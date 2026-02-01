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
	// Here you can start Change Data Capture
	// For now just return nil
	return nil
}

func (c *CDC) Shutdown(ctx context.Context) error {
	log.Println("Shutting down CDC service...")
	return nil
}
