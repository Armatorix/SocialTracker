package handlers

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"strconv"

	"github.com/Armatorix/SocialTracker/be/models"
	"github.com/Armatorix/SocialTracker/be/repository"
	"github.com/Armatorix/SocialTracker/be/twitter"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo          *repository.Repository
	twitterSyncer *twitter.Syncer
}

func NewHandler(repo *repository.Repository) *Handler {
	twitterClient := twitter.NewClient()
	return &Handler{
		repo:          repo,
		twitterSyncer: twitter.NewSyncer(twitterClient),
	}
}

// GetCurrentUser returns the current authenticated user
func (h *Handler) GetCurrentUser(c echo.Context) error {
	user, err := h.getCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// Social Account handlers
func (h *Handler) CreateSocialAccount(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req models.CreateSocialAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	account, err := h.repo.CreateSocialAccount(userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, account)
}

func (h *Handler) GetSocialAccounts(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	accounts, err := h.repo.GetSocialAccountsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if accounts == nil {
		accounts = []models.SocialAccount{}
	}

	return c.JSON(http.StatusOK, accounts)
}

func (h *Handler) DeleteSocialAccount(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid account id"})
	}

	err = h.repo.DeleteSocialAccount(accountID, userID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "account not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "account deleted"})
}

func (h *Handler) PullContentFromPlatform(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	accountID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid account id"})
	}

	// Get the social account
	account, err := h.repo.GetSocialAccountByID(accountID, userID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "account not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var response models.SyncResponse
	response.AccountID = accountID
	response.Platform = account.Platform
	response.AccountName = account.AccountName

	switch account.Platform {
	case "twitter":
		response, err = h.syncTwitterAccount(userID, account)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "auto-sync not supported for platform: " + account.Platform,
		})
	}

	// Update last pull time
	err = h.repo.UpdateSocialAccountLastPull(accountID, userID)
	if err != nil {
		log.Printf("Failed to update last pull time: %v", err)
	}

	return c.JSON(http.StatusOK, response)
}

// syncTwitterAccount syncs content from Twitter/X for the given account
func (h *Handler) syncTwitterAccount(userID int, account *models.SocialAccount) (models.SyncResponse, error) {
	response := models.SyncResponse{
		AccountID:   account.ID,
		Platform:    account.Platform,
		AccountName: account.AccountName,
	}

	// Get the latest synced tweet ID to only fetch newer tweets
	latestID, err := h.repo.GetLatestExternalPostID(account.ID)
	if err != nil {
		log.Printf("Failed to get latest external post ID: %v", err)
	}

	var sinceID string
	if latestID != nil {
		sinceID = *latestID
	}

	// Fetch tweets from Twitter
	tweets, err := h.twitterSyncer.SyncAccount(account.AccountName, account.AccountID, sinceID)
	if err != nil {
		return response, err
	}

	// If we don't have the Twitter user ID stored, fetch and save it
	if account.AccountID == nil || *account.AccountID == "" {
		twitterUserID, err := h.twitterSyncer.GetTwitterUserID(account.AccountName)
		if err != nil {
			log.Printf("Failed to get Twitter user ID: %v", err)
		} else {
			err = h.repo.UpdateSocialAccountID(account.ID, twitterUserID)
			if err != nil {
				log.Printf("Failed to update social account ID: %v", err)
			}
		}
	}

	// Store each tweet as content
	for _, tweet := range tweets {
		content, err := h.repo.CreateSyncedContent(
			userID,
			account.ID,
			"twitter",
			tweet.Link,
			tweet.Text,
			tweet.ExternalID,
			tweet.PostedAt,
		)
		if err != nil {
			log.Printf("Failed to create content for tweet %s: %v", tweet.ExternalID, err)
			response.Errors = append(response.Errors, err.Error())
			continue
		}
		if content == nil {
			// Duplicate tweet, already exists
			response.SkippedCount++
		} else {
			response.SyncedCount++
		}
	}

	response.Message = "Sync completed successfully"
	if response.SyncedCount == 0 && response.SkippedCount == 0 {
		response.Message = "No new tweets found"
	}

	log.Printf("Twitter sync for @%s: synced=%d, skipped=%d", account.AccountName, response.SyncedCount, response.SkippedCount)
	return response, nil
}

// Content handlers
func (h *Handler) CreateContent(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req models.CreateContentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	content, err := h.repo.CreateContent(userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, content)
}

func (h *Handler) GetContent(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	content, err := h.repo.GetContentByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if content == nil {
		content = []models.Content{}
	}

	return c.JSON(http.StatusOK, content)
}

func (h *Handler) DeleteContent(c echo.Context) error {
	userID, err := h.getUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	contentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid content id"})
	}

	err = h.repo.DeleteContent(contentID, userID)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "content not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "content deleted"})
}

// Admin handlers
func (h *Handler) GetAllContent(c echo.Context) error {
	// Check if user is admin
	user, err := h.getCurrentUser(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "admin access required"})
	}

	// Get query parameters for filtering
	filters := make(map[string]string)
	if platform := c.QueryParam("platform"); platform != "" {
		filters["platform"] = platform
	}
	if username := c.QueryParam("username"); username != "" {
		filters["username"] = username
	}

	content, err := h.repo.GetAllContent(filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if content == nil {
		content = []models.ContentWithUser{}
	}

	return c.JSON(http.StatusOK, content)
}

// Helper methods
func (h *Handler) getUserID(c echo.Context) (int, error) {
	user, err := h.getCurrentUser(c)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (h *Handler) getCurrentUser(c echo.Context) (*models.User, error) {
	// oauth2-proxy with PASS_USER_HEADERS=true sends X-Forwarded-* headers
	userID := c.Request().Header.Get("X-Forwarded-User")
	email := c.Request().Header.Get("X-Forwarded-Email")
	username := c.Request().Header.Get("X-Forwarded-Preferred-Username")

	// Fallback to X-Auth-Request-* headers (for nginx auth_request setups)
	if userID == "" {
		userID = c.Request().Header.Get("X-Auth-Request-User")
	}
	// attempt to decode userId
	nUserId, err := base64.RawStdEncoding.DecodeString(userID)
	if err == nil && len(nUserId) > 0 {
		userID = string(nUserId)
	}
	if email == "" {
		email = c.Request().Header.Get("X-Auth-Request-Email")
	}
	if username == "" {
		username = c.Request().Header.Get("X-Auth-Request-Preferred-Username")
	}

	// Extract username from email if not provided separately
	if idx := len(email); idx > 0 {
		if atIdx := 0; atIdx < idx {
			for i, ch := range email {
				if ch == '@' {
					atIdx = i
					break
				}
			}
			if atIdx > 0 {
				username = email[:atIdx]
			}
		}
	}

	return h.repo.GetOrCreateUser(userID, email, username)
}
