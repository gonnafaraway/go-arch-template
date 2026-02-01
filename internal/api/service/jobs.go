package service

import (
	"context"
	"log"
	"time"

	"go-arch-template/internal/api/env"
)

type Jobs struct {
	env *env.Env
}

func PrepareJobsService(env *env.Env, emailCheckerDomain interface{}) (*Jobs, error) {
	return &Jobs{
		env: env,
	}, nil
}

func (j *Jobs) Run() error {
	log.Println("Starting Jobs service...")
	// Here you can start periodic tasks
	// For now just return nil so the service doesn't block
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("Jobs service is running...")
		}
	}()
	return nil
}

func (j *Jobs) Shutdown(ctx context.Context) error {
	log.Println("Shutting down Jobs service...")
	return nil
}
