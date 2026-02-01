package app

import (
	"context"

	"go-arch-template/internal/api/env"
	"go-arch-template/internal/api/infrastructure/local/log"
	"go-arch-template/internal/api/infrastructure/local/trace"
	"go-arch-template/internal/api/integration"
	"go-arch-template/internal/api/repository"
	"go-arch-template/internal/api/service"
	"go-arch-template/internal/api/storage"
	"go-arch-template/internal/api/usecase/email_checker"

	companyUseCase "go-arch-template/internal/api/usecase/company"
	orderUseCase "go-arch-template/internal/api/usecase/order"
	userUseCase "go-arch-template/internal/api/usecase/user"
)

type Application struct {
	// Можно добавить поля для управления приложением
}

func Run() error {
	// env sections
	env, err := env.PrepareEnv()
	if err != nil {
		return err
	}

	// observability sections
	logger, err := log.NewLogger()
	if err != nil {
		// Fallback на простой logger
		logger = log.NewFallbackLogger()
	}

	tracer, err := trace.NewTracer("go-arch-template")
	if err != nil {
		// Fallback на noop tracer
		tracer = trace.NewNoopTracer()
	}
	defer tracer.Shutdown(context.Background())

	// storage sections
	storages, err := storage.PrepareStorage(env)
	if err != nil {
		return err
	}

	// integrations sections
	integrations, err := integration.PrepareIntegration(env)
	if err != nil {
		return err
	}

	// repositories sections
	repo, err := repository.PrepareRepository(storages)
	if err != nil {
		return err
	}

	// usecases section
	//api usecases
	companyUC, err := companyUseCase.PrepareCompanyUseCase(repo.CompanyRepository, integrations.CompanyIntegration, logger, tracer)
	if err != nil {
		return err
	}
	orderUC, err := orderUseCase.PrepareOrderUseCase(repo.OrderRepository, repo.UserRepository, integrations.BillingIntegration, logger, tracer)
	if err != nil {
		return err
	}
	userUC, err := userUseCase.PrepareUserUseCase(repo.UserRepository, integrations.CompanyIntegration, logger, tracer)
	if err != nil {
		return err
	}

	//jobs usecases
	emailCheckerUseCase, err := email_checker.PrepareEmailCheckerUseCase(env, integrations.CompanyIntegration, logger, tracer)
	if err != nil {
		return err
	}

	// services section
	apiService, err := service.PrepareAPIService(
		env,
		companyUC,
		userUC,
		orderUC,
	)
	if err != nil {
		return err
	}

	jobsService, err := service.PrepareJobsService(
		env,
		emailCheckerUseCase,
	)
	if err != nil {
		return err
	}

	cdcService, err := service.PrepareCDCService(env)
	if err != nil {
		return err
	}

	// run app section
	err = service.RunServices(
		apiService,
		jobsService,
		cdcService,
	)
	if err != nil {
		return err
	}

	return nil
}
