package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
	"github.com/mahimasingh04/go-website/glow-backend/models"
	"gorm.io/gorm"
	"os"
	
	
)


//Signup godoc

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) SignUpUser(c *fiber.Ctx) error {
	var payload models.SignUpInput

	if err := c.BodyParser(&payload); err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	
   hashedPassword, err:= bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
   if err != nil{
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  "fail",
		"message": "Error hashing password",
	})
   }

   newUser := models.User{
	Name : payload.Name,
	Email : payload.Email,
	Password : string(hashedPassword),
   }

   result := ac.DB.Create(&newUser) 
   if result.Error != nil {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status" :"fail",
		"message": "Error creating user",
	})
   }
   return c.Status(fiber.StatusCreated).JSON(fiber.Map{
	"status":  "success",
	"data": fiber.Map{
		"id": newUser.ID,
		"name": newUser.Name,
		"email": newUser.Email,
	},
})
}

//SignIn User 

func (ac *AuthController) SignInUser(c *fiber.Ctx) error {
	var payload models.SignInInput
	if err:= c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	var user models.User
	result := ac.DB.First(&user, "email= ?", payload.Email)
	if result.Error !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"message" :"Invalid email or password",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password),
	[]byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"message": "Invalid email or password",
		})
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail",
			"message": "JWT secret not configured",
		})
	}

	//Generate JWT token 
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["nbf"] = time.Now().Unix()



	tokenString , err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status" : "fail",
			"message": "Error signing token",
		})
	}

	cookie := fiber.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure: true,
		SameSite: "Lax",
	}
	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	"status": "success",
	"token" : tokenString,
     " data" : fiber.Map{
     "id" :user.ID,
	 "name" : user.Name,
}	,
})
}