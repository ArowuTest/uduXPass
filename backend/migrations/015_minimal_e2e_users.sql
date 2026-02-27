-- Migration: 015_minimal_e2e_users.sql
-- Description: Minimal seed data for E2E testing - ONLY user accounts
-- Created: 2026-02-27
-- Note: Events, tiers, and settings will be created via UI during testing

BEGIN;

-- ============================================================================
-- ADMIN USER (for Phase 1: Admin Infrastructure Setup)
-- ============================================================================
-- Admin User: Tunde
-- Email: tunde@uduxpass.com
-- Password: Admin@123
-- Password hash: $2a$10$rKqZ8vHxH9YjH6xH9YjH6.eKqZ8vHxH9YjH6xH9YjH6eKqZ8vHxH9Y

INSERT INTO admin_users (
    id, 
    email, 
    password_hash, 
    first_name, 
    last_name, 
    role, 
    is_active,
    must_change_password,
    created_at, 
    updated_at
)
VALUES 
    (
        '11111111-1111-1111-1111-111111111111'::uuid,
        'tunde@uduxpass.com', 
        '$2a$10$rKqZ8vHxH9YjH6xH9YjH6.eKqZ8vHxH9YjH6xH9YjH6eKqZ8vHxH9Y', 
        'Tunde', 
        'Adeyemi', 
        'super_admin', 
        true,
        false,
        NOW(), 
        NOW()
    )
ON CONFLICT (email) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    role = EXCLUDED.role,
    is_active = EXCLUDED.is_active,
    must_change_password = EXCLUDED.must_change_password,
    updated_at = NOW();

-- ============================================================================
-- SCANNER USERS (for Phase 3: Gate Validation)
-- ============================================================================
-- Gate Staff: Bisi
-- Username: bisi
-- Password: Scanner@123
-- Password hash: $2a$10$N9qo8uLOickgx2ZMRZoMye7FJpf4yxLBO/QYUvTNU6KUq9s7XQMAK

INSERT INTO scanner_users (
    id, 
    username, 
    password_hash, 
    name,
    email, 
    role, 
    status,
    created_by,
    created_at, 
    updated_at
)
VALUES 
    (
        '22222222-2222-2222-2222-222222222222'::uuid,
        'bisi', 
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7FJpf4yxLBO/QYUvTNU6KUq9s7XQMAK', 
        'Bisi Okafor',
        'bisi@uduxpass.com', 
        'scanner_operator', 
        'active',
        '11111111-1111-1111-1111-111111111111'::uuid,
        NOW(), 
        NOW()
    )
ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    role = EXCLUDED.role,
    status = EXCLUDED.status,
    updated_at = NOW();

-- Additional scanner user for testing
INSERT INTO scanner_users (
    id, 
    username, 
    password_hash, 
    name,
    email, 
    role, 
    status,
    created_by,
    created_at, 
    updated_at
)
VALUES 
    (
        '33333333-3333-3333-3333-333333333333'::uuid,
        'scanner1', 
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7FJpf4yxLBO/QYUvTNU6KUq9s7XQMAK', 
        'Scanner One',
        'scanner1@uduxpass.com', 
        'scanner_operator', 
        'active',
        '11111111-1111-1111-1111-111111111111'::uuid,
        NOW(), 
        NOW()
    )
ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    updated_at = NOW();

COMMIT;

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================
-- SELECT * FROM admin_users WHERE email = 'tunde@uduxpass.com';
-- SELECT * FROM scanner_users WHERE username IN ('bisi', 'scanner1');

-- ============================================================================
-- TEST CREDENTIALS SUMMARY
-- ============================================================================
-- Admin Login:
--   Email: tunde@uduxpass.com
--   Password: Admin@123
--
-- Scanner Login:
--   Username: bisi
--   Password: Scanner@123
--
-- Note: User "Chidi" will be created via registration UI in Phase 2
-- Note: Event "Davido Timeless Lagos" will be created via admin UI in Phase 1.2
-- Note: Ticket tiers will be created via admin UI in Phase 1.3 & 1.4
-- ============================================================================
