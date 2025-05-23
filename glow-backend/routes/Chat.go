package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

 "github.com/mahimasingh04/go-website/glow-backend/controllers"
 "github.com/mahimasingh04/go-website/glow-backend/middleware"

	
)

func ChatRoutes(app *fiber.App, db *gorm.DB)  {
	GeminiChatController := controllers.NewGeminiChatController(db)

	ai := app.Group("/api/ai")

	ai.Use(middleware.AuthMiddleware(db))
	ai.Post("/chat", GeminiChatController.ChatWithAI)
	// Add inside ChatRoutes for debugging
 ai.Get("/test", func(c *fiber.Ctx) error {
    return c.SendString("AI route is working")
})
	
}