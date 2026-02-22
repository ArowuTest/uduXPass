package qrcode

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

// Generator handles QR code generation
type Generator struct {
	size          int
	recoveryLevel qrcode.RecoveryLevel
}

// NewGenerator creates a new QR code generator with default settings
func NewGenerator() *Generator {
	return &Generator{
		size:          256,
		recoveryLevel: qrcode.High, // High error correction for better scanning
	}
}

// NewGeneratorWithOptions creates a new QR code generator with custom settings
func NewGeneratorWithOptions(size int, recoveryLevel qrcode.RecoveryLevel) *Generator {
	return &Generator{
		size:          size,
		recoveryLevel: recoveryLevel,
	}
}

// GenerateQRCode generates a QR code as PNG bytes
func (g *Generator) GenerateQRCode(data string) ([]byte, error) {
	if data == "" {
		return nil, fmt.Errorf("qr code data cannot be empty")
	}

	qr, err := qrcode.New(data, g.recoveryLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create qr code: %w", err)
	}

	png, err := qr.PNG(g.size)
	if err != nil {
		return nil, fmt.Errorf("failed to generate png: %w", err)
	}

	return png, nil
}

// GenerateQRCodeBase64 generates a QR code as base64 encoded string
func (g *Generator) GenerateQRCodeBase64(data string) (string, error) {
	png, err := g.GenerateQRCode(data)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(png)
	return fmt.Sprintf("data:image/png;base64,%s", encoded), nil
}

// GenerateQRCodeFile generates a QR code and saves it to a file
func (g *Generator) GenerateQRCodeFile(data, filepath string) error {
	if data == "" {
		return fmt.Errorf("qr code data cannot be empty")
	}

	if filepath == "" {
		return fmt.Errorf("filepath cannot be empty")
	}

	err := qrcode.WriteFile(data, g.recoveryLevel, g.size, filepath)
	if err != nil {
		return fmt.Errorf("failed to write qr code to file: %w", err)
	}

	return nil
}

// SetSize sets the QR code size
func (g *Generator) SetSize(size int) {
	if size > 0 {
		g.size = size
	}
}

// SetRecoveryLevel sets the error correction level
func (g *Generator) SetRecoveryLevel(level qrcode.RecoveryLevel) {
	g.recoveryLevel = level
}
