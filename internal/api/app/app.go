package app

import (
	"os"

	companyRepo "go-arch-template/internal/api/repository/company"
	orderRepo "go-arch-template/internal/api/repository/order"
	userRepo "go-arch-template/internal/api/repository/user"
	"go-arch-template/internal/api/service"
	"go-arch-template/internal/api/storage"
	"go-arch-template/internal/api/usecase"
)

type Application struct {
	// Можно добавить поля для управления приложением
}

type Repositories struct {
	CompanyRepository companyRepo.Repository
	UserRepository    userRepo.Repository
	OrderRepository   orderRepo.Repository
}

func Run() error {
	// env sections
	env, err := prepareEnv()
	if err != nil {
		return err
	}

	// storage sections
	stor, err := storage.PrepareStorage(env)
	if err != nil {
		return err
	}
	_ = stor // Используем моки, поэтому storage пока не нужен

	// repositories sections
	repo, err := prepareRepository()
	if err != nil {
		return err
	}

	// usecases section
	//api usecases
	companyUseCase, err := prepareCompanyUseCase(repo.CompanyRepository)
	if err != nil {
		return err
	}
	orderUseCase, err := prepareOrderUseCase(repo.OrderRepository, repo.UserRepository)
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

	// services section
	apiService, err := service.PrepareAPIService(
		env,
		companyUseCase,
		userUseCase,
		orderUseCase,
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

func prepareEnv() (*service.Env, error) {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}
	return &service.Env{
		HTTPPort: port,
	}, nil
}

func prepareRepository() (*Repositories, error) {
	return &Repositories{
		CompanyRepository: companyRepo.NewMockRepository(),
		UserRepository:    userRepo.NewMockRepository(),
		OrderRepository:   orderRepo.NewMockRepository(),
	}, nil
}

func prepareCompanyUseCase(repo companyRepo.Repository) (*usecase.CompanyUseCase, error) {
	return usecase.NewCompanyUseCase(repo), nil
}

func prepareOrderUseCase(orderRepo orderRepo.Repository, userRepo userRepo.Repository) (*usecase.OrderUseCase, error) {
	return usecase.NewOrderUseCase(orderRepo, userRepo), nil
}

func prepareUserUseCase(repo userRepo.Repository) (*usecase.UserUseCase, error) {
	return usecase.NewUserUseCase(repo), nil
}

func prepareEmailCheckerUseCase(env *service.Env) (interface{}, error) {
	// Мок для email checker usecase
	return struct{}{}, nil
}
