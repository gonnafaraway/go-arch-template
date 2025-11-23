package app

import (
	"context"
	"os"

	"go-arch-template/internal/api/integration"
	companyRepo "go-arch-template/internal/api/repository/company"
	orderRepo "go-arch-template/internal/api/repository/order"
	userRepo "go-arch-template/internal/api/repository/user"
	"go-arch-template/internal/api/observability"
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

	// observability sections
	logger, err := observability.NewLogger()
	if err != nil {
		// Fallback на простой logger
		logger = observability.NewFallbackLogger()
	}

	tracer, err := observability.NewTracer("go-arch-template")
	if err != nil {
		// Fallback на noop tracer
		tracer = observability.NewNoopTracer()
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
	repo, err := prepareRepository(storages)
	if err != nil {
		return err
	}

	// usecases section
	//api usecases
	companyUseCase, err := prepareCompanyUseCase(repo.CompanyRepository, integrations.CompanyIntegration, logger, tracer)
	if err != nil {
		return err
	}
	orderUseCase, err := prepareOrderUseCase(repo.OrderRepository, repo.UserRepository, integrations.BillingIntegration, logger, tracer)
	if err != nil {
		return err
	}
	userUseCase, err := prepareUserUseCase(repo.UserRepository, integrations.CompanyIntegration, logger, tracer)
	if err != nil {
		return err
	}

	//jobs usecases
	emailCheckerUseCase, err := prepareEmailCheckerUseCase(env, integrations.CompanyIntegration, logger, tracer)
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

func prepareRepository(storage *storage.Storage) (*Repositories, error) {
	repos := &Repositories{}

	// Используем реальные репозитории если storage доступен, иначе моки
	if storage.Postgres != nil {
		repos.CompanyRepository = companyRepo.NewPostgresRepository(storage.Postgres)
		// Можно добавить другие репозитории для PostgreSQL
	} else if storage.MongoDB != nil {
		repos.CompanyRepository = companyRepo.NewMongoDBRepository(storage.MongoDB)
		// Можно добавить другие репозитории для MongoDB
	} else {
		// Fallback на моки
		repos.CompanyRepository = companyRepo.NewMockRepository()
		repos.UserRepository = userRepo.NewMockRepository()
		repos.OrderRepository = orderRepo.NewMockRepository()
	}

	// Если не были установлены, используем моки
	if repos.UserRepository == nil {
		repos.UserRepository = userRepo.NewMockRepository()
	}
	if repos.OrderRepository == nil {
		repos.OrderRepository = orderRepo.NewMockRepository()
	}

	return repos, nil
}

func prepareCompanyUseCase(
	repo companyRepo.Repository,
	companyIntegration integration.CompanyIntegration,
	logger observability.Logger,
	tracer observability.Tracer,
) (*usecase.CompanyUseCase, error) {
	return usecase.NewCompanyUseCase(repo, companyIntegration, logger, tracer), nil
}

func prepareOrderUseCase(
	orderRepo orderRepo.Repository,
	userRepo userRepo.Repository,
	billingIntegration integration.BillingIntegration,
	logger observability.Logger,
	tracer observability.Tracer,
) (*usecase.OrderUseCase, error) {
	return usecase.NewOrderUseCase(orderRepo, userRepo, billingIntegration, logger, tracer), nil
}

func prepareUserUseCase(
	repo userRepo.Repository,
	companyIntegration integration.CompanyIntegration,
	logger observability.Logger,
	tracer observability.Tracer,
) (*usecase.UserUseCase, error) {
	return usecase.NewUserUseCase(repo, companyIntegration, logger, tracer), nil
}

func prepareEmailCheckerUseCase(
	env *service.Env,
	companyIntegration integration.CompanyIntegration,
	logger observability.Logger,
	tracer observability.Tracer,
) (interface{}, error) {
	// Мок для email checker usecase
	return struct{}{}, nil
}
