package order

import (
	"context"

	"go-arch-template/internal/api/domain/order"
)

type Repository interface {
	Save(ctx context.Context, o *order.Order) error
	FindByID(ctx context.Context, id string) (*order.Order, error)
	FindByUserID(ctx context.Context, userID string) ([]*order.Order, error)
	Update(ctx context.Context, o *order.Order) error
	Delete(ctx context.Context, id string) error
}
