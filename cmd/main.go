package main

import (
	"log"

	"github.com/okoye-dev/oss-archive/internal/config"
	"github.com/okoye-dev/oss-archive/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Create and start server
	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}