package services

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// WallArtSecretKey is the proprietary symmetric obfuscation key for web asset protection
var WallArtSecretKey = []byte("WallArtTendiesKey2026")

// EncryptTendiesBytes encrypts raw image bytes into obfuscated .dat format with WART header
func EncryptTendiesBytes(raw []byte, salt byte) []byte {
	out := make([]byte, 6+len(raw))
	out[0] = 0x57 // 'W'
	out[1] = 0x41 // 'A'
	out[2] = 0x52 // 'R'
	out[3] = 0x54 // 'T'
	out[4] = 0x01 // Version 1
	out[5] = salt
	keyLen := len(WallArtSecretKey)
	for i, b := range raw {
		k := WallArtSecretKey[(i+int(salt))%keyLen]
		shift := byte((i*7 + int(salt)) & 0xFF)
		out[6+i] = b ^ k ^ shift
	}
	return out
}

// DecryptTendiesBytes decrypts obfuscated .dat format into raw image bytes
func DecryptTendiesBytes(enc []byte) ([]byte, bool) {
	if len(enc) < 6 {
		return enc, false
	}
	if enc[0] != 0x57 || enc[1] != 0x41 || enc[2] != 0x52 || enc[3] != 0x54 {
		return enc, false // Not encrypted, return as is
	}
	salt := enc[5]
	rawLen := len(enc) - 6
	out := make([]byte, rawLen)
	keyLen := len(WallArtSecretKey)
	for i := 0; i < rawLen; i++ {
		k := WallArtSecretKey[(i+int(salt))%keyLen]
		shift := byte((i*7 + int(salt)) & 0xFF)
		out[i] = enc[6+i] ^ k ^ shift
	}
	return out, true
}

// LayerState holds transform properties for a specific state (locked, unlock, sleep)
type LayerState struct {
	PosX     float64 `json:"pos_x"`
	PosY     float64 `json:"pos_y"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Opacity  float64 `json:"opacity"`
	Rotation float64 `json:"rotation"` // in radians
	Scale    float64 `json:"scale"`
}

// LayerAnimation defines continuous loop keyframe animation parsed from CAML <animations>
type LayerAnimation struct {
	KeyPath      string    `json:"key_path"`          // "position.y", "position.x", "position", "transform.rotation.z", "opacity"
	Duration     float64   `json:"duration"`          // effective duration in seconds
	Autoreverses bool      `json:"autoreverses"`      // alternate
	RepeatCount  string    `json:"repeat_count"`      // "infinite"
	DeltaX       []float64 `json:"delta_x,omitempty"` // relative delta from base PosX
	DeltaY       []float64 `json:"delta_y,omitempty"` // relative delta from base PosY (converted to screen coords)
	RotDeg       []float64 `json:"rot_deg,omitempty"` // rotation in degrees
	Opacity      []float64 `json:"opacity,omitempty"` // opacity 0-1
}

// TendiesLayer represents an individual visual sprite layer with CAML coordinates and states
type TendiesLayer struct {
	ID         string                `json:"id"`
	Name       string                `json:"name"`
	URL        string                `json:"url"`
	Category   string                `json:"category"` // "background", "floating", "foreground"
	Order      int                   `json:"order"`
	PosX       float64               `json:"pos_x"`
	PosY       float64               `json:"pos_y"`
	Width      float64               `json:"width"`
	Height     float64               `json:"height"`
	Opacity    float64               `json:"opacity"`
	Rotation   float64               `json:"rotation"` // in radians
	Scale      float64               `json:"scale"`
	States     map[string]LayerState `json:"states"`
	Animations []LayerAnimation      `json:"animations,omitempty"`
}

// TendiesParseResult is the complete structure returned after parsing an Apple Core Animation .tendies archive
type TendiesParseResult struct {
	CanvasWidth     float64        `json:"canvas_width"`
	CanvasHeight    float64        `json:"canvas_height"`
	BackgroundURL   string         `json:"background_url"`
	ForegroundURL   string         `json:"foreground_url"`
	CoverURL        string         `json:"cover_url"`
	BGLayers        []TendiesLayer `json:"bg_layers"`
	FloatingLayers  []TendiesLayer `json:"floating_layers"`
	FGLayers        []TendiesLayer `json:"fg_layers"`
	TotalImageCount int            `json:"total_image_count"`
}

func normalizeAssetKey(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "%20", "")
	return s
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

type camlXMLNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr    `xml:",any,attr"`
	InnerXML []byte        `xml:",innerxml"`
	Children []camlXMLNode `xml:",any"`
}

func getCAMLAttr(node *camlXMLNode, name string) string {
	for _, a := range node.Attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}

type camlLayerNode struct {
	ID         string
	Name       string
	ImageSrc   string
	BoundsW    float64
	BoundsH    float64
	LocalX     float64
	LocalY     float64
	Opacity    float64
	Rotation   float64
	Scale      float64
	Flipped    bool
	Parent     *camlLayerNode
	Sublayers  []*camlLayerNode
	Animations []LayerAnimation
}

func parseCAMLLayers(camlText string, category string, assetURLMap map[string]string) ([]TendiesLayer, float64, float64) {
	var root camlXMLNode
	if err := xml.Unmarshal([]byte(camlText), &root); err != nil {
		fmt.Printf("CAML XML parse fallback error: %v\n", err)
		return nil, 390, 844
	}

	// 1. Root geometryFlipped detection
	isGeometryFlipped := false
	geomRe := regexp.MustCompile(`geometryFlipped="([^"]+)"`)
	if gm := geomRe.FindStringSubmatch(camlText); len(gm) > 1 {
		if gm[1] == "1" || strings.ToLower(gm[1]) == "true" {
			isGeometryFlipped = true
		}
	}

	// 2. Parse LKStates (<LKState name="Locked">, "Unlock", "Sleep")
	stateBlockRe := regexp.MustCompile(`(?s)<LKState\s+name="([^"]+)">(.*?)</LKState>`)
	stateMatches := stateBlockRe.FindAllStringSubmatch(camlText, -1)
	setValueRe := regexp.MustCompile(`(?s)<LKStateSetValue\s+targetId="([^"]+)"\s+keyPath="([^"]+)">(.*?)</LKStateSetValue>`)
	valRe := regexp.MustCompile(`value="([^"]+)"`)

	stateOverrides := make(map[string]map[string]map[string]float64)

	for _, sm := range stateMatches {
		sName := strings.ToLower(sm[1])
		if _, ok := stateOverrides[sName]; !ok {
			stateOverrides[sName] = make(map[string]map[string]float64)
		}
		setValues := setValueRe.FindAllStringSubmatch(sm[2], -1)
		for _, sv := range setValues {
			tid := sv[1]
			kp := sv[2]
			vStr := ""
			if vm := valRe.FindStringSubmatch(sv[3]); len(vm) > 1 {
				vStr = vm[1]
			}
			fVal, _ := strconv.ParseFloat(vStr, 64)
			if _, ok := stateOverrides[sName][tid]; !ok {
				stateOverrides[sName][tid] = make(map[string]float64)
			}
			stateOverrides[sName][tid][kp] = fVal
		}
	}

	// 3. Build CALayer Hierarchy Tree
	var buildTree func(node *camlXMLNode, parent *camlLayerNode) *camlLayerNode
	buildTree = func(node *camlXMLNode, parent *camlLayerNode) *camlLayerNode {
		if node.XMLName.Local != "CALayer" {
			for i := range node.Children {
				if n := buildTree(&node.Children[i], parent); n != nil {
					return n
				}
			}
			return nil
		}

		ln := &camlLayerNode{
			ID:       getCAMLAttr(node, "id"),
			Name:     getCAMLAttr(node, "name"),
			Opacity:  1.0,
			Scale:    1.0,
			Rotation: 0.0,
			Flipped:  isGeometryFlipped,
			Parent:   parent,
		}

		if flp := getCAMLAttr(node, "geometryFlipped"); flp != "" {
			ln.Flipped = (flp == "1" || strings.ToLower(flp) == "true")
		} else if parent != nil {
			ln.Flipped = parent.Flipped
		}

		if opStr := getCAMLAttr(node, "opacity"); opStr != "" {
			if v, err := strconv.ParseFloat(opStr, 64); err == nil {
				ln.Opacity = v
			}
		}
		if hidStr := getCAMLAttr(node, "hidden"); hidStr == "1" || strings.ToLower(hidStr) == "true" {
			ln.Opacity = 0.0
		}
		if rotStr := getCAMLAttr(node, "transform.rotation.z"); rotStr != "" {
			if v, err := strconv.ParseFloat(rotStr, 64); err == nil {
				ln.Rotation = v
			}
		}

		// Parse scale from transform attribute (e.g. scale(1.76, 1.76, 1)) or transform.scale
		trStr := getCAMLAttr(node, "transform")
		if trStr != "" {
			scaleRe := regexp.MustCompile(`scale\s*\(\s*([\d\.\-]+)`)
			if m := scaleRe.FindStringSubmatch(trStr); len(m) > 1 {
				if v, err := strconv.ParseFloat(m[1], 64); err == nil && v > 0 {
					ln.Scale = v
				}
			}
		}
		if scStr := getCAMLAttr(node, "transform.scale.xy"); scStr != "" {
			if v, err := strconv.ParseFloat(scStr, 64); err == nil && v > 0 {
				ln.Scale = v
			}
		} else if scStr := getCAMLAttr(node, "transform.scale.x"); scStr != "" {
			if v, err := strconv.ParseFloat(scStr, 64); err == nil && v > 0 {
				ln.Scale = v
			}
		}

		if bStr := getCAMLAttr(node, "bounds"); bStr != "" {
			nums := parseNumbers(bStr)
			if len(nums) >= 4 {
				ln.BoundsW = nums[2]
				ln.BoundsH = nums[3]
			}
		}

		if pStr := getCAMLAttr(node, "position"); pStr != "" {
			nums := parseNumbers(pStr)
			if len(nums) >= 2 {
				ln.LocalX = nums[0]
				ln.LocalY = nums[1]
			}
		}

		// Parse continuous animations if present
		animBlock := string(node.InnerXML)
		if animIdx := strings.Index(animBlock, "<animations>"); animIdx != -1 {
			singleAnimRe := regexp.MustCompile(`(?s)<(?:animation|p)\s+([^>]*keyPath="([^"]+)"[^>]*)>(.*?)</(?:animation|p)>`)
			animMatches := singleAnimRe.FindAllStringSubmatch(animBlock, -1)
			durRe := regexp.MustCompile(`duration="([^"]+)"`)
			spdRe := regexp.MustCompile(`speed="([^"]+)"`)
			autoRe := regexp.MustCompile(`autoreverses="([^"]+)"`)
			repRe := regexp.MustCompile(`repeatCount="([^"]+)"`)
			intRealValRe := regexp.MustCompile(`<(?:integer|real)\s+value="([^"]+)"`)
			pointValRe := regexp.MustCompile(`<CGPoint\s+value="([^"]+)"`)

			for _, am := range animMatches {
				aAttrs := am[1]
				keyPath := am[2]
				aBody := am[3]

				duration := 5.0
				if dm := durRe.FindStringSubmatch(aAttrs); len(dm) > 1 {
					if d, err := strconv.ParseFloat(dm[1], 64); err == nil && d > 0 {
						duration = d
					}
				}

				speed := 1.0
				if sm := spdRe.FindStringSubmatch(aAttrs); len(sm) > 1 {
					if sp, err := strconv.ParseFloat(sm[1], 64); err == nil && sp > 0 {
						speed = sp
					}
				}

				autoreverses := false
				if arm := autoRe.FindStringSubmatch(aAttrs); len(arm) > 1 {
					autoreverses = (arm[1] == "1" || strings.ToLower(arm[1]) == "true")
				}

				repeatCount := "infinite"
				if rm := repRe.FindStringSubmatch(aAttrs); len(rm) > 1 {
					repeatCount = rm[1]
				}

				animObj := LayerAnimation{
					KeyPath:      keyPath,
					Duration:     duration / speed,
					Autoreverses: autoreverses,
					RepeatCount:  repeatCount,
				}

				switch keyPath {
				case "position.y":
					vMatches := intRealValRe.FindAllStringSubmatch(aBody, -1)
					for _, vm := range vMatches {
						if rawVal, err := strconv.ParseFloat(vm[1], 64); err == nil {
							animObj.DeltaY = append(animObj.DeltaY, -(rawVal - ln.LocalY))
						}
					}
				case "position.x":
					vMatches := intRealValRe.FindAllStringSubmatch(aBody, -1)
					for _, vm := range vMatches {
						if rawVal, err := strconv.ParseFloat(vm[1], 64); err == nil {
							animObj.DeltaX = append(animObj.DeltaX, rawVal-ln.LocalX)
						}
					}
				case "position":
					pMatches := pointValRe.FindAllStringSubmatch(aBody, -1)
					for _, pm := range pMatches {
						nums := parseNumbers(pm[1])
						if len(nums) >= 2 {
							animObj.DeltaX = append(animObj.DeltaX, nums[0]-ln.LocalX)
							animObj.DeltaY = append(animObj.DeltaY, -(nums[1] - ln.LocalY))
						}
					}
				case "transform.rotation.z":
					vMatches := intRealValRe.FindAllStringSubmatch(aBody, -1)
					for _, vm := range vMatches {
						if rad, err := strconv.ParseFloat(vm[1], 64); err == nil {
							deg := rad * (180.0 / 3.141592653589793)
							animObj.RotDeg = append(animObj.RotDeg, deg)
						}
					}
				case "opacity":
					vMatches := intRealValRe.FindAllStringSubmatch(aBody, -1)
					for _, vm := range vMatches {
						if op, err := strconv.ParseFloat(vm[1], 64); err == nil {
							animObj.Opacity = append(animObj.Opacity, op)
						}
					}
				}

				if len(animObj.DeltaX) > 0 || len(animObj.DeltaY) > 0 || len(animObj.RotDeg) > 0 || len(animObj.Opacity) > 0 {
					ln.Animations = append(ln.Animations, animObj)
				}
			}
		}

		// Find CGImage in contents and build sublayers
		for _, c := range node.Children {
			if c.XMLName.Local == "contents" {
				for _, gc := range c.Children {
					if gc.XMLName.Local == "CGImage" {
						ln.ImageSrc = getCAMLAttr(&gc, "src")
					}
				}
			}
			if c.XMLName.Local == "sublayers" {
				for i := range c.Children {
					if childLayer := buildTree(&c.Children[i], ln); childLayer != nil {
						ln.Sublayers = append(ln.Sublayers, childLayer)
					}
				}
			}
		}

		return ln
	}

	layerTree := buildTree(&root, nil)
	if layerTree == nil {
		return nil, 390, 844
	}

	canvasW := layerTree.BoundsW
	canvasH := layerTree.BoundsH
	if canvasW == 0 {
		canvasW = 390
	}
	if canvasH == 0 {
		canvasH = 844
	}

	// 4. Compute absolute screen values for each node in each state
	type ComputedNode struct {
		ScreenX  float64
		ScreenY  float64
		Width    float64
		Height   float64
		Opacity  float64
		Rotation float64
		Scale    float64
	}

	computeState := func(stateName string) map[string]ComputedNode {
		res := make(map[string]ComputedNode)
		var traverse func(l *camlLayerNode, parentScreenX, parentScreenY, parentW, parentH, parentEffOp, parentEffRot, parentEffScale float64)
		traverse = func(l *camlLayerNode, parentScreenX, parentScreenY, parentW, parentH, parentEffOp, parentEffRot, parentEffScale float64) {
			localX := l.LocalX
			localY := l.LocalY
			w := l.BoundsW
			h := l.BoundsH
			op := l.Opacity
			rot := l.Rotation
			scale := l.Scale

			if stateName != "base" {
				if ov, ok := stateOverrides[stateName][l.ID]; ok {
					if v, has := ov["position.x"]; has {
						localX = v
					}
					if v, has := ov["position.y"]; has {
						localY = v
					}
					if v, has := ov["bounds.size.width"]; has && v > 0 {
						w = v
					}
					if v, has := ov["bounds.size.height"]; has && v > 0 {
						h = v
					}
					if v, has := ov["opacity"]; has {
						op = v
					}
					if v, has := ov["transform.rotation.z"]; has {
						rot = v
					}
					if v, has := ov["transform.scale.xy"]; has && v > 0 {
						scale = v
					} else if v, has := ov["scale"]; has && v > 0 {
						scale = v
					}
				}
			}

			var screenX, screenY float64
			if l.Parent == nil {
				screenX = localX
				if l.Flipped {
					screenY = localY
				} else {
					screenY = canvasH - localY
				}
			} else {
				// Co-registered sublayer normalization:
				// When a child layer shares identical or near-identical bounds as its parent (e.g. multi-panel comic slices),
				// it is designed as a co-registered slice overlay.
				// In rest/idle states ("locked", "base", "sleep"), creators intend it to be centered relative to parent,
				// but visual drag keyframing in CAPlayground often leaves an accidental horizontal/vertical offset.
				// If the offset is within rest range (< 35% of parent size), snap to parent center.
				// In "unlock" state, large animated displacements (flying left/right) are preserved.
				if math.Abs(w-parentW) < 4.0 && math.Abs(h-parentH) < 4.0 && (stateName == "locked" || stateName == "base" || stateName == "sleep") {
					if math.Abs(localX-parentW/2) < parentW*0.35 {
						localX = parentW / 2
					}
					if math.Abs(localY-parentH/2) < parentH*0.35 {
						localY = parentH / 2
					}
				}

				// Child relative to parent's center, scaled by parent's effective scale
				dx := localX - parentW/2
				dy := localY - parentH/2
				screenX = parentScreenX + dx*parentEffScale
				if l.Flipped {
					screenY = parentScreenY + dy*parentEffScale
				} else {
					screenY = parentScreenY - dy*parentEffScale
				}
			}

			effOp := parentEffOp * op
			effRot := parentEffRot + rot
			effScale := parentEffScale * scale

			res[l.ID] = ComputedNode{
				ScreenX:  screenX,
				ScreenY:  screenY,
				Width:    w,
				Height:   h,
				Opacity:  effOp,
				Rotation: effRot,
				Scale:    effScale,
			}

			for _, sub := range l.Sublayers {
				traverse(sub, screenX, screenY, w, h, effOp, effRot, effScale)
			}
		}

		traverse(layerTree, canvasW/2, canvasH/2, canvasW, canvasH, 1.0, 0, 1.0)
		return res
	}

	baseComputed := computeState("base")
	lockedComputed := computeState("locked")
	unlockComputed := computeState("unlock")
	sleepComputed := computeState("sleep")

	// 5. Flatten all layers that have image content
	var resultLayers []TendiesLayer
	order := 0

	var collectVisualLayers func(l *camlLayerNode)
	collectVisualLayers = func(l *camlLayerNode) {
		if l.ImageSrc != "" {
			srcBase := filepath.Base(l.ImageSrc)
			imgURL, exists := assetURLMap[strings.ToLower(srcBase)]
			if !exists {
				imgURL, exists = assetURLMap[strings.ToLower(l.Name)]
			}
			if !exists {
				imgURL, exists = assetURLMap[normalizeAssetKey(srcBase)]
			}
			if !exists {
				imgURL, exists = assetURLMap[normalizeAssetKey(l.Name)]
			}
			if !exists {
				ext := filepath.Ext(srcBase)
				if ext != "" {
					imgURL, _ = assetURLMap[normalizeAssetKey(strings.TrimSuffix(srcBase, ext))]
				}
			}

			baseVals := baseComputed[l.ID]

			tl := TendiesLayer{
				ID:         l.ID,
				Name:       srcBase,
				URL:        imgURL,
				Category:   category,
				Order:      order,
				PosX:       baseVals.ScreenX,
				PosY:       baseVals.ScreenY,
				Width:      baseVals.Width,
				Height:     baseVals.Height,
				Opacity:    baseVals.Opacity,
				Rotation:   baseVals.Rotation,
				Scale:      baseVals.Scale,
				States:     make(map[string]LayerState),
				Animations: l.Animations,
			}

			for sName, sComp := range map[string]map[string]ComputedNode{
				"locked": lockedComputed,
				"unlock": unlockComputed,
				"sleep":  sleepComputed,
			} {
				if vals, ok := sComp[l.ID]; ok {
					tl.States[sName] = LayerState{
						PosX:     vals.ScreenX,
						PosY:     vals.ScreenY,
						Width:    vals.Width,
						Height:   vals.Height,
						Opacity:  vals.Opacity,
						Rotation: vals.Rotation,
						Scale:    vals.Scale,
					}
				}
			}

			resultLayers = append(resultLayers, tl)
			order++
		}

		for _, sub := range l.Sublayers {
			collectVisualLayers(sub)
		}
	}

	collectVisualLayers(layerTree)
	return resultLayers, canvasW, canvasH
}

// ParseTendiesArchive extracts all image files from the .tendies zip archive into a dedicated directory
// and parses the CAML hierarchy for Background.ca, Floating.ca, and Foreground.ca.
func ParseTendiesArchive(zipPath string, uploadDir string, cleanID string) (*TendiesParseResult, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tendies archive: %w", err)
	}
	defer r.Close()

	// Target extracted directory
	subDirName := fmt.Sprintf("extracted_%s", cleanID)
	extractDst := filepath.Join(uploadDir, subDirName)
	if err := os.MkdirAll(extractDst, 0755); err != nil {
		return nil, fmt.Errorf("failed to create extracted directory: %w", err)
	}

	assetURLMap := make(map[string]string)
	camlContents := make(map[string]string)
	extractedFiles := make(map[string]string) // baseName -> fullURL

	// 1. Extract all visual image assets and read CAML files
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		nameLower := strings.ToLower(f.Name)
		baseName := filepath.Base(nameLower)

		// Check for image assets
		if strings.HasSuffix(nameLower, ".png") ||
			strings.HasSuffix(nameLower, ".jpg") ||
			strings.HasSuffix(nameLower, ".jpeg") ||
			strings.HasSuffix(nameLower, ".webp") {

			// Clean target filename
			cleanBase := filepath.Base(f.Name)
			cleanBase = strings.ReplaceAll(cleanBase, " ", "_")

			// Obfuscated filename: layer_<md5_hash>.dat (prevents asset names like wutian.png from leaking)
			hasher := md5.New()
			hasher.Write([]byte(cleanBase + "_" + cleanID))
			obfuscatedName := fmt.Sprintf("layer_%x.dat", hasher.Sum(nil)[:8])
			targetFilePath := filepath.Join(extractDst, obfuscatedName)

			// Extract image file to disk with encryption
			rc, err := f.Open()
			if err == nil {
				rawBytes, err := io.ReadAll(rc)
				rc.Close()
				if err == nil {
					salt := byte(0x3C + (len(cleanBase) % 100))
					encBytes := EncryptTendiesBytes(rawBytes, salt)
					_ = os.WriteFile(targetFilePath, encBytes, 0644)
				}
			}

			imgURL := fmt.Sprintf("/uploads/interactive/%s/%s", subDirName, obfuscatedName)
			assetURLMap[baseName] = imgURL
			assetURLMap[strings.ToLower(cleanBase)] = imgURL
			assetURLMap[normalizeAssetKey(baseName)] = imgURL
			assetURLMap[normalizeAssetKey(cleanBase)] = imgURL
			ext := filepath.Ext(baseName)
			if ext != "" {
				assetURLMap[normalizeAssetKey(strings.TrimSuffix(baseName, ext))] = imgURL
			}
			extractedFiles[baseName] = imgURL
		}

		// Check for CAML configuration
		if filepath.Base(nameLower) == "main.caml" {
			rc, err := f.Open()
			if err == nil {
				data, _ := io.ReadAll(rc)
				rc.Close()
				camlContents[f.Name] = string(data)
			}
		}
	}

	result := &TendiesParseResult{
		CanvasWidth:     390,
		CanvasHeight:    844,
		BGLayers:        make([]TendiesLayer, 0),
		FloatingLayers:  make([]TendiesLayer, 0),
		FGLayers:        make([]TendiesLayer, 0),
		TotalImageCount: len(extractedFiles),
	}

	claimedURLs := make(map[string]bool)

	// 2. Parse CAML layers according to directory structure
	for path, content := range camlContents {
		pathLower := strings.ToLower(path)
		if strings.Contains(pathLower, "background") {
			layers, cw, ch := parseCAMLLayers(content, "background", assetURLMap)
			result.BGLayers = layers
			if cw > 0 && ch > 0 {
				result.CanvasWidth = cw
				result.CanvasHeight = ch
			}
			for _, l := range layers {
				if l.URL != "" {
					claimedURLs[l.URL] = true
				}
			}
		} else if strings.Contains(pathLower, "floating") {
			layers, cw, ch := parseCAMLLayers(content, "floating", assetURLMap)
			result.FloatingLayers = layers
			if (result.CanvasWidth == 390 || result.CanvasWidth == 0) && cw > 0 && ch > 0 {
				result.CanvasWidth = cw
				result.CanvasHeight = ch
			}
			for _, l := range layers {
				if l.URL != "" {
					claimedURLs[l.URL] = true
				}
			}
		} else if strings.Contains(pathLower, "foreground") {
			layers, cw, ch := parseCAMLLayers(content, "foreground", assetURLMap)
			result.FGLayers = layers
			if (result.CanvasWidth == 390 || result.CanvasWidth == 0) && cw > 0 && ch > 0 {
				result.CanvasWidth = cw
				result.CanvasHeight = ch
			}
			for _, l := range layers {
				if l.URL != "" {
					claimedURLs[l.URL] = true
				}
			}
		}
	}

	// 2.5 Ensure background plate separation:
	// In some CAML packages (e.g. created via CAPlayground), all layers including the full-canvas
	// background plate (e.g. Backa.JPG, Background.JPG) were placed inside Floating.ca.
	// Any layer in FloatingLayers that has full canvas coverage (width & height >= 85% of canvas)
	// or has "back" or "bg" in its name MUST be re-routed to BGLayers so it is positioned
	// strictly BEHIND the lockscreen clock, rather than covering the clock!
	var trueFloating []TendiesLayer
	for _, l := range result.FloatingLayers {
		isBackplate := (l.Width >= result.CanvasWidth*0.85 && l.Height >= result.CanvasHeight*0.85) ||
			strings.Contains(strings.ToLower(l.Name), "background") ||
			strings.Contains(strings.ToLower(l.Name), "back") ||
			strings.Contains(strings.ToLower(l.Name), "bg")
		if isBackplate {
			l.Category = "background"
			result.BGLayers = append(result.BGLayers, l)
		} else {
			trueFloating = append(trueFloating, l)
		}
	}
	result.FloatingLayers = trueFloating

	// 3. Guarantee ZERO image loss:
	// If any extracted image was not claimed by CAML, preserve it in BGLayers or FloatingLayers
	unclaimedIdx := 100
	for base, url := range extractedFiles {
		if !claimedURLs[url] {
			isBG := strings.Contains(base, "bg") || strings.Contains(base, "background")
			layer := TendiesLayer{
				ID:       fmt.Sprintf("unclaimed_%d", unclaimedIdx),
				Name:     base,
				URL:      url,
				Category: "floating",
				Order:    unclaimedIdx,
				PosX:     result.CanvasWidth / 2,
				PosY:     result.CanvasHeight / 2,
				Width:    result.CanvasWidth,
				Height:   result.CanvasHeight,
				Opacity:  1.0,
				Rotation: 0,
				States:   make(map[string]LayerState),
			}
			if isBG {
				layer.Category = "background"
				result.BGLayers = append(result.BGLayers, layer)
			} else {
				result.FloatingLayers = append(result.FloatingLayers, layer)
			}
			unclaimedIdx++
		}
	}

	// 4. Set fallback primary background & foreground URLs for backwards compatibility
	for _, l := range result.BGLayers {
		if l.Opacity > 0 && l.URL != "" {
			result.BackgroundURL = l.URL
			break
		}
	}
	if result.BackgroundURL == "" && len(result.BGLayers) > 0 {
		result.BackgroundURL = result.BGLayers[0].URL
	}

	for _, l := range result.FloatingLayers {
		if l.Opacity > 0 && l.URL != "" {
			result.ForegroundURL = l.URL
			break
		}
	}
	if result.ForegroundURL == "" && len(result.FloatingLayers) > 0 {
		result.ForegroundURL = result.FloatingLayers[0].URL
	}

	// 5. Automatically composite authentic lock screen cover image
	coverURL := GenerateLockScreenCover(result, uploadDir, cleanID)
	if coverURL != "" {
		result.CoverURL = coverURL
	} else {
		result.CoverURL = result.BackgroundURL
	}

	return result, nil
}

func loadAndDecodeImage(path string) image.Image {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	dec, _ := DecryptTendiesBytes(data)
	img, _, err := image.Decode(bytes.NewReader(dec))
	if err != nil {
		return nil
	}
	return img
}

func scaleRGBA(src image.Image, targetW, targetH int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	if srcW <= 0 || srcH <= 0 || targetW <= 0 || targetH <= 0 {
		return dst
	}

	for y := 0; y < targetH; y++ {
		srcY := srcBounds.Min.Y + (y*srcH)/targetH
		for x := 0; x < targetW; x++ {
			srcX := srcBounds.Min.X + (x*srcW)/targetW
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func drawScaledLayer(canvas *image.RGBA, src image.Image, dstX, dstY, dstW, dstH int, opacity float64) {
	if dstW <= 0 || dstH <= 0 || opacity <= 0 {
		return
	}
	scaled := scaleRGBA(src, dstW, dstH)
	dstRect := image.Rect(dstX, dstY, dstX+dstW, dstY+dstH)

	if opacity >= 0.999 {
		draw.Draw(canvas, dstRect, scaled, image.Point{}, draw.Over)
	} else {
		mask := image.NewUniform(color.Alpha{A: uint8(opacity * 255)})
		draw.DrawMask(canvas, dstRect, scaled, image.Point{}, mask, image.Point{}, draw.Over)
	}
}

// GenerateLockScreenCover composites all layers visible in the "locked" state onto a single canvas
// and writes the result as an optimized JPEG cover image in uploadDir.
func GenerateLockScreenCover(result *TendiesParseResult, uploadDir, cleanID string) string {
	if result == nil || result.CanvasWidth <= 0 || result.CanvasHeight <= 0 {
		return ""
	}

	canvasW := int(result.CanvasWidth)
	canvasH := int(result.CanvasHeight)
	canvas := image.NewRGBA(image.Rect(0, 0, canvasW, canvasH))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	// Combine layers strictly preserving category hierarchy: Background -> Floating -> Foreground
	allLayers := make([]TendiesLayer, 0, len(result.BGLayers)+len(result.FloatingLayers)+len(result.FGLayers))
	bgCopy := append([]TendiesLayer{}, result.BGLayers...)
	sort.Slice(bgCopy, func(i, j int) bool { return bgCopy[i].Order < bgCopy[j].Order })
	flCopy := append([]TendiesLayer{}, result.FloatingLayers...)
	sort.Slice(flCopy, func(i, j int) bool { return flCopy[i].Order < flCopy[j].Order })
	fgCopy := append([]TendiesLayer{}, result.FGLayers...)
	sort.Slice(fgCopy, func(i, j int) bool { return fgCopy[i].Order < fgCopy[j].Order })

	allLayers = append(allLayers, bgCopy...)
	allLayers = append(allLayers, flCopy...)
	allLayers = append(allLayers, fgCopy...)

	for _, layer := range allLayers {
		if layer.URL == "" {
			continue
		}

		lockState, hasLock := layer.States["locked"]

		// 1. Determine Opacity
		op := layer.Opacity
		if hasLock {
			if lockState.Opacity == 0 {
				continue // explicitly hidden in locked state
			}
			op = lockState.Opacity
		}
		if op < 0.02 {
			continue
		}

		// 2. Determine bounds, position, and scale
		posX := layer.PosX
		posY := layer.PosY
		w := layer.Width
		h := layer.Height
		scale := layer.Scale
		if scale <= 0 {
			scale = 1.0
		}

		if hasLock {
			if lockState.PosX != 0 {
				posX = lockState.PosX
			}
			if lockState.PosY != 0 {
				posY = lockState.PosY
			}
			if lockState.Width > 0 {
				w = lockState.Width
			}
			if lockState.Height > 0 {
				h = lockState.Height
			}
			if lockState.Scale > 0 {
				scale = lockState.Scale
			}
		}

		renderW := w * scale
		renderH := h * scale

		// 3. Skip if completely offscreen in locked state
		if posX+renderW/2 < 0 || posX-renderW/2 > float64(canvasW) || posY+renderH/2 < 0 || posY-renderH/2 > float64(canvasH) {
			continue
		}

		relPath := strings.TrimPrefix(layer.URL, "/uploads/interactive/")
		layerDiskPath := filepath.Join(uploadDir, relPath)
		img := loadAndDecodeImage(layerDiskPath)
		if img == nil {
			continue
		}

		dstX := int(math.Round(posX - renderW/2))
		dstY := int(math.Round(posY - renderH/2))
		dstW := int(math.Round(renderW))
		dstH := int(math.Round(renderH))

		drawScaledLayer(canvas, img, dstX, dstY, dstW, dstH, op)
	}

	// Save composited cover image
	coverFilename := fmt.Sprintf("cover_lock_%s.jpg", cleanID)
	coverDiskPath := filepath.Join(uploadDir, coverFilename)
	outFile, err := os.Create(coverDiskPath)
	if err != nil {
		return ""
	}
	defer outFile.Close()

	if err := jpeg.Encode(outFile, canvas, &jpeg.Options{Quality: 92}); err != nil {
		return ""
	}

	coverURL := fmt.Sprintf("/uploads/interactive/%s", coverFilename)
	result.CoverURL = coverURL
	return coverURL
}
