-- +migrate Up
-- Add additional profile fields to users table
DO $$
BEGIN
    -- Add phone if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'phone') THEN
        ALTER TABLE users ADD COLUMN phone VARCHAR(20);
    END IF;

    -- Add date_of_birth if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'date_of_birth') THEN
        ALTER TABLE users ADD COLUMN date_of_birth DATE;
    END IF;

    -- Add gender if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'gender') THEN
        ALTER TABLE users ADD COLUMN gender VARCHAR(10);
    END IF;

    -- Add profile_image_url if not exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'profile_image_url') THEN
        ALTER TABLE users ADD COLUMN profile_image_url TEXT;
    END IF;

    -- Create index on phone
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_users_phone') THEN
        CREATE INDEX idx_users_phone ON users(phone);
    END IF;
END $$;

-- +migrate Down
DROP INDEX IF EXISTS idx_users_phone;
ALTER TABLE users DROP COLUMN IF EXISTS profile_image_url;
ALTER TABLE users DROP COLUMN IF EXISTS gender;
ALTER TABLE users DROP COLUMN IF EXISTS date_of_birth;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
