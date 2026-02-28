-- =============================================================================
-- Migration 017: Add Remaining Missing Columns
-- =============================================================================
-- Adds columns referenced in repository SQL queries that are missing from the DB.
-- All changes are idempotent (IF NOT EXISTS).
-- =============================================================================

BEGIN;

-- ---------------------------------------------------------------------------
-- 1. order_lines: Add fees, taxes, discount_amount
--    These are financial fields used in CreateBatch/UpdateBatch repository methods.
--    They represent service fees, applicable taxes, and any discount applied
--    to the line item — standard for enterprise ticketing platforms.
-- ---------------------------------------------------------------------------
ALTER TABLE order_lines
    ADD COLUMN IF NOT EXISTS fees NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS taxes NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS discount_amount NUMERIC(10,2) NOT NULL DEFAULT 0;

-- ---------------------------------------------------------------------------
-- 2. scanner_users: The entity uses 'name' (single field) but the repository
--    references 'first_name' and 'last_name'. Since the DB has 'name' and
--    the entity has 'name', we fix the repository queries (not the DB).
--    However, to support the existing repository code that uses first_name/last_name
--    in some queries, we add generated columns for backward compatibility.
-- ---------------------------------------------------------------------------
ALTER TABLE scanner_users
    ADD COLUMN IF NOT EXISTS first_name VARCHAR(100) GENERATED ALWAYS AS (
        CASE 
            WHEN position(' ' IN name) > 0 
            THEN LEFT(name, position(' ' IN name) - 1)
            ELSE name
        END
    ) STORED,
    ADD COLUMN IF NOT EXISTS last_name VARCHAR(100) GENERATED ALWAYS AS (
        CASE 
            WHEN position(' ' IN name) > 0 
            THEN SUBSTRING(name FROM position(' ' IN name) + 1)
            ELSE ''
        END
    ) STORED;

-- ---------------------------------------------------------------------------
-- 3. scanner_login_history: Add 'login_time' as alias for 'login_at'
--    Entity ScannerLoginHistory uses login_at (already in DB).
--    AdminLoginHistory uses login_time (added in migration 016).
--    No change needed here.
-- ---------------------------------------------------------------------------

-- ---------------------------------------------------------------------------
-- 4. otp_tokens: Add 'is_used' column for backward compatibility
--    Some repository methods still use is_used in queries.
--    We add it as a generated column based on status.
-- ---------------------------------------------------------------------------
ALTER TABLE otp_tokens
    ADD COLUMN IF NOT EXISTS is_used BOOLEAN GENERATED ALWAYS AS (
        status = 'used'
    ) STORED;

-- ---------------------------------------------------------------------------
-- 5. scanner_audit_logs: Verify table exists with correct columns
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS scanner_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scanner_id UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES scanner_sessions(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    details JSONB DEFAULT '{}',
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_scanner_audit_logs_scanner_id ON scanner_audit_logs(scanner_id);
CREATE INDEX IF NOT EXISTS idx_scanner_audit_logs_session_id ON scanner_audit_logs(session_id);
CREATE INDEX IF NOT EXISTS idx_scanner_audit_logs_action ON scanner_audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_scanner_audit_logs_created_at ON scanner_audit_logs(created_at);

-- ---------------------------------------------------------------------------
-- 6. ticket_validations: Verify table exists with correct columns
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ticket_validations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    scanner_id UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES scanner_sessions(id) ON DELETE CASCADE,
    validation_result VARCHAR(50) NOT NULL DEFAULT 'valid',
    validation_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ticket_validations_ticket_id ON ticket_validations(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_validations_scanner_id ON ticket_validations(scanner_id);
CREATE INDEX IF NOT EXISTS idx_ticket_validations_session_id ON ticket_validations(session_id);
CREATE INDEX IF NOT EXISTS idx_ticket_validations_timestamp ON ticket_validations(validation_timestamp);

-- ---------------------------------------------------------------------------
-- 7. scanner_event_assignments: Verify table exists with correct columns
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS scanner_event_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scanner_id UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    assigned_by UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT true,
    UNIQUE(scanner_id, event_id)
);

CREATE INDEX IF NOT EXISTS idx_scanner_event_assignments_scanner ON scanner_event_assignments(scanner_id);
CREATE INDEX IF NOT EXISTS idx_scanner_event_assignments_event ON scanner_event_assignments(event_id);

COMMIT;
