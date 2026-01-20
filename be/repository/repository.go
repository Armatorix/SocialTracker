package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Armatorix/SocialTracker/be/models"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// User operations
func (r *Repository) GetOrCreateUser(userID, email, username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(`
		INSERT INTO users (user_id, email, username, role)
		VALUES ($1, $2, $3, 'creator')
		ON CONFLICT (user_id) DO UPDATE SET email = $2, username = $3, updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, email, username, role, created_at, updated_at
	`, userID, email, username).Scan(&user.ID, &user.UserID, &user.Email, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByUserID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(`
		SELECT id, user_id, email, username, role, created_at, updated_at
		FROM users WHERE user_id = $1
	`, userID).Scan(&user.ID, &user.UserID, &user.Email, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Social Account operations
func (r *Repository) CreateSocialAccount(userID int, req models.CreateSocialAccountRequest) (*models.SocialAccount, error) {
	var account models.SocialAccount
	err := r.db.QueryRow(`
		INSERT INTO social_accounts (user_id, platform, account_name, account_id, access_token, refresh_token)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, platform, account_name, account_id, last_pull_at, created_at, updated_at
	`, userID, req.Platform, req.AccountName, req.AccountID, req.AccessToken, req.RefreshToken).
		Scan(&account.ID, &account.UserID, &account.Platform, &account.AccountName, &account.AccountID, &account.LastPullAt, &account.CreatedAt, &account.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *Repository) GetSocialAccountsByUserID(userID int) ([]models.SocialAccount, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, platform, account_name, account_id, last_pull_at, created_at, updated_at
		FROM social_accounts WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.SocialAccount
	for rows.Next() {
		var account models.SocialAccount
		err := rows.Scan(&account.ID, &account.UserID, &account.Platform, &account.AccountName, &account.AccountID, &account.LastPullAt, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *Repository) DeleteSocialAccount(accountID, userID int) error {
	result, err := r.db.Exec(`
		DELETE FROM social_accounts WHERE id = $1 AND user_id = $2
	`, accountID, userID)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) UpdateSocialAccountLastPull(accountID, userID int) error {
	_, err := r.db.Exec(`
		UPDATE social_accounts SET last_pull_at = $1, updated_at = $1
		WHERE id = $2 AND user_id = $3
	`, time.Now(), accountID, userID)
	return err
}

// Content operations
func (r *Repository) CreateContent(userID int, req models.CreateContentRequest) (*models.Content, error) {
	var content models.Content
	err := r.db.QueryRow(`
		INSERT INTO content (user_id, social_account_id, platform, link, original_text, description, tags)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, social_account_id, platform, link, original_text, description, tags, created_at, updated_at
	`, userID, req.SocialAccountID, req.Platform, req.Link, req.OriginalText, req.Description, pq.Array(req.Tags)).
		Scan(&content.ID, &content.UserID, &content.SocialAccountID, &content.Platform, &content.Link, 
			&content.OriginalText, &content.Description, pq.Array(&content.Tags), &content.CreatedAt, &content.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *Repository) GetContentByUserID(userID int) ([]models.Content, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, social_account_id, platform, link, original_text, description, tags, created_at, updated_at
		FROM content WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.Content
	for rows.Next() {
		var content models.Content
		err := rows.Scan(&content.ID, &content.UserID, &content.SocialAccountID, &content.Platform, &content.Link,
			&content.OriginalText, &content.Description, pq.Array(&content.Tags), &content.CreatedAt, &content.UpdatedAt)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (r *Repository) GetAllContent(filters map[string]string) ([]models.ContentWithUser, error) {
	query := `
		SELECT c.id, c.user_id, c.social_account_id, c.platform, c.link, c.original_text, 
		       c.description, c.tags, c.created_at, c.updated_at, u.username, u.email
		FROM content c
		JOIN users u ON c.user_id = u.id
		WHERE 1=1
	`
	
	var args []interface{}
	argCount := 1
	
	if platform, ok := filters["platform"]; ok && platform != "" {
		query += fmt.Sprintf(" AND c.platform = $%d", argCount)
		args = append(args, platform)
		argCount++
	}
	
	if username, ok := filters["username"]; ok && username != "" {
		query += fmt.Sprintf(" AND u.username ILIKE $%d", argCount)
		args = append(args, "%"+username+"%")
		argCount++
	}
	
	query += " ORDER BY c.created_at DESC"
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.ContentWithUser
	for rows.Next() {
		var content models.ContentWithUser
		err := rows.Scan(&content.ID, &content.UserID, &content.SocialAccountID, &content.Platform, &content.Link,
			&content.OriginalText, &content.Description, pq.Array(&content.Tags), &content.CreatedAt, &content.UpdatedAt,
			&content.Username, &content.Email)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	return contents, nil
}

func (r *Repository) DeleteContent(contentID, userID int) error {
	result, err := r.db.Exec(`
		DELETE FROM content WHERE id = $1 AND user_id = $2
	`, contentID, userID)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
