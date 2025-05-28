-- +migrate Up
-- Add deleted_at column for GORM soft deletes
DO $$
BEGIN
    -- Add deleted_at if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'deleted_at') THEN
        ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP NULL;
    END IF;

    -- Create index for deleted_at if not exists (improves soft delete queries)
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_users_deleted_at') THEN
        CREATE INDEX idx_users_deleted_at ON users(deleted_at);
    END IF;
END $$;

-- +migrate Down
-- Remove deleted_at column and its index
DROP INDEX IF EXISTS idx_users_deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
-- Example:
-- DROP TABLE IF EXISTS example;
