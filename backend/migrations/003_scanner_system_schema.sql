-- Migration: 003_scanner_system_schema.sql
-- Description: Create scanner authentication and management system
-- Created: 2025-09-16

BEGIN;

-- Create scanner role enum
CREATE TYPE scanner_role AS ENUM (
    'scanner_operator',
    'lead_scanner', 
    'scanner_supervisor'
);

-- Create scanner status enum
CREATE TYPE scanner_status AS ENUM (
    'active',
    'inactive',
    'locked',
    'suspended'
);

-- Scanner users table
CREATE TABLE scanner_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role scanner_role NOT NULL DEFAULT 'scanner_operator',
    permissions JSONB NOT NULL DEFAULT '[]',
    status scanner_status NOT NULL DEFAULT 'active',
    last_login TIMESTAMP WITH TIME ZONE,
    login_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    must_change_password BOOLEAN NOT NULL DEFAULT true,
    created_by UUID REFERENCES admin_users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Scanner event assignments table
CREATE TABLE scanner_event_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scanner_id UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    assigned_by UUID NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT true,
    UNIQUE(scanner_id, event_id)
);

-- Scanner sessions table
CREATE TABLE scanner_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scanner_id UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    end_time TIMESTAMP WITH TIME ZONE,
    scans_count INTEGER NOT NULL DEFAULT 0,
    valid_scans INTEGER NOT NULL DEFAULT 0,
    invalid_scans INTEGER NOT NULL DEFAULT 0,
    total_revenue DECIMAL(10, 2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    notes TEXT
);

-- Scanner audit log table
CREATE TABLE scanner_audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scanner_id UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES scanner_sessions(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    details JSONB NOT NULL DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Scanner login history table
CREATE TABLE scanner_login_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scanner_id UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN NOT NULL,
    login_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    logout_at TIMESTAMP WITH TIME ZONE
);

-- Ticket validations table (for tracking scanned tickets)
CREATE TABLE ticket_validations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    scanner_id UUID NOT NULL REFERENCES scanner_users(id) ON DELETE CASCADE,
    session_id UUID NOT NULL REFERENCES scanner_sessions(id) ON DELETE CASCADE,
    validation_result VARCHAR(20) NOT NULL CHECK (validation_result IN ('valid', 'invalid', 'duplicate', 'emergency_override')),
    validation_timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    notes TEXT,
    UNIQUE(ticket_id) -- Prevent duplicate validations
);

-- Create indexes for performance
CREATE INDEX idx_scanner_users_username ON scanner_users(username);
CREATE INDEX idx_scanner_users_email ON scanner_users(email);
CREATE INDEX idx_scanner_users_role ON scanner_users(role);
CREATE INDEX idx_scanner_users_status ON scanner_users(status);
CREATE INDEX idx_scanner_users_last_login ON scanner_users(last_login);
CREATE INDEX idx_scanner_users_created_at ON scanner_users(created_at);

CREATE INDEX idx_scanner_event_assignments_scanner_id ON scanner_event_assignments(scanner_id);
CREATE INDEX idx_scanner_event_assignments_event_id ON scanner_event_assignments(event_id);
CREATE INDEX idx_scanner_event_assignments_is_active ON scanner_event_assignments(is_active);
CREATE INDEX idx_scanner_event_assignments_assigned_at ON scanner_event_assignments(assigned_at);

CREATE INDEX idx_scanner_sessions_scanner_id ON scanner_sessions(scanner_id);
CREATE INDEX idx_scanner_sessions_event_id ON scanner_sessions(event_id);
CREATE INDEX idx_scanner_sessions_start_time ON scanner_sessions(start_time);
CREATE INDEX idx_scanner_sessions_is_active ON scanner_sessions(is_active);

CREATE INDEX idx_scanner_audit_log_scanner_id ON scanner_audit_log(scanner_id);
CREATE INDEX idx_scanner_audit_log_session_id ON scanner_audit_log(session_id);
CREATE INDEX idx_scanner_audit_log_action ON scanner_audit_log(action);
CREATE INDEX idx_scanner_audit_log_created_at ON scanner_audit_log(created_at);

CREATE INDEX idx_scanner_login_history_scanner_id ON scanner_login_history(scanner_id);
CREATE INDEX idx_scanner_login_history_login_at ON scanner_login_history(login_at);
CREATE INDEX idx_scanner_login_history_success ON scanner_login_history(success);

CREATE INDEX idx_ticket_validations_ticket_id ON ticket_validations(ticket_id);
CREATE INDEX idx_ticket_validations_scanner_id ON ticket_validations(scanner_id);
CREATE INDEX idx_ticket_validations_session_id ON ticket_validations(session_id);
CREATE INDEX idx_ticket_validations_validation_timestamp ON ticket_validations(validation_timestamp);
CREATE INDEX idx_ticket_validations_validation_result ON ticket_validations(validation_result);

-- Create triggers for updated_at timestamps
CREATE TRIGGER update_scanner_users_updated_at 
    BEFORE UPDATE ON scanner_users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to create default scanner permissions based on role
CREATE OR REPLACE FUNCTION get_default_scanner_permissions(p_role scanner_role)
RETURNS JSONB AS $$
BEGIN
    CASE p_role
        WHEN 'scanner_operator' THEN
            RETURN '["scan_tickets", "manual_entry"]'::JSONB;
        WHEN 'lead_scanner' THEN
            RETURN '["scan_tickets", "manual_entry", "bulk_scan", "view_reports"]'::JSONB;
        WHEN 'scanner_supervisor' THEN
            RETURN '["scan_tickets", "manual_entry", "emergency_override", "bulk_scan", "view_reports", "manage_settings"]'::JSONB;
        ELSE
            RETURN '["scan_tickets"]'::JSONB;
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- Function to record scanner login
CREATE OR REPLACE FUNCTION record_scanner_login(
    p_scanner_id UUID,
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL,
    p_success BOOLEAN DEFAULT true
)
RETURNS VOID AS $$
BEGIN
    INSERT INTO scanner_login_history (
        scanner_id,
        ip_address,
        user_agent,
        success,
        login_at
    ) VALUES (
        p_scanner_id,
        p_ip_address,
        p_user_agent,
        p_success,
        NOW()
    );
    
    -- Update last_login if successful
    IF p_success THEN
        UPDATE scanner_users 
        SET 
            last_login = NOW(),
            login_attempts = 0,
            locked_until = NULL,
            updated_at = NOW()
        WHERE id = p_scanner_id;
    ELSE
        -- Increment login attempts and potentially lock account
        UPDATE scanner_users 
        SET 
            login_attempts = login_attempts + 1,
            locked_until = CASE 
                WHEN login_attempts + 1 >= 5 THEN NOW() + INTERVAL '30 minutes'
                ELSE locked_until
            END,
            updated_at = NOW()
        WHERE id = p_scanner_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Function to log scanner activity
CREATE OR REPLACE FUNCTION log_scanner_activity(
    p_scanner_id UUID,
    p_action VARCHAR(100),
    p_session_id UUID DEFAULT NULL,
    p_resource_type VARCHAR(50) DEFAULT NULL,
    p_resource_id UUID DEFAULT NULL,
    p_details JSONB DEFAULT NULL,
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL
)
RETURNS VOID AS $$
BEGIN
    INSERT INTO scanner_audit_log (
        scanner_id,
        session_id,
        action,
        resource_type,
        resource_id,
        details,
        ip_address,
        user_agent,
        created_at
    ) VALUES (
        p_scanner_id,
        p_session_id,
        p_action,
        p_resource_type,
        p_resource_id,
        COALESCE(p_details, '{}'::JSONB),
        p_ip_address,
        p_user_agent,
        NOW()
    );
END;
$$ LANGUAGE plpgsql;

-- Function to start scanner session
CREATE OR REPLACE FUNCTION start_scanner_session(
    p_scanner_id UUID,
    p_event_id UUID
)
RETURNS UUID AS $$
DECLARE
    session_id UUID;
BEGIN
    -- End any active sessions for this scanner
    UPDATE scanner_sessions 
    SET 
        end_time = NOW(),
        is_active = false
    WHERE scanner_id = p_scanner_id AND is_active = true;
    
    -- Create new session
    INSERT INTO scanner_sessions (
        scanner_id,
        event_id,
        start_time,
        is_active
    ) VALUES (
        p_scanner_id,
        p_event_id,
        NOW(),
        true
    ) RETURNING id INTO session_id;
    
    -- Log the activity
    PERFORM log_scanner_activity(
        p_scanner_id,
        'session_start',
        session_id,
        'scanner_session',
        session_id,
        jsonb_build_object('event_id', p_event_id)
    );
    
    RETURN session_id;
END;
$$ LANGUAGE plpgsql;

-- Function to end scanner session
CREATE OR REPLACE FUNCTION end_scanner_session(p_session_id UUID)
RETURNS VOID AS $$
DECLARE
    scanner_id UUID;
BEGIN
    -- Get scanner_id and update session
    UPDATE scanner_sessions 
    SET 
        end_time = NOW(),
        is_active = false
    WHERE id = p_session_id AND is_active = true
    RETURNING scanner_id INTO scanner_id;
    
    -- Log the activity if session was found
    IF scanner_id IS NOT NULL THEN
        PERFORM log_scanner_activity(
            scanner_id,
            'session_end',
            p_session_id,
            'scanner_session',
            p_session_id
        );
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Function to validate ticket and record scan
CREATE OR REPLACE FUNCTION validate_and_record_ticket(
    p_ticket_id UUID,
    p_scanner_id UUID,
    p_session_id UUID,
    p_validation_result VARCHAR(20),
    p_notes TEXT DEFAULT NULL
)
RETURNS BOOLEAN AS $$
DECLARE
    ticket_exists BOOLEAN := false;
    already_validated BOOLEAN := false;
    ticket_price DECIMAL(10, 2) := 0;
BEGIN
    -- Check if ticket exists and get price
    SELECT 
        EXISTS(SELECT 1 FROM tickets WHERE id = p_ticket_id AND status = 'active'),
        COALESCE(tt.price, 0)
    INTO ticket_exists, ticket_price
    FROM tickets t
    JOIN order_lines ol ON t.order_line_id = ol.id
    JOIN ticket_tiers tt ON ol.ticket_tier_id = tt.id
    WHERE t.id = p_ticket_id;
    
    -- Check if already validated
    SELECT EXISTS(SELECT 1 FROM ticket_validations WHERE ticket_id = p_ticket_id)
    INTO already_validated;
    
    -- If ticket doesn't exist or already validated, return false
    IF NOT ticket_exists OR already_validated THEN
        RETURN false;
    END IF;
    
    -- Record the validation
    INSERT INTO ticket_validations (
        ticket_id,
        scanner_id,
        session_id,
        validation_result,
        notes
    ) VALUES (
        p_ticket_id,
        p_scanner_id,
        p_session_id,
        p_validation_result,
        p_notes
    );
    
    -- Update ticket status if valid
    IF p_validation_result = 'valid' THEN
        UPDATE tickets 
        SET 
            status = 'redeemed',
            redeemed_at = NOW(),
            updated_at = NOW()
        WHERE id = p_ticket_id;
    END IF;
    
    -- Update session statistics
    UPDATE scanner_sessions
    SET 
        scans_count = scans_count + 1,
        valid_scans = CASE WHEN p_validation_result = 'valid' THEN valid_scans + 1 ELSE valid_scans END,
        invalid_scans = CASE WHEN p_validation_result = 'invalid' THEN invalid_scans + 1 ELSE invalid_scans END,
        total_revenue = CASE WHEN p_validation_result = 'valid' THEN total_revenue + ticket_price ELSE total_revenue END
    WHERE id = p_session_id;
    
    -- Log the activity
    PERFORM log_scanner_activity(
        p_scanner_id,
        'ticket_validation',
        p_session_id,
        'ticket',
        p_ticket_id,
        jsonb_build_object(
            'validation_result', p_validation_result,
            'ticket_price', ticket_price,
            'notes', p_notes
        )
    );
    
    RETURN true;
END;
$$ LANGUAGE plpgsql;

-- Function to get scanner assigned events
CREATE OR REPLACE FUNCTION get_scanner_assigned_events(p_scanner_id UUID)
RETURNS TABLE(
    event_id UUID,
    event_name VARCHAR(255),
    event_date TIMESTAMP WITH TIME ZONE,
    venue_name VARCHAR(255),
    venue_city VARCHAR(100),
    status event_status,
    assigned_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.name,
        e.event_date,
        e.venue_name,
        e.venue_city,
        e.status,
        sea.assigned_at
    FROM scanner_event_assignments sea
    JOIN events e ON sea.event_id = e.id
    WHERE sea.scanner_id = p_scanner_id 
    AND sea.is_active = true
    AND e.is_active = true
    ORDER BY e.event_date ASC;
END;
$$ LANGUAGE plpgsql;

-- Add constraints
ALTER TABLE scanner_users ADD CONSTRAINT chk_scanner_users_username_length 
    CHECK (LENGTH(TRIM(username)) >= 3 AND LENGTH(TRIM(username)) <= 50);

ALTER TABLE scanner_users ADD CONSTRAINT chk_scanner_users_name_not_empty 
    CHECK (LENGTH(TRIM(name)) > 0);

ALTER TABLE scanner_users ADD CONSTRAINT chk_scanner_users_email_format 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE scanner_users ADD CONSTRAINT chk_scanner_users_login_attempts 
    CHECK (login_attempts >= 0 AND login_attempts <= 10);

-- Create default scanner users for testing
INSERT INTO scanner_users (
    username,
    password_hash,
    name,
    email,
    role,
    permissions,
    status,
    must_change_password,
    created_by
) VALUES 
(
    'scanner1',
    '$argon2id$v=19$m=65536,t=3,p=4$YWJjZGVmZ2hpams$+DhQmGspZct+3jqQMz+QOQ',  -- Password: Scanner123!
    'Main Gate Scanner',
    'scanner1@uduxpass.com',
    'scanner_operator',
    get_default_scanner_permissions('scanner_operator'),
    'active',
    false,
    (SELECT id FROM admin_users WHERE email = 'admin@uduxpass.com' LIMIT 1)
),
(
    'scanner2',
    '$argon2id$v=19$m=65536,t=3,p=4$YWJjZGVmZ2hpams$+DhQmGspZct+3jqQMz+QOQ',  -- Password: Scanner123!
    'VIP Gate Scanner',
    'scanner2@uduxpass.com',
    'lead_scanner',
    get_default_scanner_permissions('lead_scanner'),
    'active',
    false,
    (SELECT id FROM admin_users WHERE email = 'admin@uduxpass.com' LIMIT 1)
),
(
    'supervisor1',
    '$argon2id$v=19$m=65536,t=3,p=4$YWJjZGVmZ2hpams$+DhQmGspZct+3jqQMz+QOQ',  -- Password: Scanner123!
    'Scanner Supervisor',
    'supervisor1@uduxpass.com',
    'scanner_supervisor',
    get_default_scanner_permissions('scanner_supervisor'),
    'active',
    false,
    (SELECT id FROM admin_users WHERE email = 'admin@uduxpass.com' LIMIT 1)
);

-- Add comments
COMMENT ON TABLE scanner_users IS 'Scanner operators with role-based access control';
COMMENT ON TABLE scanner_event_assignments IS 'Assignment of scanners to specific events';
COMMENT ON TABLE scanner_sessions IS 'Active scanning sessions for tracking performance';
COMMENT ON TABLE scanner_audit_log IS 'Audit trail for all scanner activities';
COMMENT ON TABLE scanner_login_history IS 'Login history for scanner users';
COMMENT ON TABLE ticket_validations IS 'Record of all ticket validation attempts';

COMMIT;

