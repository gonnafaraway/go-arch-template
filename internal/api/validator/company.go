package validator

import (
	"context"

	"go-arch-template/internal/api/domain/company"
)

// CreateCompanyRequest request structure for validation
type CreateCompanyRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CompanyRequestValidator validator for Company requests
type CompanyRequestValidator struct{}

func NewCompanyRequestValidator() *CompanyRequestValidator {
	return &CompanyRequestValidator{}
}

// ValidateCreateRequest validates company creation request
func (v *CompanyRequestValidator) ValidateCreateRequest(ctx context.Context, req *CreateCompanyRequest) error {
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
