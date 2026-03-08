// +kubebuilder:object:generate=true
package v1beta

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=atask
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Priority",type="integer",JSONPath=".spec.priority"
// +kubebuilder:printcolumn:name="ClaimedBy",type="string",JSONPath=".status.claimedBy.agentName"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type AgentTask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentTaskSpec   `json:"spec,omitempty"`
	Status AgentTaskStatus `json:"status,omitempty"`
}

type AgentTaskSpec struct {
	TeamRef         TeamRef  `json:"teamRef"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	AssignableRoles []string `json:"assignableRoles,omitempty"`
	DependsOn       []string `json:"dependsOn,omitempty"`
	Priority        int32    `json:"priority,omitempty"`
	TimeoutSeconds  int32    `json:"timeoutSeconds,omitempty"`
}

type TeamRef struct {
	Name string `json:"name"`
}

type AgentTaskStatus struct {
	Phase       TaskPhase          `json:"phase,omitempty"`
	ClaimedBy   *ClaimInfo         `json:"claimedBy,omitempty"`
	Result      *TaskResult        `json:"result,omitempty"`
	CompletedAt *metav1.Time       `json:"completedAt,omitempty"`
	Conditions  []metav1.Condition `json:"conditions,omitempty"`
}

type TaskPhase string

const (
	TaskPhasePending    TaskPhase = "Pending"
	TaskPhaseInProgress TaskPhase = "InProgress"
	TaskPhaseSucceeded  TaskPhase = "Succeeded"
	TaskPhaseFailed     TaskPhase = "Failed"
)

type ClaimInfo struct {
	AgentName string      `json:"agentName"`
	PodName   string      `json:"podName"`
	ClaimedAt metav1.Time `json:"claimedAt"`
	LeaseRef  string      `json:"leaseRef"`
}

type TaskResult struct {
	SharedContextKey string `json:"sharedContextKey,omitempty"`
	Summary          string `json:"summary,omitempty"`
}

// +kubebuilder:object:root=true
type AgentTaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AgentTask `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AgentTask{}, &AgentTaskList{})
}
