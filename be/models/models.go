package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	Username  string    `json:"username" db:"username"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SocialAccount struct {
	ID             int        `json:"id" db:"id"`
	UserID         int        `json:"user_id" db:"user_id"`
	Platform       string     `json:"platform" db:"platform"`
	AccountName    string     `json:"account_name" db:"account_name"`
	AccountID      *string    `json:"account_id,omitempty" db:"account_id"`
	AccessToken    *string    `json:"-" db:"access_token"`
	RefreshToken   *string    `json:"-" db:"refresh_token"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty" db:"token_expires_at"`
	LastPullAt     *time.Time `json:"last_pull_at,omitempty" db:"last_pull_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type Content struct {
	ID              int       `json:"id" db:"id"`
	UserID          int       `json:"user_id" db:"user_id"`
	SocialAccountID *int      `json:"social_account_id,omitempty" db:"social_account_id"`
	Platform        string    `json:"platform" db:"platform"`
	Link            string    `json:"link" db:"link"`
	OriginalText    *string   `json:"original_text,omitempty" db:"original_text"`
	Description     *string   `json:"description,omitempty" db:"description"`
	Tags            []string  `json:"tags,omitempty" db:"tags"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSocialAccountRequest struct {
	Platform     string  `json:"platform" binding:"required"`
	AccountName  string  `json:"account_name" binding:"required"`
	AccountID    *string `json:"account_id"`
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}

type CreateContentRequest struct {
	SocialAccountID *int     `json:"social_account_id"`
	Platform        string   `json:"platform" binding:"required"`
	Link            string   `json:"link" binding:"required"`
	OriginalText    *string  `json:"original_text"`
	Description     *string  `json:"description"`
	Tags            []string `json:"tags"`
}

type ContentWithUser struct {
	Content
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}
