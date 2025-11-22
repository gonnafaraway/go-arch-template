package vo

import (
	"errors"
	"fmt"
)

type Money struct {
	amount   int64 // в копейках/центах
	currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("money amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	return Money{amount: amount, currency: currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("currencies don't match")
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(m.amount)/100, m.currency)
}
