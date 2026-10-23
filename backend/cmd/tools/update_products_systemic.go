package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

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

	targets := []string{
		"ad1f2610-02ac-4d80-a309-c4eaac122be5",
		"77055a8c-79b5-434c-b7a1-eb541e7bf7bc",
	}

	for _, id := range targets {
		var p models.Product
		if err := db.First(&p, "id = ?", id).Error; err != nil {
			fmt.Printf("Product %s not found: %v\n", id, err)
			continue
		}

		var cfg map[string]interface{}
		_ = json.Unmarshal([]byte(p.InteractiveConfig), &cfg)
		tendiesFile, _ := cfg["tendies_file"].(string)
		if tendiesFile == "" {
			fmt.Printf("Product %s has no tendies_file\n", id)
			continue
		}

		cleanID := strings.TrimSuffix(filepath.Base(tendiesFile), ".tendies")
		cleanID = strings.TrimPrefix(cleanID, "interactive_")
		if idx := strings.Index(cleanID, "_"); idx != -1 {
			cleanID = cleanID[idx+1:]
		}

		relZip := strings.TrimPrefix(tendiesFile, "/")
		fmt.Printf("Re-parsing %s (archive: %s, cleanID: %s)...\n", p.TitleZh, relZip, cleanID)

		res, err := services.ParseTendiesArchive(relZip, "uploads/interactive", cleanID)
		if err != nil {
			fmt.Println("Parse error:", err)
			continue
		}

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

		if err := db.Save(&p).Error; err != nil {
			fmt.Println("Save error:", err)
			continue
		}

		fmt.Printf("Successfully updated %s (%s)! New cover: %s\n", p.TitleZh, p.ID, p.CoverImage)
	}
}
