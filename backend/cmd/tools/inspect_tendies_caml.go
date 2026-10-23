package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"regexp"
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
		if strings.HasSuffix(f.Name, "main.caml") {
			fmt.Printf("\n===========================================\n=== CAML FILE: %s ===\n", f.Name)
			rc, err := f.Open()
			if err != nil {
				continue
			}
			data, _ := io.ReadAll(rc)
			rc.Close()
			s := string(data)

			// Print all <CALayer> lines with name, bounds, contents, position
			re := regexp.MustCompile(`(?s)<CALayer\s+([^>]+)>`)
			matches := re.FindAllStringSubmatch(s, -1)
			for idx, m := range matches {
				attrs := m[1]
				name := getAttr(attrs, "name")
				bounds := getAttr(attrs, "bounds")
				pos := getAttr(attrs, "position")
				opacity := getAttr(attrs, "opacity")
				hidden := getAttr(attrs, "hidden")
				fmt.Printf("Layer #%d: name=%q bounds=%q pos=%q opacity=%q hidden=%q\n", idx, name, bounds, pos, opacity, hidden)
			}
		}
	}
}

func getAttr(attrs, key string) string {
	re := regexp.MustCompile(fmt.Sprintf(`%s="([^"]*)"`, key))
	if m := re.FindStringSubmatch(attrs); len(m) > 1 {
		return m[1]
	}
	return ""
}
