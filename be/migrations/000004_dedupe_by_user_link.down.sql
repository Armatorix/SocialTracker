-- Remove the user_id + link unique constraint
DROP INDEX IF EXISTS idx_content_user_link_unique;
