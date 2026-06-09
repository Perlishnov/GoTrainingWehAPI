package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Perlishnov/gotrainingproject/wire"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/Perlishnov/gotrainingproject/docs" // swagger docs
)

// @title           Go REST API
// @version         1.0
// @description     Production-ready REST API with JWT auth, role-based access, and MySQL.
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
func main() {
    app, err := wire.InitializeApp()
    if err != nil {
        panic(err)
    }
    logger := app.Logger

    // Set Gin mode (release or debug)
    gin.SetMode(gin.ReleaseMode)

    router := gin.New()
    router.Use(gin.Recovery())
    router.Use(gin.Logger())

    // Swagger endpoint
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Public routes
    router.POST("/auth/login", app.AuthController.Login)
    router.POST("/auth/logout", app.AuthController.Logout)
    router.POST("/users", app.UserController.CreateUser)

    // Protected routes group
    authGroup := router.Group("/")
    authGroup.Use(app.AuthMiddleware.Authenticate())
    {
        authGroup.GET("/users/me", app.UserController.GetCurrentUser)
        authGroup.GET("/users/:id", app.UserController.GetUserByID)
        authGroup.GET("/users", app.UserController.GetAllUsers)
    }

    srv := &http.Server{
        Addr:    ":" + os.Getenv("SERVER_PORT"),
        Handler: router,
    }

    go func() {
        logger.Infof("Server starting on port %s", os.Getenv("SERVER_PORT"))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatalf("server failed: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        logger.Errorf("Server forced to shutdown: %v", err)
    }
    logger.Info("Server exited")
}