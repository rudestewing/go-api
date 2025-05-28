-- Drop roles table and its indexes
DROP INDEX IF EXISTS idx_roles_deleted_at;
DROP INDEX IF EXISTS idx_roles_code;
DROP TABLE IF EXISTS roles;
