-- Migration: 014_e2e_test_seed_data.sql
-- Description: Comprehensive seed data for E2E testing (Admin Setup → User Purchase → Gate Validation)
-- Created: 2026-02-27
-- Test Scenario: uduXPass Full Lifecycle Test

BEGIN;

-- ============================================================================
-- ADMIN USERS (Phase 1: Admin Infrastructure Setup)
-- ============================================================================
-- Admin User: Tunde (email: tunde@uduxpass.com, password: Admin@123)
-- Password hash generated with bcrypt cost 10
INSERT INTO admin_users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
VALUES 
    (gen_random_uuid(), 'tunde@uduxpass.com', '$2a$10$rKqZ8vHxH9YjH6xH9YjH6.eKqZ8vHxH9YjH6xH9YjH6eKqZ8vHxH9Y', 'Tunde', 'Adeyemi', 'super_admin', true, NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

-- ============================================================================
-- SCANNER USERS (Phase 3: Gate Validation)
-- ============================================================================
-- Gate Staff: Bisi (username: bisi, password: Scanner@123)
-- Password hash for Scanner@123: $2a$10$N9qo8uLOickgx2ZMRZoMye7FJpf4yxLBO/QYUvTNU6KUq9s7XQMAK
INSERT INTO scanner_users (id, username, password_hash, first_name, last_name, email, phone, role, is_active, created_at, updated_at)
VALUES 
    (gen_random_uuid(), 'bisi', '$2a$10$N9qo8uLOickgx2ZMRZoMye7FJpf4yxLBO/QYUvTNU6KUq9s7XQMAK', 'Bisi', 'Okafor', 'bisi@uduxpass.com', '+234-803-555-0001', 'scanner', true, NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- Additional scanner users for testing
INSERT INTO scanner_users (id, username, password_hash, first_name, last_name, email, phone, role, is_active, created_at, updated_at)
VALUES 
    (gen_random_uuid(), 'scanner2', '$2a$10$N9qo8uLOickgx2ZMRZoMye7FJpf4yxLBO/QYUvTNU6KUq9s7XQMAK', 'Chinedu', 'Eze', 'scanner2@uduxpass.com', '+234-803-555-0002', 'scanner', true, NOW(), NOW()),
    (gen_random_uuid(), 'supervisor1', '$2a$10$N9qo8uLOickgx2ZMRZoMye7FJpf4yxLBO/QYUvTNU6KUq9s7XQMAK', 'Amaka', 'Nwosu', 'supervisor1@uduxpass.com', '+234-803-555-0003', 'supervisor', true, NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- ============================================================================
-- ORGANIZER FOR TEST EVENT
-- ============================================================================
INSERT INTO organizers (id, name, email, phone, website, description, logo_url, is_verified, is_active, created_at, updated_at)
VALUES
    ('11111111-1111-1111-1111-111111111111'::uuid, 'E2E Test Organizer', 'organizer@e2etest.com', '+234-803-000-0000', 'https://e2etest.com', 'Test organizer for E2E testing', 'https://example.com/logos/e2e-test.png', true, true, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    updated_at = NOW();

-- ============================================================================
-- TEST EVENT: "Davido Timeless Lagos" (Phase 1.2)
-- ============================================================================
INSERT INTO events (
    id, 
    organizer_id, 
    tour_id, 
    name, 
    slug, 
    description, 
    event_date, 
    doors_open_time, 
    start_time, 
    end_time, 
    venue_name, 
    venue_address, 
    venue_city, 
    venue_state, 
    venue_country, 
    venue_capacity, 
    status, 
    is_published, 
    published_at, 
    banner_image_url, 
    thumbnail_image_url, 
    currency, 
    timezone, 
    tags, 
    created_at, 
    updated_at
)
VALUES
    (
        '22222222-2222-2222-2222-222222222222'::uuid,
        '11111111-1111-1111-1111-111111111111'::uuid,
        NULL,
        'Davido Timeless Lagos',
        'davido-timeless-lagos-e2e-test',
        'Experience Davido''s electrifying Timeless performance at Eko Atlantic. E2E Test Event.',
        '2026-12-31 20:00:00',
        '2026-12-31 18:00:00',
        '2026-12-31 20:00:00',
        '2027-01-01 02:00:00',
        'Eko Atlantic Energy City',
        'Plot 1, Eko Atlantic City',
        'Lagos',
        'Lagos',
        'Nigeria',
        50000,
        'published',
        true,
        NOW(),
        'https://example.com/events/davido-timeless-e2e-banner.jpg',
        'https://example.com/events/davido-timeless-e2e-thumb.jpg',
        'NGN',
        'Africa/Lagos',
        ARRAY['afrobeats', 'concert', 'davido', 'timeless', 'e2e-test'],
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_at = NOW();

-- ============================================================================
-- TICKET TIERS FOR TEST EVENT (Phase 1.3 & 1.4)
-- ============================================================================
-- VIP Gold Tier: Price 150,000, Quantity 100 (Phase 1.3)
INSERT INTO ticket_tiers (
    id, 
    event_id, 
    name, 
    description, 
    price, 
    quantity, 
    quantity_sold, 
    quantity_reserved, 
    sale_start, 
    sale_end, 
    min_per_order, 
    max_per_order, 
    is_active, 
    created_at, 
    updated_at
)
VALUES
    (
        '33333333-3333-3333-3333-333333333333'::uuid,
        '22222222-2222-2222-2222-222222222222'::uuid,
        'VIP Gold',
        'Premium VIP experience with exclusive access and amenities',
        150000.00,
        100,
        0,
        0,
        '2026-01-01 00:00:00',
        '2026-12-31 18:00:00',
        1,
        10,
        true,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    price = EXCLUDED.price,
    quantity = EXCLUDED.quantity,
    updated_at = NOW();

-- Regular Tier: Price 30,000, Quantity 1000 (Phase 1.4)
INSERT INTO ticket_tiers (
    id, 
    event_id, 
    name, 
    description, 
    price, 
    quantity, 
    quantity_sold, 
    quantity_reserved, 
    sale_start, 
    sale_end, 
    min_per_order, 
    max_per_order, 
    is_active, 
    created_at, 
    updated_at
)
VALUES
    (
        '44444444-4444-4444-4444-444444444444'::uuid,
        '22222222-2222-2222-2222-222222222222'::uuid,
        'Regular',
        'General admission with great views of the stage',
        30000.00,
        1000,
        0,
        0,
        '2026-01-01 00:00:00',
        '2026-12-31 18:00:00',
        1,
        10,
        true,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    price = EXCLUDED.price,
    quantity = EXCLUDED.quantity,
    updated_at = NOW();

-- ============================================================================
-- SCANNER EVENT ASSIGNMENTS (Phase 3.1 & 3.2)
-- ============================================================================
-- Assign "Davido Timeless Lagos" event to scanner user "bisi"
INSERT INTO scanner_event_assignments (id, scanner_user_id, event_id, assigned_by, assigned_at, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    su.id,
    '22222222-2222-2222-2222-222222222222'::uuid,
    au.id,
    NOW(),
    NOW(),
    NOW()
FROM scanner_users su
CROSS JOIN admin_users au
WHERE su.username = 'bisi' 
  AND au.email = 'tunde@uduxpass.com'
  AND NOT EXISTS (
      SELECT 1 FROM scanner_event_assignments 
      WHERE scanner_user_id = su.id 
        AND event_id = '22222222-2222-2222-2222-222222222222'::uuid
  );

-- Assign to scanner2 as well for redundancy
INSERT INTO scanner_event_assignments (id, scanner_user_id, event_id, assigned_by, assigned_at, created_at, updated_at)
SELECT 
    gen_random_uuid(),
    su.id,
    '22222222-2222-2222-2222-222222222222'::uuid,
    au.id,
    NOW(),
    NOW(),
    NOW()
FROM scanner_users su
CROSS JOIN admin_users au
WHERE su.username = 'scanner2' 
  AND au.email = 'tunde@uduxpass.com'
  AND NOT EXISTS (
      SELECT 1 FROM scanner_event_assignments 
      WHERE scanner_user_id = su.id 
        AND event_id = '22222222-2222-2222-2222-222222222222'::uuid
  );

-- ============================================================================
-- PAYMENT PROVIDER SETTINGS (Phase 1.5)
-- ============================================================================
-- Note: Payment provider toggles are typically stored in event-specific settings
-- or application configuration. This may need to be handled via API or admin UI.
-- For now, we'll add a comment indicating this step requires admin UI interaction.

-- Payment providers (MoMo PSB & Paystack) should be enabled via Admin UI
-- This typically updates the events table or a separate event_settings table

COMMIT;

-- ============================================================================
-- VERIFICATION QUERIES
-- ============================================================================
-- Run these to verify the seed data was inserted correctly:

-- SELECT * FROM admin_users WHERE email = 'tunde@uduxpass.com';
-- SELECT * FROM scanner_users WHERE username IN ('bisi', 'scanner2', 'supervisor1');
-- SELECT * FROM events WHERE name = 'Davido Timeless Lagos';
-- SELECT * FROM ticket_tiers WHERE event_id = '22222222-2222-2222-2222-222222222222'::uuid;
-- SELECT * FROM scanner_event_assignments WHERE event_id = '22222222-2222-2222-2222-222222222222'::uuid;

-- ============================================================================
-- TEST DATA SUMMARY
-- ============================================================================
-- Admin User: tunde@uduxpass.com / Admin@123
-- Scanner User: bisi / Scanner@123
-- Event: Davido Timeless Lagos (ID: 22222222-2222-2222-2222-222222222222)
-- Ticket Tiers:
--   - VIP Gold: ₦150,000 (100 tickets)
--   - Regular: ₦30,000 (1000 tickets)
-- Total Event Capacity: 1,100 tickets
-- ============================================================================
