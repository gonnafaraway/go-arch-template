package service

import (
	"context"
	"log"
	"time"
)

type Jobs struct {
	env *Env
}

func PrepareJobsService(env *Env, emailCheckerDomain interface{}) (*Jobs, error) {
	return &Jobs{
		env: env,
	}, nil
}

func (j *Jobs) Run() error {
	log.Println("Starting Jobs service...")
	// Здесь можно запустить периодические задачи
	// Пока просто возвращаем nil, чтобы сервис не блокировался
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
