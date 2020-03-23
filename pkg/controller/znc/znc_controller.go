package znc

import (
	"context"
	"fmt"
	"github.com/mitchellh/hashstructure"
	"reflect"
	"strconv"

	zncv1 "znc-operator/pkg/apis/znc/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_znc")

// Add creates a new ZNC Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileZNC{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("znc-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ZNC
	err = c.Watch(&source.Kind{Type: &zncv1.ZNC{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner ZNC
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &zncv1.ZNC{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileZNC implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileZNC{}

// ReconcileZNC reconciles a ZNC object
type ReconcileZNC struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ZNC object and makes changes based on the state read
// and what is in the ZNC.Spec
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileZNC) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ZNC")

	// Fetch the ZNC instance
	instance := &zncv1.ZNC{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var cfgHash uint64
	{
		configMap, err := newConfigMapForCR(instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err := controllerutil.SetControllerReference(instance, configMap, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		found := &corev1.ConfigMap{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, found)
		if err == nil {
			if !reflect.DeepEqual(configMap.Data, found.Data) {
				reqLogger.Info("Updating ZNC ConfigMap")
				err = r.client.Update(context.TODO(), configMap)
				if err != nil {
					return reconcile.Result{}, err
				}
				if cfgHash, err = hashstructure.Hash(configMap.Data, nil); err != nil {
					return reconcile.Result{}, err
				}
			}
		} else {
			if errors.IsNotFound(err) {
				reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name)
				if err = r.client.Create(context.TODO(), configMap); err != nil {
					return reconcile.Result{}, err
				}

			}
		}
		if cfgHash, err = hashstructure.Hash(found.Data, nil); err != nil {
			return reconcile.Result{}, err
		}
	}

	{
		cfgHashAsString := strconv.FormatUint(cfgHash, 10)
		pod := newPodForCR(instance, cfgHashAsString)
		if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
			return reconcile.Result{}, err
		}
		found := &corev1.Pod{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
		if err == nil {
			annotations := found.GetAnnotations()
			if annotations == nil || len(annotations) == 0 || annotations["config.znc.in/checksum"] != cfgHashAsString {
				// If no annotations are present or if the checksum doesn't match, we need to delete the Pod and re-create it.
				reqLogger.Info(fmt.Sprintf("Configuration updated (old checksum: %s, new checksum: %s, deleting ZNC pod", annotations["config.znc.in/checksum"], cfgHashAsString))
				return reconcile.Result{}, r.client.Delete(context.TODO(), found)
			}
		} else {
			if errors.IsNotFound(err) {
				reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
				err = r.client.Create(context.TODO(), pod)
				if err != nil {
					return reconcile.Result{}, err
				}
			} else {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}

func newConfigMapForCR(cr *zncv1.ZNC) (configMap *corev1.ConfigMap, err error) {
	zncConf, err := RenderConfiguration(&cr.Spec)
	if err != nil {
		return nil, err
	}
	configMap = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/instance":   cr.Name,
				"app.kubernetes.io/managed-by": "znc-operator",
				"app.kubernetes.io/name":       "znc",
			},
		},
		Data: map[string]string{
			"znc.conf": zncConf,
		},
	}
	return configMap, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *zncv1.ZNC, cfgHash string) *corev1.Pod {
	labels := map[string]string{
		"app.kubernetes.io/instance":   cr.Name,
		"app.kubernetes.io/managed-by": "znc-operator",
		"app.kubernetes.io/name":       "znc",
	}
	args := []string{
		"--foreground",
	}
	if cr.Spec.Debug {
		args = append(args, "--debug")
	}
	allowPrivilegeEscalation := false
	readOnlyRootFileSystem := true
	runAsNonRoot := true
	var userID int64 = 65534
	var groupID int64 = 65534
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
			Annotations: map[string]string{
				"config.znc.in/checksum": cfgHash,
			},
		},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{
				{
					Command: []string{
						"/bin/sh",
						"-c",
					},
					Args: []string{
						"mkdir -p /znc-data/configs && cp /znc-config-src/znc.conf /znc-data/configs/znc.conf",
					},
					Image:           "docker.io/alpine:3.11.3",
					ImagePullPolicy: corev1.PullIfNotPresent,
					Name:            "copy-config",
					SecurityContext: &corev1.SecurityContext{
						AllowPrivilegeEscalation: &allowPrivilegeEscalation,
						ReadOnlyRootFilesystem:   &readOnlyRootFileSystem,
						RunAsNonRoot:             &runAsNonRoot,
						RunAsUser:                &userID,
						RunAsGroup:               &groupID,
						Capabilities: &corev1.Capabilities{
							Drop: []corev1.Capability{
								"ALL",
							},
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "znc-config-src",
							MountPath: "/znc-config-src",
							ReadOnly:  true,
						},
						{
							Name:      "znc-data",
							MountPath: "/znc-data",
							ReadOnly:  false,
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Args:            args,
					Image:           fmt.Sprintf("docker.io/library/znc:%s", cr.Spec.GetVersion()),
					ImagePullPolicy: corev1.PullIfNotPresent,
					Name:            "znc",
					Ports: []corev1.ContainerPort{
						{
							Name:          "irc",
							ContainerPort: 6667,
							Protocol:      corev1.ProtocolTCP,
						},
						{
							Name:          "web",
							ContainerPort: 8080,
							Protocol:      corev1.ProtocolTCP,
						},
					},
					SecurityContext: &corev1.SecurityContext{
						AllowPrivilegeEscalation: &allowPrivilegeEscalation,
						ReadOnlyRootFilesystem:   &readOnlyRootFileSystem,
						RunAsNonRoot:             &runAsNonRoot,
						RunAsUser:                &userID,
						RunAsGroup:               &groupID,
						Capabilities: &corev1.Capabilities{
							Drop: []corev1.Capability{
								"ALL",
							},
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "znc-data",
							MountPath: "/znc-data",
							ReadOnly:  false,
						},
					},
				},
			},
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    &userID,
				RunAsGroup:   &groupID,
				RunAsNonRoot: &runAsNonRoot,
				FSGroup:      &groupID,
			},
			Volumes: []corev1.Volume{
				{
					Name: "znc-config-src",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: cr.Name,
							},
						},
					},
				},
				{
					Name: "znc-data",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}
}
