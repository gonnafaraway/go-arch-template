package usecase

import (
	"context"

	"go-arch-template/internal/api/domain/company"
	companyRepo "go-arch-template/internal/api/repository/company"
	"go-arch-template/internal/api/integration"
	"go-arch-template/internal/api/observability"
)

type CompanyUseCase struct {
	repo              companyRepo.Repository
	companyIntegration integration.CompanyIntegration
	logger            observability.Logger
	tracer            observability.Tracer
}

func NewCompanyUseCase(
	repo companyRepo.Repository,
	companyIntegration integration.CompanyIntegration,
	logger observability.Logger,
	tracer observability.Tracer,
) *CompanyUseCase {
	return &CompanyUseCase{
		repo:              repo,
		companyIntegration: companyIntegration,
		logger:            logger,
		tracer:            tracer,
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
	ctx, span := uc.tracer.Start(ctx, "CompanyUseCase.CreateCompany")
	defer span.End()

	uc.logger.Info(ctx, "Creating company", observability.Field{Key: "name", Value: req.Name})

	c := company.NewCompany(req.Name, req.Email)
	if err := uc.repo.Create(ctx, c); err != nil {
		uc.logger.Error(ctx, "Failed to create company", err, observability.Field{Key: "name", Value: req.Name})
		return nil, err
	}
	
	// Синхронизация с внешним сервисом
	if err := uc.companyIntegration.SyncCompany(ctx, c.ID); err != nil {
		uc.logger.Warn(ctx, "Failed to sync company", observability.Field{Key: "company_id", Value: c.ID}, observability.Field{Key: "error", Value: err.Error()})
	}
	
	uc.logger.Info(ctx, "Company created successfully", observability.Field{Key: "company_id", Value: c.ID})

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

