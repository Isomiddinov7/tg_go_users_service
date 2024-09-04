package models

import (
	"time"
)

type SignUpInput struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        string    `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponseByName struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Age       string    `json:"age,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
