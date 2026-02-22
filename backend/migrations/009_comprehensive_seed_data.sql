-- Migration: Comprehensive Seed Data for Production-Ready Testing
-- Purpose: Strategic seed data for complete E2E flow testing
-- Date: 2026-02-17

-- ============================================================================
-- USERS (Test Users)
-- ============================================================================

INSERT INTO users (id, email, phone_number, first_name, last_name, password_hash, auth_provider, is_active, created_at, updated_at)
VALUES
  (gen_random_uuid(), 'john.doe@test.com', '+2348012345678', 'John', 'Doe', '$2a$10$YourHashedPasswordHere', 'email', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (gen_random_uuid(), 'jane.smith@test.com', '+2348023456789', 'Jane', 'Smith', '$2a$10$YourHashedPasswordHere', 'email', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (gen_random_uuid(), 'mike.johnson@test.com', '+2348034567890', 'Mike', 'Johnson', '$2a$10$YourHashedPasswordHere', 'email', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

-- ============================================================================
-- EVENTS WITH TICKET TIERS (Using existing organizers)
-- ============================================================================

-- Get organizer IDs
DO $$
DECLARE
  org1_id UUID;
  org2_id UUID;
  event1_id UUID;
  event2_id UUID;
  event3_id UUID;
  tier1_id UUID;
  tier2_id UUID;
  tier3_id UUID;
BEGIN
  -- Get existing organizer IDs
  SELECT id INTO org1_id FROM organizers ORDER BY created_at LIMIT 1;
  SELECT id INTO org2_id FROM organizers ORDER BY created_at OFFSET 1 LIMIT 1;
  
  -- Create Event 1: Burna Boy Live in Lagos
  event1_id := gen_random_uuid();
  INSERT INTO events (
    id, organizer_id, name, slug, description, event_date, doors_open,
    venue_name, venue_address, venue_city, venue_state, venue_country,
    venue_capacity, event_image_url, sale_start, sale_end, status,
    created_at, updated_at
  ) VALUES (
    event1_id, org1_id, 'Burna Boy Live in Lagos', 'burna-boy-lagos-2026',
    'Experience the African Giant live in concert! An unforgettable night of Afrobeats, dancehall, and reggae.',
    '2026-06-15 20:00:00', '2026-06-15 18:00:00',
    'Eko Convention Centre', '1-2 Eko Hotel Roundabout, Victoria Island', 'Lagos', 'Lagos', 'Nigeria',
    15000, 'https://example.com/burna-boy.jpg',
    CURRENT_TIMESTAMP, '2026-06-14 23:59:59', 'published',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
  ) ON CONFLICT (slug) DO NOTHING;
  
  -- Ticket Tiers for Event 1
  tier1_id := gen_random_uuid();
  INSERT INTO ticket_tiers (
    id, event_id, name, description, price, quota, sold, min_purchase, max_purchase,
    sale_start, sale_end, is_active, created_at, updated_at
  ) VALUES
    (tier1_id, event1_id, 'Early Bird', 'Limited early bird tickets', 15000.00, 500, 0, 1, 10, CURRENT_TIMESTAMP, '2026-05-01 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), event1_id, 'Regular', 'Standard admission', 25000.00, 8000, 0, 1, 10, CURRENT_TIMESTAMP, '2026-06-14 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), event1_id, 'VIP', 'VIP access with premium seating', 50000.00, 1000, 0, 1, 5, CURRENT_TIMESTAMP, '2026-06-14 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), event1_id, 'VVIP', 'Exclusive backstage access and meet & greet', 150000.00, 100, 0, 1, 2, CURRENT_TIMESTAMP, '2026-06-14 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
  
  -- Create Event 2: Wizkid - Made in Lagos Tour
  event2_id := gen_random_uuid();
  INSERT INTO events (
    id, organizer_id, name, slug, description, event_date, doors_open,
    venue_name, venue_address, venue_city, venue_state, venue_country,
    venue_capacity, event_image_url, sale_start, sale_end, status,
    created_at, updated_at
  ) VALUES (
    event2_id, org1_id, 'Wizkid - Made in Lagos Tour', 'wizkid-mil-abuja-2026',
    'Starboy brings the Made in Lagos experience to Abuja! A night of Afrobeats magic.',
    '2026-07-20 21:00:00', '2026-07-20 19:00:00',
    'Abuja Stadium', 'Area 10, Garki, Abuja', 'Abuja', 'FCT', 'Nigeria',
    25000, 'https://example.com/wizkid.jpg',
    CURRENT_TIMESTAMP, '2026-07-19 23:59:59', 'published',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
  ) ON CONFLICT (slug) DO NOTHING;
  
  -- Ticket Tiers for Event 2
  INSERT INTO ticket_tiers (
    id, event_id, name, description, price, quota, sold, min_purchase, max_purchase,
    sale_start, sale_end, is_active, created_at, updated_at
  ) VALUES
    (gen_random_uuid(), event2_id, 'General Admission', 'Standing room', 10000.00, 15000, 0, 1, 10, CURRENT_TIMESTAMP, '2026-07-19 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), event2_id, 'Premium', 'Reserved seating', 30000.00, 5000, 0, 1, 8, CURRENT_TIMESTAMP, '2026-07-19 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), event2_id, 'VIP Lounge', 'VIP lounge access with complimentary drinks', 75000.00, 500, 0, 1, 4, CURRENT_TIMESTAMP, '2026-07-19 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
  
  -- Create Event 3: Davido - Timeless Concert
  event3_id := gen_random_uuid();
  INSERT INTO events (
    id, organizer_id, name, slug, description, event_date, doors_open,
    venue_name, venue_address, venue_city, venue_state, venue_country,
    venue_capacity, event_image_url, sale_start, sale_end, status,
    created_at, updated_at
  ) VALUES (
    event3_id, org2_id, 'Davido - Timeless Concert', 'davido-timeless-ph-2026',
    'OBO brings the Timeless album to life in Port Harcourt! An epic night of hits.',
    '2026-08-10 20:00:00', '2026-08-10 18:30:00',
    'Port Harcourt Pleasure Park', 'Aba Road, Port Harcourt', 'Port Harcourt', 'Rivers', 'Nigeria',
    20000, 'https://example.com/davido.jpg',
    CURRENT_TIMESTAMP, '2026-08-09 23:59:59', 'published',
    CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
  ) ON CONFLICT (slug) DO NOTHING;
  
  -- Ticket Tiers for Event 3
  INSERT INTO ticket_tiers (
    id, event_id, name, description, price, quota, sold, min_purchase, max_purchase,
    sale_start, sale_end, is_active, created_at, updated_at
  ) VALUES
    (gen_random_uuid(), event3_id, 'Student Ticket', 'Valid student ID required', 8000.00, 2000, 0, 1, 4, CURRENT_TIMESTAMP, '2026-08-09 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), event3_id, 'Regular', 'Standard admission', 20000.00, 12000, 0, 1, 10, CURRENT_TIMESTAMP, '2026-08-09 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), event3_id, 'VIP', 'VIP section with premium view', 60000.00, 1500, 0, 1, 6, CURRENT_TIMESTAMP, '2026-08-09 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), event3_id, 'Diamond', 'Ultimate VIP experience with backstage access', 200000.00, 50, 0, 1, 2, CURRENT_TIMESTAMP, '2026-08-09 23:59:59', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

END $$;

-- ============================================================================
-- GRANT PERMISSIONS
-- ============================================================================

GRANT ALL ON users TO uduxpass_user;
GRANT ALL ON events TO uduxpass_user;
GRANT ALL ON ticket_tiers TO uduxpass_user;

-- Add comments
COMMENT ON TABLE users IS 'Test users for E2E testing';
COMMENT ON TABLE events IS 'Sample events with realistic African music concerts';
COMMENT ON TABLE ticket_tiers IS 'Multiple ticket tiers per event for comprehensive testing';
