-- Add external_post_id column to content table for deduplication
ALTER TABLE content ADD COLUMN IF NOT EXISTS external_post_id VARCHAR(255);

-- Add posted_at column to store original post timestamp
ALTER TABLE content ADD COLUMN IF NOT EXISTS posted_at TIMESTAMP;

-- Add unique constraint to prevent duplicate imports
CREATE UNIQUE INDEX IF NOT EXISTS idx_content_external_post_unique 
ON content(social_account_id, external_post_id) 
WHERE external_post_id IS NOT NULL;

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_content_external_post_id ON content(external_post_id);
