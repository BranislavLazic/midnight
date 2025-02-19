package model

import "github.com/google/uuid"

type UserID struct{ uuid.UUID }

type User struct {
	ID       UserID `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"VARCHAR(255)" json:"email"`
	Password string `gorm:"VARCHAR(255)" json:"-"`
	Role     string `gorm:"VARCHAR(255)" json:"role"`
	Enabled  bool   `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
