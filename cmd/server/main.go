package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    chiMiddleware "github.com/go-chi/chi/v5/middleware"
    httpSwagger "github.com/swaggo/http-swagger"
    "github.com/Perlishnov/gotrainingproject/wire"
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

    r := chi.NewRouter()
    r.Use(chiMiddleware.RequestID)
    r.Use(chiMiddleware.RealIP)
    r.Use(chiMiddleware.Recoverer)
    r.Use(chiMiddleware.Logger)

    // Swagger UI endpoint
    r.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("/docs/swagger.json"), // point to where the JSON is served
    ))

    r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))


    // Public routes
    r.Post("/auth/login", app.AuthController.Login)
    r.Post("/auth/logout", app.AuthController.Logout)
    r.Post("/users", app.UserController.CreateUser)

    // Protected routes
    r.Group(func(r chi.Router) {
        r.Use(app.AuthMiddleware.Authenticate)
        r.Get("/users/me", app.UserController.GetCurrentUser)
        r.Get("/users/{id}", app.UserController.GetUserByID)
        r.Get("/users", app.UserController.GetAllUsers)
    })

    srv := &http.Server{
        Addr:    ":" + os.Getenv("SERVER_PORT"),
        Handler: r,
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