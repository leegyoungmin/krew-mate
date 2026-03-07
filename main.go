package main

import (
	"context"
	"fmt"
	"os"

	v1beta "github.com/leegyoungmin/krew-mate/api/v1beta"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type AgentTeamReconciler struct {
	client.Client
	Schema *runtime.Scheme
}

func (r *AgentTeamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	team := &v1beta.AgentTeam{}
	if err := r.Get(ctx, req.NamespacedName, team); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Printf("AgentTeam 감지! name=%s, goal=%s\n", team.Name, team.Spec.Goal)

	return ctrl.Result{}, nil
}

func main() {
	ctrl.SetLogger(zap.New())

	scheme := runtime.NewScheme()
	v1beta.AddToScheme(scheme)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})

	if err != nil {
		fmt.Println("Manager 생성 실패:", err)
		os.Exit(1)
	}

	if err := ctrl.NewControllerManagedBy(mgr).
		For(&v1beta.AgentTeam{}).
		Complete(&AgentTeamReconciler{
			Client: mgr.GetClient(),
			Schema: mgr.GetScheme(),
		}); err != nil {
		fmt.Println("Controller 등록 실패:", err)
		os.Exit(1)
	}

	fmt.Println("Controller 시작!")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Println("Controller 실패:", err)
		os.Exit(1)
	}
}
