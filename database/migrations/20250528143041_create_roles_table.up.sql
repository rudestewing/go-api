-- Create roles table for GORM
DO $$
BEGIN
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for better performance
IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_roles_code') THEN
    CREATE INDEX idx_roles_code ON roles(code);
END IF;

IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_roles_deleted_at') THEN
    CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);
END IF;

END $$;
