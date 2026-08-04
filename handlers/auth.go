package handlers

import (
	"net/http"
	"time"

	"iscsi-web-panel/middleware"
	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// Login handles user login
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Mock: accept admin/admin
	if req.Username != "admin" || req.Password != "admin" {
		utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := middleware.GenerateToken(1, "admin", "admin")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.Success(c, models.LoginResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		User: models.User{
			ID:       1,
			Username: "admin",
			Role:     "admin",
			Email:    "admin@localhost",
		},
	})
}

// Register handles user registration
func Register(c *gin.Context) {
	utils.Success(c, gin.H{"message": "registration disabled in mock mode"})
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(c *gin.Context) {
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	utils.Success(c, models.User{
		ID:       1,
		Username: username.(string),
		Role:     role.(string),
		Email:    "admin@localhost",
	})
}
