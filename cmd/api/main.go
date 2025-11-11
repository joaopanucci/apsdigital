package main

import (
	"log"
	"os"

	"apsdigital/internal/config"
	"apsdigital/internal/infra/db"
	"apsdigital/internal/infra/http/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database, err := db.NewPostgresConnection(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Setup Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := router.SetupRoutes(database, cfg)

	// Create uploads directory
	if err := os.MkdirAll(cfg.Upload.Path, 0755); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}