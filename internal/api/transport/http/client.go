package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"go-arch-template/internal/api/transport/http/middleware"
	companyUseCase "go-arch-template/internal/api/usecase/company"
	"go-arch-template/internal/api/usecase/order"
	"go-arch-template/internal/api/usecase/user"
)

type Server struct {
	httpServer     *http.Server
	companyHandler *CompanyHandler
	userHandler    *UserHandler
	orderHandler   *OrderHandler
}

func NewServer(port string, companyUC *companyUseCase.CompanyUseCase, userUC *user.UserUseCase, orderUC *order.OrderUseCase) *Server {
	ch := NewCompanyHandler(companyUC)
	uh := NewUserHandler(userUC)
	oh := NewOrderHandler(orderUC)

	mux := http.NewServeMux()

	// Prometheus metrics endpoint
	mux.Handle("/metrics", middleware.PrometheusHandler())

	// Company routes
	mux.HandleFunc("/api/companies", ch.HandleCompanies)
	mux.HandleFunc("/api/companies/", ch.HandleCompany)

	// User routes
	mux.HandleFunc("/api/users", uh.HandleUsers)
	mux.HandleFunc("/api/users/", uh.HandleUser)

	// Order routes
	mux.HandleFunc("/api/orders", oh.HandleOrders)
	mux.HandleFunc("/api/orders/", oh.HandleOrder)
	mux.HandleFunc("/api/orders/confirm/", oh.HandleConfirmOrder)

	// Apply middleware chain
	handler := middleware.SentryMiddleware()(middleware.PrometheusMiddleware(mux))

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer:     httpServer,
		companyHandler: ch,
		userHandler:    uh,
		orderHandler:   oh,
	}
}

func (s *Server) Run() error {
	log.Printf("HTTP server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
