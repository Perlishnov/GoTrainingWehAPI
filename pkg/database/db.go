package database

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/sirupsen/logrus"
)

type DBConfig struct {
    User     string
    Password string
    Host     string
    Port     string
    DBName   string
}

func NewMySQLConnection(cfg DBConfig, logger *logrus.Logger) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }
    logger.Info("MySQL connection established")
    return db, nil
}