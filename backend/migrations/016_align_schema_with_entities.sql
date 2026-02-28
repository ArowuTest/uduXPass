-- =============================================================================
-- Migration 016: Align Database Schema with Go Entity Structs
-- =============================================================================
-- This migration adds all columns that exist in Go entity structs but were
-- missing from the database tables. Each change is idempotent (IF NOT EXISTS).
-- =============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- 1. ticket_tiers: Add 'sold' column
--    Entity: TicketTier.Sold int db:"sold"
--    Purpose: Tracks the number of sold tickets directly on the tier for fast
--             reads. Updated by triggers or application logic on order payment.
-- ---------------------------------------------------------------------------
ALTER TABLE ticket_tiers
    ADD COLUMN IF NOT EXISTS sold INTEGER NOT NULL DEFAULT 0;

-- Backfill sold count from existing paid orders
UPDATE ticket_tiers tt
SET sold = COALESCE((
    SELECT SUM(ol.quantity)
    FROM order_lines ol
    JOIN orders o ON ol.order_id = o.id
    WHERE ol.ticket_tier_id = tt.id
      AND o.status = 'paid'
      AND o.is_active = true
), 0);

-- Add a trigger to keep sold in sync automatically
CREATE OR REPLACE FUNCTION update_ticket_tier_sold_count()
RETURNS TRIGGER AS $$
BEGIN
    -- When an order is marked paid, increment sold counts
    IF TG_OP = 'UPDATE' AND NEW.status = 'paid' AND OLD.status != 'paid' THEN
        UPDATE ticket_tiers tt
        SET sold = sold + ol.quantity
        FROM order_lines ol
        WHERE ol.order_id = NEW.id
          AND ol.ticket_tier_id = tt.id;
    END IF;
    -- When a paid order is cancelled/refunded, decrement sold counts
    IF TG_OP = 'UPDATE' AND OLD.status = 'paid' AND NEW.status IN ('cancelled', 'refunded') THEN
        UPDATE ticket_tiers tt
        SET sold = GREATEST(0, sold - ol.quantity)
        FROM order_lines ol
        WHERE ol.order_id = NEW.id
          AND ol.ticket_tier_id = tt.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_update_ticket_tier_sold ON orders;
CREATE TRIGGER trg_update_ticket_tier_sold
    AFTER UPDATE OF status ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_ticket_tier_sold_count();

-- ---------------------------------------------------------------------------
-- 2. orders: Add 'paid_at' and 'payment_reference' columns
--    Entity: Order.PaidAt *time.Time db:"paid_at"
--            Order.PaymentReference *string db:"payment_reference"
-- ---------------------------------------------------------------------------
ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS paid_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN IF NOT EXISTS payment_reference VARCHAR(255);

-- Backfill paid_at from confirmed_at for existing paid orders
UPDATE orders
SET paid_at = confirmed_at
WHERE status = 'paid'
  AND paid_at IS NULL
  AND confirmed_at IS NOT NULL;

-- ---------------------------------------------------------------------------
-- 3. order_lines: Add 'subtotal' and 'updated_at' columns
--    Entity: OrderLine.Subtotal float64 db:"subtotal"
--            OrderLine.UpdatedAt time.Time db:"updated_at"
--    Note: DB has 'total_price' which is semantically identical to 'subtotal'.
--          We add 'subtotal' and populate it from 'total_price', then keep
--          both columns for backward compatibility. Repositories will use
--          'subtotal' going forward.
-- ---------------------------------------------------------------------------
ALTER TABLE order_lines
    ADD COLUMN IF NOT EXISTS subtotal NUMERIC(10,2),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Backfill subtotal from total_price
UPDATE order_lines
SET subtotal = total_price
WHERE subtotal IS NULL;

-- Make subtotal NOT NULL after backfill
ALTER TABLE order_lines
    ALTER COLUMN subtotal SET NOT NULL,
    ALTER COLUMN subtotal SET DEFAULT 0;

-- Add trigger to keep updated_at current
DROP TRIGGER IF EXISTS trg_order_lines_updated_at ON order_lines;
CREATE TRIGGER trg_order_lines_updated_at
    BEFORE UPDATE ON order_lines
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ---------------------------------------------------------------------------
-- 4. organizers: Add 'website' column
--    Entity: Organizer.Website *string db:"website"
--    Note: DB has 'website_url'. We add 'website' and keep 'website_url' for
--          backward compatibility. Repositories will use 'website'.
-- ---------------------------------------------------------------------------
ALTER TABLE organizers
    ADD COLUMN IF NOT EXISTS website VARCHAR(500);

-- Backfill website from website_url
UPDATE organizers
SET website = website_url
WHERE website IS NULL AND website_url IS NOT NULL;

-- ---------------------------------------------------------------------------
-- 5. otp_tokens: Add missing columns
--    Entity: OTPToken.Code string db:"code"
--            OTPToken.Email *string db:"email"
--            OTPToken.Status OTPStatus db:"status"
--            OTPToken.AttemptCount int db:"attempt_count"
--            OTPToken.MaxAttempts int db:"max_attempts"
--            OTPToken.IPAddress *string db:"ip_address"
--            OTPToken.UserAgent *string db:"user_agent"
--            OTPToken.UpdatedAt time.Time db:"updated_at"
--    Note: DB has 'token' (maps to 'code') and 'attempts' (maps to 'attempt_count').
--          We add the entity-named columns and populate from existing ones.
-- ---------------------------------------------------------------------------
ALTER TABLE otp_tokens
    ADD COLUMN IF NOT EXISTS code VARCHAR(6),
    ADD COLUMN IF NOT EXISTS email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS attempt_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS max_attempts INTEGER NOT NULL DEFAULT 3,
    ADD COLUMN IF NOT EXISTS ip_address VARCHAR(45),
    ADD COLUMN IF NOT EXISTS user_agent TEXT,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Backfill code from token (existing column)
UPDATE otp_tokens
SET code = token
WHERE code IS NULL AND token IS NOT NULL;

-- Backfill attempt_count from attempts (existing column)
UPDATE otp_tokens
SET attempt_count = attempts
WHERE attempt_count = 0 AND attempts > 0;

-- Backfill status: if used_at is set → 'used', if expires_at < now → 'expired', else 'pending'
UPDATE otp_tokens
SET status = CASE
    WHEN used_at IS NOT NULL THEN 'used'
    WHEN expires_at < NOW() THEN 'expired'
    ELSE 'pending'
END
WHERE status = 'pending';

-- Make code NOT NULL after backfill (allow NULL for rows where token was NULL)
-- We leave code as nullable to avoid breaking existing rows with NULL token

-- Add trigger to keep updated_at current
DROP TRIGGER IF EXISTS trg_otp_tokens_updated_at ON otp_tokens;
CREATE TRIGGER trg_otp_tokens_updated_at
    BEFORE UPDATE ON otp_tokens
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ---------------------------------------------------------------------------
-- 6. inventory_holds: Add denormalised columns for performance
--    Entity: InventoryHold.Status InventoryHoldStatus db:"status"
--            InventoryHold.EventTitle string db:"event_title"
--            InventoryHold.EventSlug string db:"event_slug"
--            InventoryHold.TicketTierName string db:"ticket_tier_name"
--            InventoryHold.TicketTierCapacity int db:"ticket_tier_capacity"
--            InventoryHold.OrderCode string db:"order_code"
--            InventoryHold.OrderStatus string db:"order_status"
-- ---------------------------------------------------------------------------
ALTER TABLE inventory_holds
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS event_title VARCHAR(255),
    ADD COLUMN IF NOT EXISTS event_slug VARCHAR(100),
    ADD COLUMN IF NOT EXISTS ticket_tier_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS ticket_tier_capacity INTEGER,
    ADD COLUMN IF NOT EXISTS order_code VARCHAR(16),
    ADD COLUMN IF NOT EXISTS order_status VARCHAR(20);

-- Backfill denormalised columns from related tables
UPDATE inventory_holds
SET
    ticket_tier_name = sub.tier_name,
    ticket_tier_capacity = sub.tier_quota,
    event_title = sub.event_name,
    event_slug = sub.event_slug,
    order_code = sub.order_code,
    order_status = sub.order_status
FROM (
    SELECT
        ih.id AS hold_id,
        tt.name AS tier_name,
        tt.quota AS tier_quota,
        e.name AS event_name,
        e.slug AS event_slug,
        o.code AS order_code,
        o.status::VARCHAR AS order_status
    FROM inventory_holds ih
    JOIN ticket_tiers tt ON ih.ticket_tier_id = tt.id
    JOIN events e ON tt.event_id = e.id
    LEFT JOIN orders o ON ih.order_id = o.id
    WHERE ih.event_title IS NULL OR ih.ticket_tier_name IS NULL
) sub
WHERE inventory_holds.id = sub.hold_id;

-- Mark expired holds as expired
UPDATE inventory_holds
SET status = 'expired'
WHERE expires_at < NOW()
  AND status = 'active';

-- ---------------------------------------------------------------------------
-- 7. users: Add 'password' column as alias for password_hash
--    Entity: User.Password *string db:"password"
--            User.PasswordHash *string db:"password_hash"
--    Note: The entity has both fields. 'password_hash' is the canonical column.
--          We add 'password' as a separate nullable column for temporary use
--          during authentication flows (e.g., storing plaintext temporarily
--          before hashing — though this should never persist).
--          In practice, repositories should only write to password_hash.
--          We add the column to satisfy the db tag mapping.
-- ---------------------------------------------------------------------------
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS password VARCHAR(255);

-- ---------------------------------------------------------------------------
-- 8. admin_login_history: Add 'failure_reason' and 'login_time' columns
--    Entity: AdminLoginHistory.FailureReason string db:"failure_reason"
--            AdminLoginHistory.LoginTime time.Time db:"login_time"
--    Note: DB has 'login_at' (maps to 'login_time') and lacks 'failure_reason'.
-- ---------------------------------------------------------------------------
ALTER TABLE admin_login_history
    ADD COLUMN IF NOT EXISTS failure_reason TEXT,
    ADD COLUMN IF NOT EXISTS login_time TIMESTAMP WITH TIME ZONE;

-- Backfill login_time from login_at
UPDATE admin_login_history
SET login_time = login_at
WHERE login_time IS NULL;

-- ---------------------------------------------------------------------------
-- 9. events: Add missing columns to match DB (entity is missing these)
--    DB has: tour_id, venue_latitude, venue_longitude
--    These are already in the DB — we just need to update the Go entity.
--    No SQL changes needed here; entity fix is in Go code.
-- ---------------------------------------------------------------------------
-- (No SQL needed — columns already exist in DB)

-- ---------------------------------------------------------------------------
-- 10. ticket_tiers: Add 'settings' column to entity (DB has it, entity lacks)
--     The DB has settings JSONB DEFAULT '{}'. Entity doesn't have it.
--     We add MetaInfo JSONB to entity (already exists as meta_info in some
--     versions). Check if meta_info exists and add settings if not.
-- ---------------------------------------------------------------------------
ALTER TABLE ticket_tiers
    ADD COLUMN IF NOT EXISTS meta_info JSONB DEFAULT '{}';

-- Copy settings to meta_info if meta_info is empty
UPDATE ticket_tiers
SET meta_info = settings
WHERE meta_info = '{}' AND settings != '{}';

-- ---------------------------------------------------------------------------
-- 11. scanner_login_history: Align with entity
--     Entity ScannerLoginHistory uses: login_at, logout_at
--     DB has: login_at, logout_at — already in sync
-- ---------------------------------------------------------------------------
-- (No SQL needed)

-- ---------------------------------------------------------------------------
-- 12. Add indexes for new columns that will be queried frequently
-- ---------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_ticket_tiers_sold ON ticket_tiers(sold);
CREATE INDEX IF NOT EXISTS idx_orders_paid_at ON orders(paid_at) WHERE paid_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_orders_payment_reference ON orders(payment_reference) WHERE payment_reference IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_otp_tokens_code ON otp_tokens(code) WHERE code IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_otp_tokens_email ON otp_tokens(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_otp_tokens_status ON otp_tokens(status);
CREATE INDEX IF NOT EXISTS idx_inventory_holds_status ON inventory_holds(status);

COMMIT;
