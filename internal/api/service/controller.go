package service

import (
	"github.com/pkg/errors"
)

type Service interface {
	Run() error
}

func RunServices[T Service](services ...T) error {
	for _, service := range services {
		err := service.Run()
		if err != nil {
			return errors.Wrap(err, "run service")
		}
	}
	return nil
}
