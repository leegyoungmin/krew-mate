// +kubebuilder:object:generate=true
package v1beta

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,shortName=at
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Total",type="integer",JSONPath=".status.taskSummary.total"
// +kubebuilder:printcolumn:name="Done",type="integer",JSONPath=".status.taskSummary.completed"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type AgentTeam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentTeamSpec   `json:"spec,omitempty"`
	Status AgentTeamStatus `json:"status,omitempty"`
}

type AgentTeamSpec struct {
	Goal          string             `json:"goal"`
	Roles         []RoleSpec         `json:"roles"`
	Scaling       *ScalingSpec       `json:"scaling,omitempty"`
	SharedContext *SharedContextSpec `json:"sharedContext,omitempty"`
	Cleanup       *CleanupSpec       `json:"cleanup,omitempty"`
}

type RoleSpec struct {
	Name         string   `json:"name"`
	Type         RoleType `json:"type"`
	Replicas     int32    `json:"replicas,omitempty"`
	LLM          LLMSpec  `json:"llm"`
	Capabilities []string `json:"capabilities,omitempty"`
	SystemPrompt string   `json:"systemPrompt,omitempty"`
}

type RoleType string

const (
	RoleTypeOrchestrator RoleType = "orchestrator"
	RoleTypeWorker       RoleType = "worker"
)

type LLMSpec struct {
	Provider      string        `json:"provider"`
	Model         string        `json:"model"`
	CredentialRef CredentialRef `json:"credentialRef"`
}

type CredentialRef struct {
	SecretName  string `json:"secretName"`
	APIKeyField string `json:"apiKeyField"`
}

type ScalingSpec struct {
	MaxWorkers         int32 `json:"maxWorkers,omitempty"`
	IdleTimeoutSeconds int32 `json:"idleTimeoutSeconds,omitempty"`
}

type SharedContextSpec struct {
	Backend      string `json:"backend,omitempty"`
	MaxSizeBytes int64  `json:"maxSizeBytes,omitempty"`
}

type CleanupSpec struct {
	TTLAfterFinished *int32 `json:"ttlAfterFinished,omitempty"`
}

type AgentTeamStatus struct {
	Phase       TeamPhase          `json:"phase,omitempty"`
	Roles       []RoleStatus       `json:"roles,omitempty"`
	TaskSummary *TaskSummary       `json:"taskSummary,omitempty"`
	Conditions  []metav1.Condition `json:"conditions,omitempty"`
	FinishedAt  *metav1.Time       `json:"finishedAt,omitempty"`
}

type TeamPhase string

const (
	TeamPhasePending   TeamPhase = "pending"
	TeamPhaseRunning   TeamPhase = "running"
	TeamPhaseSucceeded TeamPhase = "succeeded"
	TeamPhaseFailed    TeamPhase = "failed"
)

type RoleStatus struct {
	Name       string          `json:"name"`
	PodName    string          `json:"podName,omitempty"`
	ActivePods int             `json:"activePods,omitempty"`
	Phase      corev1.PodPhase `json:"phase,omitempty"`
}

type TaskSummary struct {
	Total      int32 `json:"total"`
	Pending    int32 `json:"pending"`
	InProgress int32 `json:"inProgress"`
	Completed  int32 `json:"completed"`
	Failed     int32 `json:"failed"`
	Blocked    int32 `json:"blocked"`
}

// +kubebuilder:object:root=true
type AgentTeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AgentTeam `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AgentTeam{}, &AgentTeamList{})
}
