package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wallpaper-shop/config"
	"wallpaper-shop/database"
	"wallpaper-shop/models"
)

// MembershipPlanSetting represents a customizable membership tier
type MembershipPlanSetting struct {
	Type         string  `json:"type"`          // "monthly", "quarterly", "yearly", "lifetime"
	NameZh       string  `json:"name_zh"`       // "月度会员"
	NameEn       string  `json:"name_en"`       // "Monthly VIP"
	Price        float64 `json:"price"`         // 9.99
	Period       string  `json:"period"`        // "month", "quarter", "year", "lifetime"
	DurationDays int     `json:"duration_days"` // 30, 90, 365, 36500
	Enabled      bool    `json:"enabled"`       // true / false
	Featured     bool    `json:"featured"`      // true / false
	Badge        string  `json:"badge"`         // e.g. "热销推荐 / 省17%"
	SavePercent  int     `json:"save_percent"`  // e.g. 17
	SortOrder    int     `json:"sort_order"`
	Icon         string  `json:"icon"`          // "🌙", "⭐", "👑", "💎"
}

// PermissionSettings represents global free download and checkout permissions
type PermissionSettings struct {
	FreeDownloadRequireLogin bool   `json:"free_download_require_login"` // 是否必须登录才能下载免费壁纸
	DailyFreeLimit           int    `json:"daily_free_limit"`            // 普通注册用户每日免费下载限制（0表示不限）
	GuestBuyEnabled          bool   `json:"guest_buy_enabled"`           // 是否开启临时游客免登录购买
	MasterUnlockCodes        string `json:"master_unlock_codes"`         // 全局万能口令 (逗号分隔，如 "loveyou, 666, 1194, VIP888")
	AllowTestPay             bool   `json:"allow_test_pay"`              // 是否开启快速模拟测试支付通道
}

// AllSettingsResponse bundles all system settings for the admin panel
type AllSettingsResponse struct {
	ExchangeRate float64 `json:"exchange_rate"`
	Payment struct {
		PayPalClientID     string `json:"paypal_client_id"`
		PayPalClientSecret string `json:"paypal_client_secret"`
		PayPalMode         string `json:"paypal_mode"`
		AllowTestPay       bool   `json:"allow_test_pay"`
	} `json:"payment"`
	MembershipPlans []MembershipPlanSetting `json:"membership_plans"`
	Permissions     PermissionSettings      `json:"permissions"`
}

// UpdateAllSettingsRequest incoming payload from admin panel
type UpdateAllSettingsRequest struct {
	ExchangeRate float64 `json:"exchange_rate"`
	Payment struct {
		PayPalClientID     string `json:"paypal_client_id"`
		PayPalClientSecret string `json:"paypal_client_secret"`
		PayPalMode         string `json:"paypal_mode"`
	} `json:"payment"`
	MembershipPlans []MembershipPlanSetting `json:"membership_plans"`
	Permissions     PermissionSettings      `json:"permissions"`
}

// GetDefaultMembershipPlans provides standard fallback plans
func GetDefaultMembershipPlans() []MembershipPlanSetting {
	return []MembershipPlanSetting{
		{
			Type:         "monthly",
			NameZh:       "月度会员",
			NameEn:       "Monthly VIP",
			Price:        9.99,
			Period:       "month",
			DurationDays: 30,
			Enabled:      true,
			Featured:     false,
			Badge:        "入门尝鲜",
			SavePercent:  0,
			SortOrder:    1,
			Icon:         "🌙",
		},
		{
			Type:         "quarterly",
			NameZh:       "季度会员",
			NameEn:       "Quarterly VIP",
			Price:        24.99,
			Period:       "quarter",
			DurationDays: 90,
			Enabled:      true,
			Featured:     true,
			Badge:        "热销推荐 / 省17%",
			SavePercent:  17,
			SortOrder:    2,
			Icon:         "⭐",
		},
		{
			Type:         "yearly",
			NameZh:       "年度会员",
			NameEn:       "Yearly VIP",
			Price:        79.99,
			Period:       "year",
			DurationDays: 365,
			Enabled:      true,
			Featured:     false,
			Badge:        "超值之选 / 省33%",
			SavePercent:  33,
			SortOrder:    3,
			Icon:         "👑",
		},
		{
			Type:         "lifetime",
			NameZh:       "终生尊享",
			NameEn:       "Lifetime VIP",
			Price:        199.99,
			Period:       "lifetime",
			DurationDays: 36500, // 100 years
			Enabled:      true,
			Featured:     false,
			Badge:        "终身买断 / 限时尊享",
			SavePercent:  50,
			SortOrder:    4,
			Icon:         "💎",
		},
	}
}

// GetEffectiveMembershipPlans loads from DB, falling back to defaults
func GetEffectiveMembershipPlans() []MembershipPlanSetting {
	var s models.SystemSetting
	if err := database.DB.First(&s, "key = ?", "membership_plans").Error; err == nil && s.Value != "" && !strings.Contains(s.Value, "???") {
		var plans []MembershipPlanSetting
		if err := json.Unmarshal([]byte(s.Value), &plans); err == nil && len(plans) > 0 {
			return plans
		}
	}
	return GetDefaultMembershipPlans()
}

// GetEffectivePermissionSettings loads permission configuration from DB
func GetEffectivePermissionSettings() PermissionSettings {
	defaultPerms := PermissionSettings{
		FreeDownloadRequireLogin: false,
		DailyFreeLimit:           0,
		GuestBuyEnabled:          true,
		MasterUnlockCodes:        "loveyou, 666, 1194, VIP888, VIP666",
		AllowTestPay:             true,
	}

	var s models.SystemSetting
	if err := database.DB.First(&s, "key = ?", "permission_settings").Error; err == nil && s.Value != "" {
		var perms PermissionSettings
		if err := json.Unmarshal([]byte(s.Value), &perms); err == nil {
			return perms
		}
	}
	return defaultPerms
}

// GetMembershipPlanByType finds a specific plan
func GetMembershipPlanByType(planType string) *MembershipPlanSetting {
	plans := GetEffectiveMembershipPlans()
	for _, p := range plans {
		if strings.EqualFold(p.Type, planType) {
			return &p
		}
	}
	return nil
}

// GetAllSystemSettings handles GET /api/admin/settings/all
func GetAllSystemSettings(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, clientSecret, mode := getEffectivePayPalConfig(cfg)

		maskedSecret := ""
		if len(clientSecret) > 6 {
			maskedSecret = clientSecret[:3] + "******" + clientSecret[len(clientSecret)-3:]
		} else if clientSecret != "" {
			maskedSecret = "******"
		}

		rate := 7.25
		var rateSetting models.SystemSetting
		if err := database.DB.First(&rateSetting, "key = ?", "exchange_rate").Error; err == nil {
			if parsed, err := strconv.ParseFloat(rateSetting.Value, 64); err == nil && parsed > 0 {
				rate = parsed
			}
		}

		res := AllSettingsResponse{
			ExchangeRate:    rate,
			MembershipPlans: GetEffectiveMembershipPlans(),
			Permissions:     GetEffectivePermissionSettings(),
		}
		res.Payment.PayPalClientID = clientID
		res.Payment.PayPalClientSecret = maskedSecret
		res.Payment.PayPalMode = mode
		res.Payment.AllowTestPay = res.Permissions.AllowTestPay

		c.JSON(http.StatusOK, res)
	}
}

// UpdateAllSystemSettings handles POST /api/admin/settings/all
func UpdateAllSystemSettings(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateAllSettingsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload: " + err.Error()})
			return
		}

		now := time.Now()

		// 0. Save Exchange Rate
		if req.ExchangeRate > 0 {
			rateStr := fmt.Sprintf("%.4f", req.ExchangeRate)
			database.DB.Save(&models.SystemSetting{Key: "exchange_rate", Value: strings.TrimRight(strings.TrimRight(rateStr, "0"), "."), UpdatedAt: now})
		}

		// 1. Save Payment Gateway
		if req.Payment.PayPalClientID != "" {
			database.DB.Save(&models.SystemSetting{Key: "paypal_client_id", Value: strings.TrimSpace(req.Payment.PayPalClientID), UpdatedAt: now})
		}
		if req.Payment.PayPalClientSecret != "" && !strings.Contains(req.Payment.PayPalClientSecret, "******") {
			database.DB.Save(&models.SystemSetting{Key: "paypal_client_secret", Value: strings.TrimSpace(req.Payment.PayPalClientSecret), UpdatedAt: now})
		}
		if req.Payment.PayPalMode != "" {
			database.DB.Save(&models.SystemSetting{Key: "paypal_mode", Value: strings.TrimSpace(req.Payment.PayPalMode), UpdatedAt: now})
		}

		// 2. Save Membership Plans
		if len(req.MembershipPlans) > 0 {
			plansBytes, err := json.Marshal(req.MembershipPlans)
			if err == nil {
				database.DB.Save(&models.SystemSetting{
					Key:       "membership_plans",
					Value:     string(plansBytes),
					UpdatedAt: now,
				})
			}
		}

		// 3. Save Permission Settings
		permBytes, err := json.Marshal(req.Permissions)
		if err == nil {
			database.DB.Save(&models.SystemSetting{
				Key:       "permission_settings",
				Value:     string(permBytes),
				UpdatedAt: now,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message":          "系统配置保存成功",
			"exchange_rate":    req.ExchangeRate,
			"membership_plans": req.MembershipPlans,
			"permissions":     req.Permissions,
		})
	}
}

// ListPublicPaymentChannels returns enabled payment channels for client checkout
func ListPublicPaymentChannels() gin.HandlerFunc {
	return func(c *gin.Context) {
		var channels []models.PaymentChannel
		database.DB.Where("enabled = ?", true).Order("sort_order asc, created_at asc").Find(&channels)

		clientIP := c.ClientIP()
		_, _, _, scope, _ := detectIPRegion(c, clientIP)

		// Filter or annotate based on detected visitor region
		var filtered []models.PaymentChannel
		for _, ch := range channels {
			// Global channels show to all
			if ch.Scope == "global" || ch.Scope == "" {
				filtered = append(filtered, ch)
			} else if scope == "china" && ch.Scope == "china_only" {
				filtered = append(filtered, ch)
			} else if scope == "overseas" && ch.Scope == "overseas_only" {
				filtered = append(filtered, ch)
			}
		}

		// If filtered is empty for whatever reason, fallback to all enabled channels
		if len(filtered) == 0 {
			filtered = channels
		}

		c.JSON(http.StatusOK, gin.H{
			"channels":     filtered,
			"region_scope": scope,
		})
	}
}

// AdminListPaymentChannels returns all payment channels for admin configuration
func AdminListPaymentChannels() gin.HandlerFunc {
	return func(c *gin.Context) {
		var channels []models.PaymentChannel
		database.DB.Order("sort_order asc, created_at asc").Find(&channels)
		c.JSON(http.StatusOK, gin.H{"channels": channels})
	}
}

// AdminCreatePaymentChannel creates a new payment channel
func AdminCreatePaymentChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ch models.PaymentChannel
		if err := c.ShouldBindJSON(&ch); err != nil || ch.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Channel name is required"})
			return
		}

		ch.ID = uuid.New().String()
		if ch.Scope == "" {
			ch.Scope = "global"
		}
		ch.CreatedAt = time.Now()
		ch.UpdatedAt = time.Now()

		if err := database.DB.Create(&ch).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create channel: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "支付渠道创建成功",
			"channel": ch,
		})
	}
}

// AdminUpdatePaymentChannel updates an existing payment channel
func AdminUpdatePaymentChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var ch models.PaymentChannel
		if err := database.DB.First(&ch, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Channel not found"})
			return
		}

		var req models.PaymentChannel
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
			return
		}

		if req.Name != "" {
			ch.Name = req.Name
		}
		if req.Code != "" {
			ch.Code = req.Code
		}
		if req.Scope != "" {
			ch.Scope = req.Scope
		}
		if req.Icon != "" {
			ch.Icon = req.Icon
		}
		ch.Description = req.Description
		ch.Enabled = req.Enabled
		ch.AccountInfo = req.AccountInfo
		ch.QRCodeURL = req.QRCodeURL
		ch.SortOrder = req.SortOrder
		ch.UpdatedAt = time.Now()

		database.DB.Save(&ch)

		c.JSON(http.StatusOK, gin.H{
			"message": "支付渠道已更新",
			"channel": ch,
		})
	}
}

// AdminDeletePaymentChannel deletes a payment channel
func AdminDeletePaymentChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		database.DB.Delete(&models.PaymentChannel{}, "id = ?", id)
		c.JSON(http.StatusOK, gin.H{"message": "支付渠道已删除"})
	}
}
