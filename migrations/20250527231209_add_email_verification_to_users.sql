-- +migrate Up
-- Add columns only if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'email_verified_at') THEN
        ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMP NULL;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'email_verification_token') THEN
        ALTER TABLE users ADD COLUMN email_verification_token VARCHAR(255) NULL;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_users_email_verification_token') THEN
        CREATE INDEX idx_users_email_verification_token ON users(email_verification_token);
    END IF;
END $$;

-- +migrate Down
DROP INDEX IF EXISTS idx_users_email_verification_token;
ALTER TABLE users DROP COLUMN IF EXISTS email_verification_token;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
