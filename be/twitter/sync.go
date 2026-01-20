package twitter

import (
	"fmt"
	"log"
	"time"
)

// SyncResult contains the results of a sync operation
type SyncResult struct {
	AccountID    int       `json:"account_id"`
	TweetsSynced int       `json:"tweets_synced"`
	TweetsTotal  int       `json:"tweets_total"`
	Errors       []string  `json:"errors,omitempty"`
	SyncedAt     time.Time `json:"synced_at"`
}

// SyncedTweet represents a tweet ready to be stored
type SyncedTweet struct {
	ExternalID string
	Text       string
	Link       string
	PostedAt   time.Time
}

// Syncer handles synchronization of Twitter content
type Syncer struct {
	client       *Client
	oauthHandler *OAuthHandler
}

// NewSyncer creates a new Twitter syncer
func NewSyncer(client *Client) *Syncer {
	return &Syncer{
		client:       client,
		oauthHandler: NewOAuthHandler(),
	}
}

// GetOAuthHandler returns the OAuth handler for initiating OAuth flows
func (s *Syncer) GetOAuthHandler() *OAuthHandler {
	return s.oauthHandler
}

// SyncAccount fetches new tweets for an account (using app-level bearer token)
func (s *Syncer) SyncAccount(accountName string, accountID *string, sinceID string) ([]SyncedTweet, error) {
	var twitterUserID string

	// If we have a stored account ID, use it; otherwise look it up
	if accountID != nil && *accountID != "" {
		twitterUserID = *accountID
	} else {
		// Look up the user by username
		userResp, err := s.client.GetUserByUsername(accountName)
		if err != nil {
			return nil, fmt.Errorf("failed to lookup user %s: %w", accountName, err)
		}
		twitterUserID = userResp.Data.ID
	}

	// Fetch recent tweets
	tweetsResp, err := s.client.GetUserTweets(twitterUserID, 50, sinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tweets: %w", err)
	}

	if tweetsResp.Meta.ResultCount == 0 {
		log.Printf("No new tweets found for @%s", accountName)
		return []SyncedTweet{}, nil
	}

	var synced []SyncedTweet
	for _, tweet := range tweetsResp.Data {
		synced = append(synced, SyncedTweet{
			ExternalID: tweet.ID,
			Text:       tweet.Text,
			Link:       TweetToLink(accountName, tweet.ID),
			PostedAt:   tweet.CreatedAt,
		})
	}

	log.Printf("Fetched %d tweets for @%s", len(synced), accountName)
	return synced, nil
}

// SyncAccountWithOAuth fetches new tweets using user's OAuth access token
func (s *Syncer) SyncAccountWithOAuth(accessToken string, accountName string, accountID *string, sinceID string) ([]SyncedTweet, error) {
	userClient := NewUserClient(accessToken)

	var twitterUserID string
	var username string = accountName

	// If we have a stored account ID, use it; otherwise get it from the authenticated user
	if accountID != nil && *accountID != "" {
		twitterUserID = *accountID
	} else {
		// Get the authenticated user's ID
		userResp, err := userClient.GetAuthenticatedUser()
		if err != nil {
			return nil, fmt.Errorf("failed to get authenticated user: %w", err)
		}
		twitterUserID = userResp.Data.ID
		username = userResp.Data.Username
	}

	// Fetch recent tweets
	tweetsResp, err := userClient.GetUserTweets(twitterUserID, 50, sinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tweets: %w", err)
	}

	if tweetsResp.Meta.ResultCount == 0 {
		log.Printf("No new tweets found for @%s (OAuth)", username)
		return []SyncedTweet{}, nil
	}

	var synced []SyncedTweet
	for _, tweet := range tweetsResp.Data {
		synced = append(synced, SyncedTweet{
			ExternalID: tweet.ID,
			Text:       tweet.Text,
			Link:       TweetToLink(username, tweet.ID),
			PostedAt:   tweet.CreatedAt,
		})
	}

	log.Printf("Fetched %d tweets for @%s (OAuth)", len(synced), username)
	return synced, nil
}

// GetTwitterUserID fetches the Twitter user ID for a username
func (s *Syncer) GetTwitterUserID(username string) (string, error) {
	userResp, err := s.client.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	return userResp.Data.ID, nil
}

// GetTwitterUserIDWithOAuth fetches the Twitter user ID using OAuth token
func (s *Syncer) GetTwitterUserIDWithOAuth(accessToken string) (string, string, error) {
	userClient := NewUserClient(accessToken)
	userResp, err := userClient.GetAuthenticatedUser()
	if err != nil {
		return "", "", err
	}
	return userResp.Data.ID, userResp.Data.Username, nil
}
