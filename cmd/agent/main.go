package main

import (
	"os"

	"k8s.io/client-go/informers/coordination"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type AgentConfig struct {
	Role        string
	Name        string
	TeamName    string
	Namespace   string
	LLMProvider string
	LLMModel    string
	LLMAPIKey   string
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zap.Options{
		Development: true,
	})))
	log := ctrl.Log.WithName("agent")

	cfg, err := loadConfig()
	if err != nil {
		log.Error(err, "Failed load Agent")
		os.Exit(1)
	}

	log.Info(
		"Agent Started",
		"role", cfg.Role,
		"name", cfg.Name,
		"team", cfg.TeamName,
		"provider", cfg.LLMProvider,
		"model", cfg.LLMModel,
	)

	k8sClient, err := buildK8sClient()
	if err != nil {
		log.Error(err, "Failed to build K8s client")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM
	)
	defer cancel()

	switch cfg.Role {
	case "orchestrator":
		log.Info("Orchestrator role detected")
		if err := runOrchestrator(ctx, k8sClient, cfg); err != nil {
			log.Error(err, "Failed to run orchestrator")
			os.Exit(1)
		}
	case "worker":
		log.Info("Worker role detected")
		if err := runWorker(ctx, k8sClient, cfg); err != nil {
			log.Error(err, "Failed to run worker")
			os.Exit(1)
		}
	default:
		log.Error(nil, "Unknown role", "role", cfg.Role)
		os.Exit(1)
	}
}


func loadConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{
		Role:        os.Getenv("AGENT_ROLE"),
		Name:        os.Getenv("AGENT_NAME"),
		TeamName:    os.Getenv("AGENT_TEAM"),
		Namespace:   os.Getenv("AGENT_NAMESPACE"),
		LLMProvider: os.Getenv("LLM_PROVIDER"),
		LLMModel:    os.Getenv("LLM_MODEL"),
		LLMAPIKey:   os.Getenv("LLM_API_KEY"),
	}

	required := map[string]string{
		"AGENT_ROLE":      cfg.Role,
		"AGENT_NAME":      cfg.Name,
		"AGENT_TEAM":      cfg.TeamName,
		"AGENT_NAMESPACE": cfg.Namespace,
		"LLM_PROVIDER":    cfg.LLMProvider,
		"LLM_MODEL":       cfg.LLMModel,
	}

	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf("missing %s", key)
		}
	}

	if cfg.Namespace == "" {
		if ns, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
			cfg.Namespace = string(ns)
		}
	}

	return cfg, nil
}

func buildK8sClient() (client.Client, error) {
	scheme := runtime.NewScheme()

	_ = clientgoscheme.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = coordination.AddToScheme(scheme)
	_ = agentv1beta.AddToScheme(scheme)

	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return k8sClient, nil
}

func runOrchestrator(ctx context.Context, k8sClient client.Client) error {
	log := ctrl.Log.WithName("orchestrator").WithValues("name", cfg.Name)

	<-ctx.Done()
	return nil
}

func runWorker(ctx context.Context, k8sClient client.Client) error {
	log := ctrl.Log.WithName("worker").WithValues("name", cfg.Name)

	<-ctx.Done()
	return nil
}
