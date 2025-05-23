package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
    "os"
    "github.com/joho/godotenv"
	"github.com/mahimasingh04/go-website/glow-backend/models"
	"github.com/mahimasingh04/go-website/glow-backend/routes"
)

func main() {
	// Initialize DB

	   err := godotenv.Load()
    if err != nil {
        panic("Error loading .env file")
    }
	dsn := "host=" + os.Getenv("DB_HOST") + 
           " user=" + os.Getenv("DB_USER") + 
           " password=" + os.Getenv("DB_PASSWORD") + 
           " dbname=" + os.Getenv("DB_NAME") + 
           " port=" + os.Getenv("DB_PORT") + 
           " sslmode=" + os.Getenv("DB_SSLMODE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto-migrate User model
	db.AutoMigrate(&models.User{} , &models.Query{})

	// Initialize Fiber
	app := fiber.New()
	
	

	// Setup routes
	routes.AuthRoutes(app, db)
	routes.ChatRoutes(app, db)

	// Start server
	app.Listen(":3000")
}