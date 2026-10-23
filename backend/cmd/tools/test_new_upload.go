package main

import (
	"fmt"
	"path/filepath"
	"wallpaper-shop/services"
)

func main() {
	uploadDir := filepath.Join(".", "uploads", "interactive")
	zipPath := filepath.Join(uploadDir, "interactive_1792560062_c0140898.tendies")
	cleanID := "c0140898"

	fmt.Printf("Testing ParseTendiesArchive for %s...\n", zipPath)
	result, err := services.ParseTendiesArchive(zipPath, uploadDir, cleanID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Result: Canvas: %vx%v, BG: %d, FL: %d, FG: %d, Total: %d\n",
		result.CanvasWidth, result.CanvasHeight, len(result.BGLayers), len(result.FloatingLayers), len(result.FGLayers), result.TotalImageCount)
	fmt.Printf("BackgroundURL: %s\nForegroundURL: %s\n", result.BackgroundURL, result.ForegroundURL)
	for i, l := range result.BGLayers {
		fmt.Printf("  BG [%d]: id=%s name=%s url=%s op=%v states=%d\n", i, l.ID, l.Name, l.URL, l.Opacity, len(l.States))
	}
	for i, l := range result.FloatingLayers {
		fmt.Printf("  FL [%d]: id=%s name=%s url=%s op=%v states=%d\n", i, l.ID, l.Name, l.URL, l.Opacity, len(l.States))
	}
}
