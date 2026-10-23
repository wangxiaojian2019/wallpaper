package main

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type LayerState struct {
	PosX    float64 `json:"pos_x"`
	PosY    float64 `json:"pos_y"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Opacity float64 `json:"opacity"`
}

type TendiesLayer struct {
	ID       string                `json:"id"`
	Name     string                `json:"name"`
	URL      string                `json:"url"`
	Category string                `json:"category"`
	Order    int                   `json:"order"`
	PosX     float64               `json:"pos_x"`
	PosY     float64               `json:"pos_y"`
	Width    float64               `json:"width"`
	Height   float64               `json:"height"`
	Opacity  float64               `json:"opacity"`
	States   map[string]LayerState `json:"states"`
}

func parseNumbers(s string) []float64 {
	fields := strings.Fields(s)
	res := make([]float64, 0, len(fields))
	for _, f := range fields {
		if val, err := strconv.ParseFloat(f, 64); err == nil {
			res = append(res, val)
		}
	}
	return res
}

func parseCAMLLayers(camlText string, category string, assetURLMap map[string]string) ([]TendiesLayer, float64, float64) {
	canvasW := 390.0
	canvasH := 844.0

	boundsRe := regexp.MustCompile(`bounds="([^"]+)"`)
	if m := boundsRe.FindStringSubmatch(camlText); len(m) > 1 {
		nums := parseNumbers(m[1])
		if len(nums) >= 4 && nums[2] > 0 && nums[3] > 0 {
			canvasW = nums[2]
			canvasH = nums[3]
		}
	}

	idRe := regexp.MustCompile(`id="([^"]+)"`)
	nameRe := regexp.MustCompile(`name="([^"]+)"`)
	posRe := regexp.MustCompile(`position="([^"]+)"`)
	bndRe := regexp.MustCompile(`bounds="([^"]+)"`)
	opacityRe := regexp.MustCompile(`opacity="([^"]+)"`)
	imgRe := regexp.MustCompile(`<contents>\s*<CGImage\s+src="([^"]+)"`)

	parts := strings.Split(camlText, "<CALayer ")
	layers := make([]TendiesLayer, 0)
	claimedImages := make(map[string]bool)

	for idx, part := range parts[1:] {
		tagEnd := strings.Index(part, ">")
		if tagEnd == -1 {
			continue
		}
		attrs := part[:tagEnd]

		subLayerIdx := strings.Index(part, "<sublayers>")
		contentsIdx := strings.Index(part, "<contents>")

		if contentsIdx != -1 && (subLayerIdx == -1 || contentsIdx < subLayerIdx) {
			imgMatch := imgRe.FindStringSubmatch(part[contentsIdx:])
			if len(imgMatch) <= 1 {
				continue
			}

			srcPath := imgMatch[1]
			srcBase := filepath.Base(srcPath)
			layerID := ""
			if m := idRe.FindStringSubmatch(attrs); len(m) > 1 {
				layerID = m[1]
			}
			layerName := srcBase
			if nm := nameRe.FindStringSubmatch(attrs); len(nm) > 1 {
				layerName = nm[1]
			}

			imgURL, exists := assetURLMap[strings.ToLower(srcBase)]
			if !exists {
				imgURL, exists = assetURLMap[strings.ToLower(layerName)]
			}
			if !exists {
				for k, v := range assetURLMap {
					if strings.Contains(k, strings.ToLower(srcBase)) || strings.Contains(strings.ToLower(srcBase), k) {
						imgURL = v
						exists = true
						break
					}
				}
			}

			claimedImages[strings.ToLower(srcBase)] = true
			claimedImages[strings.ToLower(layerName)] = true

			posX := canvasW / 2
			posY := canvasH / 2
			if pm := posRe.FindStringSubmatch(attrs); len(pm) > 1 {
				nums := parseNumbers(pm[1])
				if len(nums) >= 2 {
					posX = nums[0]
					posY = nums[1]
				}
			}

			width := canvasW
			height := canvasH
			if bm := bndRe.FindStringSubmatch(attrs); len(bm) > 1 {
				nums := parseNumbers(bm[1])
				if len(nums) >= 4 {
					width = nums[2]
					height = nums[3]
				}
			}

			opacity := 1.0
			if om := opacityRe.FindStringSubmatch(attrs); len(om) > 1 {
				if v, err := strconv.ParseFloat(om[1], 64); err == nil {
					opacity = v
				}
			}

			layers = append(layers, TendiesLayer{
				ID:       layerID,
				Name:     layerName,
				URL:      imgURL,
				Category: category,
				Order:    idx,
				PosX:     posX,
				PosY:     posY,
				Width:    width,
				Height:   height,
				Opacity:  opacity,
				States:   make(map[string]LayerState),
			})
		}
	}

	stateBlockRe := regexp.MustCompile(`(?s)<LKState\s+name="([^"]+)">(.*?)</LKState>`)
	stateMatches := stateBlockRe.FindAllStringSubmatch(camlText, -1)

	setValueRe := regexp.MustCompile(`(?s)<LKStateSetValue\s+targetId="([^"]+)"\s+keyPath="([^"]+)">(.*?)</LKStateSetValue>`)
	valRe := regexp.MustCompile(`value="([^"]+)"`)

	for _, sm := range stateMatches {
		stateName := strings.ToLower(sm[1])
		stateContent := sm[2]

		setValues := setValueRe.FindAllStringSubmatch(stateContent, -1)
		for _, sv := range setValues {
			targetID := sv[1]
			keyPath := sv[2]
			valBlock := sv[3]

			valStr := ""
			if vm := valRe.FindStringSubmatch(valBlock); len(vm) > 1 {
				valStr = vm[1]
			}
			val, _ := strconv.ParseFloat(valStr, 64)

			for i := range layers {
				if layers[i].ID == targetID {
					st, exists := layers[i].States[stateName]
					if !exists {
						st = LayerState{
							PosX:    layers[i].PosX,
							PosY:    layers[i].PosY,
							Width:   layers[i].Width,
							Height:  layers[i].Height,
							Opacity: layers[i].Opacity,
						}
					}
					switch keyPath {
					case "position.x":
						st.PosX = val
					case "position.y":
						st.PosY = val
					case "bounds.size.width":
						st.Width = val
					case "bounds.size.height":
						st.Height = val
					case "opacity":
						st.Opacity = val
					case "position":
						nums := parseNumbers(valStr)
						if len(nums) >= 2 {
							st.PosX = nums[0]
							st.PosY = nums[1]
						}
					case "bounds.size", "bounds":
						nums := parseNumbers(valStr)
						if len(nums) >= 2 {
							st.Width = nums[0]
							st.Height = nums[1]
						}
					}
					layers[i].States[stateName] = st
				}
			}
		}
	}

	return layers, canvasW, canvasH
}

func main() {
	files, _ := filepath.Glob(`uploads/interactive/interactive_1792511387_3c643675.tendies`)
	for _, zipPath := range files {
		r, err := zip.OpenReader(zipPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer r.Close()

		assetURLMap := make(map[string]string)
		camlContents := make(map[string]string)

		for _, f := range r.File {
			nameLower := strings.ToLower(f.Name)
			if strings.HasSuffix(nameLower, ".png") || strings.HasSuffix(nameLower, ".jpg") || strings.HasSuffix(nameLower, ".jpeg") || strings.HasSuffix(nameLower, ".webp") {
				base := filepath.Base(nameLower)
				assetURLMap[base] = fmt.Sprintf("/uploads/interactive/extracted_test/%s", base)
			}
			if filepath.Base(nameLower) == "main.caml" {
				rc, err := f.Open()
				if err == nil {
					data, _ := io.ReadAll(rc)
					rc.Close()
					camlContents[f.Name] = string(data)
				}
			}
		}

		for path, content := range camlContents {
			pathLower := strings.ToLower(path)
			if strings.Contains(pathLower, "background") {
				layers, _, _ := parseCAMLLayers(content, "background", assetURLMap)
				fmt.Printf("\n=== BG Layers (%d) ===\n", len(layers))
				for _, l := range layers {
					lockOp := l.Opacity
					if st, ok := l.States["locked"]; ok {
						lockOp = st.Opacity
					}
					unlockOp := lockOp
					if st, ok := l.States["unlock"]; ok {
						unlockOp = st.Opacity
					}
					fmt.Printf("  BG [%d]: %s (id: %s) pos: (%.1f, %.1f) bounds: %.1fx%.1f Op: init=%.1f locked=%.1f unlock=%.1f states=%d\n",
						l.Order, l.Name, l.ID, l.PosX, l.PosY, l.Width, l.Height, l.Opacity, lockOp, unlockOp, len(l.States))
				}
			} else if strings.Contains(pathLower, "floating") {
				layers, _, _ := parseCAMLLayers(content, "floating", assetURLMap)
				fmt.Printf("\n=== Floating Layers (%d) ===\n", len(layers))
				for _, l := range layers {
					lockOp := l.Opacity
					if st, ok := l.States["locked"]; ok {
						lockOp = st.Opacity
					}
					unlockOp := lockOp
					if st, ok := l.States["unlock"]; ok {
						unlockOp = st.Opacity
					}
					fmt.Printf("  FL [%d]: %s (id: %s) pos: (%.1f, %.1f) bounds: %.1fx%.1f Op: init=%.1f locked=%.1f unlock=%.1f states=%d\n",
						l.Order, l.Name, l.ID, l.PosX, l.PosY, l.Width, l.Height, l.Opacity, lockOp, unlockOp, len(l.States))
				}
			}
		}
	}
}
