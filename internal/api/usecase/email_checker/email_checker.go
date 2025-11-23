package email_checker

import (
	"go-arch-template/internal/api/env"
	"go-arch-template/internal/api/infrastructure/internal/log"
	"go-arch-template/internal/api/infrastructure/internal/trace"
	"go-arch-template/internal/api/integration"
)

func PrepareEmailCheckerUseCase(
	env *env.Env,
	companyIntegration integration.CompanyIntegration,
	logger log.Logger,
	tracer trace.Tracer,
) (interface{}, error) {
	// Мок для email checker usecase
	return struct{}{}, nil
}
