package order

import (
	"context"
	"fmt"

	"go-arch-template/internal/api/domain/order"
)

type MockRepository struct {
	orders map[string]*order.Order
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		orders: make(map[string]*order.Order),
	}
}

func (m *MockRepository) Save(ctx context.Context, o *order.Order) error {
	if o.ID == "" {
		o.ID = fmt.Sprintf("order_%d", len(m.orders)+1)
	}
	m.orders[o.ID] = o
	return nil
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*order.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, ErrNotFound
	}
	return o, nil
}

func (m *MockRepository) FindByUserID(ctx context.Context, userID string) ([]*order.Order, error) {
	result := make([]*order.Order, 0)
	for _, o := range m.orders {
		if o.UserID == userID {
			result = append(result, o)
		}
	}
	return result, nil
}

func (m *MockRepository) Update(ctx context.Context, o *order.Order) error {
	if _, ok := m.orders[o.ID]; !ok {
		return ErrNotFound
	}
	m.orders[o.ID] = o
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.orders[id]; !ok {
		return ErrNotFound
	}
	delete(m.orders, id)
	return nil
}

