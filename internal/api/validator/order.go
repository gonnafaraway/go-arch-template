package validator

import (
	"context"
	"fmt"

	"go-arch-template/internal/api/domain/order"
)

// OrderItemRequest структура элемента заказа для валидации
type OrderItemRequest struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// CreateOrderRequest структура запроса для валидации
type CreateOrderRequest struct {
	UserID string             `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

// OrderRequestValidator валидатор для запросов Order
type OrderRequestValidator struct{}

func NewOrderRequestValidator() *OrderRequestValidator {
	return &OrderRequestValidator{}
}

// ValidateCreateRequest валидирует запрос на создание заказа
func (v *OrderRequestValidator) ValidateCreateRequest(ctx context.Context, req *CreateOrderRequest) error {
	var errs ValidationErrors

	// Валидация user_id
	if req.UserID == "" {
		errs.Add("user_id", "user_id is required")
	}

	// Валидация items
	if len(req.Items) == 0 {
		errs.Add("items", "at least one item is required")
	}

	// Валидация каждого item
	for i, item := range req.Items {
		if item.ProductID == "" {
			errs.Add(fmt.Sprintf("items[%d].product_id", i), "product_id is required")
		}
		if item.Name == "" {
			errs.Add(fmt.Sprintf("items[%d].name", i), "name is required")
		}
		if item.Quantity <= 0 {
			errs.Add(fmt.Sprintf("items[%d].quantity", i), "quantity must be greater than 0")
		}
		if item.Price < 0 {
			errs.Add(fmt.Sprintf("items[%d].price", i), "price must be non-negative")
		}
	}

	if errs.HasErrors() {
		return &errs
	}

	return nil
}

func PrepareOrderValidators() (*OrderValidators, error) {
	return &OrderValidators{
		Domain:  order.NewDomainValidator(),
		Request: NewOrderRequestValidator(),
	}, nil
}
