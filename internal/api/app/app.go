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
	// env sections
	env, err := prepareEnv()
	if err != nil {
		return err
	}

	// storage sections
	storage, err := prepareStorage(env)
	if err != nil {
		return err
	}

	// repositories sections
	repo, err := prepareRepository(storage)
	if err != nil {
		return err
	}

	// usecases section
	//api usecases
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

	//jobs usecases
	emailCheckerUseCase, err := prepareEmailCheckerUseCase(env)
	if err != nil {
		return err
	}

	// domains section
	//api domains
	companyDomain, err := prepareCompanyDomain(companyUseCase)
	if err != nil {
		return err
	}
	orderDomain, err := prepareOrderDomain(orderUseCase)
	if err != nil {
		return err
	}
	userDomain, err := prepareUserDomain(userUseCase)
	if err != nil {
		return err
	}

	//jobs domains
	emailCheckerDomain, err := prepareEmailCheckerDomain(emailCheckerUseCase)
	if err != nil {
		return err
	}

	// services section
	apiService, err := service.PrepareAPIService(
		env,
		companyDomain,
		orderDomain,
		userDomain
		)
	if err != nil {
		return err
	}

	jobsService, err := service.PrepareJobsService(
		env,
		emailCheckerDomain,
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
