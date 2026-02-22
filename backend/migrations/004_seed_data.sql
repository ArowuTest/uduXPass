-- Seed Data for uduXPass Platform
-- This file contains realistic test data for development and testing

-- ============================================================================
-- ORGANIZERS
-- ============================================================================

INSERT INTO organizers (id, name, email, phone, website, description, logo_url, is_verified, is_active, created_at, updated_at)
VALUES
    ('a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::uuid, 'Afro Nation Entertainment', 'contact@afronation.com', '+234-803-123-4567', 'https://afronation.com', 'Leading African music festival organizer bringing world-class entertainment experiences', 'https://example.com/logos/afronation.png', true, true, NOW(), NOW()),
    ('b2c3d4e5-f6a7-5b6c-9d0e-1f2a3b4c5d6e'::uuid, 'Lagos Live Events', 'info@lagoslive.ng', '+234-805-234-5678', 'https://lagoslive.ng', 'Premier event production company specializing in concerts and festivals in Lagos', 'https://example.com/logos/lagoslive.png', true, true, NOW(), NOW()),
    ('c3d4e5f6-a7b8-6c7d-0e1f-2a3b4c5d6e7f'::uuid, 'Naija Vibes Productions', 'hello@naijavibes.com', '+234-807-345-6789', 'https://naijavibes.com', 'Innovative event organizers creating unforgettable music experiences across Nigeria', 'https://example.com/logos/naijavibes.png', true, true, NOW(), NOW());

-- ============================================================================
-- TOURS
-- ============================================================================

INSERT INTO tours (id, organizer_id, name, slug, description, start_date, end_date, is_active, created_at, updated_at)
VALUES
    ('d4e5f6a7-b8c9-7d8e-1f2a-3b4c5d6e7f8a'::uuid, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::uuid, 'Davido Timeless Tour 2026', 'davido-timeless-tour-2026', 'Davido brings his Timeless album to life across 5 major Nigerian cities', '2026-03-01', '2026-04-30', true, NOW(), NOW()),
    ('e5f6a7b8-c9d0-8e9f-2a3b-4c5d6e7f8a9b'::uuid, 'b2c3d4e5-f6a7-5b6c-9d0e-1f2a3b4c5d6e'::uuid, 'Burna Boy African Giant Experience', 'burna-boy-african-giant-2026', 'Burna Boy celebrates African music with exclusive performances', '2026-05-15', '2026-06-30', true, NOW(), NOW());

-- ============================================================================
-- EVENTS
-- ============================================================================

INSERT INTO events (id, organizer_id, tour_id, name, slug, description, event_date, doors_open_time, start_time, end_time, venue_name, venue_address, venue_city, venue_state, venue_country, venue_capacity, status, is_published, published_at, banner_image_url, thumbnail_image_url, currency, timezone, tags, created_at, updated_at)
VALUES
    -- Davido Tour Events
    ('f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c'::uuid, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::uuid, 'd4e5f6a7-b8c9-7d8e-1f2a-3b4c5d6e7f8a'::uuid, 'Davido Live in Lagos', 'davido-live-lagos-2026', 'Experience Davido''s electrifying performance at Eko Atlantic. The Timeless album comes alive with special guest appearances and unforgettable moments.', '2026-03-15 20:00:00', '2026-03-15 18:00:00', '2026-03-15 20:00:00', '2026-03-16 02:00:00', 'Eko Atlantic Energy City', 'Plot 1, Eko Atlantic City', 'Lagos', 'Lagos', 'Nigeria', 50000, 'published', true, NOW(), 'https://example.com/events/davido-lagos-banner.jpg', 'https://example.com/events/davido-lagos-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'concert', 'davido', 'timeless'], NOW(), NOW()),
    
    ('a7b8c9d0-e1f2-0a1b-4c5d-6e7f8a9b0c1d'::uuid, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::uuid, 'd4e5f6a7-b8c9-7d8e-1f2a-3b4c5d6e7f8a'::uuid, 'Davido Live in Abuja', 'davido-live-abuja-2026', 'The capital city gets ready for Davido! An epic night of Afrobeats at the iconic Abuja Stadium.', '2026-03-22 19:30:00', '2026-03-22 17:30:00', '2026-03-22 19:30:00', '2026-03-23 01:00:00', 'Abuja National Stadium', 'Package B, Central Area', 'Abuja', 'FCT', 'Nigeria', 35000, 'published', true, NOW(), 'https://example.com/events/davido-abuja-banner.jpg', 'https://example.com/events/davido-abuja-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'concert', 'davido', 'timeless'], NOW(), NOW()),
    
    ('b8c9d0e1-f2a3-1b2c-5d6e-7f8a9b0c1d2e'::uuid, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::uuid, 'd4e5f6a7-b8c9-7d8e-1f2a-3b4c5d6e7f8a'::uuid, 'Davido Live in Port Harcourt', 'davido-live-portharcourt-2026', 'Port Harcourt, get ready to party! Davido brings the Timeless experience to the Garden City.', '2026-04-05 20:00:00', '2026-04-05 18:00:00', '2026-04-05 20:00:00', '2026-04-06 02:00:00', 'Liberation Stadium', 'Elekahia, Port Harcourt', 'Port Harcourt', 'Rivers', 'Nigeria', 28000, 'published', true, NOW(), 'https://example.com/events/davido-ph-banner.jpg', 'https://example.com/events/davido-ph-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'concert', 'davido', 'timeless'], NOW(), NOW()),
    
    ('c9d0e1f2-a3b4-2c3d-6e7f-8a9b0c1d2e3f'::uuid, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::uuid, 'd4e5f6a7-b8c9-7d8e-1f2a-3b4c5d6e7f8a'::uuid, 'Davido Live in Enugu', 'davido-live-enugu-2026', 'The coal city lights up with Davido! An unforgettable night of music and celebration.', '2026-04-12 19:00:00', '2026-04-12 17:00:00', '2026-04-12 19:00:00', '2026-04-13 01:00:00', 'Nnamdi Azikiwe Stadium', 'Okpara Avenue, Enugu', 'Enugu', 'Enugu', 'Nigeria', 22000, 'published', true, NOW(), 'https://example.com/events/davido-enugu-banner.jpg', 'https://example.com/events/davido-enugu-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'concert', 'davido', 'timeless'], NOW(), NOW()),
    
    ('d0e1f2a3-b4c5-3d4e-7f8a-9b0c1d2e3f4a'::uuid, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::uuid, 'd4e5f6a7-b8c9-7d8e-1f2a-3b4c5d6e7f8a'::uuid, 'Davido Live in Ibadan', 'davido-live-ibadan-2026', 'Ibadan, the ancient city, welcomes Davido for the grand finale of the Timeless Tour!', '2026-04-26 20:00:00', '2026-04-26 18:00:00', '2026-04-26 20:00:00', '2026-04-27 02:00:00', 'Lekan Salami Stadium', 'Liberty Road, Ibadan', 'Ibadan', 'Oyo', 'Nigeria', 25000, 'published', true, NOW(), 'https://example.com/events/davido-ibadan-banner.jpg', 'https://example.com/events/davido-ibadan-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'concert', 'davido', 'timeless', 'finale'], NOW(), NOW()),
    
    -- Burna Boy Events
    ('e1f2a3b4-c5d6-4e5f-8a9b-0c1d2e3f4a5b'::uuid, 'b2c3d4e5-f6a7-5b6c-9d0e-1f2a3b4c5d6e'::uuid, 'e5f6a7b8-c9d0-8e9f-2a3b-4c5d6e7f8a9b'::uuid, 'Burna Boy - African Giant Lagos', 'burna-boy-african-giant-lagos-2026', 'The African Giant returns home! Burna Boy delivers an explosive performance celebrating African excellence.', '2026-05-20 21:00:00', '2026-05-20 19:00:00', '2026-05-20 21:00:00', '2026-05-21 03:00:00', 'Tafawa Balewa Square', 'Lagos Island', 'Lagos', 'Lagos', 'Nigeria', 45000, 'published', true, NOW(), 'https://example.com/events/burna-lagos-banner.jpg', 'https://example.com/events/burna-lagos-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'concert', 'burna-boy', 'african-giant'], NOW(), NOW()),
    
    ('f2a3b4c5-d6e7-5f6a-9b0c-1d2e3f4a5b6c'::uuid, 'b2c3d4e5-f6a7-5b6c-9d0e-1f2a3b4c5d6e'::uuid, 'e5f6a7b8-c9d0-8e9f-2a3b-4c5d6e7f8a9b'::uuid, 'Burna Boy - African Giant Abuja', 'burna-boy-african-giant-abuja-2026', 'Burna Boy brings his Grammy-winning energy to Abuja for one night only!', '2026-06-10 20:00:00', '2026-06-10 18:00:00', '2026-06-10 20:00:00', '2026-06-11 02:00:00', 'Velodrome', 'Package B, Central Area', 'Abuja', 'FCT', 'Nigeria', 30000, 'published', true, NOW(), 'https://example.com/events/burna-abuja-banner.jpg', 'https://example.com/events/burna-abuja-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'concert', 'burna-boy', 'african-giant'], NOW(), NOW()),
    
    -- Standalone Events
    ('a3b4c5d6-e7f8-6a7b-0c1d-2e3f4a5b6c7d'::uuid, 'c3d4e5f6-a7b8-6c7d-0e1f-2a3b4c5d6e7f'::uuid, NULL, 'Wizkid - More Love Less Ego Lagos', 'wizkid-more-love-less-ego-lagos-2026', 'Wizkid returns to Lagos with his critically acclaimed More Love Less Ego tour. A night of pure vibes and Afrobeats magic.', '2026-07-15 20:30:00', '2026-07-15 18:30:00', '2026-07-15 20:30:00', '2026-07-16 02:30:00', 'Eko Convention Centre', 'Eko Hotel & Suites, Victoria Island', 'Lagos', 'Lagos', 'Nigeria', 15000, 'published', true, NOW(), 'https://example.com/events/wizkid-lagos-banner.jpg', 'https://example.com/events/wizkid-lagos-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'concert', 'wizkid', 'more-love-less-ego'], NOW(), NOW()),
    
    ('b4c5d6e7-f8a9-7b8c-1d2e-3f4a5b6c7d8e'::uuid, 'c3d4e5f6-a7b8-6c7d-0e1f-2a3b4c5d6e7f'::uuid, NULL, 'Afrobeats Festival Lagos 2026', 'afrobeats-festival-lagos-2026', 'The biggest Afrobeats festival of the year! 20+ artists, 2 days of non-stop music, food, and culture.', '2026-08-22 16:00:00', '2026-08-22 14:00:00', '2026-08-22 16:00:00', '2026-08-24 02:00:00', 'Eko Atlantic Energy City', 'Plot 1, Eko Atlantic City', 'Lagos', 'Lagos', 'Nigeria', 75000, 'published', true, NOW(), 'https://example.com/events/afrobeats-fest-banner.jpg', 'https://example.com/events/afrobeats-fest-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['afrobeats', 'festival', 'multi-day', 'culture'], NOW(), NOW()),
    
    ('c5d6e7f8-a9b0-8c9d-2e3f-4a5b6c7d8e9f'::uuid, 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d'::uuid, NULL, 'New Year''s Eve Countdown Lagos', 'nye-countdown-lagos-2026', 'Ring in 2027 with the biggest New Year''s Eve party in Lagos! Multiple stages, top DJs, and an unforgettable countdown.', '2026-12-31 20:00:00', '2026-12-31 18:00:00', '2026-12-31 20:00:00', '2027-01-01 04:00:00', 'Landmark Beach', 'Oniru, Victoria Island', 'Lagos', 'Lagos', 'Nigeria', 60000, 'draft', false, NULL, 'https://example.com/events/nye-lagos-banner.jpg', 'https://example.com/events/nye-lagos-thumb.jpg', 'NGN', 'Africa/Lagos', ARRAY['nye', 'countdown', 'party', 'celebration'], NOW(), NOW());

-- ============================================================================
-- TICKET TIERS
-- ============================================================================

-- Davido Lagos Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c'::uuid, 'General Admission', 'Standing area with great views of the stage', 15000.00, 30000, 12500, 150, '2026-01-15 00:00:00', '2026-03-15 18:00:00', 1, 10, true, NOW(), NOW()),
    (gen_random_uuid(), 'f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c'::uuid, 'VIP Standing', 'Premium standing area closer to the stage with dedicated bar access', 35000.00, 10000, 4200, 80, '2026-01-15 00:00:00', '2026-03-15 18:00:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c'::uuid, 'VVIP Table (4 persons)', 'Exclusive table seating with bottle service and premium food', 200000.00, 500, 320, 15, '2026-01-15 00:00:00', '2026-03-15 18:00:00', 4, 4, true, NOW(), NOW()),
    (gen_random_uuid(), 'f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c'::uuid, 'Early Bird General', 'Limited early bird tickets for general admission', 12000.00, 5000, 5000, 0, '2026-01-15 00:00:00', '2026-02-01 23:59:59', 1, 10, false, NOW(), NOW());

-- Davido Abuja Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'a7b8c9d0-e1f2-0a1b-4c5d-6e7f8a9b0c1d'::uuid, 'General Admission', 'Standing area with great views', 12000.00, 20000, 8500, 120, '2026-01-20 00:00:00', '2026-03-22 17:30:00', 1, 10, true, NOW(), NOW()),
    (gen_random_uuid(), 'a7b8c9d0-e1f2-0a1b-4c5d-6e7f8a9b0c1d'::uuid, 'VIP Standing', 'Premium standing with bar access', 28000.00, 8000, 3100, 60, '2026-01-20 00:00:00', '2026-03-22 17:30:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'a7b8c9d0-e1f2-0a1b-4c5d-6e7f8a9b0c1d'::uuid, 'VVIP Table (4 persons)', 'Exclusive table with bottle service', 180000.00, 400, 245, 10, '2026-01-20 00:00:00', '2026-03-22 17:30:00', 4, 4, true, NOW(), NOW());

-- Davido Port Harcourt Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'b8c9d0e1-f2a3-1b2c-5d6e-7f8a9b0c1d2e'::uuid, 'General Admission', 'Standing area', 10000.00, 18000, 6200, 95, '2026-02-01 00:00:00', '2026-04-05 18:00:00', 1, 10, true, NOW(), NOW()),
    (gen_random_uuid(), 'b8c9d0e1-f2a3-1b2c-5d6e-7f8a9b0c1d2e'::uuid, 'VIP Standing', 'Premium standing area', 25000.00, 6000, 2400, 45, '2026-02-01 00:00:00', '2026-04-05 18:00:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'b8c9d0e1-f2a3-1b2c-5d6e-7f8a9b0c1d2e'::uuid, 'VVIP Table (4 persons)', 'Exclusive table seating', 160000.00, 300, 180, 8, '2026-02-01 00:00:00', '2026-04-05 18:00:00', 4, 4, true, NOW(), NOW());

-- Davido Enugu Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'c9d0e1f2-a3b4-2c3d-6e7f-8a9b0c1d2e3f'::uuid, 'General Admission', 'Standing area', 10000.00, 15000, 5100, 75, '2026-02-05 00:00:00', '2026-04-12 17:00:00', 1, 10, true, NOW(), NOW()),
    (gen_random_uuid(), 'c9d0e1f2-a3b4-2c3d-6e7f-8a9b0c1d2e3f'::uuid, 'VIP Standing', 'Premium standing', 25000.00, 5000, 1850, 35, '2026-02-05 00:00:00', '2026-04-12 17:00:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'c9d0e1f2-a3b4-2c3d-6e7f-8a9b0c1d2e3f'::uuid, 'VVIP Table (4 persons)', 'Exclusive table', 150000.00, 250, 145, 6, '2026-02-05 00:00:00', '2026-04-12 17:00:00', 4, 4, true, NOW(), NOW());

-- Davido Ibadan Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'd0e1f2a3-b4c5-3d4e-7f8a-9b0c1d2e3f4a'::uuid, 'General Admission', 'Standing area', 12000.00, 16000, 5800, 85, '2026-02-10 00:00:00', '2026-04-26 18:00:00', 1, 10, true, NOW(), NOW()),
    (gen_random_uuid(), 'd0e1f2a3-b4c5-3d4e-7f8a-9b0c1d2e3f4a'::uuid, 'VIP Standing', 'Premium standing', 28000.00, 5500, 2100, 40, '2026-02-10 00:00:00', '2026-04-26 18:00:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'd0e1f2a3-b4c5-3d4e-7f8a-9b0c1d2e3f4a'::uuid, 'VVIP Table (4 persons)', 'Exclusive table', 170000.00, 300, 195, 7, '2026-02-10 00:00:00', '2026-04-26 18:00:00', 4, 4, true, NOW(), NOW());

-- Burna Boy Lagos Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'e1f2a3b4-c5d6-4e5f-8a9b-0c1d2e3f4a5b'::uuid, 'General Admission', 'Standing area', 18000.00, 25000, 10200, 140, '2026-03-01 00:00:00', '2026-05-20 19:00:00', 1, 10, true, NOW(), NOW()),
    (gen_random_uuid(), 'e1f2a3b4-c5d6-4e5f-8a9b-0c1d2e3f4a5b'::uuid, 'VIP Standing', 'Premium standing', 40000.00, 12000, 5100, 90, '2026-03-01 00:00:00', '2026-05-20 19:00:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'e1f2a3b4-c5d6-4e5f-8a9b-0c1d2e3f4a5b'::uuid, 'VVIP Table (6 persons)', 'Exclusive table', 300000.00, 600, 385, 18, '2026-03-01 00:00:00', '2026-05-20 19:00:00', 6, 6, true, NOW(), NOW());

-- Burna Boy Abuja Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'f2a3b4c5-d6e7-5f6a-9b0c-1d2e3f4a5b6c'::uuid, 'General Admission', 'Standing area', 15000.00, 18000, 7200, 110, '2026-03-15 00:00:00', '2026-06-10 18:00:00', 1, 10, true, NOW(), NOW()),
    (gen_random_uuid(), 'f2a3b4c5-d6e7-5f6a-9b0c-1d2e3f4a5b6c'::uuid, 'VIP Standing', 'Premium standing', 35000.00, 8000, 3400, 65, '2026-03-15 00:00:00', '2026-06-10 18:00:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'f2a3b4c5-d6e7-5f6a-9b0c-1d2e3f4a5b6c'::uuid, 'VVIP Table (6 persons)', 'Exclusive table', 280000.00, 450, 290, 12, '2026-03-15 00:00:00', '2026-06-10 18:00:00', 6, 6, true, NOW(), NOW());

-- Wizkid Lagos Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'a3b4c5d6-e7f8-6a7b-0c1d-2e3f4a5b6c7d'::uuid, 'General Admission', 'Standing area', 20000.00, 8000, 3200, 50, '2026-04-01 00:00:00', '2026-07-15 18:30:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'a3b4c5d6-e7f8-6a7b-0c1d-2e3f4a5b6c7d'::uuid, 'VIP Seated', 'Premium seated area', 50000.00, 4000, 1650, 35, '2026-04-01 00:00:00', '2026-07-15 18:30:00', 1, 6, true, NOW(), NOW()),
    (gen_random_uuid(), 'a3b4c5d6-e7f8-6a7b-0c1d-2e3f4a5b6c7d'::uuid, 'VVIP Table (4 persons)', 'Exclusive table', 250000.00, 200, 125, 8, '2026-04-01 00:00:00', '2026-07-15 18:30:00', 4, 4, true, NOW(), NOW());

-- Afrobeats Festival Ticket Tiers
INSERT INTO ticket_tiers (id, event_id, name, description, price, quantity, quantity_sold, quantity_reserved, sale_start, sale_end, min_per_order, max_per_order, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'b4c5d6e7-f8a9-7b8c-1d2e-3f4a5b6c7d8e'::uuid, '2-Day General Pass', '2-day access to all stages', 25000.00, 40000, 15200, 180, '2026-05-01 00:00:00', '2026-08-22 14:00:00', 1, 10, true, NOW(), NOW()),
    (gen_random_uuid(), 'b4c5d6e7-f8a9-7b8c-1d2e-3f4a5b6c7d8e'::uuid, '2-Day VIP Pass', '2-day VIP access with lounge', 60000.00, 15000, 6100, 95, '2026-05-01 00:00:00', '2026-08-22 14:00:00', 1, 8, true, NOW(), NOW()),
    (gen_random_uuid(), 'b4c5d6e7-f8a9-7b8c-1d2e-3f4a5b6c7d8e'::uuid, '2-Day VVIP Table (6 persons)', '2-day exclusive table', 400000.00, 800, 485, 22, '2026-05-01 00:00:00', '2026-08-22 14:00:00', 6, 6, true, NOW(), NOW());

-- ============================================================================
-- ADMIN USERS
-- ============================================================================

-- Password: Admin@123456 (hashed with bcrypt)
-- Note: Permissions are stored as JSONB array, not TEXT array
INSERT INTO admin_users (id, email, password_hash, first_name, last_name, role, permissions, is_active, login_attempts, must_change_password, two_factor_enabled, created_at, updated_at)
VALUES
    ('11111111-1111-1111-1111-111111111111'::uuid, 'admin@uduxpass.com', '$2b$10$nm1aqoTf4okL90HNuOLav.XmwF.KiFOqqxUHY8kLhDbvvjfCkDcM.', 'System', 'Administrator', 'super_admin', '["system_settings", "user_management", "admin_management", "event_create", "event_edit", "event_delete", "event_publish", "organizer_create", "organizer_edit", "organizer_delete", "organizer_approve", "order_view", "order_edit", "order_refund", "order_cancel", "payment_view", "payment_process", "payment_refund", "support_tickets", "customer_view", "customer_edit", "analytics_view", "reports_generate", "reports_export", "scanner_manage", "scanner_view"]'::jsonb, true, 0, false, false, NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222'::uuid, 'events@uduxpass.com', '$2b$10$nm1aqoTf4okL90HNuOLav.XmwF.KiFOqqxUHY8kLhDbvvjfCkDcM.', 'Events', 'Manager', 'event_manager', '["event_create", "event_edit", "event_delete", "event_publish", "organizer_create", "organizer_edit", "organizer_approve", "order_view", "order_edit", "analytics_view", "reports_generate", "scanner_view"]'::jsonb, true, 0, false, false, NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333'::uuid, 'support@uduxpass.com', '$2b$10$nm1aqoTf4okL90HNuOLav.XmwF.KiFOqqxUHY8kLhDbvvjfCkDcM.', 'Customer', 'Support', 'support', '["support_tickets", "customer_view", "customer_edit", "order_view", "order_edit", "order_refund", "order_cancel", "payment_view"]'::jsonb, true, 0, false, false, NOW(), NOW());

-- ============================================================================
-- SCANNER USERS
-- ============================================================================

-- Password: Scanner@123 (hashed with bcrypt)
INSERT INTO scanner_users (id, username, password_hash, first_name, last_name, email, phone, assigned_event_ids, role, status, is_active, created_at, updated_at)
VALUES
    ('44444444-4444-4444-4444-444444444444'::uuid, 'scanner_lagos_1', '$2b$12$G.9QoOWJnV1sT3bJj4qwEe8TSZPZu9K//MAdX88XqOt9l.wxMdFfy', 'John', 'Okafor', 'john.okafor@uduxpass.com', '+234-801-111-2222', ARRAY['f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c']::uuid[], 'scanner', 'active', true, NOW(), NOW()),
    ('55555555-5555-5555-5555-555555555555'::uuid, 'scanner_lagos_2', '$2b$12$G.9QoOWJnV1sT3bJj4qwEe8TSZPZu9K//MAdX88XqOt9l.wxMdFfy', 'Amina', 'Mohammed', 'amina.mohammed@uduxpass.com', '+234-802-222-3333', ARRAY['f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c', 'e1f2a3b4-c5d6-4e5f-8a9b-0c1d2e3f4a5b']::uuid[], 'scanner', 'active', true, NOW(), NOW()),
    ('66666666-6666-6666-6666-666666666666'::uuid, 'scanner_abuja_1', '$2b$12$G.9QoOWJnV1sT3bJj4qwEe8TSZPZu9K//MAdX88XqOt9l.wxMdFfy', 'Chidi', 'Nwosu', 'chidi.nwosu@uduxpass.com', '+234-803-333-4444', ARRAY['a7b8c9d0-e1f2-0a1b-4c5d-6e7f8a9b0c1d', 'f2a3b4c5-d6e7-5f6a-9b0c-1d2e3f4a5b6c']::uuid[], 'scanner', 'active', true, NOW(), NOW()),
    ('77777777-7777-7777-7777-777777777777'::uuid, 'scanner_supervisor', '$2b$12$G.9QoOWJnV1sT3bJj4qwEe8TSZPZu9K//MAdX88XqOt9l.wxMdFfy', 'Blessing', 'Adeyemi', 'blessing.adeyemi@uduxpass.com', '+234-804-444-5555', ARRAY[]::uuid[], 'supervisor', 'active', true, NOW(), NOW());

-- ============================================================================
-- REGULAR USERS
-- ============================================================================

-- Password: User@123 (hashed with bcrypt)
INSERT INTO users (id, email, password_hash, first_name, last_name, phone, auth_provider, email_verified, phone_verified, is_active, created_at, updated_at)
VALUES
    ('88888888-8888-8888-8888-888888888888'::uuid, 'adeola.williams@gmail.com', '$2b$12$fj2PsXNIFtv6sJb4WYp1t.6QIl2jKcOPFF7HhKoS6rf.Pic22.Vra', 'Adeola', 'Williams', '+234-901-111-1111', 'email', true, true, true, NOW(), NOW()),
    ('99999999-9999-9999-9999-999999999999'::uuid, 'tunde.bakare@yahoo.com', '$2b$12$fj2PsXNIFtv6sJb4WYp1t.6QIl2jKcOPFF7HhKoS6rf.Pic22.Vra', 'Tunde', 'Bakare', '+234-902-222-2222', 'email', true, false, true, NOW(), NOW()),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, 'chioma.okonkwo@outlook.com', '$2b$12$fj2PsXNIFtv6sJb4WYp1t.6QIl2jKcOPFF7HhKoS6rf.Pic22.Vra', 'Chioma', 'Okonkwo', '+234-903-333-3333', 'email', true, true, true, NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'::uuid, 'yusuf.ibrahim@gmail.com', '$2b$12$fj2PsXNIFtv6sJb4WYp1t.6QIl2jKcOPFF7HhKoS6rf.Pic22.Vra', 'Yusuf', 'Ibrahim', '+234-904-444-4444', 'email', true, true, true, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc'::uuid, 'ngozi.eze@hotmail.com', '$2b$12$fj2PsXNIFtv6sJb4WYp1t.6QIl2jKcOPFF7HhKoS6rf.Pic22.Vra', 'Ngozi', 'Eze', '+234-905-555-5555', 'email', false, false, true, NOW(), NOW());

-- ============================================================================
-- NOTES
-- ============================================================================

/*
SEED DATA SUMMARY:
- 3 Organizers (Afro Nation, Lagos Live, Naija Vibes)
- 2 Tours (Davido Timeless, Burna Boy African Giant)
- 10 Events (5 Davido, 2 Burna Boy, 3 standalone)
- 30+ Ticket Tiers across all events
- 3 Admin Users (super_admin, admin, support)
- 4 Scanner Users (3 scanners, 1 supervisor)
- 5 Regular Users

DEFAULT PASSWORDS (for testing only):
- Admin Users: Admin@123456
- Scanner Users: Scanner@123
- Regular Users: User@123

NOTE: In production, these passwords should be changed immediately and proper
password hashing should be verified.

TICKET SALES STATISTICS:
- Total tickets sold across all events: ~100,000+
- Total tickets reserved: ~1,500+
- Revenue generated: ~₦2.5 billion+

All events are published except the NYE event which is in draft status.
*/
