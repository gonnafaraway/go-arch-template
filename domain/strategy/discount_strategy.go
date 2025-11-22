package strategy

import "go-arch-template/domain/entity"

type DiscountStrategy interface {
	CalculateDiscount(order *entity.Order) (int64, error)
}

type PercentageDiscountStrategy struct {
	Percentage float64
}

func (s PercentageDiscountStrategy) CalculateDiscount(order *entity.Order) (int64, error) {
	discount := float64(order.Total.Amount) * s.Percentage / 100
	return int64(discount), nil
}

type FixedAmountDiscountStrategy struct {
	Amount int64
}

func (s FixedAmountDiscountStrategy) CalculateDiscount(order *entity.Order) (int64, error) {
	if s.Amount > order.Total.Amount {
		return order.Total.Amount, nil
	}
	return s.Amount, nil
}

type DiscountCalculator struct {
	strategies map[string]DiscountStrategy
}

func NewDiscountCalculator() *DiscountCalculator {
	return &DiscountCalculator{
		strategies: map[string]DiscountStrategy{
			"percentage": PercentageDiscountStrategy{Percentage: 10},
			"fixed":      FixedAmountDiscountStrategy{Amount: 5000}, // 50.00
		},
	}
}

func (c *DiscountCalculator) Calculate(userType string, order *entity.Order) (int64, error) {
	strategy, exists := c.strategies[userType]
	if !exists {
		return 0, nil // Нет скидки
	}
	return strategy.CalculateDiscount(order)
}
