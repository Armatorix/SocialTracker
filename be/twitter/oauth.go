package twitter

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// OAuthConfig holds Twitter OAuth 2.0 configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

// OAuthState stores temporary state for OAuth flow
type OAuthState struct {
	UserID       int
	CodeVerifier string
	CreatedAt    time.Time
}

// TokenResponse represents the OAuth token response from Twitter
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// OAuthUserResponse represents the authenticated user from Twitter
type OAuthUserResponse struct {
	Data struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"data"`
}

// OAuthHandler manages Twitter OAuth 2.0 flows
type OAuthHandler struct {
	config     OAuthConfig
	states     map[string]*OAuthState
	statesMu   sync.RWMutex
	httpClient *http.Client
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{
		config: OAuthConfig{
			ClientID:     os.Getenv("TWITTER_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITTER_CLIENT_SECRET"),
			RedirectURI:  os.Getenv("TWITTER_REDIRECT_URI"),
			Scopes:       []string{"tweet.read", "users.read", "offline.access"},
		},
		states:     make(map[string]*OAuthState),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// IsConfigured returns true if OAuth is properly configured
func (h *OAuthHandler) IsConfigured() bool {
	return h.config.ClientID != "" && h.config.ClientSecret != "" && h.config.RedirectURI != ""
}

// GetAuthorizationURL generates the authorization URL for the OAuth flow
func (h *OAuthHandler) GetAuthorizationURL(userID int) (string, error) {
	if !h.IsConfigured() {
		return "", fmt.Errorf("twitter OAuth not configured")
	}

	// Generate state and PKCE code verifier
	state, err := generateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	codeVerifier, err := generateRandomString(64)
	if err != nil {
		return "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	// Store state for verification
	h.statesMu.Lock()
	h.states[state] = &OAuthState{
		UserID:       userID,
		CodeVerifier: codeVerifier,
		CreatedAt:    time.Now(),
	}
	h.statesMu.Unlock()

	// Clean up old states
	go h.cleanupOldStates()

	// Generate code challenge from verifier (S256)
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Build authorization URL
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", h.config.ClientID)
	params.Set("redirect_uri", h.config.RedirectURI)
	params.Set("scope", strings.Join(h.config.Scopes, " "))
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	return "https://twitter.com/i/oauth2/authorize?" + params.Encode(), nil
}

// ExchangeCode exchanges the authorization code for tokens
func (h *OAuthHandler) ExchangeCode(code, state string) (*TokenResponse, int, error) {
	// Verify and retrieve state
	h.statesMu.Lock()
	oauthState, exists := h.states[state]
	if exists {
		delete(h.states, state)
	}
	h.statesMu.Unlock()

	if !exists {
		return nil, 0, fmt.Errorf("invalid or expired state")
	}

	// Check if state is too old (10 minutes)
	if time.Since(oauthState.CreatedAt) > 10*time.Minute {
		return nil, 0, fmt.Errorf("state expired")
	}

	// Exchange code for tokens
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", h.config.RedirectURI)
	data.Set("code_verifier", oauthState.CodeVerifier)

	req, err := http.NewRequest("POST", "https://api.x.com/2/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Use Basic Auth with client credentials
	auth := base64.StdEncoding.EncodeToString([]byte(h.config.ClientID + ":" + h.config.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, oauthState.UserID, nil
}

// RefreshAccessToken refreshes an expired access token
func (h *OAuthHandler) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://api.x.com/2/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(h.config.ClientID + ":" + h.config.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// GetAuthenticatedUser fetches the authenticated user's profile using the access token
func (h *OAuthHandler) GetAuthenticatedUser(accessToken string) (*OAuthUserResponse, error) {
	req, err := http.NewRequest("GET", "https://api.x.com/2/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user (status %d): %s", resp.StatusCode, string(body))
	}

	var userResp OAuthUserResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &userResp, nil
}

// cleanupOldStates removes states older than 15 minutes
func (h *OAuthHandler) cleanupOldStates() {
	h.statesMu.Lock()
	defer h.statesMu.Unlock()

	cutoff := time.Now().Add(-15 * time.Minute)
	for state, oauthState := range h.states {
		if oauthState.CreatedAt.Before(cutoff) {
			delete(h.states, state)
		}
	}
}

// generateRandomString generates a cryptographically secure random string
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length], nil
}

// generateCodeChallenge creates a S256 code challenge from the verifier
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
