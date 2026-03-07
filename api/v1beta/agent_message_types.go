package v1beta

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type AgentMessage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentMessageSpec   `json:"spec,omitempty"`
	Status AgentMessageStatus `json:"status,omitempty"`
}

type AgentMessageSpec struct {
	TeamRef    TeamRef        `json:"teamRef"`
	From       string         `json:"from"`
	To         string         `json:"to"`
	Content    string         `json:"content"`
	ActionHint ActionHintType `json:"actionHint,omitempty"`
}

type ActionHintType string

const (
	ActionHintTaskComplete ActionHintType = "TaskComplete"
	ActionHintTaskFailed   ActionHintType = "TaskFailed"
	ActionHintRequestHelp  ActionHintType = "RequestHelp"
	ActionHintStatusUpdate ActionHintType = "StatusUpdate"
	ActionHintNone         ActionHintType = ""
)

type AgentMessageStatus struct {
	Delivered   bool         `json:"delivered,omitempty"`
	DeliveredAt *metav1.Time `json:"deliveredAt,omitempty"`
	ReadAt      *metav1.Time `json:"readAt,omitempty"`
}

const (
	MessageBroadcast = "broadcast"
)

const (
	AnnotationTTLSeconds = "krewmate.io/ttl-seconds"
)

const (
	LabelTeam      = "krewmate.io/team"
	LabelRole      = "krewmate.io/role"
	LabelAgentName = "krewmate.io/agent-name"
)

// +kubebuilder:object:root=true
type AgentMessageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AgentMessage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AgentMessage{}, &AgentMessageList{})
}
