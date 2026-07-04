package repository

import (
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(dsn string) *sqlx.DB {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		slog.Error("postgres initialization failed",
			slog.String("component", "postgres"),
			slog.String("operation", "open sqlx connection"),
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		slog.Error("postgres initialization failed",
			slog.String("component", "postgres"),
			slog.String("operation", "ping sqlx connection"),
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	return db
}
