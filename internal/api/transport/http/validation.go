package http

import (
	"encoding/json"
	"net/http"

	"go-arch-template/internal/api/validator"
)

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(*validator.ValidationErrors)
	return ok
}

// RespondError sends error in JSON format
func RespondError(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]interface{}{
		"error":   message,
		"message": err.Error(),
	}

	// If it's a validation error, add details
	if validationErr, ok := err.(*validator.ValidationErrors); ok {
		response["validation_errors"] = validationErr.Errors
	}

	json.NewEncoder(w).Encode(response)
}
