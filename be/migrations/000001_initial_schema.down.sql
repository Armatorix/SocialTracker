-- Drop indexes
DROP INDEX IF EXISTS idx_content_created_at;
DROP INDEX IF EXISTS idx_content_platform;
DROP INDEX IF EXISTS idx_content_user_id;
DROP INDEX IF EXISTS idx_social_accounts_platform;
DROP INDEX IF EXISTS idx_social_accounts_user_id;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_user_id;

-- Drop tables
DROP TABLE IF EXISTS content;
DROP TABLE IF EXISTS social_accounts;
DROP TABLE IF EXISTS users;
