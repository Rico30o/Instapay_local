package main

import (
	traceroute "instapay/TraceRoute"
	"instapay/database"
	serviceroutes "instapay/serviceRoutes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Setup Endpoint routes
	serviceroutes.SetupRoutes(app)

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Middleware: Logger
	app.Use(logger.New())

	// Database connection
	err = database.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Auto-migrate models
	err = database.DB.AutoMigrate(&traceroute.NetworkAlertID{})
	if err != nil {
		log.Fatalf("Error auto-migrating models: %v", err)
	}

	// Start server
	err = app.Listen(":9090")
	if err != nil {
		log.Fatal("Error starting the server")
	}
}
