package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var body registerRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.userService.RegisterUser(c.Request.Context(), model.User{
		Email:    body.Email,
		Password: &body.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
