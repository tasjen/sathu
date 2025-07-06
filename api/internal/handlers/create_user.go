package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/tasjen/fz/db"
)

// POST /users

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (handlers *Handlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var p NewUser

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = handlers.db.CreateUser(r.Context(), db.CreateUserParams{
		Username: p.Username,
		Password: p.Password,
		Email:    &p.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
