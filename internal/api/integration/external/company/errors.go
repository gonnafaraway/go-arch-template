package company

import "errors"

var (
	ErrNotFound = errors.New("company not found")
	ErrInvalid  = errors.New("company is invalid")
)
