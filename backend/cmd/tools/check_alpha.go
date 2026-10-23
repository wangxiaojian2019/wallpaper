package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

func main() {
	dir := filepath.Join(".", "uploads", "interactive", "extracted_71efc71f")
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		fp := filepath.Join(dir, e.Name())
		f, err := os.Open(fp)
		if err != nil {
			continue
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			continue
		}
		b := img.Bounds()
		// Check opacity at top area around (w/2, 150)
		ptColor := img.At(b.Dx()/2, b.Dy()/6)
		_, _, _, a := ptColor.RGBA()
		fmt.Printf("File: %-15s Bounds: %4dx%-4d Pixel at top-center alpha: %d\n", e.Name(), b.Dx(), b.Dy(), a)
	}
}
