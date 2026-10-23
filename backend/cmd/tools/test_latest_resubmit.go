package main

import (
	"encoding/json"
	"fmt"
	"wallpaper-shop/services"
)

func main() {
	zipPath := "uploads/interactive/interactive_1792692522_962db7f3.tendies"
	uploadDir := "uploads/interactive"
	cleanID := "test_new_upload"

	res, err := services.ParseTendiesArchive(zipPath, uploadDir, cleanID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("=== NEW PARSER OUTPUT FOR USER'S RE-SUBMITTED FILE ===")
	fmt.Printf("Canvas: %vx%v, Total Images: %d\n", res.CanvasWidth, res.CanvasHeight, res.TotalImageCount)
	fmt.Println("\n--- BG LAYERS ---")
	for i, l := range res.BGLayers {
		st, _ := json.Marshal(l.States)
		fmt.Printf("[%d] %s: pos=(%v,%v) size=(%v,%v) scale=%v op=%v states=%s\n",
			i, l.Name, l.PosX, l.PosY, l.Width, l.Height, l.Scale, l.Opacity, string(st))
	}
	fmt.Println("\n--- FLOATING LAYERS ---")
	for i, l := range res.FloatingLayers {
		st, _ := json.Marshal(l.States)
		fmt.Printf("[%d] %s: pos=(%v,%v) size=(%v,%v) scale=%v op=%v states=%s\n",
			i, l.Name, l.PosX, l.PosY, l.Width, l.Height, l.Scale, l.Opacity, string(st))
	}
	fmt.Println("\nCover URL:", res.CoverURL)
}
