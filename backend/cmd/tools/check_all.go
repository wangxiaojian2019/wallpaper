package main

import (
	"encoding/json"
	"fmt"
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
	db.Find(&products)
	fmt.Printf("Total products: %d\n", len(products))
	for _, p := range products {
		fmt.Printf("ID: %s | Title: %s | IsInteractive: %v\n", p.ID, p.TitleZh, p.IsInteractive)
		if p.IsInteractive {
			var cfg map[string]interface{}
			json.Unmarshal([]byte(p.InteractiveConfig), &cfg)
			bgList, _ := cfg["bg_layers"].([]interface{})
			flList, _ := cfg["floating_layers"].([]interface{})
			fgList, _ := cfg["fg_layers"].([]interface{})
			fmt.Printf("  bg_layers: %d, floating_layers: %d, fg_layers: %d, fg_url: %v, bg_url: %v\n",
				len(bgList), len(flList), len(fgList), cfg["foreground_url"], cfg["background_url"])
			for i, l := range bgList {
				m, _ := l.(map[string]interface{})
				fmt.Printf("    BG[%d]: name=%v, category=%v, order=%v, url=%v, bounds=%vx%v\n", i, m["name"], m["category"], m["order"], m["url"], m["width"], m["height"])
			}
			for i, l := range flList {
				m, _ := l.(map[string]interface{})
				fmt.Printf("    FL[%d]: name=%v, category=%v, order=%v, url=%v, bounds=%vx%v\n", i, m["name"], m["category"], m["order"], m["url"], m["width"], m["height"])
			}
		}
	}
}
