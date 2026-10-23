package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"wallpaper-shop/models"
	"wallpaper-shop/services"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := gorm.Open(sqlite.Open("wallpaper.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}

	var p models.Product
	if err := db.Where("title_zh LIKE ?", "%程奕迅%").First(&p).Error; err != nil {
		fmt.Println("Product not found:", err)
		return
	}

	uploadDir := filepath.Join(".", "uploads", "interactive")
	tendiesBase := "interactive_1792563420_71efc71f.tendies"
	zipPath := filepath.Join(uploadDir, tendiesBase)
	cleanID := "71efc71f"

	fmt.Printf("Re-parsing %s for product %s (%s)...\n", zipPath, p.TitleZh, p.ID)
	result, err := services.ParseTendiesArchive(zipPath, uploadDir, cleanID)
	if err != nil {
		fmt.Println("Parse failed:", err)
		return
	}

	var cfg map[string]interface{}
	if p.InteractiveConfig != "" {
		json.Unmarshal([]byte(p.InteractiveConfig), &cfg)
	} else {
		cfg = make(map[string]interface{})
	}

	cfg["tendies_file"] = "/uploads/interactive/" + tendiesBase
	cfg["background_url"] = result.BackgroundURL
	cfg["foreground_url"] = result.ForegroundURL
	cfg["bg_layers"] = result.BGLayers
	cfg["floating_layers"] = result.FloatingLayers
	cfg["fg_layers"] = result.FGLayers
	cfg["canvas_width"] = result.CanvasWidth
	cfg["canvas_height"] = result.CanvasHeight
	cfg["total_image_count"] = result.TotalImageCount
	cfg["gyro_enabled"] = true
	cfg["depth_intensity"] = 16
	cfg["has_aod"] = true

	newConfigBytes, _ := json.Marshal(cfg)
	p.InteractiveConfig = string(newConfigBytes)
	p.IsInteractive = true
	p.InteractiveType = "tendies"

	if err := db.Save(&p).Error; err != nil {
		fmt.Println("Failed to save product:", err)
		return
	}

	fmt.Printf("Successfully updated %s with %d BG layers, %d Floating layers!\n",
		p.TitleZh, len(result.BGLayers), len(result.FloatingLayers))
}
