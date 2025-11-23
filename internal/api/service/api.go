package service

import (
	"context"
	"log"

	"go-arch-template/internal/api/env"
	"go-arch-template/internal/api/transport/http"
	companyUseCase "go-arch-template/internal/api/usecase/company"
	"go-arch-template/internal/api/usecase/order"
	"go-arch-template/internal/api/usecase/user"
)

type API struct {
	server *http.Server
	env    *env.Env
}

func PrepareAPIService(
	env *env.Env,
	companyUC *companyUseCase.CompanyUseCase,
	userUC *user.UserUseCase,
	orderUC *order.OrderUseCase,
) (*API, error) {
	port := env.HTTPPort
	if port == "" {
		port = "8080"
	}

	server := http.NewServer(port, companyUC, userUC, orderUC)

	return &API{
		server: server,
		env:    env,
	}, nil
}

func (a *API) Run() error {
	log.Println("Starting API service...")
	return a.server.Run()
}

func (a *API) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
