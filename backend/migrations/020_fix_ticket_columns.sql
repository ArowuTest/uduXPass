-- Migration 020: Fix ticket columns for large data storage
-- 1. Change qr_code_image_url from varchar(500) to TEXT (base64 QR images are ~6-8KB)
-- 2. Change qr_code_data from varchar(500) to TEXT (JWT tokens can exceed 500 chars)
-- 3. Fix indexes: btree indexes have a 2704-byte limit, use hash for large text columns

ALTER TABLE tickets ALTER COLUMN qr_code_image_url TYPE TEXT;
ALTER TABLE tickets ALTER COLUMN qr_code_data TYPE TEXT;

-- Drop the btree index on qr_code_image_url (stores large base64 data)
DROP INDEX IF EXISTS idx_tickets_qr_image_url;

-- Replace btree index on qr_code_data with hash index (supports unlimited key sizes)
DROP INDEX IF EXISTS idx_tickets_qr_code_data;
CREATE INDEX idx_tickets_qr_code_data ON tickets USING hash (qr_code_data);

-- Replace btree unique constraint on qr_code_data with md5-based functional unique index
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_qr_code_data_key;
CREATE UNIQUE INDEX IF NOT EXISTS tickets_qr_code_data_unique ON tickets (md5(qr_code_data));
