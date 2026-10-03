package email_checker

import (
	"go-arch-template/internal/api/env"
	"go-arch-template/internal/api/integration"
	companyIntegration "go-arch-template/internal/api/integration/external/company"
)

func PrepareEmailCheckerUseCase(
	env *env.Env,
	companies companyIntegration.Client,
	logger integration.Logger,
	tracer integration.Tracer,
) (interface{}, error) {
	_ = env
	_ = companies
	_ = logger
	_ = tracer
	return struct{}{}, nil
}
