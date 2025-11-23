package app

import (
	"go-arch-template/internal/api/service"
	http3 "go-arch-template/internal/api/transport/http"
)

type Application struct {
	OrderHandler *http3.OrderHandler
	// другие хендлеры...
}

func Run() error {
	env, err := prepareEnv()
	if err != nil {
		return err
	}

	storage, err := prepareStorage(env)
	if err != nil {
		return err
	}

	repo, err := prepareRepository(storage)
	if err != nil {
		return err
	}

	companyUseCase, err := prepareCompanyUseCase(repo.CompanyRepository)
	if err != nil {
		return err
	}
	orderUseCase, err := prepareOrderUseCase(repo.OrderRepository)
	if err != nil {
		return err
	}
	userUseCase, err := prepareUserUseCase(repo.UserRepository)
	if err != nil {
		return err
	}

	apiService, err := service.PrepareAPIService(env)
	if err != nil {
		return err
	}

	jobsService, err := service.PrepareJobsService(env)
	if err != nil {
		return err
	}

	cdcService, err := service.PrepareCDCService(env)
	if err != nil {
		return err
	}

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
