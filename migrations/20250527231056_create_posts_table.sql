-- +migrate Up
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes only if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_posts_user_id') THEN
        CREATE INDEX idx_posts_user_id ON posts(user_id);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_posts_status') THEN
        CREATE INDEX idx_posts_status ON posts(status);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_posts_published_at') THEN
        CREATE INDEX idx_posts_published_at ON posts(published_at);
    END IF;
END $$;

-- +migrate Down
DROP INDEX IF EXISTS idx_posts_published_at;
DROP INDEX IF EXISTS idx_posts_status;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP TABLE IF EXISTS posts;
