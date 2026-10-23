package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wallpaper-shop/database"
	"wallpaper-shop/models"
)

// CategoryForm represents request payload for category create/update
type CategoryForm struct {
	NameZh    string `json:"name_zh" binding:"required"`
	NameEn    string `json:"name_en"`
	Code      string `json:"code" binding:"required"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
	IsEnabled *bool  `json:"is_enabled"`
}

// ListPublicCategories returns enabled categories for public store
func ListPublicCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category
		if err := database.DB.Where("is_enabled = ?", true).
			Order("sort_order desc, created_at asc").
			Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch categories"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"items": categories,
			"total": len(categories),
		})
	}
}

// AdminListCategories returns all categories including disabled ones for admin
func AdminListCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category
		if err := database.DB.Order("sort_order desc, created_at asc").Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch categories"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"items": categories,
			"total": len(categories),
		})
	}
}

// AdminCreateCategory adds a new category
func AdminCreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form CategoryForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "分类名称(中文)与唯一英文标识编码为必填项"})
			return
		}

		code := strings.ToLower(strings.TrimSpace(form.Code))
		var existing models.Category
		if err := database.DB.Where("code = ?", code).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "已存在相同编码标识的分类，请更换编码"})
			return
		}

		icon := form.Icon
		if strings.TrimSpace(icon) == "" {
			icon = "📁"
		}

		enabled := true
		if form.IsEnabled != nil {
			enabled = *form.IsEnabled
		}

		cat := models.Category{
			ID:        uuid.New().String(),
			NameZh:    strings.TrimSpace(form.NameZh),
			NameEn:    strings.TrimSpace(form.NameEn),
			Code:      code,
			Icon:      icon,
			SortOrder: form.SortOrder,
			IsEnabled: enabled,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := database.DB.Create(&cat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create category: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "分类创建成功",
			"category": cat,
		})
	}
}

// AdminUpdateCategory edits an existing category
func AdminUpdateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var cat models.Category
		if err := database.DB.First(&cat, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "分类不存在"})
			return
		}

		var form CategoryForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "分类名称(中文)与唯一编码为必填项"})
			return
		}

		code := strings.ToLower(strings.TrimSpace(form.Code))
		var existing models.Category
		if err := database.DB.Where("code = ? AND id != ?", code, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "编码标识已被其他分类使用，请更换"})
			return
		}

		cat.NameZh = strings.TrimSpace(form.NameZh)
		cat.NameEn = strings.TrimSpace(form.NameEn)
		cat.Code = code
		if strings.TrimSpace(form.Icon) != "" {
			cat.Icon = form.Icon
		}
		cat.SortOrder = form.SortOrder
		if form.IsEnabled != nil {
			cat.IsEnabled = *form.IsEnabled
		}
		cat.UpdatedAt = time.Now()

		if err := database.DB.Save(&cat).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update category: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "分类更新成功",
			"category": cat,
		})
	}
}

// AdminDeleteCategory removes a category
func AdminDeleteCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := database.DB.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete category"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "分类删除成功"})
	}
}
