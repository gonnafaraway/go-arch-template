package specification

import (
	"go-arch-template/domain/entity"
	"time"
)

type OrderSpecification interface {
	IsSatisfiedBy(order *entity.Order) bool
}

type HighValueOrderSpec struct {
	MinAmount int64
}

func (s HighValueOrderSpec) IsSatisfiedBy(order *entity.Order) bool {
	return order.Total.Amount >= s.MinAmount
}

type RecentOrderSpec struct {
	Since time.Time
}

func (s RecentOrderSpec) IsSatisfiedBy(order *entity.Order) bool {
	return order.CreatedAt.After(s.Since)
}

type CompositeSpec struct {
	Specs []OrderSpecification
}

func (cs CompositeSpec) IsSatisfiedBy(order *entity.Order) bool {
	for _, spec := range cs.Specs {
		if !spec.IsSatisfiedBy(order) {
			return false
		}
	}
	return true
}
