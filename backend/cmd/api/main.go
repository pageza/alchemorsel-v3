package main

import (
	"log"

	"alchemorsel/backend/internal/config"
	"alchemorsel/backend/internal/interfaces/http"
)

func main() {
	cfg := config.Load()
	router := http.SetupRouter()

	log.Printf("Starting server on %s", cfg.Address)
	if err := router.Run(cfg.Address); err != nil {
		log.Fatal(err)
	}
}
