package config

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func NewDatabase(vpr *viper.Viper, logger *slog.Logger) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		vpr.GetString("DB_HOST"),
		vpr.GetString("DB_PORT"),
		vpr.GetString("DB_USER"),
		vpr.GetString("DB_PASSWORD"),
		vpr.GetString("DB_NAME"),
		vpr.GetString("DB_SSLMODE"),
	)

	logger.Info("initializing database connection")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	db.SetMaxOpenConns(vpr.GetInt("DB_MAX_OPEN_CONN"))
	db.SetMaxIdleConns(vpr.GetInt("DB_MAX_IDLE_CONN"))
	db.SetConnMaxLifetime(vpr.GetDuration("DB_CONN_MAX_LIFETIME") * time.Minute)

	logger.Info("database connected successfully")

	return db
}
