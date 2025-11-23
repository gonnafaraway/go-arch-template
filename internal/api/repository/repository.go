package repository

import (
	"go-arch-template/internal/api/repository/company"
	"go-arch-template/internal/api/repository/order"
	"go-arch-template/internal/api/repository/user"
	"go-arch-template/internal/api/storage"
)

// Repositories содержит все репозитории
type Repositories struct {
	CompanyRepository company.Repository
	UserRepository    user.Repository
	OrderRepository   order.Repository
}

func PrepareRepository(storage *storage.Storage) (*Repositories, error) {
	repos := &Repositories{}

	// Используем реальные репозитории если storage доступен, иначе моки
	if storage.Postgres != nil {
		repos.CompanyRepository = company.NewPostgresRepository(storage.Postgres)
		// Можно добавить другие репозитории для PostgreSQL
	} else if storage.MongoDB != nil {
		repos.CompanyRepository = company.NewMongoDBRepository(storage.MongoDB)
		// Можно добавить другие репозитории для MongoDB
	} else {
		// Fallback на моки
		repos.CompanyRepository = company.NewMockRepository()
		repos.UserRepository = user.NewMockRepository()
		repos.OrderRepository = order.NewMockRepository()
	}

	// Если не были установлены, используем моки
	if repos.UserRepository == nil {
		repos.UserRepository = user.NewMockRepository()
	}
	if repos.OrderRepository == nil {
		repos.OrderRepository = order.NewMockRepository()
	}

	return repos, nil
}
