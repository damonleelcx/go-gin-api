package main

import (
	"log"

	"github.com/damonleelcx/go-gin-api/controller"
	"github.com/damonleelcx/go-gin-api/entity"
	"github.com/damonleelcx/go-gin-api/repository"
	"github.com/damonleelcx/go-gin-api/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Auto migrate database tables
	err = db.AutoMigrate(
		&entity.User{},
		&entity.Session{},
		&entity.PasswordResetToken{},
	)
	if err != nil {
		log.Fatal("Database migration failed:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	passwordResetTokenRepo := repository.NewPasswordResetTokenRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo, passwordResetTokenRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)

	// Initialize routes
	router := gin.Default()

	// Register routes
	api := router.Group("/api")
	authController.RegisterRoutes(api)

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server startup failed:", err)
	}
}