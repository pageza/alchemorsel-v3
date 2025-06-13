package main

import (
	"fmt"
	"log"

	"alchemorsel/backend/internal/config"
	httpserver "alchemorsel/backend/internal/interfaces/http"
)

func main() {
	cfg := config.Load()
	router := httpserver.SetupRouter()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
