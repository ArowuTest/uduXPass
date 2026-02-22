-- Insert ticket tiers for the test event
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
        quota,
        sale_start,
        sale_end,
        position,
        is_active
    ) VALUES
    (uuid_generate_v4(), event_id, 'VIP', 'VIP seating with exclusive access', 50000.00, 'NGN', 500, CURRENT_TIMESTAMP, '2026-03-15 19:00:00+01', 1, true),
    (uuid_generate_v4(), event_id, 'Regular', 'General admission', 15000.00, 'NGN', 5000, CURRENT_TIMESTAMP, '2026-03-15 19:00:00+01', 2, true),
    (uuid_generate_v4(), event_id, 'Early Bird', 'Early bird special pricing', 10000.00, 'NGN', 1000, CURRENT_TIMESTAMP, '2026-02-28 23:59:59+01', 3, true);
END $$;
