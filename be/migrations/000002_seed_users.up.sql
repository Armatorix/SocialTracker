-- Set admin role for admin user
INSERT INTO users (user_id, email, username, role)
VALUES ('admin-001', 'admin@socialtracker.com', 'admin', 'admin')
ON CONFLICT (user_id) DO UPDATE SET role = 'admin';

-- Ensure creator users have creator role
UPDATE users SET role = 'creator' WHERE user_id LIKE 'creator-%';
