-- Migration 018: Align tickets table with repository schema
-- The tickets table was designed with only order_line_id as FK,
-- but the repository requires direct order_id, ticket_tier_id, code, 
-- attendee info, and meta_info columns for performance and correctness.
-- This migration adds all missing columns and backfills from order_lines.

BEGIN;

-- Add missing columns to tickets table
ALTER TABLE tickets
    ADD COLUMN IF NOT EXISTS order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS ticket_tier_id UUID REFERENCES ticket_tiers(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS code VARCHAR(50) UNIQUE,
    ADD COLUMN IF NOT EXISTS qr_code VARCHAR(500),
    ADD COLUMN IF NOT EXISTS attendee_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS attendee_email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS attendee_phone VARCHAR(50),
    ADD COLUMN IF NOT EXISTS meta_info JSONB;

-- Backfill order_id and ticket_tier_id from order_lines for existing tickets
UPDATE tickets t
SET 
    order_id = ol.order_id,
    ticket_tier_id = ol.ticket_tier_id
FROM order_lines ol
WHERE t.order_line_id = ol.id
  AND t.order_id IS NULL;

-- Backfill code from serial_number for existing tickets
UPDATE tickets t
SET code = t.serial_number
WHERE t.code IS NULL AND t.serial_number IS NOT NULL;

-- Backfill qr_code from qr_code_data for existing tickets
UPDATE tickets t
SET qr_code = t.qr_code_data
WHERE t.qr_code IS NULL AND t.qr_code_data IS NOT NULL;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_tickets_order_id ON tickets(order_id);
CREATE INDEX IF NOT EXISTS idx_tickets_ticket_tier_id ON tickets(ticket_tier_id);
CREATE INDEX IF NOT EXISTS idx_tickets_code ON tickets(code);
CREATE INDEX IF NOT EXISTS idx_tickets_qr_code ON tickets(qr_code);
CREATE INDEX IF NOT EXISTS idx_tickets_attendee_email ON tickets(attendee_email);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);

-- Also ensure the payments table exists with the correct schema
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_transaction_id VARCHAR(255),
    amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    provider_response JSONB,
    webhook_received_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_provider_transaction_id ON payments(provider_transaction_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

COMMIT;
