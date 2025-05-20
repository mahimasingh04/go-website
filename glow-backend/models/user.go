package models

import ("gorm.io/gorm"
"time"

)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	ID                uint      `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	// Add fields relevant to your user data here.
	SkinType          string    `gorm:"default:''" json:"skin_type"`          // User's Skin Type (e.g., "Oily", "Dry")
	ConversationHistory string    `gorm:"type:text;default:''" json:"conversation_history"` // Full Conversation History
	Location          string    `json:"location"`
	Age               int       `json:"age"`
}


type SignUpInput struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}