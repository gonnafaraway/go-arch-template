package http

import (
	"encoding/json"
	"go-arch-template/internal/api/usecase/order"
	"net/http"
	"strings"
)

// HandleOrders handles order creation
// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body usecase.CreateOrderCommand true "Order data"
// @Success 201 {object} usecase.CreateOrderResponse "Created order"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/orders [post]

type OrderHandler struct {
	useCase *order.OrderUseCase
}

func NewOrderHandler(useCase *order.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		useCase: useCase,
	}
}

func (h *OrderHandler) HandleOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodPost:
		var cmd order.CreateOrderCommand
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body", err)
			return
		}
		order, err := h.useCase.CreateOrder(ctx, cmd)
		if err != nil {
			if isValidationError(err) {
				respondError(w, http.StatusBadRequest, "validation failed", err)
			} else {
				respondError(w, http.StatusInternalServerError, "failed to create order", err)
			}
			return
		}
		respondJSON(w, http.StatusCreated, order)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleOrder handles single order operations
// @Summary Get order by ID
// @Description Get order information by ID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} usecase.OrderResponse "Order information"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/orders/{id} [get]
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

// HandleConfirmOrder handles order confirmation
// @Summary Confirm order
// @Description Confirm an existing order
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Order confirmed"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/orders/confirm/{id} [post]
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
