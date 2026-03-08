package main

import (
	"flag"
	"os"

	agentv1beta "github.com/leegyoungmin/krew-mate/operator/api/v1beta"
	"github.com/leegyoungmin/krew-mate/operator/internal/controller"
	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = coordinationv1.AddToScheme(scheme)
	_ = agentv1beta.AddToScheme(scheme)
}

func main() {
	var (
		metricAddr           string
		probeAddr            string
		enableLeaderElection bool
	)

	flag.StringVar(&metricAddr, "metric-bind-address", ":8080", "The address to bind the metric server to")
	flag.StringVar(&probeAddr, "probe-bind-address", ":8081", "The address to bind the probe server to")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false, "Enable leader election for controller manager")

	opts := zap.Options{Development: true}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	log := ctrl.Log.WithName("manager")

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: metricAddr,
		},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "krewmate.io",
	})

	if err != nil {
		log.Error(err, "failed to create manager")
		os.Exit(1)
	}

	if err := (&controller.AgentTeamReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Error(err, "failed to setup agent team reconciler")
		os.Exit(1)
	}

	if err := (&controller.AgentTaskReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Error(err, "failed to setup agent task reconciler")
		os.Exit(1)
	}

	if err := (&controller.AgentMessageReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		log.Error(err, "failed to setup agent message reconciler")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Error(err, "failed to add healthz check")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.Error(err, "failed to add readyz check")
		os.Exit(1)
	}

	log.Info("starting manager")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "failed to start manager")
		os.Exit(1)
	}
}
