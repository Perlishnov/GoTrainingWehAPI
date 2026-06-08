//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/Perlishnov/gotrainingproject/internal/config"
    "github.com/Perlishnov/gotrainingproject/internal/controller"
    "github.com/Perlishnov/gotrainingproject/internal/dao"
    "github.com/Perlishnov/gotrainingproject/internal/middleware"
    "github.com/Perlishnov/gotrainingproject/internal/service"
    "github.com/Perlishnov/gotrainingproject/internal/utils"
    "github.com/Perlishnov/gotrainingproject/pkg/database"
    "github.com/google/wire"
    "github.com/sirupsen/logrus"
)

func InitializeApp() (*App, error) {
    wire.Build(
        config.Load,
        provideLogger,

        // ------------------------------------------------------------------
        // Database connection – uncomment exactly ONE block
        // ------------------------------------------------------------------

        // --- MySQL (default) ---
        // provideDBConfig,                // converts config to DBConfig
        // database.NewMySQLConnection,   // returns *sql.DB
        // dao.NewUserDAO,                // MySQL implementation

        // --- MongoDB (bonus) ---
        database.NewMongoConnection,   // returns *mongo.Database
        dao.NewUserDAOMongo,           // MongoDB implementation

        // ------------------------------------------------------------------

        utils.NewJWTUtil,
        service.NewUserService,
        service.NewAuthService,
        controller.NewUserController,
        controller.NewAuthController,
        middleware.NewAuthMiddleware,
        NewApp,
    )
    return nil, nil
}

// provideLogger creates a logger using the log level from config.
func provideLogger(cfg *config.Config) *logrus.Logger {
    return utils.NewLogger(cfg.LogLevel)
}

// provideDBConfig creates database.DBConfig from config (needed for MySQL only).
func provideDBConfig(cfg *config.Config) database.DBConfig {
    return database.DBConfig{
        User:     cfg.DBUser,
        Password: cfg.DBPassword,
        Host:     cfg.DBHost,
        Port:     cfg.DBPort,
        DBName:   cfg.DBName,
    }
}

type App struct {
    UserController  *controller.UserController
    AuthController  *controller.AuthController
    AuthMiddleware  *middleware.AuthMiddleware
    Logger          *logrus.Logger
}

func NewApp(
    userCtrl *controller.UserController,
    authCtrl *controller.AuthController,
    authMW *middleware.AuthMiddleware,
    logger *logrus.Logger,
) *App {
    return &App{
        UserController:  userCtrl,
        AuthController:  authCtrl,
        AuthMiddleware:  authMW,
        Logger:          logger,
    }
}