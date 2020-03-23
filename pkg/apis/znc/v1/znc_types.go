package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ZNCSpec defines the desired state of ZNC
type ZNCSpec struct {

	// Version specifies the ZNC version to run.
	// +optional
	Version string `json:"version,omitempty"`

	// Debug is used to enable debug output.
	// +optional
	// +kubebuilder:validation:Default=false
	Debug bool `json:"debug,omitempty"`

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

func (in *ZNCSpec) GetConfig() ZNCSpecConfig {
	return in.Config
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

	// LoadModules controls which modules shall be loaded.
	// +optional
	// +kubebuilder:validation:MinItems=0
	LoadModules []string `json:"loadModules,omitempty"`

	// MaxBufferSize controls the maximum playback buffer size. Only admin users can exceed the limit.
	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxBufferSize int32 `json:"maxBufferSize,omitempty"`

	// Motd specifies the list of "message of the day" lines that are sent to clients on connect via notice from *status.
	// +optional
	// +kubebuilder:validation:UniqueItems=false
	// +kubebuilder:validation:MinItems=0
	Motd []string `json:"motd,omitempty"`

	// ServerThrottle controls the  number of seconds between connect attempts to the same hostname.
	// +optional
	ServerThrottle int32 `json:"serverThrottle,omitempty"`

	// StatusPrefix controls the default prefix for status and module queries. Users can override the value.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Default=*
	StatusPrefix string `json:"statusPrefix,omitempty"`

	// Users specifies the users that are allowed to interact with this ZNC instance.
	// +optional
	// +kubebuilder:validation:MinItems=0
	Users []ZNCSpecConfigUser `json:"users,omitempty"`
}

func (in ZNCSpecConfig) GetStatusPrefix() string {
	statusPrefix := in.StatusPrefix
	if len(statusPrefix) == 0 {
		return StatusPrefixDefault
	}
	return statusPrefix
}

type ZNCSpecConfigUser struct {

	// Name specifies the user's name.
	Name string `json:"name"`

	// Admin toggles whether the user has admin rights.
	// +optional
	// +kubebuilder:validation:Default=false
	Admin bool `json:"admin,omitempty"`

	// AltNick controls the default alternate nick used if the primary nick is reserved.
	// Networks can override the value.
	// +kubebuilder:validation:MinLength=1
	AltNick string `json:"altNick"`

	// AppendTimestamp controls whether Whether timestamps are appended to buffer playback messages.
	// NOTE: Only used for clients that do not support server-time.
	// +optional
	AppendTimestamp bool `json:"appendTimestamp,omitempty"`

	// AutoClearChanBuffer controls whether hether channel buffers are automatically cleared after playback.
	// When disabled, messages are buffered even while clients are attached, and already seen messages may be repeated
	// each time clients connect.
	// +optional
	AutoClearChanBuffer bool `json:"autoClearChanBuffer,omitempty"`

	// AutoClearQueryBuffer controls whether query buffers are automatically cleared after playback.
	// When disabled, messages are buffered even while clients are attached, and already seen messages may be repeated
	// each time clients connect.
	// +optional
	AutoClearQueryBuffer bool `json:"autoClearQueryBuffer,omitempty"`

	// Buffer controls the maximum amount of lines stored for each channel or query playback buffer.
	// The buffers are stored in memory, and oldest lines are discarded when the limit is reached.
	// Only admin users can exceed the maximum buffer size specified in the global section.
	// +optional
	// +kubebuilder:validation:Minimum=0
	Buffer int32 `json:"buffer,omitempty"`

	// ChanBufferSize controls the maximum amount of lines stored for each channel playback buffer.
	// The buffers are stored in memory, and oldest lines are discarded when the limit is reached.
	// Only admin users can exceed the maximum buffer size specified in the global section.
	// +optional
	// +kubebuilder:validation:Minimum=0
	ChanBufferSize int32 `json:"chanBufferSize,omitempty"`

	// ChanModes controls the default modes ZNC sets when joining an empty channel.
	// +optional
	ChanModes string `json:"chanModes,omitempty"`

	// ClientEncoding sets the client encoding.
	// +optional
	// +kubebuilder:validation:Default=UTF-8
	ClientEncoding string `json:"clientEncoding,omitempty"`

	// Ident defines the default ident. Networks can override the value.
	// +optional
	// +kubebuilder:validation:Default=znc
	Ident string `json:"ident,omitempty"`

	// JoinTries specifies the amount of times channels are attempted to join in case of a failure eg. due to channel
	// modes +i/+k/+b.
	// +optional
	// +kubebuilder:validation:Minimum=1
	JoinTries int32 `json:"joinTries,omitempty"`

	// LoadModules controls the list of user modules loaded on ZNC startup.
	// +optional
	// +kubebuilder:validation:MinItems=0
	LoadModules []string `json:"loadModules,omitempty"`

	// MaxJoins controls the maximum number of channels ZNC joins at once. Lower the value in case getting disconnected
	// for 'Excess flood'.
	// +optional
	MaxJoins int32 `json:"maxJoins,omitempty"`

	// MaxQueryBuffers controls the maximum number of query buffers that are stored. 0 is unlimited.
	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxQueryBuffers int32 `json:"maxQueryBuffers,omitempty"`

	// MultiClients controls whether multiple clients are allowed to connect simultaneously.
	// +optional
	MultiClients bool `json:"multiClients,omitempty"`

	// Nick controls the default primary nick. Networks can override the value.
	// +kubebuilder:validation:MinLength=1
	Nick string `json:"nick"`

	// NoTrafficTimeout specifies how much time ZNC waits (in seconds) until it receives something from network or
	// declares the connection timeout. This happens after attempts to ping the peer.
	// +optional
	// +kubebuilder:validation:Minimum=0
	NoTrafficTimeout int32 `json:"noTrafficTimeout,omitempty"`

	// Prependtimestamp controls whether timestamps are prepended to buffer playback messages.
	// NOTE: Only used for clients that do not support server-time.
	// +optional
	PrependTimestamp bool `json:"prependTimestamp,omitempty"`

	// QueryBufferSize controls the maximum amount of lines stored for each query playback buffer.
	// The buffers are stored in memory, and oldest lines are discarded when the limit is reached.
	// Only admin users can exceed the maximum buffer size specified in the global section.
	// +optional
	// +kubebuilder:validation:Minimum=0
	QueryBufferSize int32 `json:"queryBufferSize,omitempty"`

	// QuitMsg specifies the default quit message ZNC uses when disconnecting or shutting down.
	// Networks can override the value.
	// +optional
	QuitMsg string `json:"quitMsg,omitempty"`

	// RealName specifies the default real name. Networks can override the value.
	// +optional
	RealName string `json:"realName,omitempty"`

	// StatusPrefix controls the prefix for status and module queries.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Default=*
	StatusPrefix string `json:"statusPrefix,omitempty"`

	// TimestampFormat controls the format of the timestamps used in buffer playback messages.
	// NOTE: Only used for clients that do not support server-time.
	// +optional
	TimestampFormat string `json:"timestampFormat,omitempty"`

	// Timezone controls the timezone used for timestamps in buffer playback messages.
	// NOTE: Only used for clients that do not support server-time
	// +optional
	Timezone string `json:"timezone,omitempty"`

	// Passwords represents the definition of a password, used by clients to connect to ZNC.
	Pass string `json:"pass"`

	// Networks specifies a list of IRC networks to connect to.
	Networks []ZNCSpecConfigUserNetwork `json:"networks"`
}

func (in ZNCSpecConfigUser) GetChanModes() string {
	chanModes := in.ChanModes
	if len(chanModes) == 0 {
		return "+stn"
	}
	return chanModes
}

func (in ZNCSpecConfigUser) GetClientEncoding() string {
	clientEncoding := in.ClientEncoding
	if len(clientEncoding) == 0 {
		return "UTF-8"
	}
	return clientEncoding
}

func (in ZNCSpecConfigUser) GetIdent() string {
	ident := in.Ident
	if len(ident) == 0 {
		ident = IdentDefault
	}
	return ident
}

func (in ZNCSpecConfigUser) GetStatusPrefix() string {
	statusPrefix := in.StatusPrefix
	if len(statusPrefix) == 0 {
		return StatusPrefixDefault
	}
	return statusPrefix
}

// TODO Altough the ZNC documentation states, that separate <Password> blocks should work, this doesn't. :-(
type ZNCSpecConfigUserPass struct {

	// Hash is the hash of a salted password.
	Hash string `json:"hash"`

	// Method is the password hashing method.
	// +optional
	// +kubebuilder:validation:Default=sha256
	Method string `json:"method,omitempty"`

	// Salt is a random set of 20 characters for salting the password.
	Salt string `json:"salt,omitempty"`
}

func (in ZNCSpecConfigUserPass) GetMethod() string {
	method := in.Method
	if len(method) == 0 {
		method = PassHashMethodDefault
	}
	return method
}

type ZNCSpecConfigUserNetwork struct {

	// Name specifies the network name.
	Name string `json:"name"`

	// AltNick specifies an optional network specific alternate nick used if the primary nick is reserved.
	// +optional
	// +kubebuilder:validation:MinLength=1
	AltNick string `json:"altNick,omitempty"`

	// Encoding sets an optional network specific encoding.
	// +optional
	Encoding string `json:"encoding,omitempty"`

	// Ident defines an optional network specific ident.
	// +optional
	Ident string `json:"ident,omitempty"`

	// IRCConnectEnabled specifies whether the network is enabled ie. connects to IRC.
	// +optional
	// +kubebuilder:validation:Default=true
	IRCConnectEnabled bool `json:"ircConnectEnabled,omitempty"`

	// JoinDelay specifies the delay in seconds, until channels are joined after getting connected.
	// +optional
	JoinDelay int32 `json:"joinDelay,omitempty"`

	// LoadModules controls the list of network modules loaded on ZNC startup.
	// +optional
	// +kubebuilder:validation:MinItems=0
	LoadModules []string `json:"loadModules,omitempty"`

	// Nick specifies an optional network specific primary nick.
	// +optional
	// +kubebuilder:validation:MinLength=1
	Nick string `json:"nick,omitempty"`

	// QuitMsg specifies aA optional network specific quit message ZNC uses when disconnecting or shutting down.
	// +optional
	QuitMsg string `json:"quitMsg,omitempty"`

	// RealName specifies an optional network specific real name.
	// +optional
	RealName string `json:"realName,omitempty"`

	// Servers specifies the list of IRC servers.
	// Prefix the port number with a '+' to enable SSL. Syntax: <host> [[+]port] [password].
	Servers []string `json:"servers"`

	// Channels specifies the channels to be joined.
	// +optional
	Channels []ZNCSpecConfigUserNetworkChan `json:"channels,omitempty"`
}

type ZNCSpecConfigUserNetworkChan struct {

	// Name specifies the channel name.
	Name string `json:"name"`

	// AutoClearChanBuffer defines whether the channel specific buffer is automatically cleared after playback.
	// +optional
	AutoClearChanBuffer bool `json:"autoClearChanBuffer,omitempty"`

	// Buffer defines the maximum amount of lines stored for the channel specific playback buffer.
	// +optional
	// +kubebuilder:validation:Minimum=0
	Buffer int32 `json:"buffer,omitempty"`

	// Detached defines whether the channel is detached. Detached channels are not visible to clients.
	// +optional
	// +kubebuilder:validation:Default=false
	Detached bool `json:"detached,omitempty"`

	// Disabled defines whether the channel is disabled. ZNC does not join disabled channels.
	// +optional
	// +kubebuilder:validation:Default=false
	Disabled bool `json:"disabled,omitempty"`

	// Key is an optional channel key.
	// +optional
	Key string `json:"key,omitempty"`

	// Modes specifies an optional set of default channel modes ZNC sets when joining an empty channel.
	Modes string `json:"modes,omitempty"`
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
