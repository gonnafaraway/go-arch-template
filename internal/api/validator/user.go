package validator

import (
	"context"

	"go-arch-template/internal/api/domain/user"
)

// CreateUserRequest request structure for validation
type CreateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
}

// UserRequestValidator validator for User requests
type UserRequestValidator struct{}

func NewUserRequestValidator() *UserRequestValidator {
	return &UserRequestValidator{}
}

// ValidateCreateRequest validates user creation request
func (v *UserRequestValidator) ValidateCreateRequest(ctx context.Context, req *CreateUserRequest) error {
	var errs ValidationErrors

	// Name validation
	if req.Name == "" {
		errs.Add("name", "name is required")
	} else if len(req.Name) < 2 {
		errs.Add("name", "name must be at least 2 characters")
	} else if len(req.Name) > 100 {
		errs.Add("name", "name must be at most 100 characters")
	}

	// Email validation
	emailValidator := &EmailValidator{}
	if err := emailValidator.Validate(ctx, req.Email); err != nil {
		errs.Add("email", err.Error())
	}

	// Company ID validation
	if req.CompanyID == "" {
		errs.Add("company_id", "company_id is required")
	}

	if errs.HasErrors() {
		return &errs
	}

	return nil
}

func PrepareUserValidators() (*UserValidators, error) {
	return &UserValidators{
		Domain:  user.NewDomainValidator(),
		Request: NewUserRequestValidator(),
	}, nil
}
