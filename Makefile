THIS_FILE := $(lastword $(MAKEFILE_LIST))
PROJECT_NAME := $(shell go list -m)
SYSTEM_NAME := $(shell uname)

.PHONY: generate swagger-install swagger-generate swagger

####################
# OPENAPI
####################

swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger-generate:
	swag init -g cmd/api/main.go -o docs/api/openapi

swagger-validate:
	swagger validate docs/swagger.yaml

swagger: swagger-install swagger-generate

####################
# TBLS
####################

tbls_install:
	go install github.com/k1LoW/tbls@v1.77.0

tbls_doc_db:
	tbls doc "postgres://postgres:postgres@:5432/postgres?search_path=public&sslmode=disable" --rm-dist

tbls_diff_db:
	tbls diff "postgres://postgres:postgres@:5432/postgres?search_path=public&sslmode=disable" docs/api/db
