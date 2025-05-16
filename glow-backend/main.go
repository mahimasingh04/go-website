package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mahimasingh04/go-website/glow-backend/models"
	"github.com/mahimasingh04/go-website/glow-backend/routes"
)

func main() {
	// Initialize DB
	dsn := "host=localhost user=postgres password=Mahi91098$$ dbname=Skincare port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto-migrate User model
	db.AutoMigrate(&models.User{})

	// Initialize Fiber
	app := fiber.New()

	// Setup routes
	routes.AuthRoutes(app, db)

	// Start server
	app.Listen(":3000")
}