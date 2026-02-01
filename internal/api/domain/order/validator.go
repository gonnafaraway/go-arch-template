package order

import (
	"context"
	"fmt"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
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

// DomainValidator validator for Order domain entity
type DomainValidator struct{}

func NewDomainValidator() *DomainValidator {
	return &DomainValidator{}
}

// Validate validates the Order domain entity
func (v *DomainValidator) Validate(ctx context.Context, o *Order) error {
	var errs ValidationErrors

	// User ID validation
	if o.UserID == "" {
		errs.Add("user_id", "user_id is required")
	}

	// Items validation
	if len(o.Items) == 0 {
		errs.Add("items", "at least one item is required")
	}

	// Validate each item
	for i, item := range o.Items {
		if item.ProductID == "" {
			errs.Add(fmt.Sprintf("items[%d].product_id", i), "product_id is required")
		}
		if item.Name == "" {
			errs.Add(fmt.Sprintf("items[%d].name", i), "name is required")
		}
		if item.Quantity <= 0 {
			errs.Add(fmt.Sprintf("items[%d].quantity", i), "quantity must be greater than 0")
		}
		if item.Price < 0 {
			errs.Add(fmt.Sprintf("items[%d].price", i), "price must be non-negative")
		}
	}

	// Total validation
	if o.Total < 0 {
		errs.Add("total", "total must be non-negative")
	}

	// Status validation
	validStatuses := map[Status]bool{
		StatusPending:   true,
		StatusConfirmed: true,
		StatusCancelled: true,
	}
	if !validStatuses[o.Status] {
		errs.Add("status", "invalid order status")
	}

	if errs.HasErrors() {
		return &errs
	}

	return nil
}
