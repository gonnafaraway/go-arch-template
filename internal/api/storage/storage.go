package storage

import (
	"go-arch-template/internal/api/env"
	"os"

	mongodb "go-arch-template/internal/api/storage/mongodb"
	postgres "go-arch-template/internal/api/storage/postgres"
)

type Storage struct {
	MongoDB  *mongodb.Client
	Postgres *postgres.Client
}

func PrepareStorage(env *env.Env) (*Storage, error) {
	storage := &Storage{}

	// Initialize MongoDB if URI is specified
	if mongoURI := os.Getenv("MONGODB_URI"); mongoURI != "" {
		dbName := os.Getenv("MONGODB_DATABASE")
		if dbName == "" {
			dbName = "go_arch_template"
		}
		mongoClient, err := mongodb.NewClient(mongoURI, dbName)
		if err != nil {
			return nil, err
		}
		storage.MongoDB = mongoClient
	}

	// Initialize PostgreSQL if parameters are specified
	if pgHost := os.Getenv("POSTGRES_HOST"); pgHost != "" {
		pgPort := os.Getenv("POSTGRES_PORT")
		if pgPort == "" {
			pgPort = "5432"
		}
		pgUser := os.Getenv("POSTGRES_USER")
		pgPassword := os.Getenv("POSTGRES_PASSWORD")
		pgDB := os.Getenv("POSTGRES_DB")
		if pgDB == "" {
			pgDB = "go_arch_template"
		}

		pgClient, err := postgres.NewClient(pgHost, pgPort, pgUser, pgPassword, pgDB)
		if err != nil {
			return nil, err
		}
		storage.Postgres = pgClient
	}

	return storage, nil
}
