package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
)

func getBBox(path string) (image.Rectangle, int, int) {
	f, err := os.Open(path)
	if err != nil {
		return image.Rectangle{}, 0, 0
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return image.Rectangle{}, 0, 0
	}
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 1000 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	return image.Rect(minX, minY, maxX, maxY), b.Dx(), b.Dy()
}

func main() {
	files := []string{
		"uploads/interactive/extracted_4f800e20/1.png",
		"uploads/interactive/extracted_4f800e20/2.png",
		"uploads/interactive/extracted_4f800e20/3.png",
		"uploads/interactive/extracted_4f800e20/buu1.png",
		"uploads/interactive/extracted_4f800e20/chansa2.png",
		"uploads/interactive/extracted_4f800e20/chao3wukong.png",
		"uploads/interactive/extracted_4f800e20/buouu.png",
	}
	for _, fn := range files {
		bbox, w, h := getBBox(fn)
		fmt.Printf("%s: Image Size=(%dx%d), BBox=(%d,%d)-(%d,%d) [w=%d, h=%d]\n",
			fn, w, h, bbox.Min.X, bbox.Min.Y, bbox.Max.X, bbox.Max.Y, bbox.Dx(), bbox.Dy())
	}
}
