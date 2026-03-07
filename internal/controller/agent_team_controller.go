package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	agentv1beta "github.com/leegyoungmin/krew-mate/api/v1beta"
)

type AgentTeamReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *AgentTeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	team := &agentv1beta.AgentTeam{}

	if err := r.Get(ctx, req.NamespacedName, team); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info(
		"AgentTeam reconcile",
		"name", team.Name,
		"phase", team.Status.Phase,
	)

	return ctrl.Result{}, nil
}

func (r *AgentTeamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&agentv1beta.AgentTeam{}).
		Complete(r)
}
