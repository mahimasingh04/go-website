package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

 "github.com/mahimasingh04/go-website/glow-backend/controllers"
 "github.com/mahimasingh04/go-website/glow-backend/middleware"

	
)

func ChatRoutes(chat *fiber.App, db *gorm.DB)  {
	GeminiChatController := controllers.NewGeminiChatController(db)

	ai := chat.Group("/api/ai")

	ai.Use(middleware.AuthMiddleware(db))
	ai.Post("/chat", GeminiChatController.ChatWithAI)
	
}