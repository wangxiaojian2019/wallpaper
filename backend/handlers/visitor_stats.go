package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wallpaper-shop/database"
	"wallpaper-shop/models"
)

type TrackVisitReq struct {
	Path      string `json:"path"`
	Language  string `json:"language"`
	Device    string `json:"device"`
	OS        string `json:"os"`
	Browser   string `json:"browser"`
	UserAgent string `json:"user_agent"`
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}

type HeartbeatReq struct {
	VisitID  string `json:"visit_id" binding:"required"`
	Duration int    `json:"duration"`
}

// IP Geo cache to avoid redundant external network requests
var (
	ipGeoCache sync.Map
	httpClient = &http.Client{Timeout: 900 * time.Millisecond}
)

type ipGeoResult struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
	Status      string `json:"status"`
}

// detectDeviceAndOS parses User-Agent to determine device type, OS, and browser
func detectDeviceAndOS(ua string, reqDevice, reqOS, reqBrowser string) (string, string, string) {
	device := reqDevice
	os := reqOS
	browser := reqBrowser
	uaLower := strings.ToLower(ua)

	if device == "" {
		if strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "tablet") {
			device = "Tablet"
		} else if strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "android") || strings.Contains(uaLower, "iphone") {
			device = "Mobile"
		} else {
			device = "Desktop"
		}
	}

	if os == "" {
		if strings.Contains(uaLower, "windows") {
			os = "Windows"
		} else if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "ios") {
			os = "iOS"
		} else if strings.Contains(uaLower, "android") {
			os = "Android"
		} else if strings.Contains(uaLower, "macintosh") || strings.Contains(uaLower, "mac os") {
			os = "macOS"
		} else if strings.Contains(uaLower, "linux") {
			os = "Linux"
		} else {
			os = "Other"
		}
	}

	if browser == "" {
		if strings.Contains(uaLower, "micromessenger") {
			browser = "WeChat"
		} else if strings.Contains(uaLower, "edg") {
			browser = "Edge"
		} else if strings.Contains(uaLower, "chrome") {
			browser = "Chrome"
		} else if strings.Contains(uaLower, "safari") {
			browser = "Safari"
		} else if strings.Contains(uaLower, "firefox") {
			browser = "Firefox"
		} else {
			browser = "Browser"
		}
	}

	return device, os, browser
}

// queryIPGeo queries IP location with cache
func queryIPGeo(ipStr string) (string, string, string) {
	if cached, ok := ipGeoCache.Load(ipStr); ok {
		res := cached.(ipGeoResult)
		return res.Country, res.CountryCode, res.City
	}

	resp, err := httpClient.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,countryCode,city", ipStr))
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", "", ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", ""
	}

	var res ipGeoResult
	if err := json.Unmarshal(body, &res); err == nil && res.Status == "success" {
		ipGeoCache.Store(ipStr, res)
		return res.Country, res.CountryCode, res.City
	}

	return "", "", ""
}

// detectIPRegion determines country, country code, default language, region scope, and currency
func detectIPRegion(c *gin.Context, clientIP string) (country, code, lang, scope, currency string) {
	// Default fallbacks
	scope = "china"
	currency = "CNY"
	lang = "zh"

	// 1. Check Cloudflare / CDN headers if present
	cfCountry := strings.ToUpper(strings.TrimSpace(c.GetHeader("CF-IPCountry")))
	if cfCountry == "" {
		cfCountry = strings.ToUpper(strings.TrimSpace(c.GetHeader("X-Country-Code")))
	}
	if cfCountry != "" && cfCountry != "XX" && cfCountry != "T1" {
		code = cfCountry
		isChinese := cfCountry == "CN" || cfCountry == "HK" || cfCountry == "MO" || cfCountry == "TW"
		if isChinese {
			country = "中国"
			lang = "zh"
			scope = "china"
			currency = "CNY"
		} else {
			country = cfCountry
			lang = "en"
			scope = "overseas"
			currency = "USD"
		}
		return
	}

	// 2. Check local/LAN private IPs
	parsedIP := net.ParseIP(clientIP)
	if parsedIP != nil {
		if parsedIP.IsLoopback() || parsedIP.IsPrivate() || clientIP == "127.0.0.1" || clientIP == "::1" {
			country = "本地局域网"
			code = "LAN"
			lang = "zh"
			scope = "china"
			currency = "CNY"
			return
		}
	}

	// 3. Fast external IP lookup for real public IPs
	if parsedIP != nil && !parsedIP.IsPrivate() && !parsedIP.IsLoopback() {
		geoCountry, geoCode, _ := queryIPGeo(clientIP)
		if geoCode != "" {
			code = strings.ToUpper(geoCode)
			country = geoCountry
			isChinese := code == "CN" || code == "HK" || code == "MO" || code == "TW"
			if isChinese {
				lang = "zh"
				scope = "china"
				currency = "CNY"
				country = "中国"
			} else {
				lang = "en"
				scope = "overseas"
				currency = "USD"
			}
			return
		}
	}

	// 4. Fallback to Accept-Language header
	acceptLang := strings.ToLower(c.GetHeader("Accept-Language"))
	if strings.HasPrefix(acceptLang, "zh") || strings.Contains(acceptLang, "zh-cn") || strings.Contains(acceptLang, "zh-tw") || strings.Contains(acceptLang, "zh-hk") {
		country = "华语区 (根据浏览器语言判定)"
		code = "ZH"
		lang = "zh"
		scope = "china"
		currency = "CNY"
		return
	}

	country = "海外国际区"
	code = "GLOBAL"
	lang = "en"
	scope = "overseas"
	currency = "USD"
	return
}

// TrackVisit records incoming client page visit
func TrackVisit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TrackVisitReq
		_ = c.ShouldBindJSON(&req)

		clientIP := c.ClientIP()
		ua := req.UserAgent
		if ua == "" {
			ua = c.GetHeader("User-Agent")
		}

		device, os, browser := detectDeviceAndOS(ua, req.Device, req.OS, req.Browser)
		country, _, defaultLang, scope, currency := detectIPRegion(c, clientIP)

		lang := req.Language
		if lang == "" {
			lang = defaultLang
		}

		// Identify logged-in user if provided
		var matchedUserID, matchedUserEmail, matchedTier string
		var isMember bool
		if req.UserEmail != "" {
			matchedUserEmail = strings.ToLower(strings.TrimSpace(req.UserEmail))
			var u models.User
			if err := database.DB.Where("email = ?", matchedUserEmail).First(&u).Error; err == nil {
				matchedUserID = u.ID
				matchedTier = u.MembershipType
				isMember = u.MembershipType != "none"
			}
		} else if req.UserID != "" {
			var u models.User
			if err := database.DB.Where("id = ?", req.UserID).First(&u).Error; err == nil {
				matchedUserID = u.ID
				matchedUserEmail = u.Email
				matchedTier = u.MembershipType
				isMember = u.MembershipType != "none"
			}
		}

		visit := models.VisitLog{
			ID:             uuid.New().String(),
			IP:             clientIP,
			Country:        country,
			City:           "",
			Device:         device,
			OS:             os,
			Browser:        browser,
			Path:           req.Path,
			Language:       lang,
			Duration:       1,
			UserAgent:      ua,
			UserID:         matchedUserID,
			UserEmail:      matchedUserEmail,
			MembershipType: matchedTier,
			IsMember:       isMember,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		database.DB.Create(&visit)

		c.JSON(http.StatusOK, gin.H{
			"visit_id":     visit.ID,
			"ip":           clientIP,
			"country":      country,
			"default_lang": defaultLang,
			"region_scope": scope,
			"currency":     currency,
			"user_email":   matchedUserEmail,
		})
	}
}

// HeartbeatVisit updates dwell duration of an existing visit
func HeartbeatVisit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req HeartbeatReq
		if err := c.ShouldBindJSON(&req); err != nil || req.VisitID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "visit_id is required"})
			return
		}

		duration := req.Duration
		if duration < 0 {
			duration = 0
		} else if duration > 86400 {
			duration = 86400 // Max 24 hours cap
		}

		database.DB.Model(&models.VisitLog{}).
			Where("id = ?", req.VisitID).
			Updates(map[string]interface{}{
				"duration":   duration,
				"updated_at": time.Now(),
			})

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// GetLocaleConfig returns IP geolocation and default language for adaptive client rendering
func GetLocaleConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		country, code, defaultLang, scope, currency := detectIPRegion(c, clientIP)

		var rateSetting models.SystemSetting
		rate := 7.25
		if err := database.DB.First(&rateSetting, "key = ?", "exchange_rate").Error; err == nil {
			if parsed, err := strconv.ParseFloat(rateSetting.Value, 64); err == nil && parsed > 0 {
				rate = parsed
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"ip":            clientIP,
			"country":       country,
			"country_code":  code,
			"default_lang":  defaultLang,
			"region_scope":  scope,
			"currency":      currency,
			"exchange_rate": rate,
		})
	}
}

// AdminVisitorList provides paginated visitor logs for visitor management
func AdminVisitorList() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}

		query := database.DB.Model(&models.VisitLog{})

		ipQuery := strings.TrimSpace(c.Query("ip"))
		if ipQuery != "" {
			query = query.Where("ip LIKE ?", "%"+ipQuery+"%")
		}

		deviceQuery := strings.TrimSpace(c.Query("device"))
		if deviceQuery != "" && deviceQuery != "all" {
			query = query.Where("LOWER(device) = ?", strings.ToLower(deviceQuery))
		}

		countryQuery := strings.TrimSpace(c.Query("country"))
		if countryQuery != "" {
			query = query.Where("country LIKE ?", "%"+countryQuery+"%")
		}

		var total int64
		query.Count(&total)

		var items []models.VisitLog
		offset := (page - 1) * pageSize
		query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&items)

		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		var todayPV int64
		database.DB.Model(&models.VisitLog{}).Where("created_at >= ?", todayStart).Count(&todayPV)

		var todayUV int64
		database.DB.Model(&models.VisitLog{}).Where("created_at >= ?", todayStart).Distinct("ip").Count(&todayUV)

		var avgDuration float64
		database.DB.Model(&models.VisitLog{}).Select("COALESCE(AVG(duration), 0)").Scan(&avgDuration)

		c.JSON(http.StatusOK, gin.H{
			"items":        items,
			"total":        total,
			"page":         page,
			"page_size":    pageSize,
			"today_pv":     todayPV,
			"today_uv":     todayUV,
			"avg_duration": math.Round(avgDuration),
		})
	}
}

// AdminVisitStats returns comprehensive traffic and visitor statistics for dashboard
func AdminVisitStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		// 1. Today PV / UV
		var todayPV int64
		database.DB.Model(&models.VisitLog{}).Where("created_at >= ?", todayStart).Count(&todayPV)

		var todayUV int64
		database.DB.Model(&models.VisitLog{}).Where("created_at >= ?", todayStart).Distinct("ip").Count(&todayUV)

		// 2. Total PV / UV
		var totalPV int64
		database.DB.Model(&models.VisitLog{}).Count(&totalPV)

		var totalUV int64
		database.DB.Model(&models.VisitLog{}).Distinct("ip").Count(&totalUV)

		// 3. Average dwell duration (in seconds)
		var avgDuration float64
		database.DB.Model(&models.VisitLog{}).Select("COALESCE(AVG(duration), 0)").Scan(&avgDuration)

		// 4. Device breakdown
		type KeyCount struct {
			Name  string `json:"name"`
			Count int64  `json:"count"`
		}
		var deviceRows []KeyCount
		database.DB.Model(&models.VisitLog{}).
			Select("device as name, count(*) as count").
			Group("device").
			Scan(&deviceRows)

		var osRows []KeyCount
		database.DB.Model(&models.VisitLog{}).
			Select("os as name, count(*) as count").
			Group("os").
			Order("count desc").
			Limit(5).
			Scan(&osRows)

		var countryRows []KeyCount
		database.DB.Model(&models.VisitLog{}).
			Select("country as name, count(*) as count").
			Group("country").
			Order("count desc").
			Limit(5).
			Scan(&countryRows)

		// 5. Recent 50 visits list
		var recentLogs []models.VisitLog
		database.DB.Order("created_at desc").Limit(50).Find(&recentLogs)

		c.JSON(http.StatusOK, gin.H{
			"today_pv":      todayPV,
			"today_uv":      todayUV,
			"total_pv":      totalPV,
			"total_uv":      totalUV,
			"avg_duration":  math.Round(avgDuration),
			"device_stats":  deviceRows,
			"os_stats":      osRows,
			"country_stats": countryRows,
			"recent_visits": recentLogs,
		})
	}
}
