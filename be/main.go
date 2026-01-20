package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Armatorix/SocialTracker/be/handlers"
	"github.com/Armatorix/SocialTracker/be/migrations"
	"github.com/Armatorix/SocialTracker/be/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://socialtracker:socialtracker@localhost:5432/socialtracker?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Run migrations
	if err := migrations.Run(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Initialize repository and handlers
	repo := repository.NewRepository(db)
	h := handlers.NewHandler(repo)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API routes
	api := e.Group("/api")
	
	// User routes
	api.GET("/user", h.GetCurrentUser)
	
	// Social accounts routes
	api.GET("/social-accounts", h.GetSocialAccounts)
	api.POST("/social-accounts", h.CreateSocialAccount)
	api.DELETE("/social-accounts/:id", h.DeleteSocialAccount)
	api.POST("/social-accounts/:id/pull", h.PullContentFromPlatform)
	
	// Content routes
	api.GET("/content", h.GetContent)
	api.POST("/content", h.CreateContent)
	api.DELETE("/content/:id", h.DeleteContent)
	
	// Admin routes
	api.GET("/admin/content", h.GetAllContent)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
