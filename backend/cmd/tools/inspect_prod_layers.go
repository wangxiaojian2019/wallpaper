package main

import (
	"encoding/json"
	"fmt"
	"os"
	"wallpaper-shop/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TendiesLayerInfo struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	Category string  `json:"category"`
	Order    int     `json:"order"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Opacity  float64 `json:"opacity"`
}

type TendiesConfig struct {
	BGLayers       []TendiesLayerInfo `json:"bg_layers"`
	FloatingLayers []TendiesLayerInfo `json:"floating_layers"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("wallpaper.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Printf("db err: %v\n", err)
		os.Exit(1)
	}

	var p models.Product
	if err := db.Where("id = ?", "a39de3d8-8e52-4a68-ad57-cfb4532dc33e").First(&p).Error; err != nil {
		fmt.Printf("query err: %v\n", err)
		os.Exit(1)
	}

	var cfg TendiesConfig
	if err := json.Unmarshal([]byte(p.InteractiveConfig), &cfg); err != nil {
		fmt.Printf("json err: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("BG layers: %d\n", len(cfg.BGLayers))
	for i, l := range cfg.BGLayers {
		fmt.Printf("  BG[%d]: %s (url: %s), w=%.1f, h=%.1f, op=%.2f, cat=%s\n", i, l.Name, l.URL, l.Width, l.Height, l.Opacity, l.Category)
	}
	fmt.Printf("Floating layers: %d\n", len(cfg.FloatingLayers))
	for i, l := range cfg.FloatingLayers {
		fmt.Printf("  FL[%d]: %s (url: %s), w=%.1f, h=%.1f, op=%.2f, cat=%s\n", i, l.Name, l.URL, l.Width, l.Height, l.Opacity, l.Category)
	}
}
