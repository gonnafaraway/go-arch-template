package service

import (
	"sync"

	"github.com/pkg/errors"
)

type Service interface {
	Run() error
}

func RunServices(services ...Service) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(services))

	for _, srv := range services {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			if err := s.Run(); err != nil {
				errChan <- errors.Wrap(err, "run service")
			}
		}(srv)
	}

	// Ждем завершения всех сервисов или первую ошибку
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:
		return err
	case <-done:
		return nil
	}
}
