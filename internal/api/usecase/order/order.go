package order

import (
	"context"
	"errors"
	"go-arch-template/internal/api/infrastructure/internal/log"
	"go-arch-template/internal/api/infrastructure/internal/trace"

	"go-arch-template/internal/api/domain/order"
	"go-arch-template/internal/api/integration"
	orderRepo "go-arch-template/internal/api/repository/order"
	userRepo "go-arch-template/internal/api/repository/user"
	"go-arch-template/internal/api/validator"
)

type CreateOrderCommand struct {
	UserID string             `json:"user_id"`
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
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Items     []order.OrderItem `json:"items"`
	Total     float64           `json:"total"`
	Status    string            `json:"status"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type OrderUseCase struct {
	orderRepo          orderRepo.Repository
	userRepo           userRepo.Repository
	billingIntegration integration.BillingIntegration
	logger             log.Logger
	tracer             trace.Tracer
	validators         *validator.OrderValidators
}

func NewOrderUseCase(
	orderRepo orderRepo.Repository,
	userRepo userRepo.Repository,
	billingIntegration integration.BillingIntegration,
	logger log.Logger,
	tracer trace.Tracer,
	validators *validator.OrderValidators,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:          orderRepo,
		userRepo:           userRepo,
		billingIntegration: billingIntegration,
		logger:             logger,
		tracer:             tracer,
		validators:         validators,
	}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*CreateOrderResponse, error) {
	ctx, span := uc.tracer.Start(ctx, "OrderUseCase.CreateOrder")
	defer span.End()

	uc.logger.Info(ctx, "Creating order", log.Field{Key: "user_id", Value: cmd.UserID})

	// 1. Валидация запроса
	validatorItems := make([]validator.OrderItemRequest, len(cmd.Items))
	for i, item := range cmd.Items {
		validatorItems[i] = validator.OrderItemRequest{
			ProductID: item.ProductID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
	validatorReq := &validator.CreateOrderRequest{
		UserID: cmd.UserID,
		Items:  validatorItems,
	}
	if err := uc.validators.Request.ValidateCreateRequest(ctx, validatorReq); err != nil {
		uc.logger.Warn(ctx, "Request validation failed", log.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	// 2. Валидация пользователя
	userExists, err := uc.userRepo.Exists(ctx, cmd.UserID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to check user existence", err, log.Field{Key: "user_id", Value: cmd.UserID})
		return nil, err
	}
	if !userExists {
		uc.logger.Warn(ctx, "User not found", log.Field{Key: "user_id", Value: cmd.UserID})
		return nil, errors.New("user not found")
	}

	// 3. Создание Order Items
	orderItems := make([]order.OrderItem, len(cmd.Items))
	for i, item := range cmd.Items {
		orderItems[i] = order.OrderItem{
			ProductID: item.ProductID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	// 4. Создание заказа
	o, err := order.NewOrder(cmd.UserID, orderItems)
	if err != nil {
		uc.logger.Error(ctx, "Failed to create order entity", err)
		return nil, err
	}

	// 5. Валидация доменной сущности
	if err := uc.validators.Domain.Validate(ctx, o); err != nil {
		uc.logger.Warn(ctx, "Domain validation failed", log.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	// 6. Сохранение
	if err := uc.orderRepo.Save(ctx, o); err != nil {
		uc.logger.Error(ctx, "Failed to save order", err, log.Field{Key: "order_id", Value: o.ID})
		return nil, err
	}

	// 7. Создание инвойса через billing интеграцию
	invoiceID, err := uc.billingIntegration.CreateInvoice(ctx, o.ID, o.Total, cmd.UserID)
	if err != nil {
		uc.logger.Warn(ctx, "Failed to create invoice", log.Field{Key: "order_id", Value: o.ID}, log.Field{Key: "error", Value: err.Error()})
	} else {
		uc.logger.Info(ctx, "Invoice created", log.Field{Key: "invoice_id", Value: invoiceID}, log.Field{Key: "order_id", Value: o.ID})
	}

	uc.logger.Info(ctx, "Order created successfully", log.Field{Key: "order_id", Value: o.ID}, log.Field{Key: "total", Value: o.Total})

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

func PrepareOrderUseCase(
	orderRepo orderRepo.Repository,
	userRepo userRepo.Repository,
	billingIntegration integration.BillingIntegration,
	logger log.Logger,
	tracer trace.Tracer,
) (*OrderUseCase, error) {
	validators, err := validator.PrepareOrderValidators()
	if err != nil {
		return nil, err
	}
	return NewOrderUseCase(orderRepo, userRepo, billingIntegration, logger, tracer, validators), nil
}
