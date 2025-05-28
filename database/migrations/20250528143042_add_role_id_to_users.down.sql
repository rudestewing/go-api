-- Remove role_id column and its index
DROP INDEX IF EXISTS idx_users_role_id;
ALTER TABLE users DROP COLUMN IF EXISTS role_id;
