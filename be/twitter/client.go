package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client handles Twitter/X API interactions
type Client struct {
	httpClient  *http.Client
	bearerToken string
	baseURL     string
}

// UserClient handles Twitter/X API interactions with user OAuth tokens
type UserClient struct {
	httpClient  *http.Client
	accessToken string
	baseURL     string
}

// Tweet represents a tweet from the X API
type Tweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	AuthorID  string    `json:"author_id"`
}

// TweetsResponse represents the API response for tweets
type TweetsResponse struct {
	Data []Tweet `json:"data"`
	Meta struct {
		ResultCount   int    `json:"result_count"`
		NextToken     string `json:"next_token"`
		PreviousToken string `json:"previous_token"`
	} `json:"meta"`
	Errors []APIError `json:"errors,omitempty"`
}

// UserResponse represents the API response for user lookup
type UserResponse struct {
	Data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"data"`
	Errors []APIError `json:"errors,omitempty"`
}

// APIError represents an error from the X API
type APIError struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Type   string `json:"type"`
}

// NewClient creates a new Twitter API client
func NewClient() *Client {
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		bearerToken: bearerToken,
		baseURL:     "https://api.x.com/2",
	}
}

// NewClientWithToken creates a new Twitter API client with specific token
func NewClientWithToken(bearerToken string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		bearerToken: bearerToken,
		baseURL:     "https://api.x.com/2",
	}
}

// IsConfigured returns true if the client has necessary credentials
func (c *Client) IsConfigured() bool {
	return c.bearerToken != ""
}

// GetUserByUsername fetches a user's profile by username
func (c *Client) GetUserByUsername(username string) (*UserResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("twitter client not configured: missing bearer token")
	}

	endpoint := fmt.Sprintf("%s/users/by/username/%s", c.baseURL, url.PathEscape(username))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.bearerToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var userResp UserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(userResp.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s - %s", userResp.Errors[0].Title, userResp.Errors[0].Detail)
	}

	return &userResp, nil
}

// GetUserTweets fetches recent tweets for a user by their ID
func (c *Client) GetUserTweets(userID string, maxResults int, sinceID string) (*TweetsResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("twitter client not configured: missing bearer token")
	}

	if maxResults <= 0 || maxResults > 100 {
		maxResults = 10
	}

	endpoint := fmt.Sprintf("%s/users/%s/tweets", c.baseURL, url.PathEscape(userID))

	params := url.Values{}
	params.Set("max_results", fmt.Sprintf("%d", maxResults))
	params.Set("tweet.fields", "created_at,author_id,text")

	if sinceID != "" {
		params.Set("since_id", sinceID)
	}

	fullURL := endpoint + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.bearerToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tweetsResp TweetsResponse
	if err := json.Unmarshal(body, &tweetsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(tweetsResp.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s - %s", tweetsResp.Errors[0].Title, tweetsResp.Errors[0].Detail)
	}

	return &tweetsResp, nil
}

// TweetToLink converts a tweet ID and username to a full URL
func TweetToLink(username, tweetID string) string {
	return fmt.Sprintf("https://x.com/%s/status/%s", username, tweetID)
}

// NewUserClient creates a new Twitter API client with user OAuth access token
func NewUserClient(accessToken string) *UserClient {
	return &UserClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		accessToken: accessToken,
		baseURL:     "https://api.x.com/2",
	}
}

// IsConfigured returns true if the user client has an access token
func (c *UserClient) IsConfigured() bool {
	return c.accessToken != ""
}

// GetAuthenticatedUser fetches the authenticated user's profile
func (c *UserClient) GetAuthenticatedUser() (*UserResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("user client not configured: missing access token")
	}

	endpoint := fmt.Sprintf("%s/users/me", c.baseURL)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var userResp UserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &userResp, nil
}

// GetMyTweets fetches the authenticated user's recent tweets
func (c *UserClient) GetMyTweets(maxResults int, sinceID string) (*TweetsResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("user client not configured: missing access token")
	}

	// First get the user ID
	userResp, err := c.GetAuthenticatedUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}

	return c.GetUserTweets(userResp.Data.ID, maxResults, sinceID)
}

// GetUserTweets fetches recent tweets for a user by their ID
func (c *UserClient) GetUserTweets(userID string, maxResults int, sinceID string) (*TweetsResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("user client not configured: missing access token")
	}

	if maxResults <= 0 || maxResults > 100 {
		maxResults = 10
	}

	endpoint := fmt.Sprintf("%s/users/%s/tweets", c.baseURL, url.PathEscape(userID))

	params := url.Values{}
	params.Set("max_results", fmt.Sprintf("%d", maxResults))
	params.Set("tweet.fields", "created_at,author_id,text")

	if sinceID != "" {
		params.Set("since_id", sinceID)
	}

	fullURL := endpoint + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tweetsResp TweetsResponse
	if err := json.Unmarshal(body, &tweetsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(tweetsResp.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s - %s", tweetsResp.Errors[0].Title, tweetsResp.Errors[0].Detail)
	}

	return &tweetsResp, nil
}
