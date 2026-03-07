package podbuilder

import (
	"fmt"

	agentv1beta "github.com/leegyoungmin/krew-mate/api/v1beta"
)

const (
	AgentRuntimeImage = "ghcr.io/leegyoungmin/krew-mate/agent:latest"
	ContainerName     = "agent-runtime"
)

type Builder struct {
	team *agentv1beta.AgentTeam
}

func New(team *agentv1beta.AgentTeam) *Builder {
	return &Builder{team: team}
}

func (b *Builder) BuildOrchestratorPod(role agentv1beta.RoleSpec) *corev1.Pod {
	podName := fmt.Sprintf("%s-%s", b.team.Name, role.Name)
	return b.buildPod(podName, role)
}

func (b *Builder) BuildWorkerPod(role agentv1alpha1.RoleSpec, index int) *corev1.Pod {
	podName := fmt.Sprintf("%s-%s-%d", b.team.Name, role.Name, index)
	return b.buildPod(podName, role, "worker")
}

func (b *Builder) buildContainer(role agentv1beta.RoleSpec) *corev1.Container {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
			Namespace: b.team.Namespace,
			Labels: map[string]string{
				agentv1beta.LabelTeam:  	  b.team.Name,
				agentv1beta.LabelRole: 		  role.Name,
				agentv1beta.LabelAgentName:   podName,
				"app.kubernetes.io/part-of": "krew-mate",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: agentv1beta.GroupVersion.String(),
					Kind:       "AgentTeam",
					Name:       b.team.Name,
					UID:        b.team.UID,
					Controller: boolPtr(true),
					BlockOwnerDeletion: boolPtr(true),
				},
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			ServiceAccountName: fmt.Sprintf("%s-agent", b.team.Name),
			Containers: []corev1.Container{
				b.buildContainer(podName, role, agentRole),
			}
		},
	}
	return pod
}

func (b *Builder) buildContainer(podName string, role agentv1beta.RoleSpec, agentRole string) corev1.Container {
	return corev1.Container{
		Name: ContainerName,
		Image: AgentRuntimeImage,
		Env: []corev1.EnvVar{
			{Name: "AGENT_ROLE", Value: agentRole},
			{Name: "AGENT_NAME", Value: podName},
			{Name: "AGENT_TEAM", Value: b.team.Name},
			{Name: "AGENT_NAMESPACE", Value: b.team.Namespace},
			{Name: "LLM_PROVIDER", Value: role.LLM.Provider},
			{Name: "LLM_MODEL", Value: role.LLM.Model},
		},
		EnvFron: []corev1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LabelObjectReference: corev1.LabelObjectReference{
						Name: role.LLM.CredentialRef.SecretName,
					}
				}
			}
		},
		Resource: corev1.ResourceRequirements{

		},
	}
}

func AgentName(teamName, roleName string, index int) string {
	return fmt.Sprintf("%s-%s-%d", teamName, roleName, index)
}

func OrchestratorName(teamName, roleName string) string {
	return fmt.Sprintf("%s-%s", teamName, roleName)
}

func boolPtr(b bool) *bool {
	return &b
}
