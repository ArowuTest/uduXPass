-- Migration: Align orders table with backend Order entity
-- Description: Adds missing columns to orders table to match backend expectations
-- Version: 007
-- Date: 2026-02-17

-- Add missing columns to orders table
ALTER TABLE orders 
ADD COLUMN IF NOT EXISTS code VARCHAR(50) UNIQUE,
ADD COLUMN IF NOT EXISTS email VARCHAR(255),
ADD COLUMN IF NOT EXISTS phone VARCHAR(50),
ADD COLUMN IF NOT EXISTS first_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS last_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS customer_first_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS customer_last_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS customer_email VARCHAR(255),
ADD COLUMN IF NOT EXISTS customer_phone VARCHAR(50),
ADD COLUMN IF NOT EXISTS payment_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS notes TEXT,
ADD COLUMN IF NOT EXISTS confirmed_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS secret VARCHAR(255),
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_orders_code ON orders(code);
CREATE INDEX IF NOT EXISTS idx_orders_email ON orders(email);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_payment_reference ON orders(payment_reference);
CREATE INDEX IF NOT EXISTS idx_orders_expires_at ON orders(expires_at);
CREATE INDEX IF NOT EXISTS idx_orders_is_active ON orders(is_active);

-- Update existing orders to have codes and expiration dates
UPDATE orders 
SET 
    code = UPPER(SUBSTRING(MD5(RANDOM()::TEXT) FROM 1 FOR 8)),
    expires_at = created_at + INTERVAL '30 minutes',
    is_active = true
WHERE code IS NULL;

-- Add comment to new columns
COMMENT ON COLUMN orders.code IS 'Unique order code for customer reference';
COMMENT ON COLUMN orders.secret IS 'Secret token for order verification and access';
COMMENT ON COLUMN orders.expires_at IS 'Timestamp when unpaid order expires';
COMMENT ON COLUMN orders.customer_first_name IS 'First name of the customer (may differ from user)';
COMMENT ON COLUMN orders.customer_last_name IS 'Last name of the customer (may differ from user)';
COMMENT ON COLUMN orders.customer_email IS 'Email of the customer (may differ from user)';
COMMENT ON COLUMN orders.customer_phone IS 'Phone of the customer (may differ from user)';
