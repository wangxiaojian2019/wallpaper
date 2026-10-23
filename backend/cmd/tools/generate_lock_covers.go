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
		fmt.Println("DB error:", err)
		return
	}

	var products []models.Product
	if err := db.Where("is_interactive = ? AND interactive_type = ?", true, "tendies").Find(&products).Error; err != nil {
		fmt.Println("Error finding products:", err)
		return
	}

	uploadDir := filepath.Join(".", "uploads", "interactive")

	for _, p := range products {
		var cfg map[string]interface{}
		if err := json.Unmarshal([]byte(p.InteractiveConfig), &cfg); err != nil {
			continue
		}

		tendiesFile, ok := cfg["tendies_file"].(string)
		if !ok || tendiesFile == "" {
			continue
		}

		tendiesBase := filepath.Base(tendiesFile)
		zipPath := filepath.Join(uploadDir, tendiesBase)
		if _, err := os.Stat(zipPath); os.IsNotExist(err) {
			fmt.Printf("File not found on disk: %s\n", zipPath)
			continue
		}

		cleanID := strings.TrimSuffix(tendiesBase, filepath.Ext(tendiesBase))
		parts := strings.Split(cleanID, "_")
		if len(parts) >= 3 {
			cleanID = parts[2]
		}

		fmt.Printf("Processing %s (%s)...\n", p.TitleZh, tendiesBase)
		result, err := services.ParseTendiesArchive(zipPath, uploadDir, cleanID)
		if err != nil {
			fmt.Printf("  Parse failed: %v\n", err)
			continue
		}

		if result.CoverURL != "" {
			p.CoverImage = result.CoverURL
			cfg["cover_url"] = result.CoverURL
			newCfgBytes, _ := json.Marshal(cfg)
			p.InteractiveConfig = string(newCfgBytes)

			if err := db.Save(&p).Error; err != nil {
				fmt.Printf("  Save failed: %v\n", err)
			} else {
				fmt.Printf("  Updated cover for %s -> %s\n", p.TitleZh, result.CoverURL)
			}
		}
	}
	fmt.Println("Done generating lock covers!")
}
