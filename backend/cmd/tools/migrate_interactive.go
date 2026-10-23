package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		fmt.Println("Failed to connect database:", err)
		return
	}

	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		fmt.Println("Error finding products:", err)
		return
	}

	uploadDir := filepath.Join(".", "uploads", "interactive")

	for _, p := range products {
		fmt.Printf("Product #%d: %s (Interactive: %v)\n", p.ID, p.TitleZh, p.IsInteractive)
		if p.InteractiveConfig == "" {
			continue
		}

		var cfg map[string]interface{}
		if err := json.Unmarshal([]byte(p.InteractiveConfig), &cfg); err != nil {
			fmt.Printf("  Failed to parse config for #%d: %v\n", p.ID, err)
			continue
		}

		tendiesFile, _ := cfg["tendies_file"].(string)
		fmt.Printf("  tendies_file: %s\n", tendiesFile)

		// Check if file exists
		if tendiesFile == "" {
			continue
		}

		tendiesBase := filepath.Base(tendiesFile)
		zipPath := filepath.Join(uploadDir, tendiesBase)
		if _, err := os.Stat(zipPath); err != nil {
			// Try finding matching tendies in uploadDir
			files, _ := filepath.Glob(filepath.Join(uploadDir, "*.tendies"))
			for _, f := range files {
				if strings.Contains(f, strings.TrimSuffix(tendiesBase, filepath.Ext(tendiesBase))) {
					zipPath = f
					break
				}
			}
		}

		if _, err := os.Stat(zipPath); err != nil {
			fmt.Printf("  File not found on disk: %s\n", zipPath)
			continue
		}

		// Extract cleanID from filename
		cleanID := fmt.Sprintf("p%d", p.ID)
		parts := strings.Split(tendiesBase, "_")
		if len(parts) >= 3 {
			cleanID = strings.TrimSuffix(parts[len(parts)-1], filepath.Ext(tendiesBase))
		}

		fmt.Printf("  Parsing archive: %s (cleanID: %s)\n", zipPath, cleanID)
		result, err := services.ParseTendiesArchive(zipPath, uploadDir, cleanID)
		if err != nil {
			fmt.Printf("  Parse failed: %v\n", err)
			continue
		}

		fmt.Printf("  Successfully parsed %d BG layers, %d Floating layers, %d FG layers (total %d images)\n",
			len(result.BGLayers), len(result.FloatingLayers), len(result.FGLayers), result.TotalImageCount)

		cfg["bg_layers"] = result.BGLayers
		cfg["floating_layers"] = result.FloatingLayers
		cfg["fg_layers"] = result.FGLayers
		cfg["canvas_width"] = result.CanvasWidth
		cfg["canvas_height"] = result.CanvasHeight
		cfg["total_image_count"] = result.TotalImageCount
		if result.BackgroundURL != "" {
			cfg["background_url"] = result.BackgroundURL
		}
		if result.ForegroundURL != "" {
			cfg["foreground_url"] = result.ForegroundURL
		}

		newCfgBytes, _ := json.Marshal(cfg)
		p.InteractiveConfig = string(newCfgBytes)
		p.IsInteractive = true
		if err := db.Save(&p).Error; err != nil {
			fmt.Printf("  Failed to update product #%d: %v\n", p.ID, err)
		} else {
			fmt.Printf("  Updated product #%d interactive config successfully!\n", p.ID)
		}
	}
}
