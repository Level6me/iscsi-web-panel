package handlers

import (
	"net/http"
	"time"

	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// ListUsers returns all users
func ListUsers(c *gin.Context) {
	mockUsers := []models.User{
		{ID: 1, Username: "admin", Role: "admin", Email: "admin@localhost", CreatedAt: time.Now().Add(-30 * 24 * time.Hour)},
		{ID: 2, Username: "operator", Role: "operator", Email: "ops@localhost", CreatedAt: time.Now().Add(-10 * 24 * time.Hour)},
		{ID: 3, Username: "viewer", Role: "viewer", Email: "view@localhost", CreatedAt: time.Now().Add(-5 * 24 * time.Hour)},
	}
	utils.Success(c, mockUsers)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "missing required fields")
		return
	}
	utils.Success(c, models.User{
		ID:        10,
		Username:  req.Username,
		Role:      req.Role,
		Email:     req.Email,
		CreatedAt: time.Now(),
	})
}

// GetUser returns a specific user
func GetUser(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, models.User{
		ID:        1,
		Username:  "admin",
		Role:      "admin",
		Email:     "admin@localhost",
		CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
	})
	_ = id
}

// UpdateUser updates a user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, gin.H{"id": id, "updated": true})
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, gin.H{"id": id, "deleted": true})
}
