package main

import (
	"fmt"
	"log"

	"sales-analytics/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	fmt.Println(cfg)
}
