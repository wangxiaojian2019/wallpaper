package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"wallpaper-shop/config"
	"wallpaper-shop/database"
	"wallpaper-shop/middleware"
	"wallpaper-shop/models"
)

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password (min 6 characters)"})
			return
		}

		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		// Check if email already taken
		var existing models.User
		if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email already registered"})
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
			return
		}

		now := time.Now()
		clientIP := c.ClientIP()
		user := models.User{
			ID:             uuid.New().String(),
			Email:          req.Email,
			PasswordHash:   string(hashed),
			Role:           "user",
			MembershipType: "none",
			LastLoginAt:    &now,
			LastLoginIP:    clientIP,
			TokenVersion:   1,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
			return
		}

		token, err := middleware.GenerateToken(&user, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to issue token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Registration successful",
			"token":   token,
			"user": gin.H{
				"id":                user.ID,
				"email":             user.Email,
				"role":              user.Role,
				"membership_type":   user.MembershipType,
				"membership_expire": user.MembershipExpire,
				"last_login_at":     user.LastLoginAt,
				"last_login_ip":     user.LastLoginIP,
			},
		})
	}
}

func Login(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email and password are required"})
			return
		}

		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		var user models.User
		if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
			return
		}

		// Update last login timestamp and IP
		now := time.Now()
		clientIP := c.ClientIP()
		user.LastLoginAt = &now
		user.LastLoginIP = clientIP
		if user.TokenVersion == 0 {
			user.TokenVersion = 1
		}
		database.DB.Model(&user).Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": clientIP,
			"token_version": user.TokenVersion,
		})

		token, err := middleware.GenerateToken(&user, cfg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to issue token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"token":   token,
			"user": gin.H{
				"id":                user.ID,
				"email":             user.Email,
				"role":              user.Role,
				"membership_type":   user.MembershipType,
				"membership_expire": user.MembershipExpire,
				"last_login_at":     user.LastLoginAt,
				"last_login_ip":     user.LastLoginIP,
			},
		})
	}
}

func Me() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		user := u.(*models.User)
		// Fetch fresh data from DB
		var fresh models.User
		database.DB.First(&fresh, "id = ?", user.ID)

		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":                fresh.ID,
				"email":             fresh.Email,
				"role":              fresh.Role,
				"membership_type":   fresh.MembershipType,
				"membership_expire": fresh.MembershipExpire,
			},
		})
	}
}
