package user

import (
	"context"

	"go-arch-template/internal/api/domain/user"
)

type Repository interface {
	Create(ctx context.Context, u *user.User) error
	FindByID(ctx context.Context, id string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Exists(ctx context.Context, id string) (bool, error)
	FindAll(ctx context.Context) ([]*user.User, error)
	Update(ctx context.Context, u *user.User) error
	Delete(ctx context.Context, id string) error
}
