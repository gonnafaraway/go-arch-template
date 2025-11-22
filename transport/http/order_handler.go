package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"go-arch-template/usecase"
)

type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
}

func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID string                     `json:"user_id"`
		Items  []usecase.OrderItemRequest `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cmd := usecase.CreateOrderCommand{
		UserID: request.UserID,
		Items:  request.Items,
	}

	response, err := h.orderUseCase.CreateOrder(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	if err := h.orderUseCase.ConfirmOrder(r.Context(), orderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
