package validator

import (
	"context"
	"fmt"

	"go-arch-template/internal/api/domain/order"
)

// OrderItemRequest order item structure for validation
type OrderItemRequest struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// CreateOrderRequest request structure for validation
type CreateOrderRequest struct {
	UserID string             `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

// OrderRequestValidator validator for Order requests
type OrderRequestValidator struct{}

func NewOrderRequestValidator() *OrderRequestValidator {
	return &OrderRequestValidator{}
}

// ValidateCreateRequest validates order creation request
func (v *OrderRequestValidator) ValidateCreateRequest(ctx context.Context, req *CreateOrderRequest) error {
	var errs ValidationErrors

	// User ID validation
	if req.UserID == "" {
		errs.Add("user_id", "user_id is required")
	}

	// Items validation
	if len(req.Items) == 0 {
		errs.Add("items", "at least one item is required")
	}

	// Validate each item
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
