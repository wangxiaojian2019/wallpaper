package main

import (
	"encoding/json"
	"fmt"
	"time"
	"wallpaper-shop/models"

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
	db.Order("updated_at desc").Limit(5).Find(&products)

	fmt.Println("=== 5 MOST RECENTLY UPDATED PRODUCTS ===")
	for _, p := range products {
		fmt.Printf("ID: %s\nTitle: %s\nUpdatedAt: %s\nIsInteractive: %v\nType: %s\n",
			p.ID, p.TitleZh, p.UpdatedAt.Format(time.RFC3339), p.IsInteractive, p.InteractiveType)
		var cfg map[string]interface{}
		json.Unmarshal([]byte(p.InteractiveConfig), &cfg)
		fmt.Printf("  tendies_file: %v\n", cfg["tendies_file"])
		fmt.Printf("  background_url: %v\n", cfg["background_url"])
		fmt.Printf("  foreground_url: %v\n", cfg["foreground_url"])
		bgList, _ := cfg["bg_layers"].([]interface{})
		flList, _ := cfg["floating_layers"].([]interface{})
		fmt.Printf("  bg_layers count: %d, floating_layers count: %d\n", len(bgList), len(flList))
		if len(bgList) > 0 {
			first, _ := bgList[0].(map[string]interface{})
			fmt.Printf("  first bg: %v (url: %v)\n", first["name"], first["url"])
		}
		if len(flList) > 0 {
			first, _ := flList[0].(map[string]interface{})
			fmt.Printf("  first fl: %v (url: %v)\n", first["name"], first["url"])
		}
		fmt.Println("-------------------------------------------")
	}
}
