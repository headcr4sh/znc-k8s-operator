package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ZNCSpec defines the desired state of ZNC
type ZNCSpec struct {
	// Version specifies the ZNC version to run.
	// +optional
	Version string `json:"version,omitempty"`
	// ZNSSpecConfig is the configuration used by the ZNC instance.
	Config ZNCSpecConfig `json:"config,omitempty"`
}

func (in *ZNCSpec) GetVersion() string {
	version := in.Version
	if len(version) == 0 {
		version = VersionDefault
	}
	return version
}

type ZNCSpecConfig struct {

	// AnonIPLimit is the limit of anonymous unidentified connections per IP.
	// +optional
	// +kubebuilder:validation:Minimum=0
	AnonIPLimit int `json:"anonIPLimit,omitempty"`

	// ConnectDelay is the number of seconds every IRC connection is delayed. IRC servers may refuse a connection when reconnecting too fast. NOTE: Affects connections between ZNC and IRC servers; not connections between IRC clients and ZNC.
	// +optional
	// +kubebuilder:validation:Minimum=0
	ConnectDelay int `json:"connectDelay,omitempty"`

	// HideVersion controls whether the version number is hidden from the web interface and CTCP VERSION replies.
	// +optional
	HideVersion bool `json:"hideVersion,omitempty"`
}

// ZNCStatus defines the observed state of ZNC
type ZNCStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZNC is the Schema for the zncs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=zncs,scope=Namespaced
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="Version of this ZNC instance"
type ZNC struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZNCSpec   `json:"spec,omitempty"`
	Status ZNCStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZNCList contains a list of ZNC
type ZNCList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ZNC `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ZNC{}, &ZNCList{})
}
