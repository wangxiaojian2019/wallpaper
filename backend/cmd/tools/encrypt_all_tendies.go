package main

import (
	"encoding/json"
	"fmt"
	"os"
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

	interactiveDir := filepath.Join("uploads", "interactive")
	entries, err := os.ReadDir(interactiveDir)
	if err != nil {
		fmt.Println("ReadDir error:", err)
		return
	}

	fmt.Println("Scanning all .tendies archives in uploads/interactive...")

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".tendies") && !strings.HasSuffix(strings.ToLower(name), ".zip") {
			continue
		}

		baseNoExt := strings.TrimSuffix(name, filepath.Ext(name))
		parts := strings.Split(baseNoExt, "_")
		cleanID := ""
		if len(parts) >= 3 {
			cleanID = parts[2]
		} else if len(parts) == 2 {
			cleanID = parts[1]
		} else {
			cleanID = baseNoExt
		}

		archivePath := filepath.Join(interactiveDir, name)
		fmt.Printf("Processing archive %s (cleanID: %s)...\n", name, cleanID)

		res, err := services.ParseTendiesArchive(archivePath, interactiveDir, cleanID)
		if err != nil {
			fmt.Printf("  Parse error: %v\n", err)
			continue
		}

		// Clean up old unencrypted .png/.jpg files in extracted directory
		subDir := filepath.Join(interactiveDir, fmt.Sprintf("extracted_%s", cleanID))
		if subEntries, err := os.ReadDir(subDir); err == nil {
			for _, subEntry := range subEntries {
				subName := strings.ToLower(subEntry.Name())
				if strings.HasSuffix(subName, ".png") || strings.HasSuffix(subName, ".jpg") || strings.HasSuffix(subName, ".jpeg") {
					_ = os.Remove(filepath.Join(subDir, subEntry.Name()))
					fmt.Printf("  Removed unencrypted: %s\n", subEntry.Name())
				}
			}
		}

		fmt.Printf("  Encrypted %d layers -> saved to %s\n", res.TotalImageCount, subDir)

		// Check if any product in database references this archive or cleanID
		var products []models.Product
		db.Where("interactive_config LIKE ? OR interactive_config LIKE ?", "%"+name+"%", "%extracted_"+cleanID+"%").Find(&products)
		for _, p := range products {
			var cfg map[string]interface{}
			if err := json.Unmarshal([]byte(p.InteractiveConfig), &cfg); err == nil {
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
				_ = db.Save(&p)
				fmt.Printf("  Updated product in DB: %s (%s)\n", p.TitleZh, p.ID)
			}
		}
	}

	fmt.Println("All tendies archives processed and encrypted successfully!")
}
