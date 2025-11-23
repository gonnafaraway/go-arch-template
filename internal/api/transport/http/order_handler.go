package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-arch-template/internal/api/usecase"
)

type OrderHandler struct {
	useCase *usecase.OrderUseCase
}

func NewOrderHandler(useCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		useCase: useCase,
	}
}

func (h *OrderHandler) HandleOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost:
		var cmd usecase.CreateOrderCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		order, err := h.useCase.CreateOrder(ctx, cmd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusCreated, order)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimPrefix(r.URL.Path, "/api/orders/")

	switch r.Method {
	case http.MethodGet:
		order, err := h.useCase.GetOrder(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		respondJSON(w, http.StatusOK, order)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandler) HandleConfirmOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimPrefix(r.URL.Path, "/api/orders/confirm/")

	switch r.Method {
	case http.MethodPost:
		if err := h.useCase.ConfirmOrder(ctx, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, map[string]string{"status": "confirmed"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

