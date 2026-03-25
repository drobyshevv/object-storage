package db

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) error {
	// Берем соединение из пула
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release() // обязательно освобождаем соединение обратно в пул

	// Читаем список файлов миграций
	files, err := os.ReadDir("migrations")
	if err != nil {
		return err
	}

	// Сортируем файлы по имени, чтобы миграции шли в правильном порядке
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Выполняем каждую миграцию
	for _, file := range files {
		path := filepath.Join("migrations", file.Name())
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Выполняем SQL через conn.Exec, без лишнего .Conn()
		_, err = conn.Exec(context.Background(), string(sqlBytes))
		if err != nil {
			return err
		}
	}

	return nil
}
