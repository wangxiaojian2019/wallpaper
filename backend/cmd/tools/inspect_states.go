package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	zipPath := `d:\WALLPAPER\backend\uploads\interactive\interactive_1792717701_399ca199.tendies`
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.Contains(f.Name, "Floating") && strings.HasSuffix(f.Name, "main.caml") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			data, _ := io.ReadAll(rc)
			rc.Close()
			s := string(data)
			// Find <states> block
			if idx := strings.Index(s, "<states>"); idx != -1 {
				end := strings.Index(s, "</states>")
				if end != -1 {
					fmt.Println(s[idx : end+9])
				}
			}
		}
	}
}
