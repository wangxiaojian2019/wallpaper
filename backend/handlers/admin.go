package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"wallpaper-shop/database"
	"wallpaper-shop/models"
)

// AdminUserList returns paginated users
func AdminUserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize < 1 {
			pageSize = 20
		}

		query := database.DB.Model(&models.User{})

		membership := c.Query("membership_type")
		if membership != "" {
			query = query.Where("membership_type = ?", membership)
		}

		search := c.Query("search")
		if search != "" {
			pattern := "%" + strings.ToLower(search) + "%"
			query = query.Where("LOWER(email) LIKE ?", pattern)
		}

		var total int64
		query.Count(&total)

		var users []models.User
		offset := (page - 1) * pageSize
		query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&users)

		c.JSON(http.StatusOK, gin.H{
			"users":     users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

// AdminUpdateUser adjusts user role and membership
func AdminUpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		var req struct {
			Role           string `json:"role"`
			MembershipType string `json:"membership_type"`
			ExtendDays     int    `json:"extendDays"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid parameters"})
			return
		}

		if req.Role != "" {
			user.Role = req.Role
		}

		if req.MembershipType != "" {
			user.MembershipType = req.MembershipType
			if req.MembershipType == "none" {
				user.MembershipExpire = nil
			} else if req.ExtendDays > 0 {
				now := time.Now()
				if user.MembershipExpire != nil && user.MembershipExpire.After(now) {
					newExp := user.MembershipExpire.AddDate(0, 0, req.ExtendDays)
					user.MembershipExpire = &newExp
				} else {
					newExp := now.AddDate(0, 0, req.ExtendDays)
					user.MembershipExpire = &newExp
				}
			} else {
				// Default days per type if not set
				now := time.Now()
				var newExp time.Time
				switch req.MembershipType {
				case "monthly":
					newExp = now.AddDate(0, 1, 0)
				case "quarterly":
					newExp = now.AddDate(0, 3, 0)
				case "yearly":
					newExp = now.AddDate(1, 0, 0)
				}
				user.MembershipExpire = &newExp
			}
		}

		user.UpdatedAt = time.Now()
		database.DB.Save(&user)

		c.JSON(http.StatusOK, gin.H{
			"message": "User updated successfully",
			"user":    user,
		})
	}
}

// AdminCreateUser allows admin to manually create a user/member
func AdminCreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email          string `json:"email" binding:"required,email"`
			Password       string `json:"password" binding:"required,min=6"`
			Role           string `json:"role"`
			MembershipType string `json:"membership_type"`
			Days           int    `json:"days"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Valid email and password (min 6 chars) required"})
			return
		}

		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		var existing models.User
		if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email already exists"})
			return
		}

		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		role := req.Role
		if role == "" {
			role = "user"
		}
		mType := req.MembershipType
		if mType == "" {
			mType = "none"
		}

		var exp *time.Time
		if mType != "none" {
			now := time.Now()
			days := req.Days
			if days <= 0 {
				switch mType {
				case "monthly":
					days = 30
				case "quarterly":
					days = 90
				case "yearly":
					days = 365
				case "lifetime":
					days = 36500
				}
			}
			t := now.AddDate(0, 0, days)
			exp = &t
		}

		user := models.User{
			ID:               uuid.New().String(),
			Email:            req.Email,
			PasswordHash:     string(hashed),
			Role:             role,
			MembershipType:   mType,
			MembershipExpire: exp,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "用户创建成功",
			"user":    user,
		})
	}
}

// AdminDeleteUser removes a user
func AdminDeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		database.DB.Delete(&models.User{}, "id = ?", id)
		c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
	}
}

// AdminStats computes dashboard statistics
func AdminStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		var totalRevenue float64
		database.DB.Model(&models.Order{}).Where("status = ?", "paid").Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

		var todayRevenue float64
		todayStart := time.Now().Truncate(24 * time.Hour)
		database.DB.Model(&models.Order{}).Where("status = ? AND created_at >= ?", "paid", todayStart).Select("COALESCE(SUM(amount), 0)").Scan(&todayRevenue)

		var totalOrders int64
		database.DB.Model(&models.Order{}).Where("status = ?", "paid").Count(&totalOrders)

		var totalUsers int64
		database.DB.Model(&models.User{}).Count(&totalUsers)

		var singleRev float64
		database.DB.Model(&models.Order{}).Where("status = ? AND order_type = ?", "paid", "single").Select("COALESCE(SUM(amount), 0)").Scan(&singleRev)

		var monthlyRev float64
		database.DB.Model(&models.Order{}).Where("status = ? AND membership_type = ?", "paid", "monthly").Select("COALESCE(SUM(amount), 0)").Scan(&monthlyRev)

		var quarterlyRev float64
		database.DB.Model(&models.Order{}).Where("status = ? AND membership_type = ?", "paid", "quarterly").Select("COALESCE(SUM(amount), 0)").Scan(&quarterlyRev)

		var yearlyRev float64
		database.DB.Model(&models.Order{}).Where("status = ? AND membership_type = ?", "paid", "yearly").Select("COALESCE(SUM(amount), 0)").Scan(&yearlyRev)

		c.JSON(http.StatusOK, gin.H{
			"total_revenue":        totalRevenue,
			"today_revenue":        todayRevenue,
			"total_orders":         totalOrders,
			"total_users":          totalUsers,
			"single_sales_revenue": singleRev,
			"monthly_revenue":      monthlyRev,
			"quarterly_revenue":    quarterlyRev,
			"yearly_revenue":       yearlyRev,
		})
	}
}

// AdminResetUserPassword allows administrator to reset a user's password
func AdminResetUserPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		var req struct {
			Password string `json:"password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "新密码至少需要 6 个字符"})
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密计算失败"})
			return
		}

		user.PasswordHash = string(hashed)
		user.TokenVersion++ // Automatically invalidate user's previous sessions
		user.UpdatedAt = time.Now()
		database.DB.Save(&user)

		c.JSON(http.StatusOK, gin.H{
			"message": "用户密码已成功重置，旧登录凭据已同步强制失效",
		})
	}
}

// AdminForceLogoutUser invalidates all active sessions for a user
func AdminForceLogoutUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user models.User
		if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		user.TokenVersion++
		user.UpdatedAt = time.Now()
		database.DB.Save(&user)

		c.JSON(http.StatusOK, gin.H{
			"message": "该用户已成功强制下线，其当前所有客户端会话已被立即注销",
		})
	}
}

// AdminBatchProducts performs batch publish, unpublish, or delete on an array of product IDs
func AdminBatchProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Action string   `json:"action" binding:"required"` // "publish", "unpublish", "delete"
			IDs    []string `json:"ids" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请选择至少一个商品并指定操作指令"})
			return
		}

		switch req.Action {
		case "publish":
			database.DB.Model(&models.Product{}).Where("id IN ?", req.IDs).Update("is_published", true)
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("已成功批量上架 %d 款商品", len(req.IDs)),
				"count":   len(req.IDs),
			})
		case "unpublish":
			database.DB.Model(&models.Product{}).Where("id IN ?", req.IDs).Update("is_published", false)
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("已成功批量下架 %d 款商品", len(req.IDs)),
				"count":   len(req.IDs),
			})
		case "pin":
			now := time.Now()
			database.DB.Model(&models.Product{}).Where("id IN ?", req.IDs).
				Updates(map[string]interface{}{"is_pinned": true, "pinned_at": &now, "updated_at": now})
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("已成功批量置顶 %d 款商品", len(req.IDs)),
				"count":   len(req.IDs),
			})
		case "unpin":
			now := time.Now()
			database.DB.Model(&models.Product{}).Where("id IN ?", req.IDs).
				Updates(map[string]interface{}{"is_pinned": false, "pinned_at": nil, "updated_at": now})
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("已成功批量取消置顶 %d 款商品", len(req.IDs)),
				"count":   len(req.IDs),
			})
		case "delete":
			database.DB.Where("id IN ?", req.IDs).Delete(&models.Product{})
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("已成功批量删除 %d 款商品", len(req.IDs)),
				"count":   len(req.IDs),
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "未知的批量操作类型"})
		}
	}
}
