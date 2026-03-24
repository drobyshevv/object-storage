package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func RunMigrations(conn *pgx.Conn) error {
	sqlBytes, err := os.ReadFile("migrations/001_create_files.sql")
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), string(sqlBytes))
	return err
}
