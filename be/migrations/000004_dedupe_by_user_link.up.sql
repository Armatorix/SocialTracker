-- Add unique constraint to deduplicate content by user_id and link
-- This ensures a user cannot have the same post link twice

-- First, delete any existing duplicates (keep the oldest entry)
DELETE FROM content c1
WHERE c1.id NOT IN (
    SELECT MIN(id)
    FROM content
    GROUP BY user_id, link
);

-- Create unique index for user_id and link combination
CREATE UNIQUE INDEX IF NOT EXISTS idx_content_user_link_unique 
ON content(user_id, link);
