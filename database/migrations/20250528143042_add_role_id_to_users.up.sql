-- Add role_id column to users table (no foreign key constraint)
DO $$
BEGIN
    -- Add role_id column if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'role_id') THEN
        ALTER TABLE users ADD COLUMN role_id INTEGER NOT NULL DEFAULT 2;
    END IF;

    -- Create index for role_id for better query performance
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_users_role_id') THEN
        CREATE INDEX idx_users_role_id ON users(role_id);
    END IF;

    -- Update existing users to have default user role (role_id = 2)
    UPDATE users SET role_id = 2 WHERE role_id IS NULL OR role_id = 0;

END $$;
