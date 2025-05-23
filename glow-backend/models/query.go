package models

import ("time"
"gorm.io/gorm")

type Query struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Category    string    `json:"category"`   // e.g., "SkinType", "Disease", "Routine", "Product", "Food"
	QueryText   string    `gorm:"type:text" json:"query_text"` // User's query
	IdealResponse string `gorm:"type:text" json:"ideal_response"`   // Expected/ideal AI response
}

