package main

import (
	"log"

	"go-arch-template/internal/api/app"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal("init app", err)
	}
}
