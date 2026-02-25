package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/uduxpass/backend/internal/infrastructure/database"
	"github.com/uduxpass/backend/internal/interfaces/http/server"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get configuration from environment
	config := &server.Config{
		Host:               getEnv("HOST", "0.0.0.0"),
		Port:               getEnv("PORT", "8080"),
		Environment:        getEnv("ENV", "development"),
		JWTSecret:          getEnv("JWT_SECRET", "uduxpass-default-secret-key"),
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}

	// Initialize database manager
	dbManager, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbManager.Close()

	// Create server
	srv := server.NewServer(config, dbManager)

	// Start server
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Printf("Starting uduXPass API server on %s", addr)
	
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initializeDatabase creates and initializes the database manager
func initializeDatabase() (*database.DatabaseManager, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "ubuntu")
	password := getEnv("DB_PASSWORD", "ubuntu")
	dbname := getEnv("DB_NAME", "uduxpass_db")
	sslmode := getEnv("DB_SSLMODE", "disable")
	
	// Build PostgreSQL connection string
	databaseURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	return database.NewDatabaseManager(databaseURL)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
