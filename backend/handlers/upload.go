package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"wallpaper-shop/services"
)

// UploadImage handles local image file upload (covers, screenshots)
func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Limit to 15MB
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "No file uploaded or invalid form data"})
			return
		}

		if file.Size > 15*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "File exceeds maximum size limit (15MB)"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
			".gif":  true,
			".svg":  true,
		}

		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file type. Allowed: JPG, PNG, WEBP, GIF, SVG"})
			return
		}

		// Ensure upload directory exists
		uploadDir := filepath.Join(".", "uploads", "covers")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create upload storage directory"})
			return
		}

		// Generate clean filename
		cleanID := uuid.New().String()[:8]
		filename := fmt.Sprintf("cover_%d_%s%s", time.Now().Unix(), cleanID, ext)
		dstPath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, dstPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file"})
			return
		}

		// Return accessible URL path
		urlPath := fmt.Sprintf("/uploads/covers/%s", filename)
		c.JSON(http.StatusOK, gin.H{
			"message":  "Upload successful",
			"url":      urlPath,
			"filename": filename,
			"size":     file.Size,
		})
	}
}

// extractTendiesLayers unzips and locates PosterBoard Floating cutout & Background layers
func extractTendiesLayers(zipPath string, uploadDir string, cleanID string) (bgURL, fgURL string) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", ""
	}
	defer r.Close()

	var bgEntry *zip.File
	var fgEntry *zip.File

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		nameLower := strings.ToLower(f.Name)

		// Check for foreground floating layer (transparent cutout)
		if fgEntry == nil {
			if strings.Contains(nameLower, "floating") && strings.HasSuffix(nameLower, ".png") {
				fgEntry = f
			} else if strings.Contains(nameLower, "cutout") || strings.Contains(nameLower, "fg") {
				fgEntry = f
			}
		}

		// Check for background layer
		if bgEntry == nil {
			if strings.Contains(nameLower, "background") && (strings.HasSuffix(nameLower, ".jpg") || strings.HasSuffix(nameLower, ".jpeg") || strings.HasSuffix(nameLower, ".png")) {
				bgEntry = f
			} else if strings.HasSuffix(nameLower, "0.jpg") || strings.HasSuffix(nameLower, "0.png") {
				bgEntry = f
			}
		}
	}

	// Fallback if specific keywords not found: pick prominent png for fg and jpg for bg
	if fgEntry == nil {
		for _, f := range r.File {
			if !f.FileInfo().IsDir() && strings.HasSuffix(strings.ToLower(f.Name), ".png") && f.UncompressedSize64 > 50*1024 {
				fgEntry = f
				break
			}
		}
	}
	if bgEntry == nil {
		for _, f := range r.File {
			if !f.FileInfo().IsDir() && (strings.HasSuffix(strings.ToLower(f.Name), ".jpg") || strings.HasSuffix(strings.ToLower(f.Name), ".jpeg")) && f.UncompressedSize64 > 50*1024 {
				bgEntry = f
				break
			}
		}
	}

	// Extract foreground
	if fgEntry != nil {
		rc, err := fgEntry.Open()
		if err == nil {
			fgName := fmt.Sprintf("extracted_%s_fg.png", cleanID)
			fgPath := filepath.Join(uploadDir, fgName)
			out, err := os.Create(fgPath)
			if err == nil {
				io.Copy(out, rc)
				out.Close()
				fgURL = fmt.Sprintf("/uploads/interactive/%s", fgName)
			}
			rc.Close()
		}
	}

	// Extract background
	if bgEntry != nil {
		rc, err := bgEntry.Open()
		if err == nil {
			ext := filepath.Ext(bgEntry.Name)
			if ext == "" {
				ext = ".jpg"
			}
			bgName := fmt.Sprintf("extracted_%s_bg%s", cleanID, ext)
			bgPath := filepath.Join(uploadDir, bgName)
			out, err := os.Create(bgPath)
			if err == nil {
				io.Copy(out, rc)
				out.Close()
				bgURL = fmt.Sprintf("/uploads/interactive/%s", bgName)
			}
			rc.Close()
		}
	}

	return bgURL, fgURL
}

// UploadInteractiveFile handles .tendies archives, layers (PNG/WebP), and motion videos (MP4/WebM)
func UploadInteractiveFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Limit to 100MB for tendies packages and videos
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "No file uploaded or invalid form data"})
			return
		}

		if file.Size > 100*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "File exceeds maximum size limit (100MB)"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{
			".tendies": true,
			".zip":     true,
			".png":     true,
			".webp":    true,
			".jpg":     true,
			".jpeg":    true,
			".mp4":     true,
			".webm":    true,
			".mov":     true,
		}

		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid file type. Allowed: .tendies, .zip, PNG, WebP, JPG, MP4, WebM, MOV"})
			return
		}

		// Ensure upload directory exists
		uploadDir := filepath.Join(".", "uploads", "interactive")
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create interactive storage directory"})
			return
		}

		// Generate clean filename
		cleanID := uuid.New().String()[:8]
		filename := fmt.Sprintf("interactive_%d_%s%s", time.Now().Unix(), cleanID, ext)
		dstPath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, dstPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save file"})
			return
		}

		var tendiesData *services.TendiesParseResult
		var bgURL, fgURL, coverURL string
		if ext == ".tendies" || ext == ".zip" {
			parsed, err := services.ParseTendiesArchive(dstPath, uploadDir, cleanID)
			if err == nil && parsed != nil {
				tendiesData = parsed
				bgURL = parsed.BackgroundURL
				fgURL = parsed.ForegroundURL
				coverURL = parsed.CoverURL
			} else {
				bgURL, fgURL = extractTendiesLayers(dstPath, uploadDir, cleanID)
			}
		}

		if coverURL == "" && bgURL != "" {
			coverURL = bgURL
		}

		// Return accessible URL path
		urlPath := fmt.Sprintf("/uploads/interactive/%s", filename)
		resp := gin.H{
			"message":        "Interactive file upload successful",
			"url":            urlPath,
			"filename":       filename,
			"ext":            ext,
			"size":           file.Size,
			"background_url": bgURL,
			"foreground_url": fgURL,
			"cover_url":      coverURL,
		}
		if tendiesData != nil {
			resp["tendies_data"] = tendiesData
		}
		c.JSON(http.StatusOK, resp)
	}
}

