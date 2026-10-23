package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
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

// generateRandomCardCode generates an 8-character clean alphanumeric code
func generateRandomCardCode(prefix string) string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 8)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	if prefix != "" {
		return fmt.Sprintf("%s-%s", prefix, string(b))
	}
	return string(b)
}

// AdminListRedeemCodes returns list of redeem codes with filters
func AdminListRedeemCodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize < 1 {
			pageSize = 20
		}

		query := database.DB.Model(&models.RedeemCode{})

		purpose := c.Query("purpose")
		if purpose != "" {
			query = query.Where("purpose = ?", purpose)
		}

		search := c.Query("search")
		if search != "" {
			pattern := "%" + strings.ToLower(search) + "%"
			query = query.Where("LOWER(code) LIKE ? OR LOWER(remark) LIKE ?", pattern, pattern)
		}

		var total int64
		query.Count(&total)

		var codes []models.RedeemCode
		offset := (page - 1) * pageSize
		query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&codes)

		c.JSON(http.StatusOK, gin.H{
			"codes":     codes,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

type CreateRedeemCodeReq struct {
	Count       int        `json:"count"`       // Number of codes to batch generate (1-100)
	Code        string     `json:"code"`        // Custom code if count == 1
	Prefix      string     `json:"prefix"`      // Prefix for random codes e.g. "VIP", "FREE"
	Purpose     string     `json:"purpose"`     // "price_limit", "membership", "product"
	TargetTier  string     `json:"target_tier"` // "monthly", "quarterly", "yearly", "lifetime"
	MaxPrice    float64    `json:"max_price"`   // e.g. 5.0 (for 5元及以下商品免费兑换)
	Currency    string     `json:"currency"`    // "CNY", "USD"
	ProductID   string     `json:"product_id"`
	MaxUses     int        `json:"max_uses"`    // e.g. 1 or 100
	ExpiresAt   *time.Time `json:"expires_at"`  // Optional expiry date
	IsPermanent bool       `json:"is_permanent"`
	Remark      string     `json:"remark"`
}

// AdminCreateRedeemCode creates a custom or batch of redeem codes
func AdminCreateRedeemCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateRedeemCodeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
			return
		}

		if req.Purpose == "" {
			req.Purpose = "price_limit"
		}
		if req.MaxUses <= 0 {
			req.MaxUses = 1
		}
		if req.Currency == "" {
			req.Currency = "CNY"
		}
		if req.IsPermanent {
			req.ExpiresAt = nil
		}

		count := req.Count
		if count <= 0 {
			count = 1
		}
		if count > 200 {
			count = 200
		}

		var createdCodes []models.RedeemCode
		now := time.Now()

		for i := 0; i < count; i++ {
			codeStr := strings.TrimSpace(req.Code)
			if count > 1 || codeStr == "" {
				codeStr = generateRandomCardCode(req.Prefix)
			} else {
				codeStr = strings.ToUpper(codeStr)
			}

			// Ensure uniqueness
			var countExisting int64
			database.DB.Model(&models.RedeemCode{}).Where("code = ?", codeStr).Count(&countExisting)
			if countExisting > 0 {
				if count == 1 {
					c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("该口令【%s】已存在，请换一个口令", codeStr)})
					return
				}
				codeStr = generateRandomCardCode(req.Prefix)
			}

			rc := models.RedeemCode{
				ID:         uuid.New().String(),
				Code:       codeStr,
				Purpose:    req.Purpose,
				TargetTier: req.TargetTier,
				MaxPrice:   req.MaxPrice,
				Currency:   req.Currency,
				ProductID:  req.ProductID,
				MaxUses:    req.MaxUses,
				UsedCount:  0,
				ExpiresAt:  req.ExpiresAt,
				IsActive:   true,
				Remark:     req.Remark,
				CreatedAt:  now,
				UpdatedAt:  now,
			}

			if err := database.DB.Create(&rc).Error; err == nil {
				createdCodes = append(createdCodes, rc)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("已成功生成 %d 个兑换口令", len(createdCodes)),
			"codes":   createdCodes,
		})
	}
}

// AdminToggleRedeemCode toggles active status
func AdminToggleRedeemCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var code models.RedeemCode
		if err := database.DB.First(&code, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Redeem code not found"})
			return
		}

		code.IsActive = !code.IsActive
		code.UpdatedAt = time.Now()
		database.DB.Save(&code)

		c.JSON(http.StatusOK, gin.H{
			"message":   fmt.Sprintf("口令【%s】状态已更新为：%s", code.Code, map[bool]string{true: "启用", false: "禁用"}[code.IsActive]),
			"is_active": code.IsActive,
		})
	}
}

// AdminDeleteRedeemCode removes a redeem code
func AdminDeleteRedeemCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		database.DB.Delete(&models.RedeemCode{}, "id = ?", id)
		c.JSON(http.StatusOK, gin.H{"message": "口令已删除"})
	}
}

// AdminListRedeemLogs returns redeem logs
func AdminListRedeemLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize < 1 {
			pageSize = 20
		}

		query := database.DB.Model(&models.RedeemLog{})

		code := c.Query("code")
		if code != "" {
			query = query.Where("code = ?", strings.ToUpper(strings.TrimSpace(code)))
		}

		email := c.Query("email")
		if email != "" {
			query = query.Where("LOWER(user_email) LIKE ?", "%"+strings.ToLower(email)+"%")
		}

		var total int64
		query.Count(&total)

		var logs []models.RedeemLog
		offset := (page - 1) * pageSize
		query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&logs)

		c.JSON(http.StatusOK, gin.H{
			"logs":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

type ApplyRedeemReq struct {
	Code      string `json:"code" binding:"required"`
	ProductID string `json:"product_id"`
	UserEmail string `json:"user_email"`
}

// ApplyRedeemCode processes client redemption
func ApplyRedeemCode(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ApplyRedeemReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请输入兑换口令"})
			return
		}

		inputCode := strings.ToUpper(strings.TrimSpace(req.Code))
		if inputCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请输入有效的口令"})
			return
		}

		// 1. Check in RedeemCode table
		var rc models.RedeemCode
		if err := database.DB.Where("UPPER(code) = ?", inputCode).First(&rc).Error; err != nil {
			// Check fallback master unlock codes in permissions for backwards compatibility
			perms := GetEffectivePermissionSettings()
			rawCodes := strings.Split(perms.MasterUnlockCodes, ",")
			isMasterMatch := false
			for _, mc := range rawCodes {
				if strings.EqualFold(inputCode, strings.TrimSpace(mc)) {
					isMasterMatch = true
					break
				}
			}

			if !isMasterMatch {
				c.JSON(http.StatusBadRequest, gin.H{"message": "兑换口令不存在或已失效，请检查后重新输入"})
				return
			}

			// Master unlock code matched: grant direct product unlock
			rc = models.RedeemCode{
				ID:        uuid.New().String(),
				Code:      inputCode,
				Purpose:   "price_limit",
				MaxPrice:  999999,
				MaxUses:   999999,
				IsActive:  true,
				Remark:    "系统全局万能解锁码",
			}
		} else {
			// Validate code status
			if !rc.IsActive {
				c.JSON(http.StatusForbidden, gin.H{"message": "该兑换口令已被管理员停用"})
				return
			}
			if rc.ExpiresAt != nil && time.Now().After(*rc.ExpiresAt) {
				c.JSON(http.StatusForbidden, gin.H{
					"message": fmt.Sprintf("该口令已于 %s 过期失效", rc.ExpiresAt.Format("2006-01-02 15:04")),
				})
				return
			}
			if rc.MaxUses > 0 && rc.UsedCount >= rc.MaxUses {
				c.JSON(http.StatusForbidden, gin.H{
					"message": fmt.Sprintf("该口令已被兑换完毕（上限 %d 次）", rc.MaxUses),
				})
				return
			}
		}

		// 2. Resolve User (logged-in or by provided email)
		var user models.User
		var tokenStr string
		if u, exists := c.Get("user"); exists {
			user = *(u.(*models.User))
			tokenStr, _ = middleware.GenerateToken(&user, cfg)
		} else {
			targetEmail := strings.ToLower(strings.TrimSpace(req.UserEmail))
			if targetEmail == "" {
				targetEmail = fmt.Sprintf("guest_%s@wallart.vip", uuid.New().String()[:8])
			}
			if err := database.DB.Where("email = ?", targetEmail).First(&user).Error; err != nil {
				hashed, _ := bcrypt.GenerateFromPassword([]byte("guest123456"), bcrypt.DefaultCost)
				now := time.Now()
				user = models.User{
					ID:             uuid.New().String(),
					Email:          targetEmail,
					PasswordHash:   string(hashed),
					Role:           "user",
					MembershipType: "none",
					LastLoginAt:    &now,
					LastLoginIP:    c.ClientIP(),
					TokenVersion:   1,
					CreatedAt:      now,
					UpdatedAt:      now,
				}
				database.DB.Create(&user)
			}
			tokenStr, _ = middleware.GenerateToken(&user, cfg)
		}

		// 3. Process Redemption based on Purpose
		switch rc.Purpose {
		case "membership":
			// Grant membership
			tier := rc.TargetTier
			if tier == "" {
				tier = "monthly"
			}
			now := time.Now()
			var exp time.Time
			switch tier {
			case "monthly":
				exp = now.AddDate(0, 1, 0)
			case "quarterly":
				exp = now.AddDate(0, 3, 0)
			case "yearly":
				exp = now.AddDate(1, 0, 0)
			case "lifetime":
				exp = now.AddDate(100, 0, 0)
			default:
				exp = now.AddDate(0, 1, 0)
			}

			user.MembershipType = tier
			user.MembershipExpire = &exp
			user.UpdatedAt = now
			database.DB.Save(&user)

			// Update code count if real code
			if rc.ID != "" && rc.Remark != "系统全局万能解锁码" {
				database.DB.Model(&models.RedeemCode{}).Where("id = ?", rc.ID).Update("used_count", rc.UsedCount+1)
			}

			// Log redemption
			logItem := models.RedeemLog{
				ID:          uuid.New().String(),
				CodeID:      rc.ID,
				Code:        rc.Code,
				UserID:      user.ID,
				UserEmail:   user.Email,
				Purpose:     rc.Purpose,
				GrantedTier: tier,
				IP:          c.ClientIP(),
				CreatedAt:   now,
			}
			database.DB.Create(&logItem)

			tierNames := map[string]string{
				"monthly":   "月度 VIP 会员 (30天)",
				"quarterly": "季度 VIP 会员 (90天)",
				"yearly":    "年度 VIP 会员 (365天)",
				"lifetime":  "终身尊享 VIP (永久全免)",
			}

			c.JSON(http.StatusOK, gin.H{
				"success":           true,
				"message":           fmt.Sprintf("恭喜！已成功兑换【%s】特权", tierNames[tier]),
				"membership_type":   tier,
				"membership_expire": exp,
				"token":             tokenStr,
				"user": gin.H{
					"id":                user.ID,
					"email":             user.Email,
					"membership_type":   user.MembershipType,
					"membership_expire": user.MembershipExpire,
				},
			})
			return

		case "price_limit", "product":
			// Wallpaper unlock
			if req.ProductID == "" && rc.ProductID != "" {
				req.ProductID = rc.ProductID
			}
			if req.ProductID == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "该口令用于兑换壁纸，请在壁纸详情页收银台兑换，或指定壁纸商品 ID",
				})
				return
			}

			var product models.Product
			if err := database.DB.First(&product, "id = ?", req.ProductID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": "目标壁纸不存在"})
				return
			}

			if rc.Purpose == "product" && rc.ProductID != "" && rc.ProductID != product.ID {
				c.JSON(http.StatusBadRequest, gin.H{"message": "此口令仅限兑换指定的专属壁纸，与当前壁纸不匹配"})
				return
			}

			// Price limit check
			if rc.Purpose == "price_limit" && rc.MaxPrice > 0 {
				var prodPrice float64 = product.Price
				if rc.Currency == "CNY" && product.PriceCNY > 0 {
					prodPrice = product.PriceCNY
				} else if rc.Currency == "USD" && product.PriceUSD > 0 {
					prodPrice = product.PriceUSD
				}

				if prodPrice > rc.MaxPrice {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": fmt.Sprintf("该口令最高支持免费兑换 %.2f %s 及以下商品（当前商品为 %.2f %s）", rc.MaxPrice, rc.Currency, prodPrice, rc.Currency),
					})
					return
				}
			}

			now := time.Now()
			// Create paid order record
			orderID := uuid.New().String()
			order := models.Order{
				ID:           orderID,
				UserID:       user.ID,
				UserEmail:    user.Email,
				OrderType:    "single",
				ProductID:    product.ID,
				ProductTitle: product.TitleZh,
				Amount:       0,
				Status:       "paid",
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			database.DB.Create(&order)

			// Record download
			rec := models.DownloadRecord{
				ID:        uuid.New().String(),
				UserID:    user.ID,
				UserEmail: user.Email,
				ProductID: product.ID,
				IP:        c.ClientIP(),
				CreatedAt: now,
			}
			database.DB.Create(&rec)

			// Increment used count
			if rc.ID != "" && rc.Remark != "系统全局万能解锁码" {
				database.DB.Model(&models.RedeemCode{}).Where("id = ?", rc.ID).Update("used_count", rc.UsedCount+1)
			}

			// Log redemption
			logItem := models.RedeemLog{
				ID:           uuid.New().String(),
				CodeID:       rc.ID,
				Code:         rc.Code,
				UserID:       user.ID,
				UserEmail:    user.Email,
				Purpose:      rc.Purpose,
				ProductID:    product.ID,
				ProductTitle: product.TitleZh,
				IP:           c.ClientIP(),
				CreatedAt:    now,
			}
			database.DB.Create(&logItem)

			c.JSON(http.StatusOK, gin.H{
				"success":      true,
				"message":      "口令兑换成功！已为您免费解锁本款高清壁纸",
				"order_id":     orderID,
				"product_id":   product.ID,
				"token":        tokenStr,
				"download_url": fmt.Sprintf("/pages/download/index?id=%s&title=%s", product.ID, product.TitleZh),
			})
			return

		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "未知的口令用途类型"})
			return
		}
	}
}
