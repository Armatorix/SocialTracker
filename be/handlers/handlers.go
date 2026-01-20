package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Armatorix/SocialTracker/be/models"
	"github.com/Armatorix/SocialTracker/be/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// GetCurrentUser returns the current authenticated user
func (h *Handler) GetCurrentUser(c echo.Context) error {
	// Get user info from OAuth2 headers
	userID := c.Request().Header.Get("X-Auth-Request-User")
	email := c.Request().Header.Get("X-Auth-Request-Email")
	
	if userID == "" {
		userID = "2137" // Default for testing
		email = "a@e.com"
	}
	
	username := email
	if userID == "2137" {
		username = "admin"
	}
	
	// Get or create user in database
	user, err := h.repo.GetOrCreateUser(userID, email, username)
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
	
	// Update last pull time
	err = h.repo.UpdateSocialAccountLastPull(accountID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	
	// In a real implementation, this would call the actual social media APIs
	// For now, we just update the last pull time
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Content pull initiated. In production, this would fetch content from the platform API.",
	})
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
	userID := c.Request().Header.Get("X-Auth-Request-User")
	email := c.Request().Header.Get("X-Auth-Request-Email")
	
	if userID == "" {
		userID = "admin-001" // Default for testing
		email = "admin@socialtracker.com"
	}
	
	username := email
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
