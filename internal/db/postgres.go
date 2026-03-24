package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func New(dbURL string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), dbURL)
}
