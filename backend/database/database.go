package database

import (
	"encoding/json"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"wallpaper-shop/config"
	"wallpaper-shop/models"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto-migrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.SystemSetting{},
		&models.DownloadRecord{},
		&models.VisitLog{},
		&models.PaymentChannel{},
		&models.RedeemCode{},
		&models.RedeemLog{},
		&models.Category{},
	)
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	seedData()
	return DB
}

func seedData() {
	// 1. Seed Admin user if not present
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
		exp := time.Now().AddDate(10, 0, 0)
		admin := models.User{
			ID:               uuid.New().String(),
			Email:            "admin@wallpaper.com",
			PasswordHash:     string(hashed),
			Role:             "admin",
			MembershipType:   "yearly",
			MembershipExpire: &exp,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		DB.Create(&admin)
		log.Println("Seeded default admin: admin@wallpaper.com / admin123456")
	}

	// Seed default exchange rate (1 USD = 7.25 CNY) if not present
	var rateSetting models.SystemSetting
	if err := DB.First(&rateSetting, "key = ?", "exchange_rate").Error; err != nil {
		DB.Create(&models.SystemSetting{
			Key:       "exchange_rate",
			Value:     "7.25",
			UpdatedAt: time.Now(),
		})
	}

	// Seed default Payment Channels if none exist
	var channelCount int64
	DB.Model(&models.PaymentChannel{}).Count(&channelCount)
	if channelCount == 0 {
		defaultChannels := []models.PaymentChannel{
			{
				ID:          uuid.New().String(),
				Name:        "支付宝",
				Code:        "alipay",
				Scope:       "china_only",
				Icon:        "🔵",
				Description: "推荐国内用户使用，支持支付宝手机扫码与花呗快捷支付",
				Enabled:     true,
				SortOrder:   1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          uuid.New().String(),
				Name:        "微信支付",
				Code:        "wechat",
				Scope:       "china_only",
				Icon:        "🟢",
				Description: "微信扫码一秒支付，极速下发原画网盘提取链接",
				Enabled:     true,
				SortOrder:   2,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          uuid.New().String(),
				Name:        "PayPal 国际支付",
				Code:        "paypal",
				Scope:       "overseas_only",
				Icon:        "💳",
				Description: "支持国际 Visa / MasterCard / 信用卡与 PayPal 美元结算",
				Enabled:     true,
				SortOrder:   3,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          uuid.New().String(),
				Name:        "沙盒免密秒测",
				Code:        "test",
				Scope:       "global",
				Icon:        "⚡",
				Description: "全球免密沙盒测试体验，秒级解锁下载",
				Enabled:     true,
				SortOrder:   4,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}
		for _, ch := range defaultChannels {
			DB.Create(&ch)
		}
		log.Println("Seeded default payment channels (Alipay, WeChat, PayPal, Test)")
	}

	// 2. Seed initial wallpapers if empty
	var prodCount int64
	DB.Model(&models.Product{}).Count(&prodCount)
	if prodCount == 0 {
		sampleProducts := []struct {
			titleZh     string
			titleEn     string
			desc        string
			cover       string
			video       string
			price       float64
			category    string
			tags        []string
			sales       int64
			views       int64
			downloads   []models.DownloadItem
		}{
			{
				titleZh:  "赛博霓虹雨夜 8K 超清壁纸",
				titleEn:  "Cyberpunk Rainy Night 8K",
				desc:     "<h3>赛博朋克未来都市 · 动态视觉</h3><p>本作品由专业 3D 渲染引擎制作，呈现未来都市霓虹灯光在雨水湿润路面上的高精度反射。内含 8K (7680x4320) 超宽分辨率壁纸包，以及 Wallpaper Engine 动态交互壁纸包。</p><ul><li>分辨率：8K 超清 / 4K / 2K</li><li>色域：DCI-P3 广色域</li><li>格式：无损 PNG + MP4 动效</li></ul>",
				cover:    "https://images.unsplash.com/photo-1542751371-adc38448a05e?w=1200&q=80",
				video:    "https://www.bilibili.com/video/BV1xx411c7mD",
				price:    3.99,
				category: "Cyberpunk",
				tags:     []string{"8K", "赛博朋克", "霓虹", "雨夜", "动态"},
				sales:    86,
				views:    1420,
				downloads: []models.DownloadItem{
					{Platform: "baidu", URL: "https://pan.baidu.com/s/1demoCyberBaidu88", ExtractCode: "wp88"},
					{Platform: "quark", URL: "https://pan.quark.cn/s/demoQuarkCyber99", ExtractCode: "99aa"},
					{Platform: "google", URL: "https://drive.google.com/drive/folders/demoGoogleCyber77", ExtractCode: ""},
				},
			},
			{
				titleZh:  "星际深邃黑洞 4K 动态视界",
				titleEn:  "Deep Space Black Hole 4K",
				desc:     "<h3>爱因斯坦引力透镜 · 引力奇点</h3><p>基于广义相对论光线追踪演算生成的黑洞吸积盘模拟。极致深邃纯黑底色，专为 OLED 屏幕优化，省电且对比度拉满。</p><blockquote>“穿越事件视界，感受宇宙最极致的宁静。”</blockquote>",
				cover:    "https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1200&q=80",
				video:    "https://www.youtube.com/watch?v=demoBlackHole",
				price:    4.99,
				category: "Space",
				tags:     []string{"4K", "深空", "黑洞", "OLED", "宇宙"},
				sales:    64,
				views:    1105,
				downloads: []models.DownloadItem{
					{Platform: "baidu", URL: "https://pan.baidu.com/s/1demoSpaceBaidu55", ExtractCode: "sp55"},
					{Platform: "quark", URL: "https://pan.quark.cn/s/demoQuarkSpace33", ExtractCode: "33bb"},
					{Platform: "google", URL: "https://drive.google.com/drive/folders/demoGoogleSpace11", ExtractCode: ""},
				},
			},
			{
				titleZh:  "极简北欧晨雾雪山 5K",
				titleEn:  "Minimal Nordic Misty Mountain 5K",
				desc:     "<h3>北欧冷峻纯净美学</h3><p>拍摄于挪威罗弗敦群岛冬季晨曦，柔和的晨雾环绕雪山尖峰，极简色调适合全天候办公与桌面整理。</p>",
				cover:    "https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=1200&q=80",
				video:    "",
				price:    2.99,
				category: "Minimal",
				tags:     []string{"5K", "极简", "雪山", "自然", "冷淡风"},
				sales:    42,
				views:    890,
				downloads: []models.DownloadItem{
					{Platform: "baidu", URL: "https://pan.baidu.com/s/1demoNordicBaidu66", ExtractCode: "no66"},
					{Platform: "quark", URL: "https://pan.quark.cn/s/demoQuarkNordic22", ExtractCode: "22cc"},
					{Platform: "google", URL: "https://drive.google.com/drive/folders/demoGoogleNordic44", ExtractCode: ""},
				},
			},
			{
				titleZh:  "日系盛夏青空与海浪",
				titleEn:  "Japanese Summer Blue Sky & Waves",
				desc:     "<h3>新海诚风格 · 盛夏风物诗</h3><p>碧海蓝天，阳光穿透浪花。温暖治愈的色调，为你的桌面带来满满夏日清爽活力。</p>",
				cover:    "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=1200&q=80",
				video:    "https://www.bilibili.com/video/BV1demoSummerWaves",
				price:    2.49,
				category: "Nature",
				tags:     []string{"治愈", "动漫风", "海洋", "蓝天", "夏日"},
				sales:    51,
				views:    980,
				downloads: []models.DownloadItem{
					{Platform: "baidu", URL: "https://pan.baidu.com/s/1demoSummerBaidu77", ExtractCode: "su77"},
					{Platform: "quark", URL: "https://pan.quark.cn/s/demoQuarkSummer88", ExtractCode: "88dd"},
					{Platform: "google", URL: "https://drive.google.com/drive/folders/demoGoogleSummer99", ExtractCode: ""},
				},
			},
			{
				titleZh:  "东京秋夜流光街道 4K",
				titleEn:  "Tokyo Autumn Night Streets 4K",
				desc:     "<h3>街头摄影 · 都会孤寂感</h3><p>东京涩谷十字路口长曝光车轨，雨后银杏与都市光晕交织，沉浸感极强。</p>",
				cover:    "https://images.unsplash.com/photo-1503899036084-c55cdd92da26?w=1200&q=80",
				video:    "",
				price:    3.49,
				category: "City",
				tags:     []string{"东京", "城市", "夜景", "街道", "长曝光"},
				sales:    37,
				views:    740,
				downloads: []models.DownloadItem{
					{Platform: "baidu", URL: "https://pan.baidu.com/s/1demoTokyoBaidu12", ExtractCode: "tk12"},
					{Platform: "quark", URL: "https://pan.quark.cn/s/demoQuarkTokyo34", ExtractCode: "tk34"},
					{Platform: "google", URL: "https://drive.google.com/drive/folders/demoGoogleTokyo56", ExtractCode: ""},
				},
			},
		}

		for _, sp := range sampleProducts {
			tagsJson, _ := json.Marshal(sp.tags)
			dlJson, _ := json.Marshal(sp.downloads)
			p := models.Product{
				ID:              uuid.New().String(),
				TitleZh:         sp.titleZh,
				TitleEn:         sp.titleEn,
				Description:     sp.desc,
				CoverImage:      sp.cover,
				PreviewVideoURL: sp.video,
				Price:           sp.price,
				Category:        sp.category,
				Tags:            string(tagsJson),
				Downloads:       string(dlJson),
				IsPublished:     true,
				ViewsCount:      sp.views,
				SalesCount:      sp.sales,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			DB.Create(&p)
		}
		log.Println("Seeded initial wallpaper products with cloud drive links.")
	}

	// 3. Seed default categories if none exist
	var catCount int64
	DB.Model(&models.Category{}).Count(&catCount)
	if catCount == 0 {
		defaultCategories := []models.Category{
			{
				ID:        uuid.New().String(),
				NameZh:    "iOS 风格",
				NameEn:    "iOS Theme",
				Code:      "ios",
				Icon:      "📱",
				SortOrder: 100,
				IsEnabled: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				NameZh:    "动漫二次元",
				NameEn:    "Anime",
				Code:      "anime",
				Icon:      "✨",
				SortOrder: 90,
				IsEnabled: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				NameZh:    "赛博科技",
				NameEn:    "Cyberpunk",
				Code:      "cyberpunk",
				Icon:      "🏙️",
				SortOrder: 80,
				IsEnabled: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				NameZh:    "治愈风景",
				NameEn:    "Nature",
				Code:      "nature",
				Icon:      "🌿",
				SortOrder: 70,
				IsEnabled: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				NameZh:    "极简冷淡",
				NameEn:    "Minimal",
				Code:      "minimal",
				Icon:      "🌸",
				SortOrder: 60,
				IsEnabled: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				NameZh:    "游戏同人",
				NameEn:    "Gaming",
				Code:      "game",
				Icon:      "🎮",
				SortOrder: 50,
				IsEnabled: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				NameZh:    "赛车超跑",
				NameEn:    "Supercars",
				Code:      "cars",
				Icon:      "🏎️",
				SortOrder: 40,
				IsEnabled: true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		for _, cat := range defaultCategories {
			DB.Create(&cat)
		}
		log.Println("Seeded default wallpaper categories.")
	}

	// 4. Seed default Universal Tutorial Template if not present
	var tutSetting models.SystemSetting
	if err := DB.First(&tutSetting, "key = ?", "tutorial_template").Error; err != nil {
		defaultTutorial := models.TutorialTemplate{
			Enabled:  true,
			TitleZh:  "iOS 动态壁纸实机交互设置教程",
			TitleEn:  "iOS Live Wallpaper Setup Guide",
			NoticeZh: "💡 提示：依据当前 iPhone 系统版本选择对应指引与视频教程，后续随 iOS 系统大版本统一在此更新维护",
			NoticeEn: "💡 Tip: Follow the instructions and video link according to your iOS version.",
			Items: []models.TutorialItem{
				{
					ID:            "tut_ios18_26_4",
					VersionRange:  "iOS 18.0 - iOS 26.4",
					TitleZh:       "经典锁屏与景深联动设置法",
					TitleEn:       "Classic Lock Screen & Depth Effect Setup",
					BadgeZh:       "主流系统推荐",
					BadgeEn:       "Recommended",
					BadgeColor:    "blue",
					VideoURL:      "https://www.bilibili.com/video/BV1xx411c7mD",
					VideoPlatform: "bilibili",
					Steps: []string{
						"解压下载的原画及 .tendies 动画资源包到手机相册或「文件」App",
						"在「设置」-「墙纸」中新增壁纸，选择实况照片并开启「实况」联动",
						"长按锁屏自定义，点击时钟样式调整景深，使人物或主体产生穿透层叠效果",
					},
				},
				{
					ID:            "tut_ios26_5_27_0",
					VersionRange:  "iOS 26.5 - iOS 27.0 Beta",
					TitleZh:       "新版锁屏与交互层穿透设置法",
					TitleEn:       "New Interaction Layer Penetration Setup",
					BadgeZh:       "Beta 测试版专用",
					BadgeEn:       "Beta Version",
					BadgeColor:    "purple",
					VideoURL:      "https://www.bilibili.com/video/BV1xx411c7mD",
					VideoPlatform: "bilibili",
					Steps: []string{
						"更新 iPhone 至支持全屏浮动交互的 iOS 版本",
						"进入「设置」-「显示与亮度」，确保「始终显示」与「轻触唤醒」开启",
						"导入交互动作并绑定对应快捷指令，滑动或点击锁屏即可触发专属动态",
					},
				},
			},
		}
		rawBytes, _ := json.Marshal(defaultTutorial)
		DB.Create(&models.SystemSetting{
			Key:       "tutorial_template",
			Value:     string(rawBytes),
			UpdatedAt: time.Now(),
		})
		log.Println("Seeded default universal tutorial template.")
	}
}
