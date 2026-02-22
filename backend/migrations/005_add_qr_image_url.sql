-- Migration: Add QR code image URL to tickets table
-- This allows storing pre-generated QR code images for better performance

-- Add qr_code_image_url column to tickets table
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS qr_code_image_url VARCHAR(500);

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_tickets_qr_image_url ON tickets(qr_code_image_url);

-- Add comment
COMMENT ON COLUMN tickets.qr_code_image_url IS 'URL to pre-generated QR code image (optional, falls back to client-side generation)';
