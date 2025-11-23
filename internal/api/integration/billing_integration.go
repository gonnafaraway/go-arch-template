package integration

import (
	"context"

	"go-arch-template/internal/api/infrastructure/external/billing"
)

type billingIntegration struct {
	billingClient billing.Client
}

func NewBillingIntegration(billingClient billing.Client) BillingIntegration {
	return &billingIntegration{
		billingClient: billingClient,
	}
}

func (i *billingIntegration) CreateInvoice(ctx context.Context, orderID string, amount float64, userID string) (string, error) {
	req := billing.CreateInvoiceRequest{
		OrderID: orderID,
		Amount:  amount,
		UserID:  userID,
	}
	
	invoice, err := i.billingClient.CreateInvoice(ctx, req)
	if err != nil {
		return "", err
	}
	return invoice.ID, nil
}

func (i *billingIntegration) GetInvoice(ctx context.Context, invoiceID string) (*billing.InvoiceResponse, error) {
	return i.billingClient.GetInvoice(ctx, invoiceID)
}

func (i *billingIntegration) CancelInvoice(ctx context.Context, invoiceID string) error {
	return i.billingClient.CancelInvoice(ctx, invoiceID)
}

