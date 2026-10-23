package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"wallpaper-shop/config"
	"wallpaper-shop/database"
	"wallpaper-shop/handlers"
	"wallpaper-shop/middleware"
)

func main() {
	cfg := config.LoadConfig()
	database.InitDB(cfg)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Serve static uploads (local covers)
	r.Static("/uploads", "./uploads")

	// API Group
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "wallpaper-shop-api"})
		})

		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register(cfg))
			auth.POST("/login", handlers.Login(cfg))
			auth.GET("/me", middleware.JWTAuthMiddleware(cfg), handlers.Me())
		}

		// Public Products (with optional auth to detect if current user unlocked)
		api.GET("/products", handlers.ListProducts())
		api.GET("/products/:id", middleware.OptionalAuthMiddleware(cfg), handlers.GetProduct())

		// Downloads (Optional auth: free wallpapers don't require login, paid check order/member)
		api.GET("/products/:id/download", middleware.OptionalAuthMiddleware(cfg), handlers.GetDownloadLinks())
		api.POST("/products/:id/unlock-code", middleware.OptionalAuthMiddleware(cfg), handlers.UnlockWithCode(cfg))

		// Orders & Payment (Optional auth: supports guest checkout and sandbox test payment)
		api.POST("/payment/paypal/create", middleware.OptionalAuthMiddleware(cfg), handlers.CreatePayPalOrder(cfg))
		api.POST("/payment/paypal/capture/:order_id", middleware.OptionalAuthMiddleware(cfg), handlers.CapturePayPalOrder(cfg))
		api.POST("/payment/paypal/capture", middleware.OptionalAuthMiddleware(cfg), handlers.CapturePayPalOrder(cfg))
		api.POST("/payment/paypal/subscription", middleware.OptionalAuthMiddleware(cfg), handlers.CreatePayPalSubscription(cfg))
		api.POST("/payment/paypal/webhook", handlers.HandlePayPalWebhook(cfg))

		// Public order lookup (Allows buyers & guests to retrieve paid wallpapers by email)
		api.POST("/orders/lookup", handlers.LookupOrders())

		// Public site configuration (PayPal client ID, mode, currency)
		api.GET("/config/public", handlers.GetPublicConfig(cfg))
		api.GET("/config/locale", handlers.GetLocaleConfig())

		// Visitor Tracking & Dwell Time Analytics
		api.POST("/stats/track", handlers.TrackVisit())
		api.POST("/stats/heartbeat", handlers.HeartbeatVisit())

		// Payment Channels (public)
		api.GET("/payment/channels", handlers.ListPublicPaymentChannels())

		// Redeem Codes (public application)
		api.POST("/redeem/apply", middleware.OptionalAuthMiddleware(cfg), handlers.ApplyRedeemCode(cfg))

		// Categories & Universal Tutorial (Public)
		api.GET("/categories", handlers.ListPublicCategories())
		api.GET("/config/tutorial", handlers.GetPublicTutorial())

		orders := api.Group("/orders", middleware.JWTAuthMiddleware(cfg))
		{
			orders.GET("/my", handlers.UserOrders())
			orders.GET("/:id", handlers.GetOrder())
		}

		// Admin routes (Requires auth + admin role)
		admin := api.Group("/admin", middleware.JWTAuthMiddleware(cfg), middleware.AdminRequired())
		{
			// Products management & batch operations
			admin.GET("/products", handlers.AdminProductList())
			admin.GET("/products/:id", handlers.AdminGetProduct())
			admin.POST("/products", handlers.AdminCreateProduct())
			admin.POST("/products/batch", handlers.AdminBatchProducts())
			admin.PUT("/products/:id", handlers.AdminUpdateProduct())
			admin.PUT("/products/:id/pin", handlers.AdminToggleProductPin())
			admin.PUT("/products/:id/sort-order", handlers.AdminUpdateProductSort())
			admin.DELETE("/products/:id", handlers.AdminDeleteProduct())
			admin.POST("/upload/image", handlers.UploadImage())
			admin.POST("/upload/interactive", handlers.UploadInteractiveFile())

			// Categories management
			admin.GET("/categories", handlers.AdminListCategories())
			admin.POST("/categories", handlers.AdminCreateCategory())
			admin.PUT("/categories/:id", handlers.AdminUpdateCategory())
			admin.DELETE("/categories/:id", handlers.AdminDeleteCategory())

			// Universal Tutorial Template
			admin.GET("/tutorial", handlers.AdminGetTutorial())
			admin.PUT("/tutorial", handlers.AdminUpdateTutorial())

			// Orders management
			admin.GET("/orders", handlers.AdminOrderList())
			admin.PUT("/orders/:id/status", handlers.AdminUpdateOrderStatus(cfg))

			// Users management & Security operations
			admin.GET("/users", handlers.AdminUserList())
			admin.POST("/users", handlers.AdminCreateUser())
			admin.PUT("/users/:id", handlers.AdminUpdateUser())
			admin.DELETE("/users/:id", handlers.AdminDeleteUser())
			admin.POST("/users/:id/reset-password", handlers.AdminResetUserPassword())
			admin.POST("/users/:id/force-logout", handlers.AdminForceLogoutUser())

			// Visitors management & Analytics stats
			admin.GET("/visitors", handlers.AdminVisitorList())
			admin.GET("/stats", handlers.AdminStats())
			admin.GET("/stats/visits", handlers.AdminVisitStats())

			// Payment & System Settings
			admin.GET("/settings/payment", handlers.GetPaymentSettings(cfg))
			admin.POST("/settings/payment", handlers.UpdatePaymentSettings(cfg))
			admin.GET("/settings/all", handlers.GetAllSystemSettings(cfg))
			admin.POST("/settings/all", handlers.UpdateAllSystemSettings(cfg))

			// Payment Channels Configuration
			admin.GET("/payment/channels", handlers.AdminListPaymentChannels())
			admin.POST("/payment/channels", handlers.AdminCreatePaymentChannel())
			admin.PUT("/payment/channels/:id", handlers.AdminUpdatePaymentChannel())
			admin.DELETE("/payment/channels/:id", handlers.AdminDeletePaymentChannel())

			// Redeem Codes (Promotional Tokens) & Logs
			admin.GET("/redeem-codes", handlers.AdminListRedeemCodes())
			admin.POST("/redeem-codes", handlers.AdminCreateRedeemCode())
			admin.PUT("/redeem-codes/:id", handlers.AdminToggleRedeemCode())
			admin.DELETE("/redeem-codes/:id", handlers.AdminDeleteRedeemCode())
			admin.GET("/redeem-logs", handlers.AdminListRedeemLogs())
		}
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting Wallpaper Shop server on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
