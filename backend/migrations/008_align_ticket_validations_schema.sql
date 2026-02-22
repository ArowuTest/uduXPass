-- Migration: Align ticket_validations table with backend entity
-- Purpose: Strategic schema alignment for enterprise-grade validation tracking
-- Date: 2026-02-17

-- Add missing columns to ticket_validations table
ALTER TABLE ticket_validations
ADD COLUMN IF NOT EXISTS scanner_id UUID,
ADD COLUMN IF NOT EXISTS session_id UUID,
ADD COLUMN IF NOT EXISTS validation_result VARCHAR(50) DEFAULT 'valid',
ADD COLUMN IF NOT EXISTS notes TEXT,
ADD COLUMN IF NOT EXISTS device_info JSONB DEFAULT '{}',
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add foreign key constraints
ALTER TABLE ticket_validations
ADD CONSTRAINT IF NOT EXISTS ticket_validations_scanner_id_fkey 
FOREIGN KEY (scanner_id) REFERENCES scanner_users(id);

ALTER TABLE ticket_validations
ADD CONSTRAINT IF NOT EXISTS ticket_validations_session_id_fkey 
FOREIGN KEY (session_id) REFERENCES scanner_sessions(id);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_ticket_validations_scanner 
ON ticket_validations(scanner_id);

CREATE INDEX IF NOT EXISTS idx_ticket_validations_session 
ON ticket_validations(session_id);

CREATE INDEX IF NOT EXISTS idx_ticket_validations_result 
ON ticket_validations(validation_result);

CREATE INDEX IF NOT EXISTS idx_ticket_validations_validated_at 
ON ticket_validations(validated_at);

-- Update existing records to have validation_result
UPDATE ticket_validations 
SET validation_result = 'valid' 
WHERE validation_result IS NULL;

-- Grant permissions
GRANT ALL ON ticket_validations TO uduxpass_user;

-- Add comment
COMMENT ON TABLE ticket_validations IS 'Enterprise-grade ticket validation tracking with full audit trail';
COMMENT ON COLUMN ticket_validations.scanner_id IS 'ID of the scanner who validated the ticket';
COMMENT ON COLUMN ticket_validations.session_id IS 'ID of the scanning session';
COMMENT ON COLUMN ticket_validations.validation_result IS 'Result of validation: valid, invalid, expired, used, etc.';
COMMENT ON COLUMN ticket_validations.notes IS 'Additional notes about the validation';
COMMENT ON COLUMN ticket_validations.device_info IS 'JSON object containing device information';
