package service

import (
	"context"
	"log"

	"go-arch-template/internal/api/env"
)

type CDC struct {
	env *env.Env
}

func PrepareCDCService(env *env.Env) (*CDC, error) {
	return &CDC{
		env: env,
	}, nil
}

func (c *CDC) Run() error {
	log.Println("Starting CDC service...")
	// Здесь можно запустить Change Data Capture
	// Пока просто возвращаем nil
	return nil
}

func (c *CDC) Shutdown(ctx context.Context) error {
	log.Println("Shutting down CDC service...")
	return nil
}
