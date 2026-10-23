package main

import (
	"encoding/json"
	"fmt"
	"wallpaper-shop/models"
	"wallpaper-shop/services"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("wallpaper.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}

	prodID := "a39de3d8-8e52-4a68-ad57-cfb4532dc33e"
	var p models.Product
	if err := db.First(&p, "id = ?", prodID).Error; err != nil {
		fmt.Println("Find product error:", err)
		return
	}

	zipPath := "uploads/interactive/interactive_1792717701_399ca199.tendies"
	cleanID := "399ca199"
	res, err := services.ParseTendiesArchive(zipPath, "uploads/interactive", cleanID)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	var cfg map[string]interface{}
	_ = json.Unmarshal([]byte(p.InteractiveConfig), &cfg)
	if cfg == nil {
		cfg = make(map[string]interface{})
	}
	cfg["tendies_file"] = "/" + zipPath
	cfg["bg_layers"] = res.BGLayers
	cfg["floating_layers"] = res.FloatingLayers
	cfg["fg_layers"] = res.FGLayers
	cfg["background_url"] = res.BackgroundURL
	cfg["foreground_url"] = res.ForegroundURL
	cfg["canvas_width"] = res.CanvasWidth
	cfg["canvas_height"] = res.CanvasHeight
	cfg["total_image_count"] = res.TotalImageCount

	newCfgBytes, _ := json.Marshal(cfg)
	p.InteractiveConfig = string(newCfgBytes)
	if res.CoverURL != "" {
		p.CoverImage = res.CoverURL
	}
	db.Save(&p)

	fmt.Printf("Successfully updated %s (%s)!\n", p.TitleZh, p.ID)
	fmt.Printf("BGLayers: %d\n", len(res.BGLayers))
	for i, l := range res.BGLayers {
		fmt.Printf("  BG[%d]: %s, w=%.1f, h=%.1f\n", i, l.Name, l.Width, l.Height)
	}
	fmt.Printf("FloatingLayers: %d\n", len(res.FloatingLayers))
	for i, l := range res.FloatingLayers {
		fmt.Printf("  FL[%d]: %s, w=%.1f, h=%.1f\n", i, l.Name, l.Width, l.Height)
	}
	fmt.Printf("BackgroundURL: %s\n", res.BackgroundURL)
	fmt.Printf("ForegroundURL: %s\n", res.ForegroundURL)
}
