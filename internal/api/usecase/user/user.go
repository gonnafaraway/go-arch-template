package user

import (
	"context"
	"errors"

	"go-arch-template/internal/api/domain/user"
	"go-arch-template/internal/api/infrastructure/local/log"
	"go-arch-template/internal/api/infrastructure/local/trace"
	"go-arch-template/internal/api/integration"
	"go-arch-template/internal/api/validator"

	userRepo "go-arch-template/internal/api/repository/user"
)

type UserUseCase struct {
	repo               userRepo.Repository
	companyIntegration integration.CompanyIntegration
	logger             log.Logger
	tracer             trace.Tracer
	validators         *validator.UserValidators
}

func NewUserUseCase(
	repo userRepo.Repository,
	companyIntegration integration.CompanyIntegration,
	logger log.Logger,
	tracer trace.Tracer,
	validators *validator.UserValidators,
) *UserUseCase {
	return &UserUseCase{
		repo:               repo,
		companyIntegration: companyIntegration,
		logger:             logger,
		tracer:             tracer,
		validators:         validators,
	}
}

type CreateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	ctx, span := uc.tracer.Start(ctx, "UserUseCase.CreateUser")
	defer span.End()

	uc.logger.Info(ctx, "Creating user", log.Field{Key: "email", Value: req.Email})

	// Request validation
	validatorReq := &validator.CreateUserRequest{
		Name:      req.Name,
		Email:     req.Email,
		CompanyID: req.CompanyID,
	}
	if err := uc.validators.Request.ValidateCreateRequest(ctx, validatorReq); err != nil {
		uc.logger.Warn(ctx, "Request validation failed", log.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	// Company validation through integration
	valid, err := uc.companyIntegration.ValidateCompany(ctx, req.CompanyID)
	if err != nil {
		uc.logger.Error(ctx, "Failed to validate company", err, log.Field{Key: "company_id", Value: req.CompanyID})
		return nil, err
	}
	if !valid {
		uc.logger.Warn(ctx, "Invalid company", log.Field{Key: "company_id", Value: req.CompanyID})
		return nil, errors.New("invalid company")
	}

	// Create domain entity
	u := user.NewUser(req.Name, req.Email, req.CompanyID)

	// Domain entity validation
	if err := uc.validators.Domain.Validate(ctx, u); err != nil {
		uc.logger.Warn(ctx, "Domain validation failed", log.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	if err := uc.repo.Create(ctx, u); err != nil {
		uc.logger.Error(ctx, "Failed to create user", err, log.Field{Key: "email", Value: req.Email})
		return nil, err
	}

	uc.logger.Info(ctx, "User created successfully", log.Field{Key: "user_id", Value: u.ID})

	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CompanyID: u.CompanyID,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*UserResponse, error) {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CompanyID: u.CompanyID,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (uc *UserUseCase) ListUsers(ctx context.Context) ([]*UserResponse, error) {
	users, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*UserResponse, len(users))
	for i, u := range users {
		result[i] = &UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CompanyID: u.CompanyID,
			CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return result, nil
}

func PrepareUserUseCase(
	repo userRepo.Repository,
	companyIntegration integration.CompanyIntegration,
	logger log.Logger,
	tracer trace.Tracer,
) (*UserUseCase, error) {
	validators, err := validator.PrepareUserValidators()
	if err != nil {
		return nil, err
	}
	return NewUserUseCase(repo, companyIntegration, logger, tracer, validators), nil
}
