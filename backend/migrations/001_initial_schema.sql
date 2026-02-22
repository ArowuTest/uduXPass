-- Migration: 001_initial_schema.sql
-- Description: Create initial database schema for uduXPass ticketing platform
-- Created: 2024-01-01

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable JSONB operators
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- Create custom types
CREATE TYPE auth_provider AS ENUM ('email', 'momo');
CREATE TYPE event_status AS ENUM ('draft', 'published', 'on_sale', 'sold_out', 'cancelled', 'completed');
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'expired', 'cancelled', 'refunded');
CREATE TYPE payment_method AS ENUM ('momo', 'paystack');
CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'failed', 'cancelled', 'refunded');
CREATE TYPE ticket_status AS ENUM ('active', 'redeemed', 'voided');
CREATE TYPE otp_purpose AS ENUM ('login', 'registration', 'password_reset');

-- Organizers table
CREATE TABLE organizers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    website_url VARCHAR(500),
    logo_url VARCHAR(500),
    description TEXT,
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100) NOT NULL DEFAULT 'Nigeria',
    is_active BOOLEAN NOT NULL DEFAULT true,
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Tours table
CREATE TABLE tours (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organizer_id UUID NOT NULL REFERENCES organizers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    artist_name VARCHAR(255) NOT NULL,
    description TEXT,
    tour_image_url VARCHAR(500),
    start_date DATE,
    end_date DATE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(organizer_id, slug)
);

-- Events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organizer_id UUID NOT NULL REFERENCES organizers(id) ON DELETE CASCADE,
    tour_id UUID REFERENCES tours(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL,
    doors_open TIMESTAMP WITH TIME ZONE,
    venue_name VARCHAR(255) NOT NULL,
    venue_address TEXT NOT NULL,
    venue_city VARCHAR(100) NOT NULL,
    venue_state VARCHAR(100),
    venue_country VARCHAR(100) NOT NULL DEFAULT 'Nigeria',
    venue_capacity INTEGER,
    venue_latitude DECIMAL(10, 8),
    venue_longitude DECIMAL(11, 8),
    event_image_url VARCHAR(500),
    status event_status NOT NULL DEFAULT 'draft',
    sale_start TIMESTAMP WITH TIME ZONE,
    sale_end TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(organizer_id, slug)
);

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255),
    phone VARCHAR(20),
    password_hash VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    auth_provider auth_provider NOT NULL,
    momo_id VARCHAR(100),
    email_verified BOOLEAN NOT NULL DEFAULT false,
    phone_verified BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login TIMESTAMP WITH TIME ZONE,
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(email),
    UNIQUE(phone),
    UNIQUE(momo_id)
);

-- Ticket tiers table
CREATE TABLE ticket_tiers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'NGN',
    quota INTEGER,
    max_per_order INTEGER NOT NULL DEFAULT 10,
    min_per_order INTEGER NOT NULL DEFAULT 1,
    sale_start TIMESTAMP WITH TIME ZONE,
    sale_end TIMESTAMP WITH TIME ZONE,
    position INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    settings JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Orders table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(16) NOT NULL UNIQUE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    status order_status NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'NGN',
    payment_method payment_method,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    secret VARCHAR(128) NOT NULL,
    locale VARCHAR(10) NOT NULL DEFAULT 'en',
    comment TEXT,
    meta_info JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Order lines table
CREATE TABLE order_lines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    ticket_tier_id UUID NOT NULL REFERENCES ticket_tiers(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Tickets table
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_line_id UUID NOT NULL REFERENCES order_lines(id) ON DELETE CASCADE,
    serial_number VARCHAR(50) NOT NULL UNIQUE,
    qr_code_data VARCHAR(500) NOT NULL UNIQUE,
    status ticket_status NOT NULL DEFAULT 'active',
    redeemed_at TIMESTAMP WITH TIME ZONE,
    redeemed_by VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    provider payment_method NOT NULL,
    provider_transaction_id VARCHAR(255),
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'NGN',
    status payment_status NOT NULL DEFAULT 'pending',
    provider_response JSONB NOT NULL DEFAULT '{}',
    webhook_received_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Inventory holds table
CREATE TABLE inventory_holds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_tier_id UUID NOT NULL REFERENCES ticket_tiers(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- OTP tokens table
CREATE TABLE otp_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(20) NOT NULL,
    token VARCHAR(6) NOT NULL,
    purpose otp_purpose NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    attempts INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_organizers_slug ON organizers(slug);
CREATE INDEX idx_organizers_email ON organizers(email);
CREATE INDEX idx_organizers_active ON organizers(is_active);

CREATE INDEX idx_tours_organizer_id ON tours(organizer_id);
CREATE INDEX idx_tours_slug ON tours(organizer_id, slug);
CREATE INDEX idx_tours_active ON tours(is_active);
CREATE INDEX idx_tours_dates ON tours(start_date, end_date);

CREATE INDEX idx_events_organizer_id ON events(organizer_id);
CREATE INDEX idx_events_tour_id ON events(tour_id);
CREATE INDEX idx_events_slug ON events(organizer_id, slug);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_active ON events(is_active);
CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_city ON events(venue_city);
CREATE INDEX idx_events_country ON events(venue_country);
CREATE INDEX idx_events_sale_period ON events(sale_start, sale_end);
CREATE INDEX idx_events_public ON events(status, is_active, event_date) WHERE status IN ('published', 'on_sale');

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_momo_id ON users(momo_id);
CREATE INDEX idx_users_auth_provider ON users(auth_provider);
CREATE INDEX idx_users_active ON users(is_active);

CREATE INDEX idx_ticket_tiers_event_id ON ticket_tiers(event_id);
CREATE INDEX idx_ticket_tiers_active ON ticket_tiers(is_active);
CREATE INDEX idx_ticket_tiers_position ON ticket_tiers(event_id, position);
CREATE INDEX idx_ticket_tiers_sale_period ON ticket_tiers(sale_start, sale_end);

CREATE INDEX idx_orders_code ON orders(code);
CREATE INDEX idx_orders_secret ON orders(secret);
CREATE INDEX idx_orders_event_id ON orders(event_id);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_email ON orders(email);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_expires_at ON orders(expires_at);
CREATE INDEX idx_orders_created_at ON orders(created_at);

CREATE INDEX idx_order_lines_order_id ON order_lines(order_id);
CREATE INDEX idx_order_lines_ticket_tier_id ON order_lines(ticket_tier_id);

CREATE INDEX idx_tickets_order_line_id ON tickets(order_line_id);
CREATE INDEX idx_tickets_serial_number ON tickets(serial_number);
CREATE INDEX idx_tickets_qr_code_data ON tickets(qr_code_data);
CREATE INDEX idx_tickets_status ON tickets(status);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_provider ON payments(provider);
CREATE INDEX idx_payments_provider_transaction_id ON payments(provider, provider_transaction_id);
CREATE INDEX idx_payments_status ON payments(status);

CREATE INDEX idx_inventory_holds_ticket_tier_id ON inventory_holds(ticket_tier_id);
CREATE INDEX idx_inventory_holds_session_id ON inventory_holds(session_id);
CREATE INDEX idx_inventory_holds_expires_at ON inventory_holds(expires_at);

CREATE INDEX idx_otp_tokens_phone_purpose ON otp_tokens(phone, purpose);
CREATE INDEX idx_otp_tokens_expires_at ON otp_tokens(expires_at);

-- Create triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_organizers_updated_at BEFORE UPDATE ON organizers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tours_updated_at BEFORE UPDATE ON tours FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ticket_tiers_updated_at BEFORE UPDATE ON ticket_tiers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tickets_updated_at BEFORE UPDATE ON tickets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create constraints
ALTER TABLE users ADD CONSTRAINT users_auth_provider_check 
    CHECK (
        (auth_provider = 'email' AND email IS NOT NULL AND password_hash IS NOT NULL) OR
        (auth_provider = 'momo' AND phone IS NOT NULL AND momo_id IS NOT NULL)
    );

ALTER TABLE ticket_tiers ADD CONSTRAINT ticket_tiers_order_limits_check 
    CHECK (min_per_order <= max_per_order AND min_per_order > 0);

ALTER TABLE ticket_tiers ADD CONSTRAINT ticket_tiers_quota_check 
    CHECK (quota IS NULL OR quota > 0);

ALTER TABLE order_lines ADD CONSTRAINT order_lines_quantity_check 
    CHECK (quantity > 0);

ALTER TABLE order_lines ADD CONSTRAINT order_lines_price_check 
    CHECK (unit_price >= 0 AND total_price >= 0);

ALTER TABLE inventory_holds ADD CONSTRAINT inventory_holds_quantity_check 
    CHECK (quantity > 0);

ALTER TABLE otp_tokens ADD CONSTRAINT otp_tokens_attempts_check 
    CHECK (attempts >= 0 AND attempts <= 10);

-- Create functions for business logic
CREATE OR REPLACE FUNCTION get_ticket_tier_availability(tier_id UUID)
RETURNS TABLE(
    available_quantity INTEGER,
    total_sold INTEGER,
    total_held INTEGER
) AS $$
BEGIN
    RETURN QUERY
    WITH tier_info AS (
        SELECT quota FROM ticket_tiers WHERE id = tier_id
    ),
    sold_count AS (
        SELECT COALESCE(SUM(ol.quantity), 0)::INTEGER as sold
        FROM order_lines ol
        JOIN orders o ON ol.order_id = o.id
        WHERE ol.ticket_tier_id = tier_id 
        AND o.status = 'paid'
    ),
    held_count AS (
        SELECT COALESCE(SUM(quantity), 0)::INTEGER as held
        FROM inventory_holds
        WHERE ticket_tier_id = tier_id 
        AND expires_at > NOW()
    )
    SELECT 
        CASE 
            WHEN ti.quota IS NULL THEN 999999
            ELSE GREATEST(0, ti.quota - sc.sold - hc.held)
        END as available_quantity,
        sc.sold as total_sold,
        hc.held as total_held
    FROM tier_info ti, sold_count sc, held_count hc;
END;
$$ LANGUAGE plpgsql;

-- Create function to clean up expired data
CREATE OR REPLACE FUNCTION cleanup_expired_data()
RETURNS INTEGER AS $$
DECLARE
    expired_count INTEGER := 0;
BEGIN
    -- Clean up expired inventory holds
    DELETE FROM inventory_holds WHERE expires_at <= NOW();
    GET DIAGNOSTICS expired_count = ROW_COUNT;
    
    -- Clean up expired OTP tokens
    DELETE FROM otp_tokens WHERE expires_at <= NOW();
    
    -- Mark expired orders
    UPDATE orders 
    SET status = 'expired', updated_at = NOW()
    WHERE status = 'pending' AND expires_at <= NOW();
    
    RETURN expired_count;
END;
$$ LANGUAGE plpgsql;

-- Create views for common queries
CREATE VIEW public_events AS
SELECT 
    e.id,
    e.name,
    e.slug,
    e.description,
    e.event_date,
    e.doors_open,
    e.venue_name,
    e.venue_address,
    e.venue_city,
    e.venue_state,
    e.venue_country,
    e.venue_capacity,
    e.venue_latitude,
    e.venue_longitude,
    e.event_image_url,
    e.status,
    e.sale_start,
    e.sale_end,
    e.created_at,
    o.name as organizer_name,
    o.slug as organizer_slug,
    t.name as tour_name,
    t.artist_name
FROM events e
JOIN organizers o ON e.organizer_id = o.id
LEFT JOIN tours t ON e.tour_id = t.id
WHERE e.status IN ('published', 'on_sale') 
AND e.is_active = true 
AND o.is_active = true;

-- Insert initial data
INSERT INTO organizers (name, slug, email, description, country) VALUES
('uduXPass Demo', 'uduxpass-demo', 'demo@uduxpass.com', 'Demo organizer for testing', 'Nigeria');

-- Create admin user for testing
INSERT INTO users (email, password_hash, first_name, last_name, auth_provider, email_verified, is_active) VALUES
('admin@uduxpass.com', 'hashed_admin123', 'Admin', 'User', 'email', true, true);

COMMIT;

