package company

import (
	"context"

	"go-arch-template/internal/api/domain/company"
)

type Repository interface {
	Create(ctx context.Context, c *company.Company) error
	FindByID(ctx context.Context, id string) (*company.Company, error)
	FindAll(ctx context.Context) ([]*company.Company, error)
	Update(ctx context.Context, c *company.Company) error
	Delete(ctx context.Context, id string) error
}

