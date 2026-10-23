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

	var p models.Product
	if err := db.Where("title_zh LIKE ?", "%布欧%").First(&p).Error; err != nil {
		fmt.Println("Product not found:", err)
		return
	}

	fmt.Printf("ID: %s\nTitleZh: %s\nCoverImage: %s\nIsInteractive: %v\nType: %s\n", p.ID, p.TitleZh, p.CoverImage, p.IsInteractive, p.InteractiveType)
	var cfg map[string]interface{}
	json.Unmarshal([]byte(p.InteractiveConfig), &cfg)
	fmt.Printf("tendies_file: %v\n", cfg["tendies_file"])
	fmt.Printf("background_url: %v\n", cfg["background_url"])
	fmt.Printf("foreground_url: %v\n", cfg["foreground_url"])
	fmt.Printf("total_image_count: %v\n", cfg["total_image_count"])
	fmt.Printf("canvas_width: %v, canvas_height: %v\n", cfg["canvas_width"], cfg["canvas_height"])
	bgList, _ := cfg["bg_layers"].([]interface{})
	flList, _ := cfg["floating_layers"].([]interface{})
	fmt.Printf("bg_layers count: %d\n", len(bgList))
	for i, l := range bgList {
		m, _ := l.(map[string]interface{})
		statesB, _ := json.Marshal(m["states"])
		fmt.Printf("  BG [%d]: name=%v scale=%v bounds=%vx%v states=%s\n", i, m["name"], m["scale"], m["width"], m["height"], string(statesB))
	}
	fmt.Printf("floating_layers count: %d\n", len(flList))
	for i, l := range flList {
		m, _ := l.(map[string]interface{})
		statesB, _ := json.Marshal(m["states"])
		fmt.Printf("  FL [%d]: name=%v scale=%v bounds=%vx%v states=%s\n", i, m["name"], m["scale"], m["width"], m["height"], string(statesB))
	}
}
