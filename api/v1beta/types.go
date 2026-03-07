package v1beta

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type AgentTeamSpec struct {
	Goal  string `json:"goal,omitempty"`
	Model string `json:"model,omitempty"`
}

type AgentTeamStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
}

type AgentTeam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentTeamSpec   `json:"spec,omitempty"`
	Status AgentTeamStatus `json:"status,omitempty"`
}

type AgentTeamList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Items []AgentTeam `json:"items"`
}

var (
	SchemeGroupVersion = schema.GroupVersion{Group: "agents.krewmate.io", Version: "v1beta"}
	SchemeBuilder      = &scheme.Builder{GroupVersion: SchemeGroupVersion}
	AddToScheme        = SchemeBuilder.AddToScheme
)

func init() {
	SchemeBuilder.Register(&AgentTeam{}, &AgentTeamList{})
}

func (a *AgentTeam) DeepCopyObject() runtime.Object {
	copy := a.DeepCopy()
	return copy
}

func (a *AgentTeam) DeepCopy() *AgentTeam {
	if a == nil {
		return nil
	}

	copy := *a
	return &copy
}

func (a *AgentTeamList) DeepCopyObject() runtime.Object {
	copy := a.DeepCopy()
	return copy
}

func (a *AgentTeamList) DeepCopy() *AgentTeamList {
	if a == nil {
		return nil
	}
	copy := *a
	return &copy
}
