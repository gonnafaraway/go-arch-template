package billing

import (
	"context"
	"fmt"
)

type Client interface {
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*Invoice, error)
	GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error)
	CancelInvoice(ctx context.Context, invoiceID string) error
}

type client struct {
	invoices map[string]*Invoice
}

func NewClient() Client {
	return &client{
		invoices: make(map[string]*Invoice),
	}
}

func (c *client) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*Invoice, error) {
	invoice := &Invoice{
		ID:      fmt.Sprintf("invoice_%d", len(c.invoices)+1),
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  "pending",
	}
	c.invoices[invoice.ID] = invoice
	return invoice, nil
}

func (c *client) GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error) {
	invoice, ok := c.invoices[invoiceID]
	if !ok {
		return nil, ErrNotFound
	}
	return invoice, nil
}

func (c *client) CancelInvoice(ctx context.Context, invoiceID string) error {
	invoice, ok := c.invoices[invoiceID]
	if !ok {
		return ErrNotFound
	}
	invoice.Status = "cancelled"
	return nil
}
