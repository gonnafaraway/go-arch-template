package http

import (
	"encoding/json"
	"net/http"

	"go-arch-template/internal/api/validator"
)

// isValidationError проверяет является ли ошибка ошибкой валидации
func isValidationError(err error) bool {
	_, ok := err.(*validator.ValidationErrors)
	return ok
}

// respondError отправляет ошибку в формате JSON
func respondError(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]interface{}{
		"error":   message,
		"message": err.Error(),
	}

	// Если это ошибка валидации, добавляем детали
	if validationErr, ok := err.(*validator.ValidationErrors); ok {
		response["validation_errors"] = validationErr.Errors
	}

	json.NewEncoder(w).Encode(response)
}

