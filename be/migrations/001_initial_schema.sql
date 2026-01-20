-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'creator',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_role CHECK (role IN ('admin', 'creator'))
);

-- Create social_accounts table
CREATE TABLE IF NOT EXISTS social_accounts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    account_id VARCHAR(255),
    access_token TEXT,
    refresh_token TEXT,
    token_expires_at TIMESTAMP,
    last_pull_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_platform CHECK (platform IN ('twitter', 'facebook', 'instagram', 'youtube', 'tiktok'))
);

-- Create content table
CREATE TABLE IF NOT EXISTS content (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    social_account_id INTEGER REFERENCES social_accounts(id) ON DELETE SET NULL,
    platform VARCHAR(50) NOT NULL,
    link TEXT NOT NULL,
    original_text TEXT,
    description TEXT,
    tags TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_content_platform CHECK (platform IN ('twitter', 'facebook', 'instagram', 'youtube', 'tiktok'))
);

-- Create indexes for better query performance
CREATE INDEX idx_users_user_id ON users(user_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_social_accounts_user_id ON social_accounts(user_id);
CREATE INDEX idx_social_accounts_platform ON social_accounts(platform);
CREATE INDEX idx_content_user_id ON content(user_id);
CREATE INDEX idx_content_platform ON content(platform);
CREATE INDEX idx_content_created_at ON content(created_at);
