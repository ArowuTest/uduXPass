-- Migration 019: Add missing users columns and seed scanner event assignments
-- Author: Enterprise schema alignment
-- Date: 2026-02-28

BEGIN;

-- ============================================================
-- 1. Add missing columns to users table
-- ============================================================
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_login_ip VARCHAR(45),
    ADD COLUMN IF NOT EXISTS email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- ============================================================
-- 2. Seed scanner event assignments for all scanners to both published events
-- Assign all scanners to both published events so E2E tests work
-- ============================================================

-- Get admin user ID for assigned_by (use the first admin user)
DO $$
DECLARE
    admin_id UUID;
    davido_event_id UUID := '3b9e540e-07ef-4f64-a4e6-252d45f93315';
    wizkid_event_id UUID := '8e1edfda-163f-47dd-9f6e-b5e4202eaff8';
    scanner_rec RECORD;
BEGIN
    -- Get admin user ID
    SELECT id INTO admin_id FROM admin_users LIMIT 1;
    
    -- Assign all active scanners to both published events
    FOR scanner_rec IN SELECT id FROM scanner_users WHERE status = 'active' LOOP
        -- Assign to Davido event
        INSERT INTO scanner_event_assignments (id, scanner_id, event_id, assigned_by, assigned_at, is_active)
        VALUES (gen_random_uuid(), scanner_rec.id, davido_event_id, admin_id, NOW(), TRUE)
        ON CONFLICT DO NOTHING;
        
        -- Assign to Wizkid event
        INSERT INTO scanner_event_assignments (id, scanner_id, event_id, assigned_by, assigned_at, is_active)
        VALUES (gen_random_uuid(), scanner_rec.id, wizkid_event_id, admin_id, NOW(), TRUE)
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;

-- ============================================================
-- 3. Verify the assignments were created
-- ============================================================
DO $$
DECLARE
    assignment_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO assignment_count FROM scanner_event_assignments;
    IF assignment_count = 0 THEN
        RAISE EXCEPTION 'No scanner event assignments were created';
    END IF;
    RAISE NOTICE 'Created % scanner event assignments', assignment_count;
END $$;

COMMIT;
