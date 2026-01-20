package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{
				"http://localhost:5173",
				"http://localhost",
			},
			AllowHeaders: []string{
				echo.HeaderAccept,
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAuthorization,
				echo.HeaderAccessControlAllowOrigin,
			},
			AllowCredentials: true,
		},
	),
	)

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

	fe := e.Group("", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderCacheControl, "no-store, max-age=0")

			return next(c)
		}
	})

	staticPath := "/app/public"

	err = filepath.Walk(staticPath,
		func(path string, _ os.FileInfo, _ error) error {
			routePath := path[len(staticPath):]
			fe.GET(routePath, func(c echo.Context) error {
				return c.File(path)
			})

			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	fe.Any("*", func(c echo.Context) error {
		return c.File("/app/public/index.html")
	})
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
