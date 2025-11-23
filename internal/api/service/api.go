package service

import (
	"context"
	"log"

	"go-arch-template/internal/api/transport/http"
	"go-arch-template/internal/api/usecase"
)

type API struct {
	server *http.Server
	env    *Env
}

func PrepareAPIService(
	env *Env,
	companyUseCase *usecase.CompanyUseCase,
	userUseCase *usecase.UserUseCase,
	orderUseCase *usecase.OrderUseCase,
) (*API, error) {
	port := env.HTTPPort
	if port == "" {
		port = "8080"
	}

	server := http.NewServer(port, companyUseCase, userUseCase, orderUseCase)

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
