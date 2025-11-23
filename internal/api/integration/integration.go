package integration

import (
	"context"

	"go-arch-template/internal/api/infrastructure/external/billing"
	usersservice "go-arch-template/internal/api/infrastructure/external/users-service"
)

type Integrations struct {
	CompanyIntegration CompanyIntegration
	BillingIntegration BillingIntegration
	OAuthIntegration   OAuthIntegration
}

type CompanyIntegration interface {
	ValidateCompany(ctx context.Context, companyID string) (bool, error)
	SyncCompany(ctx context.Context, companyID string) error
}

type BillingIntegration interface {
	CreateInvoice(ctx context.Context, orderID string, amount float64, userID string) (string, error)
	GetInvoice(ctx context.Context, invoiceID string) (*billing.InvoiceResponse, error)
	CancelInvoice(ctx context.Context, invoiceID string) error
}

type OAuthIntegration interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)
}

func PrepareIntegration(env interface{}) (*Integrations, error) {
	billingClient := billing.NewMockClient()
	usersClient := usersservice.NewMockClient()

	return &Integrations{
		CompanyIntegration: NewCompanyIntegration(usersClient),
		BillingIntegration: NewBillingIntegration(billingClient),
		OAuthIntegration:   NewOAuthIntegration(),
	}, nil
}

