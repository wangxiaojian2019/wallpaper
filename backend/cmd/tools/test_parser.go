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
	Category string                `json:"category"` // "background", "floating", "foreground"
	Order    int                   `json:"order"`
	PosX     float64               `json:"pos_x"`
	PosY     float64               `json:"pos_y"`
	Width    float64               `json:"width"`
	Height   float64               `json:"height"`
	Opacity  float64               `json:"opacity"`
	States   map[string]LayerState `json:"states"`
}

type TendiesParseResult struct {
	CanvasWidth     float64        `json:"canvas_width"`
	CanvasHeight    float64        `json:"canvas_height"`
	BackgroundURL   string         `json:"background_url"`
	ForegroundURL   string         `json:"foreground_url"`
	BGLayers        []TendiesLayer `json:"bg_layers"`
	FloatingLayers  []TendiesLayer `json:"floating_layers"`
	FGLayers        []TendiesLayer `json:"fg_layers"`
	TotalImageCount int            `json:"total_image_count"`
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

	// Check root bounds
	boundsRe := regexp.MustCompile(`bounds="([^"]+)"`)
	if m := boundsRe.FindStringSubmatch(camlText); len(m) > 1 {
		nums := parseNumbers(m[1])
		if len(nums) >= 4 && nums[2] > 0 && nums[3] > 0 {
			canvasW = nums[2]
			canvasH = nums[3]
		}
	}

	// Match every CALayer with a CGImage contents
	// Each layer with CGImage is a visual sprite
	cgImageLayerRe := regexp.MustCompile(`(?s)<CALayer\s+([^>]*id="([^"]+)"[^>]*)>.*?<contents>\s*<CGImage\s+src="([^"]+)"\s*/>\s*</contents>`)
	nameRe := regexp.MustCompile(`name="([^"]+)"`)
	posRe := regexp.MustCompile(`position="([^"]+)"`)
	bndRe := regexp.MustCompile(`bounds="([^"]+)"`)
	opacityRe := regexp.MustCompile(`opacity="([^"]+)"`)

	matches := cgImageLayerRe.FindAllStringSubmatch(camlText, -1)
	layers := make([]TendiesLayer, 0)
	claimedImages := make(map[string]bool)

	for idx, m := range matches {
		attrs := m[1]
		layerID := m[2]
		srcPath := m[3]

		srcBase := filepath.Base(srcPath)
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

	// Parse States (<LKState name="Locked">, "Unlock", "Sleep")
	stateBlockRe := regexp.MustCompile(`(?s)<LKState\s+name="([^"]+)">(.*?)</LKState>`)
	stateMatches := stateBlockRe.FindAllStringSubmatch(camlText, -1)

	setValueRe := regexp.MustCompile(`(?s)<LKStateSetValue\s+targetId="([^"]+)"\s+keyPath="([^"]+)">(.*?)</LKStateSetValue>`)
	valRe := regexp.MustCompile(`value="([^"]+)"`)

	for _, sm := range stateMatches {
		stateName := strings.ToLower(sm[1]) // "locked", "unlock", "sleep"
		stateContent := sm[2]

		setValues := setValueRe.FindAllStringSubmatch(stateContent, -1)
		for _, sv := range setValues {
			targetID := sv[1]
			keyPath := sv[2]
			valBlock := sv[3]

			val := 0.0
			if vm := valRe.FindStringSubmatch(valBlock); len(vm) > 1 {
				val, _ = strconv.ParseFloat(vm[1], 64)
			}

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
					}
					layers[i].States[stateName] = st
				}
			}
		}
	}

	return layers, canvasW, canvasH
}

func main() {
	zipPath := `F:\WALLPAPER\backend\uploads\interactive\interactive_1792566650_5a9d4b61.tendies`
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "main.caml") {
			fmt.Println("==================================================")
			fmt.Println("CAML file:", f.Name)
			rc, _ := f.Open()
			b, _ := io.ReadAll(rc)
			rc.Close()
			content := string(b)

			// Print LKState: Unlock
			stRe := regexp.MustCompile(`(?s)<LKState\s+name="Unlock">(.*?)</LKState>`)
			if m := stRe.FindStringSubmatch(content); len(m) > 1 {
				fmt.Println("\n=== LKState: Unlock ===")
				svRe := regexp.MustCompile(`(?s)<LKStateSetValue\s+targetId="([^"]+)"\s+keyPath="([^"]+)">(.*?)</LKStateSetValue>`)
				valRe := regexp.MustCompile(`value="([^"]+)"`)
				for _, sv := range svRe.FindAllStringSubmatch(m[1], -1) {
					val := ""
					if vm := valRe.FindStringSubmatch(sv[3]); len(vm) > 1 {
						val = vm[1]
					}
					fmt.Printf("  Target: %s  KeyPath: %s  Value: %s\n", sv[1], sv[2], val)
				}
			}
		}
	}
}
