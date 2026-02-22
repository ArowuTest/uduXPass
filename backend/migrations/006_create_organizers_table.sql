-- Migration: Create organizers table
-- Description: Creates the organizers table to support event organizer management
-- Version: 006
-- Date: 2026-02-17

-- Create organizers table
CREATE TABLE IF NOT EXISTS organizers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    website VARCHAR(500),
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    logo_url TEXT,
    banner_url TEXT,
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    postal_code VARCHAR(20),
    timezone VARCHAR(100) DEFAULT 'Africa/Lagos',
    currency VARCHAR(3) DEFAULT 'NGN',
    settings JSONB DEFAULT '{}',
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_organizers_slug ON organizers(slug);
CREATE INDEX IF NOT EXISTS idx_organizers_email ON organizers(email);
CREATE INDEX IF NOT EXISTS idx_organizers_is_active ON organizers(is_active);
CREATE INDEX IF NOT EXISTS idx_organizers_is_verified ON organizers(is_verified);
CREATE INDEX IF NOT EXISTS idx_organizers_created_at ON organizers(created_at);

-- Add foreign key constraint to events table (if not already exists)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'events_organizer_id_fkey'
    ) THEN
        ALTER TABLE events 
        ADD CONSTRAINT events_organizer_id_fkey 
        FOREIGN KEY (organizer_id) REFERENCES organizers(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Add foreign key constraint to tours table (if not already exists)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'tours_organizer_id_fkey'
    ) THEN
        ALTER TABLE tours 
        ADD CONSTRAINT tours_organizer_id_fkey 
        FOREIGN KEY (organizer_id) REFERENCES organizers(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Insert seed data for testing
INSERT INTO organizers (
    id,
    name,
    slug,
    description,
    website,
    email,
    phone,
    timezone,
    currency,
    is_verified,
    is_active
) VALUES
(
    '11111111-1111-1111-1111-111111111111',
    'uduXPass Events',
    'uduxpass-events',
    'Official uduXPass event organizer for testing and demo events',
    'https://uduxpass.com',
    'events@uduxpass.com',
    '+2348012345678',
    'Africa/Lagos',
    'NGN',
    true,
    true
),
(
    '22222222-2222-2222-2222-222222222222',
    'Lagos Entertainment Ltd',
    'lagos-entertainment',
    'Premier entertainment and event management company in Lagos',
    'https://lagosent.com',
    'info@lagosent.com',
    '+2348098765432',
    'Africa/Lagos',
    'NGN',
    true,
    true
),
(
    '33333333-3333-3333-3333-333333333333',
    'Abuja Live Events',
    'abuja-live-events',
    'Leading event organizer in Abuja specializing in concerts and festivals',
    'https://abujalive.com',
    'contact@abujalive.com',
    '+2349012345678',
    'Africa/Lagos',
    'NGN',
    true,
    true
)
ON CONFLICT (id) DO NOTHING;

-- Grant permissions to uduxpass_user
GRANT SELECT, INSERT, UPDATE, DELETE ON organizers TO uduxpass_user;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO uduxpass_user;

-- Add comment to table
COMMENT ON TABLE organizers IS 'Stores information about event organizers who can create and manage events';
COMMENT ON COLUMN organizers.slug IS 'URL-friendly unique identifier for the organizer';
COMMENT ON COLUMN organizers.is_verified IS 'Indicates if the organizer has been verified by platform administrators';
COMMENT ON COLUMN organizers.settings IS 'JSON object storing organizer-specific settings and preferences';
