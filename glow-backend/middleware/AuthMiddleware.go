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
	UserID uint `json:"user_id"` // Store user ID in claims
	jwt.RegisteredClaims
}

var jwtSecret []byte // Declare JWT secret as a byte slice

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	err := godotenv.Load() // Load .env file to get secret
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecretStr := os.Getenv("JWT_SECRET") //Get secret from .env
	if jwtSecretStr == "" {
		log.Fatal("JWT_SECRET not set in .env")
	}

	jwtSecret = []byte(jwtSecretStr) // Convert secret to byte slice

	return func(c *fiber.Ctx) error {
		// 1. Get token from the cookie
		tokenString := c.Cookies("token")

		// 2. If no token, check authorization header
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ") // Remove "Bearer " prefix
			}
		}

		// 3. Check if the token is missing
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Missing token"})
		}

		// 4. Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid signature"})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid token"})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid token"})
		}

		// 5. Check if the token is expired
		if claims.ExpiresAt != nil { // Check if expiry claim exists
			if !claims.ExpiresAt.Time.After(time.Now()) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Token has expired"})
			}
		} else {
			log.Println("Warning: Token missing expiration claim.  This is insecure.")
			// Consider this invalid in a real application, or extend expiry
		}

		// 6. Attach user ID to context
		c.Locals("user_id", claims.UserID) // Store user ID in context
		c.Locals("db", db)                // Store db in context

		// 7. Call the next handler
		return c.Next()
	}
}

