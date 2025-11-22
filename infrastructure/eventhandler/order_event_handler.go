package eventhandler

import (
	"context"
	"log"

	"go-arch-template/domain/event"
	"go-arch-template/ports/notification"
)

type OrderEventHandler struct {
	emailService  notification.EmailService
	analyticsRepo notification.AnalyticsRepository
}

func NewOrderEventHandler(
	emailService notification.EmailService,
	analyticsRepo notification.AnalyticsRepository,
) *OrderEventHandler {
	return &OrderEventHandler{
		emailService:  emailService,
		analyticsRepo: analyticsRepo,
	}
}

func (h *OrderEventHandler) Handle(ctx context.Context, event interface{}) error {
	switch e := event.(type) {
	case event.OrderCreatedEvent:
		return h.handleOrderCreated(ctx, e)
	case event.OrderConfirmedEvent:
		return h.handleOrderConfirmed(ctx, e)
	case event.OrderCancelledEvent:
		return h.handleOrderCancelled(ctx, e)
	default:
		return nil // Игнорируем неизвестные события
	}
}

func (h *OrderEventHandler) handleOrderCreated(ctx context.Context, e event.OrderCreatedEvent) error {
	// Отправка email асинхронно
	go func() {
		if err := h.emailService.SendOrderConfirmation(e.OrderID, e.UserID); err != nil {
			log.Printf("Failed to send order confirmation: %v", err)
		}
	}()

	// Запись в аналитику
	if err := h.analyticsRepo.TrackOrderCreated(e); err != nil {
		log.Printf("Failed to track order creation: %v", err)
	}

	return nil
}

func (h *OrderEventHandler) HandleError(ctx context.Context, event interface{}, err error) {
	log.Printf("Failed to handle event %T: %v", event, err)
}
