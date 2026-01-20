-- Remove indexes and column
DROP INDEX IF EXISTS idx_content_external_post_unique;
DROP INDEX IF EXISTS idx_content_external_post_id;
ALTER TABLE content DROP COLUMN IF EXISTS posted_at;
ALTER TABLE content DROP COLUMN IF EXISTS external_post_id;
