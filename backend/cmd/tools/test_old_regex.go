package main

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	uploadDir := filepath.Join(".", "uploads", "interactive")
	zipPath := filepath.Join(uploadDir, "interactive_1792560062_c0140898.tendies")

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer r.Close()

	cgImageLayerRe := regexp.MustCompile(`(?s)<CALayer\s+([^>]*id="([^"]+)"[^>]*)>.*?<contents>\s*<CGImage\s+src="([^"]+)"\s*/>\s*</contents>`)

	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), "main.caml") {
			rc, _ := f.Open()
			data, _ := io.ReadAll(rc)
			rc.Close()
			camlText := string(data)
			matches := cgImageLayerRe.FindAllStringSubmatch(camlText, -1)
			fmt.Printf("File %s: %d matches with OLD regex\n", f.Name, len(matches))
			for i, m := range matches {
				fmt.Printf("  Match [%d]: id=%s src=%s\n", i, m[2], m[3])
			}
		}
	}
}
