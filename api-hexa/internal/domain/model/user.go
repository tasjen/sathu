package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Avatar   *string   `json:"avatar"`
	Email    *string   `json:"email"`
}
