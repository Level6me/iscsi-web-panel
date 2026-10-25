package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"iscsi-web-panel/middleware"
	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// 管理口令来源：环境变量 ADMIN_PASSWORD 优先；未设置时首启生成随机口令并落盘，
// 拒绝任何硬编码默认口令（修复 mock 弱口令 admin/admin123）
var (
	adminPasswordOnce sync.Once
	adminPassword     string
)

func resolveAdminPassword() string {
	adminPasswordOnce.Do(func() {
		if p := os.Getenv("ADMIN_PASSWORD"); p != "" {
			adminPassword = p
			return
		}
		const path = "data/admin_password"
		if b, err := os.ReadFile(path); err == nil && len(b) > 0 {
			adminPassword = string(b)
			return
		}
		buf := make([]byte, 16)
		if _, err := rand.Read(buf); err != nil {
			adminPassword = ""
			return
		}
		adminPassword = hex.EncodeToString(buf)
		_ = os.MkdirAll("data", 0o700)
		if err := os.WriteFile(path, []byte(adminPassword), 0o600); err != nil {
			log.Printf("[auth] failed to persist generated admin password: %v", err)
			return
		}
		log.Printf("[auth] generated initial admin password: %s (stored in %s, change via ADMIN_PASSWORD env)", adminPassword, path)
	})
	return adminPassword
}

// Login handles user login
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	expected := resolveAdminPassword()
	if expected == "" || req.Username != "admin" || req.Password != expected {
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
