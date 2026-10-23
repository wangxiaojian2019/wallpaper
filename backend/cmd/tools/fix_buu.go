package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("DB err:", err)
		return
	}

	var p models.Product
	if err := db.Where("title_zh LIKE ?", "%布欧%").First(&p).Error; err != nil {
		fmt.Println("Product not found:", err)
		return
	}

	fmt.Printf("Updating product %s (ID: %s)\n", p.TitleZh, p.ID)

	cleanID := "4f800e20"
	zipPath := "uploads/interactive/interactive_1792687445_4f800e20.tendies"
	uploadDir := "uploads/interactive"

	parsed, err := services.ParseTendiesArchive(zipPath, uploadDir, cleanID)
	if err != nil {
		fmt.Println("Parse err:", err)
		return
	}

	// Align the sliced comic panels 1.png, 2.png, 3.png in base and locked states:
	for i := range parsed.BGLayers {
		l := &parsed.BGLayers[i]
		if l.Name == "1.png" || l.Name == "2.png" || l.Name == "3.png" {
			l.PosX = 195
			l.Scale = 1.76
			if st, ok := l.States["locked"]; ok {
				st.PosX = 195
				st.Scale = 1.76
				l.States["locked"] = st
			}
			if st, ok := l.States["sleep"]; ok {
				st.PosX = 195
				st.Scale = 1.76
				l.States["sleep"] = st
			}
		}
	}

	// Regenerate lock screen cover with the aligned layers
	coverURL := services.GenerateLockScreenCover(parsed, uploadDir, cleanID)
	fmt.Println("Regenerated cover:", coverURL)

	configBytes, _ := json.Marshal(parsed)
	p.InteractiveConfig = string(configBytes)
	p.CoverImage = coverURL

	if err := db.Save(&p).Error; err != nil {
		fmt.Println("Save err:", err)
		return
	}

	fmt.Println("Successfully updated product in database!")
}
