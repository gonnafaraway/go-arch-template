package app

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-arch-template/infrastructure/eventhandler"
	"go-arch-template/infrastructure/repository"
	"go-arch-template/transport/http"
	"go-arch-template/usecase"
)

type Application struct {
	OrderHandler *http.OrderHandler
	// другие хендлеры...
}

func NewApplication() (*Application, error) {
	// 1. Инициализация инфраструктуры
	mongoClient, err := initMongoDB()
	if err != nil {
		return nil, err
	}

	db := mongoClient.Database("orders_db")

	// 2. Инициализация репозиториев
	orderRepo := repository.NewMongoOrderRepository(db)
	userRepo := repository.NewMongoUserRepository(db)
	productRepo := repository.NewMongoProductRepository(db)

	// 3. Инициализация внешних сервисов
	emailService := infrastructure.NewEmailService()
	analyticsRepo := infrastructure.NewAnalyticsRepository()

	// 4. Инициализация event handlers
	orderEventHandler := eventhandler.NewOrderEventHandler(emailService, analyticsRepo)

	// 5. Инициализация use cases
	orderUseCase := usecase.NewOrderUseCase(
		orderRepo,
		userRepo,
		productRepo,
		orderEventHandler,
	)

	// 6. Инициализация хендлеров
	orderHandler := http.NewOrderHandler(orderUseCase)

	return &Application{
		OrderHandler: orderHandler,
	}, nil
}

func initMongoDB() (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return client, nil
}
