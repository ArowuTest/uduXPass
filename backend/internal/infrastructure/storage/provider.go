// Package storage provides a pluggable file storage abstraction.
// Currently ships with LocalStorage (files on disk, served as static assets).
// To switch to GCP Cloud Storage in production, implement the StorageProvider
// interface in a new gcs.go file and set STORAGE_PROVIDER=gcs in your env.
package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StorageProvider defines the contract for file storage backends.
// Implementations: LocalStorage (default), GCSStorage (production).
type StorageProvider interface {
	// Upload saves the file and returns its public URL.
	Upload(file multipart.File, header *multipart.FileHeader, folder string) (string, error)
	// Delete removes a file by its public URL.
	Delete(publicURL string) error
}

// AllowedImageTypes lists MIME types accepted for event images.
var AllowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// AllowedVideoTypes lists MIME types accepted for promo videos.
var AllowedVideoTypes = map[string]string{
	"video/mp4":       ".mp4",
	"video/webm":      ".webm",
	"video/quicktime": ".mov",
}

const (
	MaxImageSize = 10 * 1024 * 1024  // 10 MB
	MaxVideoSize = 100 * 1024 * 1024 // 100 MB
)

// ValidateFile checks MIME type and size constraints.
func ValidateFile(header *multipart.FileHeader, allowedTypes map[string]string, maxSize int64) error {
	if header.Size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum of %d bytes", header.Size, maxSize)
	}
	contentType := header.Header.Get("Content-Type")
	// Normalise: strip charset suffix if present
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = strings.TrimSpace(contentType[:idx])
	}
	if _, ok := allowedTypes[contentType]; !ok {
		allowed := make([]string, 0, len(allowedTypes))
		for k := range allowedTypes {
			allowed = append(allowed, k)
		}
		return fmt.Errorf("content type %q is not allowed; accepted: %s", contentType, strings.Join(allowed, ", "))
	}
	return nil
}

// generateFilename creates a collision-resistant filename preserving the extension.
func generateFilename(header *multipart.FileHeader) string {
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		// Fall back to MIME-derived extension
		ct := header.Header.Get("Content-Type")
		if idx := strings.Index(ct, ";"); idx != -1 {
			ct = strings.TrimSpace(ct[:idx])
		}
		if e, ok := AllowedImageTypes[ct]; ok {
			ext = e
		} else if e, ok := AllowedVideoTypes[ct]; ok {
			ext = e
		} else {
			ext = ".bin"
		}
	}
	return fmt.Sprintf("%d-%s%s", time.Now().UnixMilli(), uuid.New().String()[:8], ext)
}

// ─── LocalStorage ─────────────────────────────────────────────────────────────

// LocalStorage stores files on the local filesystem and serves them via a
// static file handler mounted at /uploads/.
// In production, replace with GCSStorage by setting STORAGE_PROVIDER=gcs.
type LocalStorage struct {
	// BaseDir is the absolute path to the directory where uploads are stored.
	BaseDir string
	// BaseURL is the public URL prefix used to construct file URLs,
	// e.g. "http://localhost:3000/uploads".
	BaseURL string
}

// NewLocalStorage creates a LocalStorage instance, creating BaseDir if needed.
func NewLocalStorage(baseDir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("storage: failed to create upload directory %s: %w", baseDir, err)
	}
	return &LocalStorage{BaseDir: baseDir, BaseURL: strings.TrimRight(baseURL, "/")}, nil
}

// Upload saves the multipart file to disk and returns its public URL.
func (s *LocalStorage) Upload(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	dir := filepath.Join(s.BaseDir, folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("storage: failed to create folder %s: %w", folder, err)
	}

	filename := generateFilename(header)
	destPath := filepath.Join(dir, filename)

	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("storage: failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("storage: failed to write file: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s/%s", s.BaseURL, folder, filename)
	return publicURL, nil
}

// Delete removes the file identified by its public URL.
func (s *LocalStorage) Delete(publicURL string) error {
	// Strip base URL prefix to get relative path
	rel := strings.TrimPrefix(publicURL, s.BaseURL+"/")
	if rel == publicURL {
		return fmt.Errorf("storage: URL %q does not belong to this storage provider", publicURL)
	}
	fullPath := filepath.Join(s.BaseDir, rel)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage: failed to delete %s: %w", fullPath, err)
	}
	return nil
}
