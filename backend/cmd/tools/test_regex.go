package main

import (
	"archive/zip"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func main() {
	zipPath := `uploads/interactive/interactive_1792511387_3c643675.tendies`
	r, _ := zip.OpenReader(zipPath)
	defer r.Close()

	idRe := regexp.MustCompile(`id="([^"]+)"`)
	nameRe := regexp.MustCompile(`name="([^"]+)"`)
	opacityRe := regexp.MustCompile(`opacity="([^"]+)"`)
	posRe := regexp.MustCompile(`position="([^"]+)"`)
	bndRe := regexp.MustCompile(`bounds="([^"]+)"`)
	imgRe := regexp.MustCompile(`<contents>\s*<CGImage\s+src="([^"]+)"`)

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "main.caml") && !strings.Contains(f.Name, "Foreground") {
			fmt.Printf("\n=== %s ===\n", f.Name)
			rc, _ := f.Open()
			data, _ := io.ReadAll(rc)
			rc.Close()

			camlText := string(data)
			parts := strings.Split(camlText, "<CALayer ")
			for idx, part := range parts[1:] {
				tagEnd := strings.Index(part, ">")
				if tagEnd == -1 {
					continue
				}
				attrs := part[:tagEnd]
				// Does this layer block have an immediate <contents> before another <CALayer?
				// Note: if there is another <sublayers> or <CALayer, it will be after
				subLayerIdx := strings.Index(part, "<sublayers>")
				contentsIdx := strings.Index(part, "<contents>")

				// If it has <contents> and it comes before any <sublayers>
				if contentsIdx != -1 && (subLayerIdx == -1 || contentsIdx < subLayerIdx) {
					imgMatch := imgRe.FindStringSubmatch(part[contentsIdx:])
					if len(imgMatch) > 1 {
						id := ""
						if m := idRe.FindStringSubmatch(attrs); len(m) > 1 {
							id = m[1]
						}
						name := ""
						if m := nameRe.FindStringSubmatch(attrs); len(m) > 1 {
							name = m[1]
						}
						pos := ""
						if m := posRe.FindStringSubmatch(attrs); len(m) > 1 {
							pos = m[1]
						}
						bnd := ""
						if m := bndRe.FindStringSubmatch(attrs); len(m) > 1 {
							bnd = m[1]
						}
						op := "1"
						if m := opacityRe.FindStringSubmatch(attrs); len(m) > 1 {
							op = m[1]
						}
						fmt.Printf("Layer [%d]: id=%s name=%s src=%s pos=(%s) bnd=(%s) op=%s\n",
							idx, id, name, imgMatch[1], pos, bnd, op)
					}
				}
			}
		}
	}
}
