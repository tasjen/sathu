package handler

import (
	"encoding/json"
	"net/http"

	"github.com/tasjen/fz/api-hexa/internal/domain/model"
	"github.com/tasjen/fz/api-hexa/internal/domain/port"
)

type UserHandler struct {
	userService port.UserService
}

func NewUserHandler(userService port.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// POST /users
type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body registerRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userService.RegisterUser(r.Context(), &model.User{
		Username: newUser.Username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
