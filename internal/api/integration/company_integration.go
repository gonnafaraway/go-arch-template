package integration

import (
	"context"

	usersservice "go-arch-template/internal/api/infrastructure/external/users-service"
)

type companyIntegration struct {
	usersClient usersservice.Client
}

func NewCompanyIntegration(usersClient usersservice.Client) CompanyIntegration {
	return &companyIntegration{
		usersClient: usersClient,
	}
}

func (i *companyIntegration) ValidateCompany(ctx context.Context, companyID string) (bool, error) {
	// Валидация компании через внешний сервис
	// Пока используем мок
	return true, nil
}

func (i *companyIntegration) SyncCompany(ctx context.Context, companyID string) error {
	// Синхронизация компании
	// Пока используем мок
	return nil
}

