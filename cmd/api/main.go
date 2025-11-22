// cmd/api/main.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go-arch-template/app"
)

func main() {
	// Инициализация приложения
	application, err := app.NewApplication()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	// Настройка маршрутов
	router := mux.NewRouter()

	// Order routes
	router.HandleFunc("/api/orders", application.OrderHandler.CreateOrder).Methods("POST")
	router.HandleFunc("/api/orders/{id}/confirm", application.OrderHandler.ConfirmOrder).Methods("POST")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed:", err)
	}
}
