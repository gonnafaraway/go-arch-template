package validator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
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

// Validator интерфейс для валидации
type Validator interface {
	Validate(ctx context.Context, value interface{}) error
}

// DomainValidator валидатор для доменных сущностей
type DomainValidator interface {
	Validate(ctx context.Context, entity interface{}) error
}

// RequestValidator валидатор для HTTP запросов
type RequestValidator interface {
	ValidateRequest(ctx context.Context, req interface{}) error
}

// StringValidator валидатор для строк
type StringValidator struct {
	Required bool
	MinLen   int
	MaxLen   int
	Pattern  *regexp.Regexp
}

func (v *StringValidator) Validate(ctx context.Context, value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}

	if v.Required && str == "" {
		return errors.New("field is required")
	}

	if str == "" && !v.Required {
		return nil // Пустая строка допустима если не требуется
	}

	if v.MinLen > 0 && len(str) < v.MinLen {
		return fmt.Errorf("minimum length is %d", v.MinLen)
	}

	if v.MaxLen > 0 && len(str) > v.MaxLen {
		return fmt.Errorf("maximum length is %d", v.MaxLen)
	}

	if v.Pattern != nil && !v.Pattern.MatchString(str) {
		return errors.New("value does not match required pattern")
	}

	return nil
}

// EmailValidator валидатор для email
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type EmailValidator struct{}

func (v *EmailValidator) Validate(ctx context.Context, value interface{}) error {
	email, ok := value.(string)
	if !ok {
		return errors.New("value is not a string")
	}

	if email == "" {
		return errors.New("email is required")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

// NumberValidator валидатор для чисел
type NumberValidator struct {
	Required bool
	Min      *float64
	Max      *float64
}

func (v *NumberValidator) Validate(ctx context.Context, value interface{}) error {
	var num float64
	switch n := value.(type) {
	case float64:
		num = n
	case int:
		num = float64(n)
	case int64:
		num = float64(n)
	default:
		return errors.New("value is not a number")
	}

	if v.Min != nil && num < *v.Min {
		return fmt.Errorf("value must be at least %f", *v.Min)
	}

	if v.Max != nil && num > *v.Max {
		return fmt.Errorf("value must be at most %f", *v.Max)
	}

	return nil
}

// CompositeValidator композитный валидатор
type CompositeValidator struct {
	Validators []Validator
}

func (v *CompositeValidator) Validate(ctx context.Context, value interface{}) error {
	var errs ValidationErrors
	for _, validator := range v.Validators {
		if err := validator.Validate(ctx, value); err != nil {
			errs.Add("", err.Error())
		}
	}
	if errs.HasErrors() {
		return &errs
	}
	return nil
}

