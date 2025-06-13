package main

import (
	"fmt"
	"log"

	"alchemorsel/backend/internal/config"
	httpserver "alchemorsel/backend/internal/interfaces/http"
)

func main() {
	cfg := config.Load()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)


	router := httpserver.SetupRouter()
	if err := router.Run(addr); err != nil {

		log.Fatal(err)
	}
}
