package billing

type CreateInvoiceRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	UserID  string  `json:"user_id"`
}

type Invoice struct {
	ID      string  `json:"id"`
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}
