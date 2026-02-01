package order

import (
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
)

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    Status      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func NewOrder(userID string, items []OrderItem) (*Order, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if len(items) == 0 {
		return nil, ErrEmptyItems
	}

	total := 0.0
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		if item.Price < 0 {
			return nil, ErrInvalidPrice
		}
		total += item.Price * float64(item.Quantity)
	}

	now := time.Now()
	return &Order{
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (o *Order) Confirm() error {
	if o.Status != StatusPending {
		return ErrInvalidStatus
	}
	o.Status = StatusConfirmed
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == StatusCancelled {
		return ErrInvalidStatus
	}
	o.Status = StatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}
