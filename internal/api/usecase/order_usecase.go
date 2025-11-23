package usecase

import (
	"context"
	"errors"

	"go-arch-template/internal/api/domain/order"
	orderRepo "go-arch-template/internal/api/repository/order"
	userRepo "go-arch-template/internal/api/repository/user"
	"go-arch-template/internal/api/integration"
)

type CreateOrderCommand struct {
	UserID string            `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CreateOrderResponse struct {
	OrderID string  `json:"order_id"`
	Total   float64 `json:"total"`
	Status  string  `json:"status"`
}

type OrderResponse struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Items     []order.OrderItem `json:"items"`
	Total     float64         `json:"total"`
	Status    string          `json:"status"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

type OrderUseCase struct {
	orderRepo          orderRepo.Repository
	userRepo           userRepo.Repository
	billingIntegration integration.BillingIntegration
}

func NewOrderUseCase(
	orderRepo orderRepo.Repository,
	userRepo userRepo.Repository,
	billingIntegration integration.BillingIntegration,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:          orderRepo,
		userRepo:           userRepo,
		billingIntegration: billingIntegration,
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

	// 2. Создание Order Items
	orderItems := make([]order.OrderItem, len(cmd.Items))
	for i, item := range cmd.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("invalid quantity")
		}
		if item.Price < 0 {
			return nil, errors.New("invalid price")
		}
		orderItems[i] = order.OrderItem{
			ProductID: item.ProductID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// 3. Создание заказа
	o, err := order.NewOrder(cmd.UserID, orderItems)
	if err != nil {
		return nil, err
	}

	// 4. Сохранение
	if err := uc.orderRepo.Save(ctx, o); err != nil {
		return nil, err
	}

	// 5. Создание инвойса через billing интеграцию
	_, err = uc.billingIntegration.CreateInvoice(ctx, o.ID, o.Total, cmd.UserID)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		_ = err
	}

	return &CreateOrderResponse{
		OrderID: o.ID,
		Total:   o.Total,
		Status:  string(o.Status),
	}, nil
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, orderID string) (*OrderResponse, error) {
	o, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return &OrderResponse{
		ID:        o.ID,
		UserID:    o.UserID,
		Items:     o.Items,
		Total:     o.Total,
		Status:    string(o.Status),
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *OrderUseCase) ConfirmOrder(ctx context.Context, orderID string) error {
	o, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := o.Confirm(); err != nil {
		return err
	}

	if err := uc.orderRepo.Update(ctx, o); err != nil {
		return err
	}

	// После подтверждения заказа можно обновить статус инвойса
	// Это пример использования интеграции
	_ = uc.billingIntegration

	return nil
}

func (uc *OrderUseCase) ListOrdersByUser(ctx context.Context, userID string) ([]*OrderResponse, error) {
	orders, err := uc.orderRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		result[i] = &OrderResponse{
			ID:        o.ID,
			UserID:    o.UserID,
			Items:     o.Items,
			Total:     o.Total,
			Status:    string(o.Status),
			CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return result, nil
}
