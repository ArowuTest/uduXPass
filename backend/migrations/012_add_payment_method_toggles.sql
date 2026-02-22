-- Migration: Add payment method toggles to events table
-- Created: 2026-02-22
-- Purpose: Allow admins to enable/disable specific payment methods per event

-- Add payment method toggle columns to events table
ALTER TABLE events 
ADD COLUMN enable_momo BOOLEAN DEFAULT true,
ADD COLUMN enable_paystack BOOLEAN DEFAULT true;

-- Update existing events to have both payment methods enabled by default
UPDATE events 
SET enable_momo = true, enable_paystack = true
WHERE enable_momo IS NULL OR enable_paystack IS NULL;

-- Add comment for documentation
COMMENT ON COLUMN events.enable_momo IS 'Whether MoMo PSB payment is enabled for this event';
COMMENT ON COLUMN events.enable_paystack IS 'Whether Paystack payment is enabled for this event';
