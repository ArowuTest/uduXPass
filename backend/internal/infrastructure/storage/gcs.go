// Package storage — GCP Cloud Storage implementation.
//
// To activate in production:
//  1. Add the GCS client library: go get cloud.google.com/go/storage
//  2. Set environment variables:
//       STORAGE_PROVIDER=gcs
//       GCS_BUCKET=your-bucket-name
//       GCS_BASE_URL=https://storage.googleapis.com/your-bucket-name
//       GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
//  3. Uncomment the implementation below and remove the stub.
//
// The StorageProvider interface is identical for both LocalStorage and GCSStorage,
// so no other code changes are required when switching providers.

package storage

import (
	"fmt"
	"mime/multipart"
)

// GCSStorage is a placeholder for the GCP Cloud Storage implementation.
// Uncomment and implement when deploying to production on GCP.
type GCSStorage struct {
	BucketName string
	BaseURL    string
}

// NewGCSStorage creates a GCSStorage instance.
// Returns an error until the full implementation is activated.
func NewGCSStorage(bucketName, baseURL string) (*GCSStorage, error) {
	return nil, fmt.Errorf("GCSStorage is not yet activated — set STORAGE_PROVIDER=local or implement GCSStorage")
}

func (s *GCSStorage) Upload(_ multipart.File, _ *multipart.FileHeader, _ string) (string, error) {
	return "", fmt.Errorf("GCSStorage.Upload not implemented")
}

func (s *GCSStorage) Delete(_ string) error {
	return fmt.Errorf("GCSStorage.Delete not implemented")
}
