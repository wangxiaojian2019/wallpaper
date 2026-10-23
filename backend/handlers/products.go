package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
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
	"wallpaper-shop/services"
)

type ProductForm struct {
	TitleZh         string                `json:"title_zh" binding:"required"`
	TitleEn         string                `json:"title_en"`
	Description     string                `json:"description"`
	CoverImage      string                `json:"cover_image"`
	PreviewVideoURL string                `json:"preview_video_url"`
	PreviewPlatform string                `json:"preview_platform"`
	PreviewVideos   []models.VideoItem    `json:"preview_videos"`
	Price           float64               `json:"price"`
	PriceUSD        float64               `json:"price_usd"`
	PriceCNY        float64               `json:"price_cny"`
	VerifyCode      string                `json:"verify_code"`
	Category        string                `json:"category"`
	Tags              []string              `json:"tags"`
	AllowedVIPLevels  []string              `json:"allowed_vip_levels"`
	Downloads         []models.DownloadItem `json:"downloads"`
	IsPublished       bool                  `json:"is_published"`
	IsPinned          bool                  `json:"is_pinned"`
	SortOrder         int                   `json:"sort_order"`
	IsInteractive     bool                  `json:"is_interactive"`
	InteractiveType   string                `json:"interactive_type"`
	InteractiveConfig string                `json:"interactive_config"`
}

// computeDualPrices calculates dual prices based on system exchange rate
func computeDualPrices(price, priceUSD, priceCNY float64) (finalBase, finalUSD, finalCNY float64) {
	rate := 7.25
	var rateSetting models.SystemSetting
	if err := database.DB.First(&rateSetting, "key = ?", "exchange_rate").Error; err == nil {
		if parsed, err := strconv.ParseFloat(rateSetting.Value, 64); err == nil && parsed > 0 {
			rate = parsed
		}
	}

	// 1. Both explicitly provided -> respect both
	if priceUSD > 0 && priceCNY > 0 {
		return priceUSD, priceUSD, priceCNY
	}

	// 2. Only USD provided -> auto-calc CNY
	if priceUSD > 0 && priceCNY == 0 {
		cny := math.Round(priceUSD*rate*100) / 100
		return priceUSD, priceUSD, cny
	}

	// 3. Only CNY provided -> auto-calc USD
	if priceCNY > 0 && priceUSD == 0 {
		usd := math.Round((priceCNY/rate)*100) / 100
		return usd, usd, priceCNY
	}

	// 4. Fallback to legacy price field
	if price > 0 {
		cny := math.Round(price*rate*100) / 100
		return price, price, cny
	}

	return 0, 0, 0
}

// IsVIPLevelAllowed checks if a given user's membership type is permitted to download the product
func IsVIPLevelAllowed(allowedJSON string, userMembershipType string) bool {
	if userMembershipType == "" || userMembershipType == "none" {
		return false
	}
	trimmed := strings.TrimSpace(allowedJSON)
	// If empty/null/unset, default to all VIP levels allowed for backward compatibility
	if trimmed == "" || trimmed == "null" {
		return true
	}
	var allowed []string
	if err := json.Unmarshal([]byte(trimmed), &allowed); err != nil {
		return true
	}
	if len(allowed) == 0 {
		return false
	}
	for _, lvl := range allowed {
		if strings.EqualFold(lvl, userMembershipType) {
			return true
		}
	}
	return false
}

// GetAllowedVIPDisplayNames returns human readable Chinese names for the allowed tiers
func GetAllowedVIPDisplayNames(allowedJSON string) string {
	trimmed := strings.TrimSpace(allowedJSON)
	if trimmed == "" || trimmed == "null" {
		return "全等级 VIP 会员"
	}
	var allowed []string
	if err := json.Unmarshal([]byte(trimmed), &allowed); err != nil || len(allowed) == 0 {
		return "全等级 VIP 会员"
	}
	plans := GetEffectiveMembershipPlans()
	planMap := make(map[string]string)
	for _, p := range plans {
		planMap[strings.ToLower(p.Type)] = p.NameZh
	}
	var names []string
	for _, lvl := range allowed {
		if name, ok := planMap[strings.ToLower(lvl)]; ok {
			names = append(names, name)
		} else {
			names = append(names, lvl)
		}
	}
	if len(names) == 0 {
		return "限单款购买"
	}
	return strings.Join(names, "、")
}

// ListProducts returns wallpapers with filters and pagination
func ListProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}

		query := database.DB.Model(&models.Product{}).Where("is_published = ?", true)

		category := c.Query("category")
		if category != "" && category != "all" {
			query = query.Where("LOWER(category) = ?", strings.ToLower(category))
		}

		search := c.Query("search")
		if search != "" {
			pattern := "%" + strings.ToLower(search) + "%"
			query = query.Where("LOWER(title_zh) LIKE ? OR LOWER(title_en) LIKE ? OR LOWER(tags) LIKE ?", pattern, pattern, pattern)
		}

		sort := c.DefaultQuery("sort", "newest")
		switch sort {
		case "popular":
			query = query.Order("is_pinned desc, pinned_at desc, sort_order desc, sales_count desc, views_count desc")
		case "price_asc":
			query = query.Order("is_pinned desc, pinned_at desc, sort_order desc, price asc")
		case "price_desc":
			query = query.Order("is_pinned desc, pinned_at desc, sort_order desc, price desc")
		default:
			query = query.Order("is_pinned desc, pinned_at desc, sort_order desc, created_at desc")
		}

		var total int64
		query.Count(&total)

		var products []models.Product
		query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

		c.JSON(http.StatusOK, gin.H{
			"items":     products,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

// GetProduct returns single wallpaper details
func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Wallpaper not found"})
			return
		}

		// Auto-heal / synchronize interactive tendies CAML layers if out of sync
		if product.IsInteractive && product.InteractiveConfig != "" {
			syncedConfig := syncInteractiveTendiesConfig(product.InteractiveConfig)
			if syncedConfig != product.InteractiveConfig {
				database.DB.Model(&product).UpdateColumn("interactive_config", syncedConfig)
				product.InteractiveConfig = syncedConfig
			}
		}

		// Increment view count
		database.DB.Model(&product).UpdateColumn("views_count", product.ViewsCount+1)

		// Check if user is unlocked
		isUnlocked := false
		requireLoginForFree := false
		dailyLimitReached := false
		var downloadList []models.DownloadItem

		perms := GetEffectivePermissionSettings()
		var currentUser *models.User
		if u, exists := c.Get("user"); exists {
			currentUser = u.(*models.User)
		}

		isVIPUser := currentUser != nil && currentUser.MembershipExpire != nil && currentUser.MembershipExpire.After(time.Now())
		vipAllowed := false
		if isVIPUser {
			vipAllowed = IsVIPLevelAllowed(product.AllowedVIPLevels, currentUser.MembershipType)
		}

		// 1. Free wallpaper check
		if product.Price <= 0 {
			if perms.FreeDownloadRequireLogin && currentUser == nil {
				// Login required by admin setting
				isUnlocked = false
				requireLoginForFree = true
			} else {
				// User is either anonymous (with login requirement turned off) or logged in
				if currentUser != nil && currentUser.Role != "admin" && !isVIPUser {
					// Regular user: check daily download limit
					if perms.DailyFreeLimit > 0 {
						var todayCount int64
						todayStart := time.Now().Truncate(24 * time.Hour)
						database.DB.Model(&models.DownloadRecord{}).
							Where("user_id = ? AND created_at >= ?", currentUser.ID, todayStart).
							Count(&todayCount)
						if todayCount >= int64(perms.DailyFreeLimit) {
							isUnlocked = false
							dailyLimitReached = true
						} else {
							isUnlocked = true
						}
					} else {
						isUnlocked = true
					}
				} else {
					isUnlocked = true
				}
			}
		}

		orderID := strings.TrimSpace(c.Query("order_id"))
		email := strings.ToLower(strings.TrimSpace(c.Query("email")))

		if !isUnlocked && orderID != "" {
			var orderCount int64
			database.DB.Model(&models.Order{}).
				Where("(id = ? OR paypal_order_id = ?) AND (product_id = ? OR order_type = 'membership') AND status = ?", orderID, orderID, product.ID, "paid").
				Count(&orderCount)
			if orderCount > 0 {
				isUnlocked = true
			}
		}

		if !isUnlocked && email != "" {
			var orderCount int64
			database.DB.Model(&models.Order{}).
				Where("LOWER(user_email) = ? AND (product_id = ? OR order_type = 'membership') AND status = ?", email, product.ID, "paid").
				Count(&orderCount)
			if orderCount > 0 {
				isUnlocked = true
			}
		}

		if !isUnlocked && currentUser != nil {
			// Admin always unlocked
			if currentUser.Role == "admin" {
				isUnlocked = true
			} else if isVIPUser {
				// Member unlocked only if tier is allowed
				if vipAllowed {
					isUnlocked = true
				}
			} else {
				// Check single order purchase by user_id or user's email
				var orderCount int64
				userEmail := strings.ToLower(strings.TrimSpace(currentUser.Email))
				database.DB.Model(&models.Order{}).
					Where("(user_id = ? OR (user_email != '' AND LOWER(user_email) = ?)) AND (product_id = ? OR order_type = 'membership') AND status = ?", currentUser.ID, userEmail, product.ID, "paid").
					Count(&orderCount)
				if orderCount > 0 {
					isUnlocked = true
				}
			}
		}

		if isUnlocked && product.Downloads != "" {
			_ = json.Unmarshal([]byte(product.Downloads), &downloadList)
		}

		c.JSON(http.StatusOK, gin.H{
			"product":                product,
			"is_unlocked":            isUnlocked,
			"require_login_for_free": requireLoginForFree,
			"daily_limit_reached":    dailyLimitReached,
			"daily_free_limit":       perms.DailyFreeLimit,
			"download_links":         downloadList,
			"downloads":              downloadList,
			"is_vip_user":            isVIPUser,
			"vip_allowed":            vipAllowed,
			"allowed_vip_names":      GetAllowedVIPDisplayNames(product.AllowedVIPLevels),
		})
	}
}

// GetDownloadLinks allows unlocked users (or any user if free) to access cloud drive download links
func GetDownloadLinks() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var product models.Product
		if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}

		// Authorization verification:
		authorized := false
		perms := GetEffectivePermissionSettings()
		var currentUser *models.User
		if u, exists := c.Get("user"); exists {
			currentUser = u.(*models.User)
		}

		orderID := strings.TrimSpace(c.Query("order_id"))
		email := strings.ToLower(strings.TrimSpace(c.Query("email")))
		code := strings.TrimSpace(c.Query("code"))

		if product.Price <= 0 {
			if perms.FreeDownloadRequireLogin && currentUser == nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message":                "本款免费壁纸需要登录账号后方可下载，请先登录或注册！",
					"require_login_for_free": true,
					"is_unlocked":            false,
				})
				return
			}

			// Check daily free limit for regular logged in users
			if currentUser != nil && currentUser.Role != "admin" && (currentUser.MembershipExpire == nil || currentUser.MembershipExpire.Before(time.Now())) {
				if perms.DailyFreeLimit > 0 {
					var todayCount int64
					todayStart := time.Now().Truncate(24 * time.Hour)
					database.DB.Model(&models.DownloadRecord{}).
						Where("user_id = ? AND created_at >= ?", currentUser.ID, todayStart).
						Count(&todayCount)
					if todayCount >= int64(perms.DailyFreeLimit) {
						c.JSON(http.StatusTooManyRequests, gin.H{
							"message":             fmt.Sprintf("您今日的免费壁纸下载已达上限（每日限额 %d 次），请明日再来或升级 VIP 会员畅享无限下载！", perms.DailyFreeLimit),
							"daily_limit_reached": true,
							"is_unlocked":         false,
						})
						return
					}
				}
			}

			authorized = true
		} else {
			// Paid product authorization checks:
			// 1. Verify code match
			if code != "" && product.VerifyCode != "" && strings.EqualFold(code, product.VerifyCode) {
				authorized = true
			}

			// 2. Admin role
			if !authorized && currentUser != nil && currentUser.Role == "admin" {
				authorized = true
			}

			// 3. User VIP membership
			vipTierMismatch := false
			var vipPlanName string
			if !authorized && currentUser != nil && currentUser.MembershipExpire != nil && currentUser.MembershipExpire.After(time.Now()) {
				if IsVIPLevelAllowed(product.AllowedVIPLevels, currentUser.MembershipType) {
					authorized = true
				} else {
					vipTierMismatch = true
					vipPlanName = currentUser.MembershipType
					if p := GetMembershipPlanByType(currentUser.MembershipType); p != nil && p.NameZh != "" {
						vipPlanName = p.NameZh
					}
				}
			}

			// 4. Order ID verification (works for guests, direct buyers, paypal)
			if !authorized && orderID != "" {
				var count int64
				database.DB.Model(&models.Order{}).
					Where("(id = ? OR paypal_order_id = ?) AND (product_id = ? OR order_type = 'membership') AND status = ?", orderID, orderID, product.ID, "paid").
					Count(&count)
				if count > 0 {
					authorized = true
				}
			}

			// 5. Email verification (works for guests who paid and filled email)
			if !authorized && email != "" {
				var count int64
				database.DB.Model(&models.Order{}).
					Where("LOWER(user_email) = ? AND (product_id = ? OR order_type = 'membership') AND status = ?", email, product.ID, "paid").
					Count(&count)
				if count > 0 {
					authorized = true
				}
			}

			// 6. Current logged-in user order check (by user_id or user's email)
			if !authorized && currentUser != nil {
				var count int64
				uEmail := strings.ToLower(strings.TrimSpace(currentUser.Email))
				database.DB.Model(&models.Order{}).
					Where("(user_id = ? OR (user_email != '' AND LOWER(user_email) = ?)) AND (product_id = ? OR order_type = 'membership') AND status = ?", currentUser.ID, uEmail, product.ID, "paid").
					Count(&count)
				if count > 0 {
					authorized = true
				}
			}

			// If still not authorized, check if it's a VIP tier mismatch or missing purchase
			if !authorized {
				if vipTierMismatch {
					c.JSON(http.StatusForbidden, gin.H{
						"message": fmt.Sprintf("本壁纸仅限【%s】免费下载，您当前会员等级为【%s】，请升级会员或单独购买！",
							GetAllowedVIPDisplayNames(product.AllowedVIPLevels), vipPlanName),
						"vip_tier_mismatch": true,
						"allowed_vip_names": GetAllowedVIPDisplayNames(product.AllowedVIPLevels),
						"is_unlocked":       false,
					})
					return
				}

				c.JSON(http.StatusForbidden, gin.H{
					"message":     "请先完成支付或开通相应 VIP 会员即可解锁网盘高速下载",
					"is_unlocked": false,
				})
				return
			}
		}

		// Record download record
		uid := ""
		uemail := email
		if currentUser != nil {
			uid = currentUser.ID
			if uemail == "" {
				uemail = currentUser.Email
			}
		}
		database.DB.Create(&models.DownloadRecord{
			ID:        uuid.New().String(),
			UserID:    uid,
			UserEmail: uemail,
			ProductID: product.ID,
			IP:        c.ClientIP(),
			CreatedAt: time.Now(),
		})

		var downloadList []models.DownloadItem
		if product.Downloads != "" {
			_ = json.Unmarshal([]byte(product.Downloads), &downloadList)
		}

		c.JSON(http.StatusOK, gin.H{
			"product_id":     product.ID,
			"links":          downloadList,
			"download_links": downloadList,
			"downloads":      downloadList,
			"is_unlocked":    true,
		})
	}
}

type UnlockCodeReq struct {
	Code      string `json:"code" binding:"required"`
	UserEmail string `json:"user_email"`
}

// UnlockWithCode verifies an unlock code and grants instant access
func UnlockWithCode(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req UnlockCodeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请输入解锁验证码"})
			return
		}

		var product models.Product
		if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}

		code := strings.TrimSpace(req.Code)
		isMatch := false
		var matchedRC *models.RedeemCode

		// 1. Check in custom RedeemCode table
		var rc models.RedeemCode
		if err := database.DB.Where("UPPER(code) = ? AND is_active = ?", strings.ToUpper(code), true).First(&rc).Error; err == nil {
			if rc.ExpiresAt == nil || rc.ExpiresAt.After(time.Now()) {
				if rc.MaxUses <= 0 || rc.UsedCount < rc.MaxUses {
					if rc.Purpose == "price_limit" {
						var prodPrice float64 = product.Price
						if rc.Currency == "CNY" && product.PriceCNY > 0 {
							prodPrice = product.PriceCNY
						} else if rc.Currency == "USD" && product.PriceUSD > 0 {
							prodPrice = product.PriceUSD
						}
						if rc.MaxPrice == 0 || prodPrice <= rc.MaxPrice {
							isMatch = true
							matchedRC = &rc
						}
					} else if rc.Purpose == "product" {
						if rc.ProductID == "" || rc.ProductID == product.ID {
							isMatch = true
							matchedRC = &rc
						}
					} else if rc.Purpose == "membership" {
						isMatch = true
						matchedRC = &rc
					}
				}
			}
		}

		// 2. Check product verify code
		if !isMatch && product.VerifyCode != "" && strings.EqualFold(code, product.VerifyCode) {
			isMatch = true
		}

		// 3. Master / author unlock codes from customizable system settings
		if !isMatch {
			perms := GetEffectivePermissionSettings()
			rawCodes := strings.Split(perms.MasterUnlockCodes, ",")
			for _, mc := range rawCodes {
				trimmed := strings.TrimSpace(mc)
				if trimmed != "" && strings.EqualFold(code, trimmed) {
					isMatch = true
					break
				}
			}
		}

		if !isMatch {
			c.JSON(http.StatusBadRequest, gin.H{"message": "解锁码无效或不正确，请重新输入"})
			return
		}

		// Resolve user
		var user models.User
		var tokenStr string
		if u, exists := c.Get("user"); exists {
			user = *(u.(*models.User))
			tokenStr, _ = middleware.GenerateToken(&user, cfg)
		} else {
			guestEmail := strings.ToLower(strings.TrimSpace(req.UserEmail))
			if guestEmail == "" {
				guestEmail = fmt.Sprintf("guest_%s@wallpaper.com", uuid.New().String()[:8])
			}
			if err := database.DB.Where("email = ?", guestEmail).First(&user).Error; err != nil {
				hashed, _ := bcrypt.GenerateFromPassword([]byte("guest123456"), bcrypt.DefaultCost)
				user = models.User{
					ID:             uuid.New().String(),
					Email:          guestEmail,
					PasswordHash:   string(hashed),
					Role:           "user",
					MembershipType: "none",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				database.DB.Create(&user)
			}
			tokenStr, _ = middleware.GenerateToken(&user, cfg)
		}

		// Record paid order
		order := models.Order{
			ID:            uuid.New().String(),
			UserID:        user.ID,
			UserEmail:     user.Email,
			OrderType:     "single",
			ProductID:     product.ID,
			ProductTitle:  product.TitleZh,
			Amount:        0,
			Currency:      "USD",
			PaypalOrderID: "CODE-" + code,
			Status:        "paid",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		database.DB.Create(&order)

		// Record download record
		rec := models.DownloadRecord{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			UserEmail: user.Email,
			ProductID: product.ID,
			IP:        c.ClientIP(),
			CreatedAt: time.Now(),
		}
		database.DB.Create(&rec)

		// Record redeem log and update code usage if from custom RedeemCode
		if matchedRC != nil {
			database.DB.Model(&models.RedeemCode{}).Where("id = ?", matchedRC.ID).Update("used_count", matchedRC.UsedCount+1)
			database.DB.Create(&models.RedeemLog{
				ID:           uuid.New().String(),
				CodeID:       matchedRC.ID,
				Code:         matchedRC.Code,
				UserID:       user.ID,
				UserEmail:    user.Email,
				Purpose:      matchedRC.Purpose,
				ProductID:    product.ID,
				ProductTitle: product.TitleZh,
				GrantedTier:  matchedRC.TargetTier,
				IP:           c.ClientIP(),
				CreatedAt:    time.Now(),
			})

			if matchedRC.Purpose == "membership" && matchedRC.TargetTier != "" {
				now := time.Now()
				var exp time.Time
				switch matchedRC.TargetTier {
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
				user.MembershipType = matchedRC.TargetTier
				user.MembershipExpire = &exp
				database.DB.Model(&user).Updates(map[string]interface{}{
					"membership_type":   matchedRC.TargetTier,
					"membership_expire": exp,
				})
			}
		}

		var downloadList []models.DownloadItem
		if product.Downloads != "" {
			_ = json.Unmarshal([]byte(product.Downloads), &downloadList)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":         "paid",
			"message":        "🎉 验证码正确，壁纸已成功解锁！",
			"token":          tokenStr,
			"is_unlocked":    true,
			"download_links": downloadList,
		})
	}
}

// AdminProductList allows admin to view both published and draft products
func AdminProductList() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
		if pageSize < 1 || pageSize > 100 {
			pageSize = 15
		}

		query := database.DB.Model(&models.Product{})

		category := c.Query("category")
		if category != "" {
			query = query.Where("LOWER(category) = ?", strings.ToLower(category))
		}

		search := c.Query("search")
		if search != "" {
			pattern := "%" + strings.ToLower(search) + "%"
			query = query.Where("LOWER(title_zh) LIKE ? OR LOWER(title_en) LIKE ? OR LOWER(tags) LIKE ?", pattern, pattern, pattern)
		}

		var total int64
		query.Count(&total)

		var products []models.Product
		query.Order("is_pinned desc, pinned_at desc, sort_order desc, created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)

		c.JSON(http.StatusOK, gin.H{
			"items":     products,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		})
	}
}

// AdminGetProduct returns complete wallpaper details for admin editing (including downloads)
func AdminGetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}

		var downloads []models.DownloadItem
		if product.Downloads != "" {
			_ = json.Unmarshal([]byte(product.Downloads), &downloads)
		}

		var tags []string
		if product.Tags != "" {
			_ = json.Unmarshal([]byte(product.Tags), &tags)
		}

		var previewVideos []models.VideoItem
		if product.PreviewVideos != "" {
			_ = json.Unmarshal([]byte(product.PreviewVideos), &previewVideos)
		}

		c.JSON(http.StatusOK, gin.H{
			"product": gin.H{
				"id":                 product.ID,
				"title_zh":           product.TitleZh,
				"title_en":           product.TitleEn,
				"description":        product.Description,
				"cover_image":        product.CoverImage,
				"preview_video_url":  product.PreviewVideoURL,
				"preview_platform":   product.PreviewPlatform,
				"preview_videos":     previewVideos,
				"price":              product.Price,
				"price_usd":          product.PriceUSD,
				"price_cny":          product.PriceCNY,
				"verify_code":        product.VerifyCode,
				"category":           product.Category,
				"tags":               tags,
				"allowed_vip_levels": product.AllowedVIPLevels,
				"downloads":          downloads,
				"is_published":       product.IsPublished,
				"is_pinned":          product.IsPinned,
				"pinned_at":          product.PinnedAt,
				"sort_order":         product.SortOrder,
				"is_interactive":     product.IsInteractive,
				"interactive_type":   product.InteractiveType,
				"interactive_config": product.InteractiveConfig,
				"views_count":        product.ViewsCount,
				"sales_count":        product.SalesCount,
				"created_at":         product.CreatedAt,
				"updated_at":         product.UpdatedAt,
			},
			"downloads":      downloads,
			"download_links": downloads,
		})
	}
}

// AdminCreateProduct adds new wallpaper item
func AdminCreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form ProductForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Title is required"})
			return
		}

		tagsJson, _ := json.Marshal(form.Tags)
		dlJson, _ := json.Marshal(form.Downloads)
		videosJson, _ := json.Marshal(form.PreviewVideos)

		platform := form.PreviewPlatform
		if platform == "" {
			platform = "douyin"
		}

		interactiveConfig := form.InteractiveConfig
		if form.IsInteractive && interactiveConfig != "" {
			interactiveConfig = syncInteractiveTendiesConfig(interactiveConfig)
		}

		allowedVIPJson := ""
		if len(form.AllowedVIPLevels) > 0 {
			b, _ := json.Marshal(form.AllowedVIPLevels)
			allowedVIPJson = string(b)
		} else {
			// Default to all enabled VIP levels on creation
			plans := GetEffectiveMembershipPlans()
			var defaultTiers []string
			for _, p := range plans {
				if p.Enabled {
					defaultTiers = append(defaultTiers, p.Type)
				}
			}
			if len(defaultTiers) == 0 {
				defaultTiers = []string{"monthly", "quarterly", "yearly", "lifetime"}
			}
			b, _ := json.Marshal(defaultTiers)
			allowedVIPJson = string(b)
		}

		basePrice, usdPrice, cnyPrice := computeDualPrices(form.Price, form.PriceUSD, form.PriceCNY)

		var pinnedAt *time.Time
		if form.IsPinned {
			now := time.Now()
			pinnedAt = &now
		}

		product := models.Product{
			ID:               uuid.New().String(),
			TitleZh:          form.TitleZh,
			TitleEn:          form.TitleEn,
			Description:      form.Description,
			CoverImage:       form.CoverImage,
			PreviewVideoURL:  form.PreviewVideoURL,
			PreviewPlatform:  platform,
			PreviewVideos:    string(videosJson),
			Price:            basePrice,
			PriceUSD:         usdPrice,
			PriceCNY:         cnyPrice,
			VerifyCode:       form.VerifyCode,
			Category:         form.Category,
			Tags:             string(tagsJson),
			AllowedVIPLevels: allowedVIPJson,
			Downloads:         string(dlJson),
			IsPublished:       form.IsPublished,
			IsPinned:          form.IsPinned,
			PinnedAt:          pinnedAt,
			SortOrder:         form.SortOrder,
			IsInteractive:     form.IsInteractive,
			InteractiveType:   form.InteractiveType,
			InteractiveConfig: interactiveConfig,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := database.DB.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Product created successfully",
			"product": product,
		})
	}
}

// AdminUpdateProduct edits existing wallpaper item
func AdminUpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}

		var form ProductForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		tagsJson, _ := json.Marshal(form.Tags)
		dlJson, _ := json.Marshal(form.Downloads)
		videosJson, _ := json.Marshal(form.PreviewVideos)

		platform := form.PreviewPlatform
		if platform == "" {
			platform = product.PreviewPlatform
			if platform == "" {
				platform = "douyin"
			}
		}

		product.TitleZh = form.TitleZh
		product.TitleEn = form.TitleEn
		product.Description = form.Description
		product.CoverImage = form.CoverImage
		product.PreviewVideoURL = form.PreviewVideoURL
		product.PreviewPlatform = platform
		product.PreviewVideos = string(videosJson)
		basePrice, usdPrice, cnyPrice := computeDualPrices(form.Price, form.PriceUSD, form.PriceCNY)
		product.Price = basePrice
		product.PriceUSD = usdPrice
		product.PriceCNY = cnyPrice
		product.VerifyCode = form.VerifyCode
		product.Category = form.Category
		product.Tags = string(tagsJson)
		if form.AllowedVIPLevels != nil {
			b, _ := json.Marshal(form.AllowedVIPLevels)
			product.AllowedVIPLevels = string(b)
		}
		product.Downloads = string(dlJson)
		product.IsPublished = form.IsPublished

		if form.IsPinned != product.IsPinned {
			product.IsPinned = form.IsPinned
			if form.IsPinned {
				now := time.Now()
				product.PinnedAt = &now
			} else {
				product.PinnedAt = nil
			}
		}
		product.SortOrder = form.SortOrder

		product.IsInteractive = form.IsInteractive
		product.InteractiveType = form.InteractiveType
		if form.IsInteractive && form.InteractiveConfig != "" {
			product.InteractiveConfig = syncInteractiveTendiesConfig(form.InteractiveConfig)
		} else {
			product.InteractiveConfig = form.InteractiveConfig
		}
		product.UpdatedAt = time.Now()

		if err := database.DB.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":        "Product updated successfully",
			"product":        product,
			"downloads":      form.Downloads,
			"download_links": form.Downloads,
		})
	}
}

// AdminToggleProductPin toggles is_pinned for a product
func AdminToggleProductPin() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var product models.Product
		if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "商品不存在"})
			return
		}

		product.IsPinned = !product.IsPinned
		if product.IsPinned {
			now := time.Now()
			product.PinnedAt = &now
		} else {
			product.PinnedAt = nil
		}
		product.UpdatedAt = time.Now()

		if err := database.DB.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "置顶状态更新失败"})
			return
		}

		statusMsg := "已成功置顶该商品 (排在最前)"
		if !product.IsPinned {
			statusMsg = "已取消置顶该商品"
		}

		c.JSON(http.StatusOK, gin.H{
			"message":   statusMsg,
			"is_pinned": product.IsPinned,
			"pinned_at": product.PinnedAt,
		})
	}
}

// AdminUpdateProductSort updates sort_order for a product
func AdminUpdateProductSort() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			SortOrder int `json:"sort_order"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请输入有效的排序数字权重"})
			return
		}

		if err := database.DB.Model(&models.Product{}).Where("id = ?", id).
			Updates(map[string]interface{}{"sort_order": req.SortOrder, "updated_at": time.Now()}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "排序权重更新失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "排序权重已更新",
			"sort_order": req.SortOrder,
		})
	}
}

// AdminDeleteProduct removes product
func AdminDeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := database.DB.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}

// syncInteractiveTendiesConfig verifies that CAML layers in interactive_config match the current tendies_file archive.
// If the archive is replaced or layers are missing or out of sync, it parses CAML and auto-updates all layer data.
func syncInteractiveTendiesConfig(configJSON string) string {
	if strings.TrimSpace(configJSON) == "" {
		return configJSON
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return configJSON
	}
	tendiesFile, _ := cfg["tendies_file"].(string)
	if strings.TrimSpace(tendiesFile) == "" {
		return configJSON
	}

	cleanFile := strings.TrimPrefix(tendiesFile, "/")
	zipPath := filepath.Join(".", filepath.FromSlash(cleanFile))
	if _, err := os.Stat(zipPath); err != nil {
		return configJSON
	}

	baseName := filepath.Base(tendiesFile)
	nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	parts := strings.Split(nameWithoutExt, "_")
	cleanID := ""
	if len(parts) >= 3 {
		cleanID = parts[2]
	} else if len(parts) == 2 {
		cleanID = parts[1]
	} else {
		cleanID = uuid.New().String()[:8]
	}

	// Check if already synced with this cleanID
	needsSync := true
	if bgLayers, ok := cfg["bg_layers"].([]interface{}); ok && len(bgLayers) > 0 {
		if firstLayer, ok := bgLayers[0].(map[string]interface{}); ok {
			layerID, _ := firstLayer["id"].(string)
			if layerID == "custom_bg" {
				// User manually configured custom single background
				needsSync = false
			} else if url, ok := firstLayer["url"].(string); ok && strings.Contains(url, "extracted_"+cleanID) {
				needsSync = false
			}
		}
	}

	if needsSync {
		uploadDir := filepath.Join(".", "uploads", "interactive")
		parsed, err := services.ParseTendiesArchive(zipPath, uploadDir, cleanID)
		if err == nil && parsed != nil {
			cfg["canvas_width"] = parsed.CanvasWidth
			cfg["canvas_height"] = parsed.CanvasHeight
			if parsed.BackgroundURL != "" {
				cfg["background_url"] = parsed.BackgroundURL
			}
			if parsed.ForegroundURL != "" {
				cfg["foreground_url"] = parsed.ForegroundURL
			}
			cfg["bg_layers"] = parsed.BGLayers
			cfg["floating_layers"] = parsed.FloatingLayers
			cfg["fg_layers"] = parsed.FGLayers
			cfg["total_image_count"] = parsed.TotalImageCount

			newBytes, err := json.Marshal(cfg)
			if err == nil {
				return string(newBytes)
			}
		}
	}
	return configJSON
}

