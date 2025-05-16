package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/mahimasingh04/go-website/glow-backend/controllers"
	
)

func AuthRoutes(app *fiber.App, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	api := app.Group("/api/auth")
	api.Post("/register", authController.SignUpUser)
	api.Post("/login", authController.SignInUser)
	
}