package validator

import (
	"context"

	"go-arch-template/internal/api/domain/company"
)

// CreateCompanyRequest структура запроса для валидации
type CreateCompanyRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CompanyRequestValidator валидатор для запросов Company
type CompanyRequestValidator struct{}

func NewCompanyRequestValidator() *CompanyRequestValidator {
	return &CompanyRequestValidator{}
}

// ValidateCreateRequest валидирует запрос на создание компании
func (v *CompanyRequestValidator) ValidateCreateRequest(ctx context.Context, req *CreateCompanyRequest) error {
	var errs ValidationErrors

	// Валидация имени
	if req.Name == "" {
		errs.Add("name", "name is required")
	} else if len(req.Name) < 2 {
		errs.Add("name", "name must be at least 2 characters")
	} else if len(req.Name) > 100 {
		errs.Add("name", "name must be at most 100 characters")
	}

	// Валидация email
	emailValidator := &EmailValidator{}
	if err := emailValidator.Validate(ctx, req.Email); err != nil {
		errs.Add("email", err.Error())
	}

	if errs.HasErrors() {
		return &errs
	}

	return nil
}

func PrepareCompanyValidators() (*CompanyValidators, error) {
	return &CompanyValidators{
		Domain:  company.NewDomainValidator(),
		Request: NewCompanyRequestValidator(),
	}, nil
}
