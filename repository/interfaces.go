package repository

import (
	"context"

	"go-arch-template/domain/entity"
	"go-arch-template/domain/specification"
)

// OrderRepository - порт для работы с заказами
type OrderRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Order, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.Order, error)
	FindBySpec(ctx context.Context, spec specification.OrderSpecification) ([]*entity.Order, error)
	Save(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
}

// UserRepository - порт для работы с пользователями
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	Exists(ctx context.Context, id string) (bool, error)
}

// ProductRepository - порт для работы с товарами
type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	FindByIDs(ctx context.Context, ids []string) (map[string]*entity.Product, error)
}
