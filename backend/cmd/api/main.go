package main

import (
	"fmt"

	"alchemorsel/backend/internal/config"
	httpserver "alchemorsel/backend/internal/interfaces/http"
	"alchemorsel/backend/internal/pkg/logger"
)

func main() {
	cfg := config.Load()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Starting server on %s", addr)

	router := httpserver.SetupRouter()
	if err := router.Run(addr); err != nil {
		logger.Fatal(err)
	}
}
