package repository

import (
	"go-arch-template/internal/api/repository/company"
	"go-arch-template/internal/api/repository/order"
	"go-arch-template/internal/api/repository/user"
	"go-arch-template/internal/api/storage"
)

// Repositories contains all repositories
type Repositories struct {
	CompanyRepository company.Repository
	UserRepository    user.Repository
	OrderRepository   order.Repository
}

func PrepareRepository(storage *storage.Storage) (*Repositories, error) {
	repos := &Repositories{}

	// Use real repositories if storage is available, otherwise use mocks
	if storage.Postgres != nil {
		repos.CompanyRepository = company.NewPostgresRepository(storage.Postgres)
		// Can add other repositories for PostgreSQL
	} else if storage.MongoDB != nil {
		repos.CompanyRepository = company.NewMongoDBRepository(storage.MongoDB)
		// Can add other repositories for MongoDB
	} else {
		// Fallback to mocks
		repos.CompanyRepository = company.NewMockRepository()
		repos.UserRepository = user.NewMockRepository()
		repos.OrderRepository = order.NewMockRepository()
	}

	// If not set, use mocks
	if repos.UserRepository == nil {
		repos.UserRepository = user.NewMockRepository()
	}
	if repos.OrderRepository == nil {
		repos.OrderRepository = order.NewMockRepository()
	}

	return repos, nil
}
