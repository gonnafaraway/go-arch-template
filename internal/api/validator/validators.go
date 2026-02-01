package validator

import (
	companyDomain "go-arch-template/internal/api/domain/company"
	orderDomain "go-arch-template/internal/api/domain/order"
	userDomain "go-arch-template/internal/api/domain/user"
)

// CompanyValidators contains all validators for Company
type CompanyValidators struct {
	Domain  *companyDomain.DomainValidator
	Request *CompanyRequestValidator
}

// UserValidators contains all validators for User
type UserValidators struct {
	Domain  *userDomain.DomainValidator
	Request *UserRequestValidator
}

// OrderValidators contains all validators for Order
type OrderValidators struct {
	Domain  *orderDomain.DomainValidator
	Request *OrderRequestValidator
}

