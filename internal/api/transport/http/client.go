package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"go-arch-template/internal/api/usecase"
)

type Server struct {
	httpServer     *http.Server
	companyHandler *CompanyHandler
	userHandler    *UserHandler
	orderHandler   *OrderHandler
}

func NewServer(port string, companyUseCase *usecase.CompanyUseCase, userUseCase *usecase.UserUseCase, orderUseCase *usecase.OrderUseCase) *Server {
	ch := NewCompanyHandler(companyUseCase)
	uh := NewUserHandler(userUseCase)
	oh := NewOrderHandler(orderUseCase)

	mux := http.NewServeMux()
	
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

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer:    httpServer,
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
