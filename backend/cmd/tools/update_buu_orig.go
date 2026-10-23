package main

import (
	"encoding/json"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"wallpaper-shop/models"
	"wallpaper-shop/services"
)

func main() {
	db, err := gorm.Open(sqlite.Open("wallpaper.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}

	var p models.Product
	if err := db.First(&p, "id = ?", "77055a8c-79b5-434c-b7a1-eb541e7bf7bc").Error; err == nil {
		zipPath := "uploads/interactive/interactive_1792687445_4f800e20.tendies"
		cleanID := "4f800e20"
		res, err := services.ParseTendiesArchive(zipPath, "uploads/interactive", cleanID)
		if err == nil {
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
			p.CoverImage = res.CoverURL
			db.Save(&p)
			fmt.Println("Successfully updated 布欧超三 (77055a8c-79b5-434c-b7a1-eb541e7bf7bc)!")
		}
	}
}
