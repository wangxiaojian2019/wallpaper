package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wallpaper-shop/config"
	"wallpaper-shop/database"
	"wallpaper-shop/models"
)

// UserOrders returns purchase history for current user
func UserOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		user := u.(*models.User)

		var orders []models.Order
		if err := database.DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch orders"})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

// GetOrder returns single order
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var order models.Order
		if err := database.DB.First(&order, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
			return
		}

		// Security: verify requester is owner or admin
		u, exists := c.Get("user")
		if exists {
			user := u.(*models.User)
			if user.Role != "admin" && user.ID != order.UserID {
				c.JSON(http.StatusForbidden, gin.H{"message": "Access denied"})
				return
			}
		}

		c.JSON(http.StatusOK, order)
	}
}

// AdminOrderList returns paginated list of all orders
func AdminOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page < 1 {
			page = 1
		}
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		if pageSize < 1 {
			pageSize = 20
		}

		query := database.DB.Model(&models.Order{})

		status := c.Query("status")
		if status != "" {
			query = query.Where("status = ?", status)
		}

		orderType := c.Query("order_type")
		if orderType != "" {
			query = query.Where("order_type = ?", orderType)
		}

		search := c.Query("search")
		if search != "" {
			pattern := "%" + strings.ToLower(search) + "%"
			query = query.Where("LOWER(id) LIKE ? OR LOWER(paypal_order_id) LIKE ? OR LOWER(user_email) LIKE ?", pattern, pattern, pattern)
		}

		var total int64
		query.Count(&total)

		var orders []models.Order
		offset := (page - 1) * pageSize
		if err := query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch orders"})
			return
		}

		for i := range orders {
			if orders[i].Currency == "" || orders[i].Currency == "USD" {
				ref := strings.ToUpper(orders[i].PaypalOrderID)
				if strings.HasPrefix(ref, "ALIPAY-") || strings.HasPrefix(ref, "WECHAT-") || strings.HasPrefix(ref, "CUSTOM-") {
					orders[i].Currency = "CNY"
				}
			}
		}

		var pendingReviewTotal int64
		database.DB.Model(&models.Order{}).Where("status = ?", "pending_review").Count(&pendingReviewTotal)

		c.JSON(http.StatusOK, gin.H{
			"orders":               orders,
			"total":                total,
			"page":                 page,
			"page_size":            pageSize,
			"pending_review_total": pendingReviewTotal,
		})
	}
}

// AdminUpdateOrderStatus updates order status manually
func AdminUpdateOrderStatus(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var order models.Order
		if err := database.DB.First(&order, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
			return
		}

		var req struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Status required"})
			return
		}

		// If admin marks order as "paid", call FulfillOrder to provision VIP or sales count cleanly
		if req.Status == "paid" {
			_, _, err := FulfillOrder(&order, cfg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "审核通过下发权益失败: " + err.Error()})
				return
			}
		} else {
			order.Status = req.Status
			order.UpdatedAt = time.Now()
			database.DB.Save(&order)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "订单状态已更新",
			"order":   order,
		})
	}
}

type OrderLookupReq struct {
	Email   string `json:"email" binding:"required"`
	OrderID string `json:"order_id"`
}

// LookupOrders allows buyers (including guests) to retrieve their paid and reviewing orders
func LookupOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OrderLookupReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请输入购买时填写的联系邮箱"})
			return
		}

		email := strings.ToLower(strings.TrimSpace(req.Email))
		if email == "" || !strings.Contains(email, "@") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "请输入有效的电子邮箱地址"})
			return
		}

		query := database.DB.Model(&models.Order{}).Where("LOWER(user_email) = ? AND status IN ('paid', 'pending_review')", email)
		orderID := strings.TrimSpace(req.OrderID)
		if orderID != "" {
			query = query.Where("id = ? OR paypal_order_id = ?", orderID, orderID)
		}

		var orders []models.Order
		if err := query.Order("created_at desc").Limit(20).Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单记录失败"})
			return
		}

		if len(orders) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "未查询到该邮箱名下的订单，请核对邮箱或订单编号"})
			return
		}

		type OrderResult struct {
			models.Order
			CoverImage    string                `json:"cover_image,omitempty"`
			DownloadLinks []models.DownloadItem `json:"download_links,omitempty"`
		}

		var results []OrderResult
		for _, o := range orders {
			res := OrderResult{Order: o}
			if o.Status == "paid" && o.ProductID != "" {
				var p models.Product
				if err := database.DB.Select("id, cover_image, downloads").First(&p, "id = ?", o.ProductID).Error; err == nil {
					res.CoverImage = p.CoverImage
					if p.Downloads != "" {
						_ = json.Unmarshal([]byte(p.Downloads), &res.DownloadLinks)
					}
				}
			} else if o.ProductID != "" {
				var p models.Product
				if err := database.DB.Select("id, cover_image").First(&p, "id = ?", o.ProductID).Error; err == nil {
					res.CoverImage = p.CoverImage
				}
			}
			results = append(results, res)
		}

		c.JSON(http.StatusOK, gin.H{
			"total":  len(results),
			"orders": results,
		})
	}
}

