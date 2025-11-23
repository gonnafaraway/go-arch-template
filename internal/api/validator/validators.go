package validator

import (
	companyDomain "go-arch-template/internal/api/domain/company"
	orderDomain "go-arch-template/internal/api/domain/order"
	userDomain "go-arch-template/internal/api/domain/user"
)

// CompanyValidators содержит все валидаторы для Company
type CompanyValidators struct {
	Domain  *companyDomain.DomainValidator
	Request *CompanyRequestValidator
}

// UserValidators содержит все валидаторы для User
type UserValidators struct {
	Domain  *userDomain.DomainValidator
	Request *UserRequestValidator
}

// OrderValidators содержит все валидаторы для Order
type OrderValidators struct {
	Domain  *orderDomain.DomainValidator
	Request *OrderRequestValidator
}

