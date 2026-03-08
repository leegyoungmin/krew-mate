package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	agentv1beta "github.com/leegyoungmin/krew-mate/operator/api/v1beta"
)

type AgentTaskReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *AgentTaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	task := &agentv1beta.AgentTask{}
	if err := r.Get(ctx, req.NamespacedName, task); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info(
		"AgentTask reconciled",
		"name", task.Name,
		"phase", task.Status.Phase,
	)

	return ctrl.Result{}, nil
}

func (r *AgentTaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&agentv1beta.AgentTask{}).
		Complete(r)
}
