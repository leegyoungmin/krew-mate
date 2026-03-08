package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	agentv1beta "github.com/leegyoungmin/krew-mate/operator/api/v1beta"
)

type AgentMessageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *AgentMessageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	msg := &agentv1beta.AgentMessage{}
	if err := r.Get(ctx, req.NamespacedName, msg); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info(
		"Reconciling AgentMessage",
		"name", msg.Name,
		"to", msg.Spec.To,
		"from", msg.Spec.From,
	)

	return ctrl.Result{}, nil
}

func (r *AgentMessageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&agentv1beta.AgentMessage{}).
		Complete(r)
}
