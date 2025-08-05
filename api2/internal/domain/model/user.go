package model

import "github.com/google/uuid"

type User struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Password        *string   `json:"-"`
	IsEmailVerified bool      `json:"is_email_verified"`
	Username        string    `json:"username"`
	Avatar          *string   `json:"avatar"`
	CreatedAt       string    `json:"created_at"`
}
