package order

import "errors"

var (
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrEmptyItems      = errors.New("empty items")
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrInvalidPrice    = errors.New("invalid price")
	ErrInvalidStatus   = errors.New("invalid status")
)

