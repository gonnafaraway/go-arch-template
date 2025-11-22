package entity

import (
	"errors"
	"time"

	"go-arch-template/domain/vo"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderConfirmed OrderStatus = "confirmed"
	OrderCompleted OrderStatus = "completed"
	OrderCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        string
	UserID    string
	Items     []OrderItem
	Status    OrderStatus
	Total     vo.Money
	CreatedAt time.Time
	UpdatedAt time.Time

	events []interface{} // Domain Events
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     vo.Money
	Name      string
}

func NewOrder(userID string, items []OrderItem) (*Order, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	order := &Order{
		ID:        generateID(),
		UserID:    userID,
		Items:     items,
		Status:    OrderPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := order.calculateTotal(); err != nil {
		return nil, err
	}

	order.AddEvent(OrderCreatedEvent{
		OrderID:    order.ID,
		UserID:     order.UserID,
		Total:      order.Total,
		OccurredAt: time.Now(),
	})

	return order, nil
}

func (o *Order) Confirm() error {
	if o.Status != OrderPending {
		return errors.New("only pending orders can be confirmed")
	}
	o.Status = OrderConfirmed
	o.UpdatedAt = time.Now()

	o.AddEvent(OrderConfirmedEvent{
		OrderID:    o.ID,
		OccurredAt: time.Now(),
	})

	return nil
}

func (o *Order) Cancel(reason string) error {
	if o.Status == OrderCompleted || o.Status == OrderCancelled {
		return errors.New("cannot cancel completed or cancelled order")
	}
	o.Status = OrderCancelled
	o.UpdatedAt = time.Now()

	o.AddEvent(OrderCancelledEvent{
		OrderID:    o.ID,
		Reason:     reason,
		OccurredAt: time.Now(),
	})

	return nil
}

func (o *Order) calculateTotal() error {
	total, err := vo.NewMoney(0, "USD") // базовая валюта
	if err != nil {
		return err
	}

	for _, item := range o.Items {
		itemTotal, err := item.Price.Multiply(item.Quantity)
		if err != nil {
			return err
		}
		total, err = total.Add(itemTotal)
		if err != nil {
			return err
		}
	}

	o.Total = total
	return nil
}

func (o *Order) AddEvent(event interface{}) {
	o.events = append(o.events, event)
}

func (o *Order) GetEvents() []interface{} {
	events := o.events
	o.events = nil // clear events after reading
	return events
}
