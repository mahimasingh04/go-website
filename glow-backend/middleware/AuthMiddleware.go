

package middleware

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// Initialize JWT secret once when the package is loaded
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file (might be okay if using environment variables)", err)
	}

	jwtSecretStr := os.Getenv("JWT_SECRET")
	if jwtSecretStr == "" {
		log.Fatal("JWT_SECRET not set in environment variables")
	}

	jwtSecret = []byte(jwtSecretStr)
}

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from cookie or header
		tokenString := c.Cookies("token")
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" {
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenString = strings.TrimPrefix(authHeader, "Bearer ")
				} else {
					
					tokenString = authHeader
				}
			}
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Missing token",
			})
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			log.Printf("Token parsing error: %v", err) // Log the actual error for debugging
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized: Invalid token",
				"details": err.Error(), // Include the actual error details for debugging
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid token",
			})
		}

		// Check token expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Token has expired",
			})
		}

		// Store user ID and DB in context
		c.Locals("user_id", claims.UserID)
		c.Locals("db", db)

		return c.Next()
	}
}