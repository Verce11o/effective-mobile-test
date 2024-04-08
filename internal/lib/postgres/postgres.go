package postgres

import (
	"context"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func New(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Name, cfg.Postgres.SSLMode))

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	if err := db.Ping(ctx); err != nil {
		log.Fatal("Ping error connecting to database: ", err)
	}

	return db
}
