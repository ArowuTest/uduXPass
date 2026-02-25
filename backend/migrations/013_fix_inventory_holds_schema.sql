-- Migration 013: Fix inventory_holds schema
-- Add order_id column and make session_id nullable

-- Add order_id column with foreign key to orders table
ALTER TABLE inventory_holds
ADD COLUMN order_id UUID REFERENCES orders(id) ON DELETE CASCADE;

-- Make session_id nullable (it was NOT NULL before)
ALTER TABLE inventory_holds
ALTER COLUMN session_id DROP NOT NULL;

-- Add index on order_id for better query performance
CREATE INDEX IF NOT EXISTS idx_inventory_holds_order_id ON inventory_holds(order_id);

-- Add comment explaining the change
COMMENT ON COLUMN inventory_holds.order_id IS 'References the order that created this inventory hold';
COMMENT ON COLUMN inventory_holds.session_id IS 'Session ID for temporary holds before order creation (nullable)';
