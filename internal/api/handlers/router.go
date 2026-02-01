package handlers

import (
	"context"
	"log"
	"net/http"

	"go-arch-template/internal/api/handlers/company"
	"go-arch-template/internal/api/transport/http/middleware"
	"go-arch-template/internal/api/usecase/order"
	"go-arch-template/internal/api/usecase/user"

	orderhandler "go-arch-template/internal/api/handlers/order"
	userhandler "go-arch-template/internal/api/handlers/user"
	httptransport "go-arch-template/internal/api/transport/http"
	companyUseCase "go-arch-template/internal/api/usecase/company"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string) *Server {
	httpServer := httptransport.NewHTTPTransport(port)
	return &Server{httpServer: httpServer}
}

func BindRoutes(server *Server, companyUC *companyUseCase.CompanyUseCase, userUC *user.UserUseCase, orderUC *order.OrderUseCase) {
	ch := company.NewCompanyHandler(companyUC)
	uh := userhandler.NewUserHandler(userUC)
	oh := orderhandler.NewOrderHandler(orderUC)

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

	server.httpServer.Handler = handler
}

func (s *Server) Run() error {
	log.Printf("HTTP server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
