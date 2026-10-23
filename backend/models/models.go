package models

import (
	"time"
)

// User represents a registered customer or administrator
type User struct {
	ID               string     `json:"id" gorm:"primaryKey;size:36"`
	Email            string     `json:"email" gorm:"uniqueIndex;not null;size:128"`
	PasswordHash     string     `json:"-" gorm:"size:255;not null"`
	Role             string     `json:"role" gorm:"default:'user';size:20"`                         // "user", "admin"
	MembershipType   string     `json:"membership_type" gorm:"default:'none';size:20"`              // "none", "monthly", "quarterly", "yearly", "lifetime"
	MembershipExpire *time.Time `json:"membership_expire"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	LastLoginIP      string     `json:"last_login_ip" gorm:"size:64"`
	TokenVersion     int        `json:"token_version" gorm:"default:1"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// DownloadItem represents a cloud storage link (Baidu, Quark, Google Drive)
type DownloadItem struct {
	Platform    string `json:"platform"`     // "baidu", "quark", "google"
	URL         string `json:"url"`          // Drive URL
	ExtractCode string `json:"extract_code"` // Extract/access code if any
}

// VideoItem represents a preview video link with platform specification
type VideoItem struct {
	Platform string `json:"platform"` // "tiktok", "douyin", "xiaohongshu", "bilibili", "youtube", "direct"
	URL      string `json:"url"`
	Title    string `json:"title"`
}

// Product represents a wallpaper package
type Product struct {
	ID              string    `json:"id" gorm:"primaryKey;size:36"`
	TitleZh         string    `json:"title_zh" gorm:"size:255;not null"`
	TitleEn         string    `json:"title_en" gorm:"size:255"`
	Description     string    `json:"description" gorm:"type:text"` // HTML Rich text content
	CoverImage      string    `json:"cover_image" gorm:"size:512"`  // Cover screenshot image
	PreviewVideoURL string    `json:"preview_video_url" gorm:"size:512"` // Hyperlink to preview video
	PreviewPlatform string    `json:"preview_platform" gorm:"size:64;default:'douyin'"` // "tiktok", "douyin", "xiaohongshu", "bilibili", "youtube", "direct"
	PreviewVideos   string    `json:"preview_videos" gorm:"type:text"` // JSON array of VideoItem
	Price           float64   `json:"price" gorm:"default:0"`       // Base price (fallback)
	PriceUSD        float64   `json:"price_usd" gorm:"default:0"`   // USD price ($)
	PriceCNY        float64   `json:"price_cny" gorm:"default:0"`   // CNY price (¥)
	VerifyCode      string    `json:"verify_code" gorm:"size:64"`   // Unlock/redeem code (e.g. from legacy site or custom)
	Category        string    `json:"category" gorm:"size:64;default:'general'"`
	Tags            string    `json:"tags" gorm:"type:text"`        // JSON array string
	AllowedVIPLevels string   `json:"allowed_vip_levels" gorm:"type:text"` // JSON array of allowed VIP levels: ["monthly","quarterly","yearly","lifetime"]
	Downloads       string    `json:"-" gorm:"type:text"`           // JSON array of DownloadItem, hidden from public
	IsPublished       bool      `json:"is_published" gorm:"default:true"`
	IsInteractive     bool      `json:"is_interactive" gorm:"default:false"`
	InteractiveType   string    `json:"interactive_type" gorm:"size:32;default:'tendies'"` // "tendies", "gyro_depth", "live_motion"
	InteractiveConfig string    `json:"interactive_config" gorm:"type:text"`              // JSON string of interactive layers & settings
	IsPinned          bool      `json:"is_pinned" gorm:"default:false;index"`
	PinnedAt          *time.Time `json:"pinned_at" gorm:"index"`
	SortOrder         int       `json:"sort_order" gorm:"default:0;index"`
	ViewsCount        int64     `json:"views_count" gorm:"default:0"`
	SalesCount        int64     `json:"sales_count" gorm:"default:0"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Category represents a wallpaper category
type Category struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	NameZh    string    `json:"name_zh" gorm:"size:128;not null"`
	NameEn    string    `json:"name_en" gorm:"size:128"`
	Code      string    `json:"code" gorm:"size:64;uniqueIndex;not null"`
	Icon      string    `json:"icon" gorm:"size:32;default:'📁'"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	IsEnabled bool      `json:"is_enabled" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TutorialItem represents a version-specific tutorial guide with a video link
type TutorialItem struct {
	ID            string   `json:"id"`
	VersionRange  string   `json:"version_range"`  // e.g. "iOS 18.0 - iOS 26.4"
	TitleZh       string   `json:"title_zh"`       // e.g. "经典锁屏与景深联动设置法"
	TitleEn       string   `json:"title_en"`       // e.g. "Classic Lock Screen & Depth Effect Setup"
	BadgeZh       string   `json:"badge_zh"`       // e.g. "主流推荐"
	BadgeEn       string   `json:"badge_en"`       // e.g. "Recommended"
	BadgeColor    string   `json:"badge_color"`    // "blue", "purple", "green", "orange"
	VideoURL      string   `json:"video_url"`      // Video URL
	VideoPlatform string   `json:"video_platform"` // "bilibili", "douyin", "youtube", "direct"
	Steps         []string `json:"steps"`          // Step-by-step instructions
}

// TutorialTemplate represents the universal detail-page tutorial configuration
type TutorialTemplate struct {
	Enabled  bool           `json:"enabled"`
	TitleZh  string         `json:"title_zh"`
	TitleEn  string         `json:"title_en"`
	NoticeZh string         `json:"notice_zh"`
	NoticeEn string         `json:"notice_en"`
	Items    []TutorialItem `json:"items"`
}

// SystemSetting represents key-value store for site configuration (PayPal, etc.)
type SystemSetting struct {
	Key       string    `json:"key" gorm:"primaryKey;size:64"`
	Value     string    `json:"value" gorm:"type:text"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Order represents an order for a single wallpaper or membership plan
type Order struct {
	ID             string    `json:"id" gorm:"primaryKey;size:36"`
	UserID         string    `json:"user_id" gorm:"size:36;index"`
	UserEmail      string    `json:"user_email" gorm:"size:128"`
	OrderType      string    `json:"order_type" gorm:"size:20"` // "single", "membership"
	ProductID      string    `json:"product_id" gorm:"size:36"`
	ProductTitle   string    `json:"product_title" gorm:"size:255"`
	MembershipType string    `json:"membership_type" gorm:"size:20"` // "monthly", "quarterly", "yearly"
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency" gorm:"default:'USD';size:10"`
	PaypalOrderID  string    `json:"paypal_order_id" gorm:"size:64;index"`
	PaypalPayerID  string    `json:"paypal_payer_id" gorm:"size:64"`
	Status         string    `json:"status" gorm:"default:'pending';size:20"` // "pending", "paid", "cancelled"
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// DownloadRecord tracks wallpaper downloads for user history and daily limit enforcement
type DownloadRecord struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	UserID    string    `json:"user_id" gorm:"size:36;index"`
	UserEmail string    `json:"user_email" gorm:"size:128"`
	ProductID string    `json:"product_id" gorm:"size:36;index"`
	IP        string    `json:"ip" gorm:"size:64"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
}

// VisitLog tracks visitor traffic, IP, device, duration and page activity
type VisitLog struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	IP        string    `json:"ip" gorm:"size:64;index"`
	Country   string    `json:"country" gorm:"size:64"`
	City      string    `json:"city" gorm:"size:64"`
	Device    string    `json:"device" gorm:"size:32"`  // "Desktop", "Mobile", "Tablet"
	OS        string    `json:"os" gorm:"size:64"`      // "Windows", "iOS", "Android", "macOS", "Linux"
	Browser   string    `json:"browser" gorm:"size:64"` // "Chrome", "Safari", "Firefox", "Edge", "WeChat"
	Path      string    `json:"path" gorm:"size:255"`
	Language  string    `json:"language" gorm:"size:32"`
	Duration  int       `json:"duration" gorm:"default:0"` // Dwell time in seconds
	UserAgent string    `json:"user_agent" gorm:"size:512"`
	UserID         string    `json:"user_id" gorm:"size:36;index"`
	UserEmail      string    `json:"user_email" gorm:"size:128;index"`
	MembershipType string    `json:"membership_type" gorm:"size:32"`
	IsMember       bool      `json:"is_member" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaymentChannel represents a configured payment method with scope
type PaymentChannel struct {
	ID          string    `json:"id" gorm:"primaryKey;size:36"`
	Name        string    `json:"name" gorm:"size:128"`          // e.g. "支付宝", "微信支付", "PayPal"
	Code        string    `json:"code" gorm:"size:64;index"`     // "alipay", "wechat", "paypal", "stripe", "custom"
	Scope       string    `json:"scope" gorm:"size:32;default:'global'"` // "global", "china_only", "overseas_only"
	Icon        string    `json:"icon" gorm:"size:64"`          // "🔵", "🟢", "💳", "⚡"
	Description string    `json:"description" gorm:"size:255"`  // e.g. "支持手机支付宝扫码、快捷支付"
	Enabled     bool      `json:"enabled" gorm:"default:true"`
	AccountInfo string    `json:"account_info" gorm:"type:text"` // JSON config or QR code url / merchant id
	QRCodeURL   string    `json:"qr_code_url" gorm:"size:512"`  // Direct QR code image URL (e.g. /uploads/covers/...)
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RedeemCode represents a custom promotional token or redemption voucher
type RedeemCode struct {
	ID         string     `json:"id" gorm:"primaryKey;size:36"`
	Code       string     `json:"code" gorm:"size:64;uniqueIndex;not null"`
	Purpose    string     `json:"purpose" gorm:"size:32"` // "price_limit" (限额免单), "membership" (兑换会员), "product" (兑换指定壁纸)
	TargetTier string     `json:"target_tier" gorm:"size:32"` // "monthly", "quarterly", "yearly", "lifetime" if Purpose=="membership"
	MaxPrice   float64    `json:"max_price" gorm:"default:0"` // e.g. 5.0 (5元及以下免费兑换) if Purpose=="price_limit"
	Currency   string     `json:"currency" gorm:"size:10;default:'CNY'"` // "CNY" or "USD"
	ProductID  string     `json:"product_id" gorm:"size:36"` // specific product ID if Purpose=="product"
	MaxUses    int        `json:"max_uses" gorm:"default:1"` // Custom max redeem uses, e.g. 1, 10, 100, 99999
	UsedCount  int        `json:"used_count" gorm:"default:0"`
	ExpiresAt  *time.Time `json:"expires_at"` // Custom expiration date, null if permanent
	IsActive   bool       `json:"is_active" gorm:"default:true"`
	Remark     string     `json:"remark" gorm:"size:255"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// RedeemLog records who redeemed what code and when
type RedeemLog struct {
	ID           string    `json:"id" gorm:"primaryKey;size:36"`
	CodeID       string    `json:"code_id" gorm:"size:36;index"`
	Code         string    `json:"code" gorm:"size:64;index"`
	UserID       string    `json:"user_id" gorm:"size:36;index"`
	UserEmail    string    `json:"user_email" gorm:"size:128;index"`
	Purpose      string    `json:"purpose" gorm:"size:32"`
	ProductID    string    `json:"product_id" gorm:"size:36"`
	ProductTitle string    `json:"product_title" gorm:"size:255"`
	GrantedTier  string    `json:"granted_tier" gorm:"size:32"`
	IP           string    `json:"ip" gorm:"size:64"`
	CreatedAt    time.Time `json:"created_at"`
}


