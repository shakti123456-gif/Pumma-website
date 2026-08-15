package handlers

import (
	"ecommerce-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

type AuthHandler struct {
	Email    string
	Password string
	Secret   string
}

func NewAuthHandler() AuthHandler {
	return AuthHandler{Email: env("ADMIN_EMAIL", "admin@example.com"), Password: env("ADMIN_PASSWORD", "change-me-before-production"), Secret: env("JWT_SECRET", "development-secret-change-before-production")}
}
func (h AuthHandler) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if c.ShouldBindJSON(&request) != nil || request.Email != h.Email || request.Password != h.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	token, err := auth.Issue(h.Email, h.Secret)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not create login token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "expiresIn": 86400})
}
func (h AuthHandler) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "admin login required"})
			return
		}
		if _, err := auth.Verify(token, h.Secret); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired login token"})
			return
		}
		c.Next()
	}
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
