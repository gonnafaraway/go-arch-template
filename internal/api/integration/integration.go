package integration

import (
	"go-arch-template/internal/api/integration/external/billing"
	"go-arch-template/internal/api/integration/external/company"
	"go-arch-template/internal/api/integration/external/prometheus"
	"go-arch-template/internal/api/integration/external/sentry"
	usersservice "go-arch-template/internal/api/integration/external/usersservice"
	"go-arch-template/internal/api/integration/internal/log"
	"go-arch-template/internal/api/integration/internal/oauth"
	"go-arch-template/internal/api/integration/internal/trace"
)

// Aliases for packages under integration/internal
type (
	Logger = log.Client
	Field  = log.Field
	Tracer = trace.Client
)

type Integrations struct {
	Billing    billing.Client
	Company    company.Client
	Users      usersservice.Client
	OAuth      oauth.Client
	Prometheus prometheus.Client
	Sentry     sentry.Client
	Log        Logger
	Trace      Tracer
}

func PrepareIntegration(env interface{}) (*Integrations, error) {
	_ = env

	logger, err := log.NewClient()
	if err != nil {
		logger = log.NewFallbackClient()
	}

	tracer, err := trace.NewClient("go-arch-template")
	if err != nil {
		tracer = trace.NewNoopClient()
	}

	users := usersservice.NewClient()

	return &Integrations{
		Billing:    billing.NewClient(),
		Company:    company.NewClient(users),
		Users:      users,
		OAuth:      oauth.NewClient(),
		Prometheus: prometheus.NewClient(),
		Sentry:     sentry.NewClient(),
		Log:        logger,
		Trace:      tracer,
	}, nil
}
