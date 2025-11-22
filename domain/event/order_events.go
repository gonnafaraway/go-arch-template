package event

import (
	"time"

	"go-arch-template/domain/vo"
)

type OrderCreatedEvent struct {
	OrderID    string
	UserID     string
	Total      vo.Money
	OccurredAt time.Time
}

type OrderConfirmedEvent struct {
	OrderID    string
	OccurredAt time.Time
}

type OrderCancelledEvent struct {
	OrderID    string
	Reason     string
	OccurredAt time.Time
}
