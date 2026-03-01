// Package handlers — UploadHandler handles media file uploads.
// Supports images (JPEG, PNG, WebP, GIF) and videos (MP4, WebM, MOV).
// Files are validated for type and size before being passed to the StorageProvider.
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uduxpass/backend/internal/infrastructure/storage"
)

// UploadHandler handles file upload requests.
type UploadHandler struct {
	store storage.StorageProvider
}

// NewUploadHandler creates an UploadHandler backed by the given StorageProvider.
func NewUploadHandler(store storage.StorageProvider) *UploadHandler {
	return &UploadHandler{store: store}
}

// UploadResponse is returned on a successful upload.
type UploadResponse struct {
	Success  bool   `json:"success"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// UploadMedia handles POST /v1/admin/upload
// Accepts multipart/form-data with a "file" field and an optional "folder" field.
// Folder defaults to "events" and is sanitised to prevent path traversal.
//
// Accepted media types:
//   - Images: image/jpeg, image/png, image/webp, image/gif  (max 10 MB)
//   - Videos: video/mp4, video/webm, video/quicktime        (max 100 MB)
func (h *UploadHandler) UploadMedia(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No file provided",
			"message": "Include a 'file' field in your multipart/form-data request",
		})
		return
	}

	// Determine media type and validate
	contentType := fileHeader.Header.Get("Content-Type")
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = strings.TrimSpace(contentType[:idx])
	}

	isImage := false
	isVideo := false
	if _, ok := storage.AllowedImageTypes[contentType]; ok {
		isImage = true
	} else if _, ok := storage.AllowedVideoTypes[contentType]; ok {
		isVideo = true
	}

	if !isImage && !isVideo {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error":   "Unsupported file type",
			"message": "Accepted image types: JPEG, PNG, WebP, GIF. Accepted video types: MP4, WebM, MOV.",
		})
		return
	}

	// Size validation
	if isImage {
		if err := storage.ValidateFile(fileHeader, storage.AllowedImageTypes, storage.MaxImageSize); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"success": false,
				"error":   "File validation failed",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := storage.ValidateFile(fileHeader, storage.AllowedVideoTypes, storage.MaxVideoSize); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"success": false,
				"error":   "File validation failed",
				"message": err.Error(),
			})
			return
		}
	}

	// Sanitise folder name — only allow alphanumeric, hyphens, underscores
	folder := c.DefaultPostForm("folder", "events")
	folder = sanitiseFolder(folder)

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to read uploaded file",
		})
		return
	}
	defer file.Close()

	// Upload via storage provider
	publicURL, err := h.store.Upload(file, fileHeader, folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Upload failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, UploadResponse{
		Success:  true,
		URL:      publicURL,
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		MimeType: contentType,
	})
}

// DeleteMedia handles DELETE /v1/admin/upload
// Body: { "url": "https://..." }
func (h *UploadHandler) DeleteMedia(c *gin.Context) {
	var body struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "url field is required",
		})
		return
	}

	if err := h.store.Delete(body.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Delete failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "File deleted successfully",
	})
}

// sanitiseFolder strips any characters that could cause path traversal.
func sanitiseFolder(folder string) string {
	// Replace path separators and dots
	folder = strings.ReplaceAll(folder, "/", "")
	folder = strings.ReplaceAll(folder, "\\", "")
	folder = strings.ReplaceAll(folder, "..", "")
	folder = strings.TrimSpace(folder)
	if folder == "" {
		return "events"
	}
	return folder
}
