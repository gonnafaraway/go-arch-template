package usecase

import (
	"context"

	"go-arch-template/internal/api/domain/company"
	companyRepo "go-arch-template/internal/api/repository/company"
	"go-arch-template/internal/api/integration"
)

type CompanyUseCase struct {
	repo              companyRepo.Repository
	companyIntegration integration.CompanyIntegration
}

func NewCompanyUseCase(repo companyRepo.Repository, companyIntegration integration.CompanyIntegration) *CompanyUseCase {
	return &CompanyUseCase{
		repo:              repo,
		companyIntegration: companyIntegration,
	}
}

type CreateCompanyRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CompanyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (uc *CompanyUseCase) CreateCompany(ctx context.Context, req CreateCompanyRequest) (*CompanyResponse, error) {
	c := company.NewCompany(req.Name, req.Email)
	if err := uc.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	
	// Синхронизация с внешним сервисом
	if err := uc.companyIntegration.SyncCompany(ctx, c.ID); err != nil {
		// Логируем ошибку, но не прерываем выполнение
		_ = err
	}
	
	return &CompanyResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *CompanyUseCase) GetCompany(ctx context.Context, id string) (*CompanyResponse, error) {
	c, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &CompanyResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *CompanyUseCase) ListCompanies(ctx context.Context) ([]*CompanyResponse, error) {
	companies, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*CompanyResponse, len(companies))
	for i, c := range companies {
		result[i] = &CompanyResponse{
			ID:        c.ID,
			Name:      c.Name,
			Email:     c.Email,
			CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return result, nil
}

