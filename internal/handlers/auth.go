package handlers

import (
	"github.com/gin-gonic/gin"
	"library-backend/internal/dto/request"
	"library-backend/internal/dto/response"
	"library-backend/internal/helpers"
	"net/http"
)

type AuthHandler struct {
	admins []struct {
		Username string
		Password string
	}
}

func NewAuthHandler() *AuthHandler {
	adminPassword, err := helpers.HashPassword("password")
	if err != nil {
		panic("failed to hash password for admin" + err.Error())
	}

	return &AuthHandler{
		admins: []struct {
			Username string
			Password string
		}{
			{"admin", adminPassword},
		},
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials request.LoginRequestDTO
	if err := c.ShouldBindJSON(&credentials); err != nil {
		helpers.HandleError(c, http.StatusBadRequest, "invalid payload")
		return
	}

	for _, admin := range h.admins {
		if admin.Username == credentials.Username && helpers.CheckPassword(credentials.Password, admin.Password) {
			token, err := helpers.GenerateToken(credentials.Username)
			if err != nil {
				helpers.HandleError(c, http.StatusInternalServerError, "failed to generate token")
				return
			}

			responseData := response.LoginResponseDTO{
				Token: token,
			}

			c.JSON(http.StatusOK, responseData)
			return
		}
	}

	helpers.HandleError(c, http.StatusUnauthorized, "invalid credentials")
}
