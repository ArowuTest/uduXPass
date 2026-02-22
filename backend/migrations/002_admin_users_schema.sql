-- Description: Create admin users and role-based access control schema
-- Created: 2024-01-02

BEGIN;

-- Create admin role enum
CREATE TYPE admin_role AS ENUM (
    'super_admin',
    'admin', 
    'event_manager',
    'support',
    'analyst'
);

-- Admin users table
CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role admin_role NOT NULL DEFAULT 'admin',
    permissions JSONB NOT NULL DEFAULT '[]',
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login TIMESTAMP WITH TIME ZONE,
    login_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    must_change_password BOOLEAN NOT NULL DEFAULT true,
    two_factor_enabled BOOLEAN NOT NULL DEFAULT false,
    two_factor_secret VARCHAR(255),
    created_by UUID REFERENCES admin_users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Admin login history table
CREATE TABLE admin_login_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN NOT NULL,
    login_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    logout_at TIMESTAMP WITH TIME ZONE
);

-- Admin activity log table
CREATE TABLE admin_activity_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL REFERENCES admin_users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for admin_users
CREATE INDEX idx_admin_users_email ON admin_users(email);
CREATE INDEX idx_admin_users_role ON admin_users(role);
CREATE INDEX idx_admin_users_is_active ON admin_users(is_active);
CREATE INDEX idx_admin_users_last_login ON admin_users(last_login);
CREATE INDEX idx_admin_users_locked_until ON admin_users(locked_until);
CREATE INDEX idx_admin_users_created_at ON admin_users(created_at);

-- Indexes for admin_login_history
CREATE INDEX idx_admin_login_history_admin_id ON admin_login_history(admin_id);
CREATE INDEX idx_admin_login_history_login_at ON admin_login_history(login_at);
CREATE INDEX idx_admin_login_history_success ON admin_login_history(success);
CREATE INDEX idx_admin_login_history_ip_address ON admin_login_history(ip_address);

-- Indexes for admin_activity_log
CREATE INDEX idx_admin_activity_log_admin_id ON admin_activity_log(admin_id);
CREATE INDEX idx_admin_activity_log_action ON admin_activity_log(action);
CREATE INDEX idx_admin_activity_log_resource_type ON admin_activity_log(resource_type);
CREATE INDEX idx_admin_activity_log_resource_id ON admin_activity_log(resource_id);
CREATE INDEX idx_admin_activity_log_created_at ON admin_activity_log(created_at);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_admin_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_admin_users_updated_at
    BEFORE UPDATE ON admin_users
    FOR EACH ROW
    EXECUTE FUNCTION update_admin_users_updated_at();

-- Function to create default super admin
CREATE OR REPLACE FUNCTION create_default_super_admin()
RETURNS VOID AS $$
DECLARE
    super_admin_permissions JSONB;
BEGIN
    -- Define super admin permissions
    super_admin_permissions := '[
        "system_settings",
        "user_management", 
        "admin_management",
        "event_create",
        "event_edit",
        "event_delete",
        "event_publish",
        "organizer_create",
        "organizer_edit",
        "organizer_delete",
        "organizer_approve",
        "order_view",
        "order_edit",
        "order_refund",
        "order_cancel",
        "payment_view",
        "payment_process",
        "payment_refund",
        "support_tickets",
        "customer_view",
        "customer_edit",
        "analytics_view",
        "reports_generate",
        "reports_export",
        "scanner_manage",
        "scanner_view"
    ]'::JSONB;
    
    -- Insert super admin if not exists
    INSERT INTO admin_users (
        email,
        password_hash,
        first_name,
        last_name,
        role,
        permissions,
        is_active,
        must_change_password,
        created_at,
        updated_at
    ) VALUES (
        'admin@uduxpass.com',
        '$argon2id$v=19$m=65536,t=3,p=4$YWJjZGVmZ2hpams$+DhQmGspZct+3jqQMz+QOQ',  -- Password: Admin123!
        'Super',
        'Administrator',
        'super_admin',
        super_admin_permissions,
        true,
        true,  -- Must change password on first login
        NOW(),
        NOW()
    ) ON CONFLICT (email) DO NOTHING;
    
    -- Create additional test admin accounts
    INSERT INTO admin_users (
        email,
        password_hash,
        first_name,
        last_name,
        role,
        permissions,
        is_active,
        must_change_password,
        created_at,
        updated_at
    ) VALUES 
    (
        'eventmanager@uduxpass.com',
        '$argon2id$v=19$m=65536,t=3,p=4$YWJjZGVmZ2hpams$+DhQmGspZct+3jqQMz+QOQ',  -- Password: Admin123!
        'Event',
        'Manager',
        'event_manager',
        '["event_create", "event_edit", "event_publish", "organizer_edit", "order_view", "order_edit", "analytics_view", "reports_generate", "scanner_view"]'::JSONB,
        true,
        true,
        NOW(),
        NOW()
    ),
    (
        'support@uduxpass.com',
        '$argon2id$v=19$m=65536,t=3,p=4$YWJjZGVmZ2hpams$+DhQmGspZct+3jqQMz+QOQ',  -- Password: Admin123!
        'Support',
        'Agent',
        'support',
        '["support_tickets", "customer_view", "customer_edit", "order_view", "order_refund", "analytics_view"]'::JSONB,
        true,
        true,
        NOW(),
        NOW()
    ),
    (
        'analyst@uduxpass.com',
        '$argon2id$v=19$m=65536,t=3,p=4$YWJjZGVmZ2hpams$+DhQmGspZct+3jqQMz+QOQ',  -- Password: Admin123!
        'Data',
        'Analyst',
        'analyst',
        '["analytics_view", "reports_generate", "reports_export", "order_view", "customer_view", "scanner_view"]'::JSONB,
        true,
        true,
        NOW(),
        NOW()
    )
    ON CONFLICT (email) DO NOTHING;
END;
$$ LANGUAGE plpgsql;

-- Execute the function to create default admin users
SELECT create_default_super_admin();

-- Function to log admin activity
CREATE OR REPLACE FUNCTION log_admin_activity(
    p_admin_id UUID,
    p_action VARCHAR(100),
    p_resource_type VARCHAR(50) DEFAULT NULL,
    p_resource_id UUID DEFAULT NULL,
    p_details JSONB DEFAULT NULL,
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL
)
RETURNS VOID AS $$
BEGIN
    INSERT INTO admin_activity_log (
        admin_id,
        action,
        resource_type,
        resource_id,
        details,
        ip_address,
        user_agent,
        created_at
    ) VALUES (
        p_admin_id,
        p_action,
        p_resource_type,
        p_resource_id,
        p_details,
        p_ip_address,
        p_user_agent,
        NOW()
    );
END;
$$ LANGUAGE plpgsql;

-- Function to record admin login
CREATE OR REPLACE FUNCTION record_admin_login(
    p_admin_id UUID,
    p_ip_address INET DEFAULT NULL,
    p_user_agent TEXT DEFAULT NULL,
    p_success BOOLEAN DEFAULT true
)
RETURNS VOID AS $$
BEGIN
    INSERT INTO admin_login_history (
        admin_id,
        ip_address,
        user_agent,
        success,
        login_at
    ) VALUES (
        p_admin_id,
        p_ip_address,
        p_user_agent,
        p_success,
        NOW()
    );
    
    -- Update last_login if successful
    IF p_success THEN
        UPDATE admin_users 
        SET 
            last_login = NOW(),
            login_attempts = 0,
            locked_until = NULL,
            updated_at = NOW()
        WHERE id = p_admin_id;
    ELSE
        -- Increment login attempts and potentially lock account
        UPDATE admin_users 
        SET 
            login_attempts = login_attempts + 1,
            locked_until = CASE 
                WHEN login_attempts + 1 >= 5 THEN NOW() + INTERVAL '30 minutes'
                ELSE locked_until
            END,
            updated_at = NOW()
        WHERE id = p_admin_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Function to clean up old login history (keep last 1000 records per admin)
CREATE OR REPLACE FUNCTION cleanup_admin_login_history()
RETURNS VOID AS $$
BEGIN
    DELETE FROM admin_login_history
    WHERE id NOT IN (
        SELECT id FROM (
            SELECT id, ROW_NUMBER() OVER (PARTITION BY admin_id ORDER BY login_at DESC) as rn
            FROM admin_login_history
        ) ranked
        WHERE rn <= 1000
    );
END;
$$ LANGUAGE plpgsql;

-- Function to clean up old activity logs (keep last 6 months)
CREATE OR REPLACE FUNCTION cleanup_admin_activity_log()
RETURNS VOID AS $$
BEGIN
    DELETE FROM admin_activity_log
    WHERE created_at < NOW() - INTERVAL '6 months';
END;
$$ LANGUAGE plpgsql;

-- Add constraints
ALTER TABLE admin_users ADD CONSTRAINT chk_admin_users_email_format 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE admin_users ADD CONSTRAINT chk_admin_users_login_attempts 
    CHECK (login_attempts >= 0 AND login_attempts <= 10);

ALTER TABLE admin_users ADD CONSTRAINT chk_admin_users_names_not_empty 
    CHECK (LENGTH(TRIM(first_name)) > 0 AND LENGTH(TRIM(last_name)) > 0);

-- Add comments
COMMENT ON TABLE admin_users IS 'Administrative users with role-based access control';
COMMENT ON TABLE admin_login_history IS 'Login history tracking for admin users';
COMMENT ON TABLE admin_activity_log IS 'Activity log for admin actions and changes';

COMMENT ON COLUMN admin_users.permissions IS 'JSON array of specific permissions for this admin user';
COMMENT ON COLUMN admin_users.locked_until IS 'Account lock expiration timestamp after failed login attempts';
COMMENT ON COLUMN admin_users.must_change_password IS 'Forces password change on next login';
COMMENT ON COLUMN admin_users.two_factor_enabled IS 'Whether two-factor authentication is enabled';
COMMENT ON COLUMN admin_users.two_factor_secret IS 'TOTP secret for two-factor authentication';

COMMIT;

