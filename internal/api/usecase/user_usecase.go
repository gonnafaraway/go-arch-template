package usecase

import (
	"context"

	"go-arch-template/internal/api/domain/user"
	userRepo "go-arch-template/internal/api/repository/user"
)

type UserUseCase struct {
	repo userRepo.Repository
}

func NewUserUseCase(repo userRepo.Repository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
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
	u := user.NewUser(req.Name, req.Email, req.CompanyID)
	if err := uc.repo.Create(ctx, u); err != nil {
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

