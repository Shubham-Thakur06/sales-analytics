package main

import (
	"fmt"
	"log"

	"sales-analytics/internal/config"
	"sales-analytics/internal/container"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize container
	container, err := container.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Error initializing container: %v", err)
	}
	defer container.Stop()

	// Start background services
	container.Start()

	// Setup and start HTTP server
	router := container.SetupHTTPServer()
	addr := fmt.Sprintf(":%d", container.Config.AppPort)
	container.Logger.Infof("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		container.Logger.Fatalf("Error starting server: %v", err)
	}
}
