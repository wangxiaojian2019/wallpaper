package main

import (
	"archive/zip"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	zipPath := `uploads/interactive/interactive_1792687445_4f800e20.tendies`
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer r.Close()

	out, _ := os.Create("buu_caml.txt")
	defer out.Close()

	fmt.Println("=== IMAGES IN BUU ===")
	for _, f := range r.File {
		name := strings.ToLower(f.Name)
		if strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg") {
			rc, err := f.Open()
			if err == nil {
				cfg, _, err := image.DecodeConfig(rc)
				if err == nil {
					fmt.Printf("Image: %s (%dx%d)\n", filepath.Base(f.Name), cfg.Width, cfg.Height)
				}
				rc.Close()
			}
		}
	}

	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), "main.caml") {
			rc, _ := f.Open()
			data, _ := io.ReadAll(rc)
			rc.Close()
			out.WriteString(fmt.Sprintf("\n=========================================\nCAML: %s\n=========================================\n", f.Name))
			out.Write(data)
		}
	}
	fmt.Println("Wrote buu_caml.txt")
}
