package app

import (
	companyDomain "go-arch-template/internal/api/domain/company"
	orderDomain "go-arch-template/internal/api/domain/order"
	userDomain "go-arch-template/internal/api/domain/user"
	"go-arch-template/internal/api/validator"
)

func prepareCompanyValidators() (*validator.CompanyValidators, error) {
	return &validator.CompanyValidators{
		Domain:  companyDomain.NewDomainValidator(),
		Request: validator.NewCompanyRequestValidator(),
	}, nil
}

func prepareUserValidators() (*validator.UserValidators, error) {
	return &validator.UserValidators{
		Domain:  userDomain.NewDomainValidator(),
		Request: validator.NewUserRequestValidator(),
	}, nil
}

func prepareOrderValidators() (*validator.OrderValidators, error) {
	return &validator.OrderValidators{
		Domain:  orderDomain.NewDomainValidator(),
		Request: validator.NewOrderRequestValidator(),
	}, nil
}

