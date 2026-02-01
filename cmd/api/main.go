package main

import (
	"log"

	"go-arch-template/internal/api/app"
)

// @title Go Arch Template API
// @version 1.0
// @description This is a sample server for Go Architecture Template
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @schemes http https
func main() {
	err := app.Run()
	if err != nil {
		log.Fatal("run app", err)
	}
}
