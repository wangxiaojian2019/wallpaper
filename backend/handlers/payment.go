package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wallpaper-shop/config"
	"wallpaper-shop/database"
	"wallpaper-shop/middleware"
	"wallpaper-shop/models"
)

type PayPalTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type CreateOrderReq struct {
	ProductID      string  `json:"product_id"`
	OrderType      string  `json:"order_type"`      // "single", "membership"
	MembershipType string  `json:"membership_type"` // "monthly", "quarterly", "yearly"
	UserEmail      string  `json:"user_email"`      // Optional for guest checkout
	PaymentMethod  string  `json:"payment_method"`  // "alipay", "wechat", "paypal", "test", etc.
	IsTest         bool    `json:"is_test"`         // Sandbox fast test checkout
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
}

// getEffectivePayPalConfig retrieves credentials, prioritizing database SystemSettings
func getEffectivePayPalConfig(cfg *config.Config) (clientID, clientSecret, mode string) {
	clientID = cfg.PayPalClientID
	clientSecret = cfg.PayPalClientSecret
	mode = cfg.PayPalMode

	var sID models.SystemSetting
	if err := database.DB.First(&sID, "key = ?", "paypal_client_id").Error; err == nil && sID.Value != "" {
		clientID = sID.Value
	}
	var sSecret models.SystemSetting
	if err := database.DB.First(&sSecret, "key = ?", "paypal_client_secret").Error; err == nil && sSecret.Value != "" {
		clientSecret = sSecret.Value
	}
	var sMode models.SystemSetting
	if err := database.DB.First(&sMode, "key = ?", "paypal_mode").Error; err == nil && sMode.Value != "" {
		mode = sMode.Value
	}
	if mode == "" {
		mode = "sandbox"
	}
	return
}

func getPayPalBaseURL(mode string) string {
	if mode == "live" {
		return "https://api-m.paypal.com"
	}
	return "https://api-m.sandbox.paypal.com"
}

func getPayPalAccessToken(cfg *config.Config) (string, error) {
	clientID, clientSecret, mode := getEffectivePayPalConfig(cfg)
	if clientSecret == "" || clientID == "" || clientID == "sb" {
		return "", fmt.Errorf("paypal credentials not configured")
	}

	url := fmt.Sprintf("%s/v1/oauth2/token", getPayPalBaseURL(mode))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("grant_type=client_credentials")))
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("paypal auth failed: %s", string(body))
	}

	var tokenRes PayPalTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return "", err
	}

	return tokenRes.AccessToken, nil
}

// CreatePayPalOrder creates pending/test order and PayPal transaction
func CreatePayPalOrder(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateOrderReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
			return
		}

		perms := GetEffectivePermissionSettings()

		// 1. Resolve User (Logged-in user or guest user by email)
		var user models.User
		var tokenStr string

		if u, exists := c.Get("user"); exists {
			user = *(u.(*models.User))
			tokenStr, _ = middleware.GenerateToken(&user, cfg)
		} else {
			// Check if guest checkout is permitted
			if !perms.GuestBuyEnabled {
				c.JSON(http.StatusForbidden, gin.H{
					"message": "临时游客免登录购买通道当前已关闭，请先注册或登录您的账号！",
					"require_login": true,
				})
				return
			}

			// Guest checkout
			guestEmail := strings.ToLower(strings.TrimSpace(req.UserEmail))
			if guestEmail == "" || !strings.Contains(guestEmail, "@") || !strings.Contains(guestEmail, ".") {
				c.JSON(http.StatusBadRequest, gin.H{"message": "请填写有效的联系邮箱，以便接收壁纸网盘链接与查询订单凭证"})
				return
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

		// 2. Determine Order Type, Currency and Amount
		currency := strings.ToUpper(strings.TrimSpace(req.Currency))
		if currency == "" {
			if strings.HasPrefix(strings.ToLower(req.PaymentMethod), "alipay") ||
				strings.HasPrefix(strings.ToLower(req.PaymentMethod), "wechat") ||
				req.PaymentMethod == "custom" {
				currency = "CNY"
			} else {
				currency = "USD"
			}
		}

		var amount float64
		var title string
		var matchedPlan *MembershipPlanSetting
		orderType := req.OrderType
		if orderType == "" {
			if req.MembershipType != "" {
				orderType = "membership"
			} else {
				orderType = "single"
			}
		}

		if req.Amount > 0 {
			// Client explicitly passed the price displayed on the cashier
			amount = req.Amount
			if orderType == "membership" {
				matchedPlan = GetMembershipPlanByType(req.MembershipType)
				if matchedPlan != nil {
					title = fmt.Sprintf("VIP %s (%s)", matchedPlan.NameZh, matchedPlan.NameEn)
				} else {
					title = "VIP Membership"
				}
			} else {
				var product models.Product
				if err := database.DB.First(&product, "id = ?", req.ProductID).Error; err == nil {
					title = product.TitleZh
					if title == "" {
						title = product.TitleEn
					}
				}
			}
		} else if orderType == "membership" {
			matchedPlan = GetMembershipPlanByType(req.MembershipType)
			if matchedPlan == nil || !matchedPlan.Enabled {
				c.JSON(http.StatusBadRequest, gin.H{"message": "该会员方案不存在或已被下架"})
				return
			}
			if currency == "CNY" {
				rate := 7.25
				var rateSetting models.SystemSetting
				if err := database.DB.First(&rateSetting, "key = ?", "exchange_rate").Error; err == nil {
					if parsed, err := strconv.ParseFloat(rateSetting.Value, 64); err == nil && parsed > 0 {
						rate = parsed
					}
				}
				amount = math.Round(matchedPlan.Price*rate*100) / 100
			} else {
				amount = matchedPlan.Price
			}
			title = fmt.Sprintf("VIP %s (%s)", matchedPlan.NameZh, matchedPlan.NameEn)
		} else {
			// Single wallpaper order
			var product models.Product
			if err := database.DB.First(&product, "id = ?", req.ProductID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
				return
			}
			title = product.TitleZh
			if title == "" {
				title = product.TitleEn
			}
			if currency == "CNY" && product.PriceCNY > 0 {
				amount = product.PriceCNY
			} else if currency == "USD" && product.PriceUSD > 0 {
				amount = product.PriceUSD
			} else {
				amount = product.Price
			}
		}

		orderID := uuid.New().String()
		paypalOrderId := "PAYPAL-" + uuid.New().String()[:12]
		approvalURL := fmt.Sprintf("%s/payment/success?order_id=%s&token=%s", cfg.AppBaseURL, orderID, paypalOrderId)

		// 3. Fast Test / Sandbox Checkout Mode
		if req.IsTest {
			if !perms.AllowTestPay {
				c.JSON(http.StatusForbidden, gin.H{"message": "模拟测试支付通道已被管理员关闭"})
				return
			}

			order := models.Order{
				ID:             orderID,
				UserID:         user.ID,
				UserEmail:      user.Email,
				OrderType:      orderType,
				ProductID:      req.ProductID,
				ProductTitle:   title,
				MembershipType: req.MembershipType,
				Amount:         amount,
				Currency:       currency,
				PaypalOrderID:  "TEST-" + uuid.New().String()[:8],
				Status:         "paid",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			database.DB.Create(&order)

			// Fulfill immediate access
			if orderType == "membership" {
				now := time.Now()
				var newExp time.Time
				if user.MembershipExpire != nil && user.MembershipExpire.After(now) {
					newExp = *user.MembershipExpire
				} else {
					newExp = now
				}

				if matchedPlan != nil && matchedPlan.DurationDays > 0 {
					newExp = newExp.AddDate(0, 0, matchedPlan.DurationDays)
				} else {
					switch req.MembershipType {
					case "monthly":
						newExp = newExp.AddDate(0, 1, 0)
					case "quarterly":
						newExp = newExp.AddDate(0, 3, 0)
					case "yearly":
						newExp = newExp.AddDate(1, 0, 0)
					case "lifetime":
						newExp = newExp.AddDate(99, 0, 0)
					default:
						newExp = newExp.AddDate(0, 1, 0)
					}
				}
				user.MembershipType = req.MembershipType
				user.MembershipExpire = &newExp
				database.DB.Save(&user)
			} else if req.ProductID != "" {
				database.DB.Model(&models.Product{}).Where("id = ?", req.ProductID).
					UpdateColumn("sales_count", gorm.Expr("sales_count + 1"))
			}

			// Issue fresh token with updated membership
			tokenStr, _ = middleware.GenerateToken(&user, cfg)

			c.JSON(http.StatusOK, gin.H{
				"order_id": order.ID,
				"status":   "paid",
				"is_test":  true,
				"token":    tokenStr,
				"user": gin.H{
					"id":                user.ID,
					"email":             user.Email,
					"membership_type":   user.MembershipType,
					"membership_expire": user.MembershipExpire,
				},
				"message": "⚡ 模拟测试支付成功，已解锁对应权益！",
			})
			return
		}

		// 3.5 Alipay, WeChat, or custom payment channels (creates pending order with secondary confirmation support)
		if req.PaymentMethod != "paypal" && req.PaymentMethod != "" && !req.IsTest {
			var ch models.PaymentChannel
			database.DB.Where("code = ?", req.PaymentMethod).First(&ch)
			channelName := ch.Name
			if channelName == "" {
				channelName = req.PaymentMethod
			}

			customRefID := fmt.Sprintf("%s-%s", strings.ToUpper(req.PaymentMethod), orderID[:8])

			order := models.Order{
				ID:             orderID,
				UserID:         user.ID,
				UserEmail:      user.Email,
				OrderType:      orderType,
				ProductID:      req.ProductID,
				ProductTitle:   title,
				MembershipType: req.MembershipType,
				Amount:         amount,
				Currency:       currency,
				PaypalOrderID:  customRefID,
				Status:         "pending",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := database.DB.Create(&order).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "创建订单记录失败: " + err.Error()})
				return
			}

			tokenStr, _ = middleware.GenerateToken(&user, cfg)

			c.JSON(http.StatusOK, gin.H{
				"order_id":         order.ID,
				"status":           "pending",
				"payment_method":   req.PaymentMethod,
				"channel_name":     channelName,
				"account_info":     ch.AccountInfo,
				"qr_code_url":      ch.QRCodeURL,
				"amount":           amount,
				"currency":         currency,
				"token":            tokenStr,
				"user": gin.H{
					"id":                user.ID,
					"email":             user.Email,
					"membership_type":   user.MembershipType,
					"membership_expire": user.MembershipExpire,
				},
				"message": fmt.Sprintf("【%s】支付订单已生成，请核对信息并完成支付确认", channelName),
			})
			return
		}

		// 4. Real PayPal Orders v2 API Call
		clientID, clientSecret, mode := getEffectivePayPalConfig(cfg)
		if clientID == "" || clientID == "sb" || clientSecret == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "PayPal 商业收款账号未配置 (缺少 Client ID 或 Secret)。请前往管理员后台【系统设置 - PayPal 支付网关】填入您的沙盒或正式账户凭据，或勾选【沙盒免密秒测】体验！",
			})
			return
		}

		token, err := getPayPalAccessToken(cfg)
		if err != nil || token == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("PayPal 身份凭据核验未通过: %v。请检查【系统设置 - PayPal 支付网关】中的 Client ID 与 Secret 是否正确，并确认模式为 %s", err, mode),
			})
			return
		}

		baseURL := cfg.AppBaseURL
		if origin := c.Request.Header.Get("Origin"); origin != "" {
			baseURL = strings.TrimRight(origin, "/")
		} else if referer := c.Request.Header.Get("Referer"); referer != "" {
			if u, err := url.Parse(referer); err == nil && u.Host != "" {
				baseURL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
			}
		}

		ppURL := fmt.Sprintf("%s/v2/checkout/orders", getPayPalBaseURL(mode))
		orderPayload := map[string]interface{}{
			"intent": "CAPTURE",
			"purchase_units": []map[string]interface{}{
				{
					"reference_id": orderID,
					"description":  title,
					"amount": map[string]string{
						"currency_code": "USD",
						"value":         fmt.Sprintf("%.2f", amount),
					},
				},
			},
			"application_context": map[string]string{
				"return_url": fmt.Sprintf("%s/payment/success?order_id=%s", baseURL, orderID),
				"cancel_url": fmt.Sprintf("%s/payment/cancel?order_id=%s", baseURL, orderID),
			},
		}

		payloadBytes, _ := json.Marshal(orderPayload)
		httpReq, _ := http.NewRequest("POST", ppURL, bytes.NewBuffer(payloadBytes))
		httpReq.Header.Set("Authorization", "Bearer "+token)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, doErr := client.Do(httpReq)
		if doErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("无法连接至 PayPal 官方支付网关: %v", doErr),
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("PayPal 网关创建交易失败: %s", string(body)),
			})
			return
		}

		var ppRes struct {
			ID    string `json:"id"`
			Links []struct {
				Href string `json:"href"`
				Rel  string `json:"rel"`
			} `json:"links"`
		}
		if json.NewDecoder(resp.Body).Decode(&ppRes) == nil && ppRes.ID != "" {
			paypalOrderId = ppRes.ID
			for _, link := range ppRes.Links {
				if link.Rel == "approve" {
					approvalURL = link.Href
					break
				}
			}
		}

		order := models.Order{
			ID:             orderID,
			UserID:         user.ID,
			UserEmail:      user.Email,
			OrderType:      orderType,
			ProductID:      req.ProductID,
			ProductTitle:   title,
			MembershipType: req.MembershipType,
			Amount:         amount,
			Currency:       currency,
			PaypalOrderID:  paypalOrderId,
			Status:         "pending",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := database.DB.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to record order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"order_id":         order.ID,
			"paypal_order_id":  paypalOrderId,
			"approval_url":     approvalURL,
			"amount":           amount,
			"currency":         "USD",
			"paypal_client_id": clientID,
			"token":            tokenStr,
		})
	}
}

// FulfillOrder marks order paid and provisions user membership or sales count idempotently
func FulfillOrder(order *models.Order, cfg *config.Config) (*models.User, string, error) {
	var user models.User
	hasUser := false
	if order.UserID != "" && database.DB.First(&user, "id = ?", order.UserID).Error == nil {
		hasUser = true
	} else if order.UserEmail != "" && database.DB.First(&user, "LOWER(email) = ?", strings.ToLower(strings.TrimSpace(order.UserEmail))).Error == nil {
		hasUser = true
		order.UserID = user.ID
	}

	if order.Status == "paid" {
		if hasUser {
			tokenStr, _ := middleware.GenerateToken(&user, cfg)
			return &user, tokenStr, nil
		}
		return nil, "", nil
	}

	order.Status = "paid"
	order.UpdatedAt = time.Now()
	if err := database.DB.Save(order).Error; err != nil {
		return nil, "", err
	}

	if order.OrderType == "membership" && hasUser {
		now := time.Now()
		var newExp time.Time
		if user.MembershipExpire != nil && user.MembershipExpire.After(now) {
			newExp = *user.MembershipExpire
		} else {
			newExp = now
		}

		matchedPlan := GetMembershipPlanByType(order.MembershipType)
		if matchedPlan != nil && matchedPlan.DurationDays > 0 {
			newExp = newExp.AddDate(0, 0, matchedPlan.DurationDays)
		} else {
			switch order.MembershipType {
			case "monthly":
				newExp = newExp.AddDate(0, 1, 0)
			case "quarterly":
				newExp = newExp.AddDate(0, 3, 0)
			case "yearly":
				newExp = newExp.AddDate(1, 0, 0)
			case "lifetime":
				newExp = newExp.AddDate(99, 0, 0)
			default:
				newExp = newExp.AddDate(0, 1, 0)
			}
		}

		user.MembershipType = order.MembershipType
		user.MembershipExpire = &newExp
		database.DB.Save(&user)
	} else if order.ProductID != "" {
		database.DB.Model(&models.Product{}).Where("id = ?", order.ProductID).
			UpdateColumn("sales_count", gorm.Expr("sales_count + 1"))
	}

	if hasUser {
		tokenStr, _ := middleware.GenerateToken(&user, cfg)
		return &user, tokenStr, nil
	}
	return nil, "", nil
}

// CapturePayPalOrder captures payment on PayPal with strict verification and fulfills order
func CapturePayPalOrder(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenOrId := c.Param("order_id")
		if tokenOrId == "" || tokenOrId == "undefined" {
			tokenOrId = c.Query("order_id")
		}
		if tokenOrId == "" || tokenOrId == "undefined" {
			tokenOrId = c.Query("token")
		}
		if tokenOrId == "" || tokenOrId == "undefined" {
			var bodyReq struct {
				OrderID       string `json:"order_id"`
				PaypalOrderID string `json:"paypal_order_id"`
				Token         string `json:"token"`
			}
			if c.ShouldBindJSON(&bodyReq) == nil {
				if bodyReq.OrderID != "" {
					tokenOrId = bodyReq.OrderID
				} else if bodyReq.PaypalOrderID != "" {
					tokenOrId = bodyReq.PaypalOrderID
				} else if bodyReq.Token != "" {
					tokenOrId = bodyReq.Token
				}
			}
		}

		var order models.Order
		if err := database.DB.Where("paypal_order_id = ? OR id = ?", tokenOrId, tokenOrId).First(&order).Error; err != nil {
			if u, exists := c.Get("user"); exists {
				user := u.(*models.User)
				database.DB.Where("user_id = ? AND status = ?", user.ID, "pending").Order("created_at desc").First(&order)
			}
		}

		if order.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "未查询到对应订单，请确认订单号是否正确",
			})
			return
		}

		// If order is already paid, return fulfilled details immediately
		if order.Status == "paid" {
			user, tokenStr, _ := FulfillOrder(&order, cfg)
			c.JSON(http.StatusOK, gin.H{
				"status":   "paid",
				"message":  "订单已完成支付",
				"order_id": order.ID,
				"order":    order,
				"token":    tokenStr,
				"user": gin.H{
					"id":                user.ID,
					"email":             user.Email,
					"membership_type":   user.MembershipType,
					"membership_expire": user.MembershipExpire,
				},
			})
			return
		}

		// Handle simulated test payment vs manual QR code payment (Alipay, WeChat, Custom)
		if strings.HasPrefix(order.PaypalOrderID, "TEST-") ||
			strings.HasPrefix(order.PaypalOrderID, "ALIPAY-") ||
			strings.HasPrefix(order.PaypalOrderID, "WECHAT-") ||
			strings.HasPrefix(order.PaypalOrderID, "CUSTOM-") ||
			!strings.HasPrefix(order.PaypalOrderID, "PAYPAL-") {
			if strings.HasPrefix(order.PaypalOrderID, "TEST-") {
				perms := GetEffectivePermissionSettings()
				if !perms.AllowTestPay {
					c.JSON(http.StatusForbidden, gin.H{"message": "模拟测试支付通道已被管理员关闭"})
					return
				}
				user, tokenStr, _ := FulfillOrder(&order, cfg)
				c.JSON(http.StatusOK, gin.H{
					"status":   "paid",
					"is_test":  true,
					"message":  "⚡ 模拟测试支付成功，壁纸与会员权益已成功解锁！",
					"order_id": order.ID,
					"order":    order,
					"token":    tokenStr,
					"user": gin.H{
						"id":                user.ID,
						"email":             user.Email,
						"membership_type":   user.MembershipType,
						"membership_expire": user.MembershipExpire,
					},
				})
				return
			}

			// For offline / direct QR payment channels (WeChat, Alipay, Custom transfer):
			// User clicks "我已付款" -> marks order as "pending_review" for admin manual verification
			order.Status = "pending_review"
			order.UpdatedAt = time.Now()
			if err := database.DB.Save(&order).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "更新订单状态失败: " + err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":         "pending_review",
				"pending_review": true,
				"message":        "您的付款凭据已提交审核！管理员核对微信/支付宝转账账单后将立即审核通过并下发提取码（通常1-5分钟）。",
				"order_id":       order.ID,
				"order":          order,
			})
			return
		}

		// Real PayPal payment verification
		token, err := getPayPalAccessToken(cfg)
		if err != nil || token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "PayPal 商业授权凭据未配置，无法向官方核验支付状态",
			})
			return
		}

		_, _, mode := getEffectivePayPalConfig(cfg)
		baseURL := getPayPalBaseURL(mode)

		isPaid := false
		var payerID string

		// Step 1: Attempt to capture order
		captureURL := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", baseURL, order.PaypalOrderID)
		req, _ := http.NewRequest("POST", captureURL, bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, doErr := client.Do(req)

		if doErr == nil && resp != nil {
			defer resp.Body.Close()
			var ppRes struct {
				Status string `json:"status"` // "COMPLETED"
				Payer  struct {
					PayerID string `json:"payer_id"`
				} `json:"payer"`
			}
			bodyBytes, _ := io.ReadAll(resp.Body)
			_ = json.Unmarshal(bodyBytes, &ppRes)

			if (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) && ppRes.Status == "COMPLETED" {
				isPaid = true
				payerID = ppRes.Payer.PayerID
			} else {
				log.Printf("[PayPal Capture] POST capture status: %d, response: %s", resp.StatusCode, string(bodyBytes))
			}
		}

		// Step 2: Fallback check order status via GET (e.g. if already captured by client JS SDK or webhook)
		if !isPaid {
			getOrderURL := fmt.Sprintf("%s/v2/checkout/orders/%s", baseURL, order.PaypalOrderID)
			getReq, _ := http.NewRequest("GET", getOrderURL, nil)
			getReq.Header.Set("Authorization", "Bearer "+token)
			if getResp, getErr := client.Do(getReq); getErr == nil && getResp != nil {
				defer getResp.Body.Close()
				var getRes struct {
					Status string `json:"status"`
					Payer  struct {
						PayerID string `json:"payer_id"`
					} `json:"payer"`
				}
				if json.NewDecoder(getResp.Body).Decode(&getRes) == nil {
					if getRes.Status == "COMPLETED" {
						isPaid = true
						payerID = getRes.Payer.PayerID
					}
				}
			}
		}

		if !isPaid {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "pending",
				"message": "PayPal 订单尚未完成扣款或付款未通过，请检查 PayPal 付款窗口",
			})
			return
		}

		if payerID != "" {
			order.PaypalPayerID = payerID
		}

		user, tokenStr, fulfillErr := FulfillOrder(&order, cfg)
		if fulfillErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "订单权益履约异常"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "paid",
			"message":  "支付成功，已解锁对应权益！",
			"order_id": order.ID,
			"order":    order,
			"token":    tokenStr,
			"user": gin.H{
				"id":                user.ID,
				"email":             user.Email,
				"membership_type":   user.MembershipType,
				"membership_expire": user.MembershipExpire,
			},
		})
	}
}

// HandlePayPalWebhook processes incoming asynchronous notifications from PayPal
func HandlePayPalWebhook(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		var event struct {
			ID           string `json:"id"`
			EventType    string `json:"event_type"`
			ResourceType string `json:"resource_type"`
			Resource     struct {
				ID          string `json:"id"`
				Status      string `json:"status"`
				CustomID    string `json:"custom_id"`
				SupplementaryData struct {
					RelatedIDs struct {
						OrderID string `json:"order_id"`
					} `json:"related_ids"`
				} `json:"supplementary_data"`
			} `json:"resource"`
		}

		if err := json.Unmarshal(bodyBytes, &event); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		log.Printf("[PayPal Webhook] Received %s (Event ID: %s)", event.EventType, event.ID)

		targetOrderID := event.Resource.SupplementaryData.RelatedIDs.OrderID
		if targetOrderID == "" {
			targetOrderID = event.Resource.ID
		}

		switch event.EventType {
		case "PAYMENT.CAPTURE.COMPLETED":
			if event.Resource.Status == "COMPLETED" {
				var order models.Order
				if err := database.DB.Where("paypal_order_id = ? OR id = ?", targetOrderID, event.Resource.CustomID).First(&order).Error; err == nil {
					if order.Status != "paid" {
						_, _, _ = FulfillOrder(&order, cfg)
						log.Printf("[PayPal Webhook] Fulfilled order via capture completed: %s", order.ID)
					}
				}
			}

		case "CHECKOUT.ORDER.APPROVED":
			var order models.Order
			if err := database.DB.Where("paypal_order_id = ? OR id = ?", targetOrderID, event.Resource.CustomID).First(&order).Error; err == nil {
				if order.Status != "paid" {
					if token, err := getPayPalAccessToken(cfg); err == nil && token != "" {
						_, _, mode := getEffectivePayPalConfig(cfg)
						captureURL := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", getPayPalBaseURL(mode), order.PaypalOrderID)
						req, _ := http.NewRequest("POST", captureURL, bytes.NewBuffer([]byte("{}")))
						req.Header.Set("Authorization", "Bearer "+token)
						req.Header.Set("Content-Type", "application/json")
						client := &http.Client{Timeout: 10 * time.Second}
						if resp, doErr := client.Do(req); doErr == nil && resp != nil {
							defer resp.Body.Close()
							if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
								_, _, _ = FulfillOrder(&order, cfg)
								log.Printf("[PayPal Webhook] Captured & fulfilled order via approved: %s", order.ID)
							}
						}
					}
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "received"})
	}
}

// CreatePayPalSubscription for membership plans
func CreatePayPalSubscription(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			PlanType  string `json:"plan_type"` // "monthly", "quarterly", "yearly"
			UserEmail string `json:"user_email"`
			IsTest    bool   `json:"is_test"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Plan type required"})
			return
		}

		// Build forward payload
		bodyBytes, _ := json.Marshal(CreateOrderReq{
			OrderType:      "membership",
			MembershipType: req.PlanType,
			UserEmail:      req.UserEmail,
			IsTest:         req.IsTest,
		})
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		CreatePayPalOrder(cfg)(c)
	}
}

// GetPublicConfig returns non-sensitive site configuration
func GetPublicConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _, mode := getEffectivePayPalConfig(cfg)
		plans := GetEffectiveMembershipPlans()
		perms := GetEffectivePermissionSettings()

		c.JSON(http.StatusOK, gin.H{
			"paypal_client_id": clientID,
			"paypal_mode":      mode,
			"allow_test_pay":   perms.AllowTestPay,
			"currency":         "USD",
			"membership_plans": plans,
			"permissions": gin.H{
				"free_download_require_login": perms.FreeDownloadRequireLogin,
				"daily_free_limit":           perms.DailyFreeLimit,
				"guest_buy_enabled":          perms.GuestBuyEnabled,
				"allow_test_pay":             perms.AllowTestPay,
			},
		})
	}
}

// GetPaymentSettings for Admin
func GetPaymentSettings(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, clientSecret, mode := getEffectivePayPalConfig(cfg)
		// Mask secret
		maskedSecret := ""
		if len(clientSecret) > 6 {
			maskedSecret = clientSecret[:3] + "******" + clientSecret[len(clientSecret)-3:]
		} else if clientSecret != "" {
			maskedSecret = "******"
		}

		c.JSON(http.StatusOK, gin.H{
			"paypal_client_id":     clientID,
			"paypal_client_secret": maskedSecret,
			"paypal_mode":          mode,
			"allow_test_pay":       true,
		})
	}
}

// UpdatePaymentSettings for Admin
func UpdatePaymentSettings(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			PayPalClientID     string `json:"paypal_client_id"`
			PayPalClientSecret string `json:"paypal_client_secret"`
			PayPalMode         string `json:"paypal_mode"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid settings payload"})
			return
		}

		if body.PayPalClientID != "" {
			database.DB.Save(&models.SystemSetting{Key: "paypal_client_id", Value: strings.TrimSpace(body.PayPalClientID), UpdatedAt: time.Now()})
		}
		if body.PayPalClientSecret != "" && !strings.Contains(body.PayPalClientSecret, "******") {
			database.DB.Save(&models.SystemSetting{Key: "paypal_client_secret", Value: strings.TrimSpace(body.PayPalClientSecret), UpdatedAt: time.Now()})
		}
		if body.PayPalMode != "" {
			database.DB.Save(&models.SystemSetting{Key: "paypal_mode", Value: strings.TrimSpace(body.PayPalMode), UpdatedAt: time.Now()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "Payment settings updated successfully"})
	}
}
