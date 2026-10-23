package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wallpaper-shop/database"
	"wallpaper-shop/models"
)

// GetDefaultTutorialTemplate returns fallback tutorial content
func GetDefaultTutorialTemplate() models.TutorialTemplate {
	return models.TutorialTemplate{
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
}

// GetPublicTutorial provides the universal tutorial guide to frontends
func GetPublicTutorial() gin.HandlerFunc {
	return func(c *gin.Context) {
		var setting models.SystemSetting
		if err := database.DB.First(&setting, "key = ?", "tutorial_template").Error; err != nil {
			c.JSON(http.StatusOK, GetDefaultTutorialTemplate())
			return
		}

		var tut models.TutorialTemplate
		if err := json.Unmarshal([]byte(setting.Value), &tut); err != nil {
			c.JSON(http.StatusOK, GetDefaultTutorialTemplate())
			return
		}

		c.JSON(http.StatusOK, tut)
	}
}

// AdminGetTutorial provides the tutorial template to admin for editing
func AdminGetTutorial() gin.HandlerFunc {
	return func(c *gin.Context) {
		var setting models.SystemSetting
		if err := database.DB.First(&setting, "key = ?", "tutorial_template").Error; err != nil {
			c.JSON(http.StatusOK, GetDefaultTutorialTemplate())
			return
		}

		var tut models.TutorialTemplate
		if err := json.Unmarshal([]byte(setting.Value), &tut); err != nil {
			c.JSON(http.StatusOK, GetDefaultTutorialTemplate())
			return
		}

		c.JSON(http.StatusOK, tut)
	}
}

// AdminUpdateTutorial saves updated tutorial template
func AdminUpdateTutorial() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload models.TutorialTemplate
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tutorial payload: " + err.Error()})
			return
		}

		// Ensure items have IDs
		for i := range payload.Items {
			if payload.Items[i].ID == "" {
				payload.Items[i].ID = "tut_" + uuid.New().String()[:8]
			}
		}

		rawBytes, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to serialize tutorial data"})
			return
		}

		var setting models.SystemSetting
		if err := database.DB.First(&setting, "key = ?", "tutorial_template").Error; err == nil {
			setting.Value = string(rawBytes)
			setting.UpdatedAt = time.Now()
			database.DB.Save(&setting)
		} else {
			setting = models.SystemSetting{
				Key:       "tutorial_template",
				Value:     string(rawBytes),
				UpdatedAt: time.Now(),
			}
			database.DB.Create(&setting)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "教程通用模板保存成功",
			"tutorial": payload,
		})
	}
}
