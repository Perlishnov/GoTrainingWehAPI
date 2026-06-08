package config

import (
    "os"
    "strconv"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
)

type Config struct {
    ServerPort         string
    DBUser             string
    DBPassword         string
    DBHost             string
    DBPort             string
    DBName             string
    JWTSecret          string
    JWTExpirationHours int
    LogLevel           logrus.Level
}

func Load() (*Config, error) {
    _ = godotenv.Load() // ignore if missing

    level, err := logrus.ParseLevel(getEnv("LOG_LEVEL", "info"))
    if err != nil {
        level = logrus.InfoLevel
    }

    jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))

    return &Config{
        ServerPort:         getEnv("SERVER_PORT", "8080"),
        DBUser:             getEnv("DB_USER", "root"),
        DBPassword:         getEnv("DB_PASSWORD", ""),
        DBHost:             getEnv("DB_HOST", "127.0.0.1"),
        DBPort:             getEnv("DB_PORT", "3306"),
        DBName:             getEnv("DB_NAME", "go_api_db"),
        JWTSecret:          getEnv("JWT_SECRET", "change-me"),
        JWTExpirationHours: jwtExpHours,
        LogLevel:           level,
    }, nil
}

func getEnv(key, defaultValue string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultValue
}