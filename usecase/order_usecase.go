package usecase

import (
	"context"
	"errors"

	"go-arch-template/domain/entity"
	"go-arch-template/ports/eventhandler"
	"go-arch-template/ports/repository"
)

type CreateOrderCommand struct {
	UserID string
	Items  []OrderItemRequest
}

type OrderItemRequest struct {
	ProductID string
	Quantity  int
}

type CreateOrderResponse struct {
	OrderID string
	Total   string
	Status  string
}

type OrderUseCase struct {
	orderRepo    repository.OrderRepository
	userRepo     repository.UserRepository
	productRepo  repository.ProductProductRepository
	eventHandler eventhandler.EventHandler
}

func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
	productRepo repository.ProductRepository,
	eventHandler eventhandler.EventHandler,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:    orderRepo,
		userRepo:     userRepo,
		productRepo:  productRepo,
		eventHandler: eventHandler,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*CreateOrderResponse, error) {
	// 1. Валидация пользователя
	userExists, err := uc.userRepo.Exists(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	// 2. Получение информации о товарах
	productIDs := make([]string, len(cmd.Items))
	for i, item := range cmd.Items {
		productIDs[i] = item.ProductID
	}

	products, err := uc.productRepo.FindByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	// 3. Создание Order Items
	orderItems := make([]entity.OrderItem, len(cmd.Items))
	for i, item := range cmd.Items {
		product, exists := products[item.ProductID]
		if !exists {
			return nil, errors.New("product not found: " + item.ProductID)
		}

		orderItems[i] = entity.OrderItem{
			ProductID: product.ID,
			Name:      product.Name,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
	}

	// 4. Создание заказа
	order, err := entity.NewOrder(cmd.UserID, orderItems)
	if err != nil {
		return nil, err
	}

	// 5. Сохранение
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	// 6. Обработка доменных событий
	events := order.GetEvents()
	for _, e := range events {
		if err := uc.eventHandler.Handle(ctx, e); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			// В реальном приложении нужно ретраи и dead letter queue
			uc.eventHandler.HandleError(ctx, e, err)
		}
	}

	return &CreateOrderResponse{
		OrderID: order.ID,
		Total:   order.Total.String(),
		Status:  string(order.Status),
	}, nil
}

func (uc *OrderUseCase) ConfirmOrder(ctx context.Context, orderID string) error {
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := order.Confirm(); err != nil {
		return err
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return err
	}

	// Обработка событий подтверждения
	events := order.GetEvents()
	for _, e := range events {
		uc.eventHandler.Handle(ctx, e)
	}

	return nil
}
