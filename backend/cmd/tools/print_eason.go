package main

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

func main() {
	zipPath := filepath.Join(".", "uploads", "interactive", "interactive_1792563420_71efc71f.tendies")
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.Contains(f.Name, "Floating") && strings.HasSuffix(f.Name, "main.caml") {
			rc, _ := f.Open()
			data, _ := io.ReadAll(rc)
			rc.Close()
			lines := strings.Split(string(data), "\n")
			for i := 0; i < len(lines) && i < 150; i++ {
				fmt.Printf("%3d: %s\n", i+1, lines[i])
			}
		}
	}
}
