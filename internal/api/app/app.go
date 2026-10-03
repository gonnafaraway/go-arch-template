package app

import (
	"context"

	"go-arch-template/internal/api/env"
	"go-arch-template/internal/api/integration"
	"go-arch-template/internal/api/repository"
	"go-arch-template/internal/api/service"
	"go-arch-template/internal/api/storage"
	"go-arch-template/internal/api/usecase/email_checker"

	companyUseCase "go-arch-template/internal/api/usecase/company"
	orderUseCase "go-arch-template/internal/api/usecase/order"
	userUseCase "go-arch-template/internal/api/usecase/user"
)

type Application struct{}

func Run() error {
	env, err := env.PrepareEnv()
	if err != nil {
		return err
	}

	integrations, err := integration.PrepareIntegration(env)
	if err != nil {
		return err
	}
	defer func() { _ = integrations.Trace.Shutdown(context.Background()) }()

	storages, err := storage.PrepareStorage(env)
	if err != nil {
		return err
	}

	repo, err := repository.PrepareRepository(storages)
	if err != nil {
		return err
	}

	companyUC, err := companyUseCase.PrepareCompanyUseCase(
		repo.CompanyRepository,
		integrations.Company,
		integrations.Log,
		integrations.Trace,
	)
	if err != nil {
		return err
	}

	orderUC, err := orderUseCase.PrepareOrderUseCase(
		repo.OrderRepository,
		repo.UserRepository,
		integrations.Billing,
		integrations.Log,
		integrations.Trace,
	)
	if err != nil {
		return err
	}

	userUC, err := userUseCase.PrepareUserUseCase(
		repo.UserRepository,
		integrations.Company,
		integrations.Log,
		integrations.Trace,
	)
	if err != nil {
		return err
	}

	emailCheckerUseCase, err := email_checker.PrepareEmailCheckerUseCase(
		env,
		integrations.Company,
		integrations.Log,
		integrations.Trace,
	)
	if err != nil {
		return err
	}

	apiService, err := service.PrepareAPIService(env, companyUC, userUC, orderUC)
	if err != nil {
		return err
	}

	jobsService, err := service.PrepareJobsService(env, emailCheckerUseCase)
	if err != nil {
		return err
	}

	cdcService, err := service.PrepareCDCService(env)
	if err != nil {
		return err
	}

	return service.RunServices(apiService, jobsService, cdcService)
}
