package user

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	nameRegex  = regexp.MustCompile(`^[a-zA-Z0-9\s\-_]{2,100}$`)
)

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors представляет множество ошибок валидации
type ValidationErrors struct {
	Errors []ValidationError
}

func (e *ValidationErrors) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func (e *ValidationErrors) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{Field: field, Message: message})
}

func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// DomainValidator валидатор для доменной сущности User
type DomainValidator struct{}

func NewDomainValidator() *DomainValidator {
	return &DomainValidator{}
}

// Validate валидирует доменную сущность User
func (v *DomainValidator) Validate(ctx context.Context, u *User) error {
	var errs ValidationErrors

	// Валидация имени
	if u.Name == "" {
		errs.Add("name", "name is required")
	} else if len(u.Name) < 2 {
		errs.Add("name", "name must be at least 2 characters")
	} else if len(u.Name) > 100 {
		errs.Add("name", "name must be at most 100 characters")
	} else if !nameRegex.MatchString(u.Name) {
		errs.Add("name", "name contains invalid characters")
	}

	// Валидация email
	if u.Email == "" {
		errs.Add("email", "email is required")
	} else if !emailRegex.MatchString(u.Email) {
		errs.Add("email", "invalid email format")
	}

	// Валидация company_id
	if u.CompanyID == "" {
		errs.Add("company_id", "company_id is required")
	}

	if errs.HasErrors() {
		return &errs
	}

	return nil
}
