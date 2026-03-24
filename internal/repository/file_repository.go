package repository

import (
	"context"

	"github.com/drobyshevv/object-storage/internal/model"
	"github.com/jackc/pgx/v5"
)

type FileRepository struct {
	conn *pgx.Conn
}

func NewFileRepository(conn *pgx.Conn) *FileRepository {
	return &FileRepository{conn: conn}
}

// --- CREATE ---
func (r *FileRepository) Create(ctx context.Context, f *model.File) error {
	query := `
		INSERT INTO files (filename, size, content_type, s3_key)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.conn.QueryRow(ctx, query,
		f.Filename,
		f.Size,
		f.ContentType,
		f.S3Key, // 🔥 обязательно
	).Scan(&f.ID, &f.CreatedAt)
}

// --- GET ALL ---
func (r *FileRepository) GetAll(ctx context.Context) ([]model.File, error) {
	query := `
		SELECT id, filename, size, content_type, s3_key, created_at
		FROM files
		ORDER BY created_at DESC
	`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []model.File

	for rows.Next() {
		var f model.File

		err := rows.Scan(
			&f.ID,
			&f.Filename,
			&f.Size,
			&f.ContentType,
			&f.S3Key, // 🔥 обязательно
			&f.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		files = append(files, f)
	}

	return files, nil
}

// --- GET BY ID ---
func (r *FileRepository) GetByID(ctx context.Context, id int) (*model.File, error) {
	query := `
		SELECT id, filename, size, content_type, s3_key, created_at
		FROM files
		WHERE id = $1
	`

	var f model.File

	err := r.conn.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.Filename,
		&f.Size,
		&f.ContentType,
		&f.S3Key, // 🔥 обязательно
		&f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// --- DELETE ---
func (r *FileRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := r.conn.Exec(ctx, query, id)
	return err
}
