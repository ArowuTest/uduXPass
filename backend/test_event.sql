-- Insert test event with category
INSERT INTO events (
    id,
    organizer_id,
    category_id,
    name,
    slug,
    description,
    event_date,
    doors_open,
    venue_name,
    venue_address,
    venue_city,
    venue_state,
    venue_country,
    venue_capacity,
    event_image_url,
    status,
    sale_start,
    sale_end,
    is_active
) VALUES (
    uuid_generate_v4(),
    '56c5c757-eed7-43e2-8eaf-ed5c92154075',
    '1f8085f8-b735-47df-a39b-b379142cd0e6',
    'Burna Boy Live in Lagos',
    'burna-boy-live-lagos-2026',
    'Experience an unforgettable night with Grammy-winning artist Burna Boy live in concert at Eko Atlantic.',
    '2026-03-15 20:00:00+01',
    '2026-03-15 18:00:00+01',
    'Eko Atlantic Energy City',
    'Plot 1, Eko Atlantic City',
    'Lagos',
    'Lagos State',
    'Nigeria',
    50000,
    'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=800',
    'published',
    CURRENT_TIMESTAMP,
    '2026-03-15 19:00:00+01',
    true
) RETURNING id;

-- Get the event ID for ticket tiers
DO $$
DECLARE
    event_id UUID;
BEGIN
    SELECT id INTO event_id FROM events WHERE slug = 'burna-boy-live-lagos-2026';
    
    -- Insert ticket tiers
    INSERT INTO ticket_tiers (
        id,
        event_id,
        name,
        description,
        price,
        currency,
        quantity_total,
        quantity_available,
        sale_start,
        sale_end,
        is_active
    ) VALUES
    (uuid_generate_v4(), event_id, 'VIP', 'VIP seating with exclusive access', 50000.00, 'NGN', 500, 500, CURRENT_TIMESTAMP, '2026-03-15 19:00:00+01', true),
    (uuid_generate_v4(), event_id, 'Regular', 'General admission', 15000.00, 'NGN', 5000, 5000, CURRENT_TIMESTAMP, '2026-03-15 19:00:00+01', true),
    (uuid_generate_v4(), event_id, 'Early Bird', 'Early bird special pricing', 10000.00, 'NGN', 1000, 1000, CURRENT_TIMESTAMP, '2026-02-28 23:59:59+01', true);
END $$;
