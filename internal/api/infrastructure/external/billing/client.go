package billing

import (
	"context"
	"fmt"
)

type Client interface {
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*InvoiceResponse, error)
	GetInvoice(ctx context.Context, invoiceID string) (*InvoiceResponse, error)
	CancelInvoice(ctx context.Context, invoiceID string) error
}

type CreateInvoiceRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	UserID  string  `json:"user_id"`
}

type InvoiceResponse struct {
	ID      string  `json:"id"`
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}

type MockClient struct {
	invoices map[string]*InvoiceResponse
}

func NewMockClient() *MockClient {
	return &MockClient{
		invoices: make(map[string]*InvoiceResponse),
	}
}

func (m *MockClient) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*InvoiceResponse, error) {
	invoice := &InvoiceResponse{
		ID:      fmt.Sprintf("invoice_%d", len(m.invoices)+1),
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  "pending",
	}
	m.invoices[invoice.ID] = invoice
	return invoice, nil
}

func (m *MockClient) GetInvoice(ctx context.Context, invoiceID string) (*InvoiceResponse, error) {
	invoice, ok := m.invoices[invoiceID]
	if !ok {
		return nil, ErrNotFound
	}
	return invoice, nil
}

func (m *MockClient) CancelInvoice(ctx context.Context, invoiceID string) error {
	invoice, ok := m.invoices[invoiceID]
	if !ok {
		return ErrNotFound
	}
	invoice.Status = "cancelled"
	return nil
}
