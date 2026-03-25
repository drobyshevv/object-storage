package repository

import (
	"context"

	"github.com/drobyshevv/object-storage/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	pool *pgxpool.Pool
}

func NewFileRepository(pool *pgxpool.Pool) *FileRepository {
	return &FileRepository{pool: pool}
}

func (r *FileRepository) Create(ctx context.Context, f *model.File) error {
	query := `
	INSERT INTO files (filename, size, content_type, s3_key, folder)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		f.Filename,
		f.Size,
		f.ContentType,
		f.S3Key,
		f.Folder,
	).Scan(&f.ID, &f.CreatedAt)
}

func (r *FileRepository) GetAll(ctx context.Context, folder string) ([]model.File, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if folder == "" {
		rows, err = r.pool.Query(ctx, `
			SELECT id, filename, size, content_type, s3_key, folder, created_at
			FROM files
			ORDER BY created_at DESC
		`)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, filename, size, content_type, s3_key, folder, created_at
			FROM files
			WHERE folder = $1
			ORDER BY created_at DESC
		`, folder)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]model.File, 0)
	for rows.Next() {
		var f model.File
		if err := rows.Scan(
			&f.ID,
			&f.Filename,
			&f.Size,
			&f.ContentType,
			&f.S3Key,
			&f.Folder,
			&f.CreatedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (r *FileRepository) GetByID(ctx context.Context, id int) (*model.File, error) {
	var f model.File
	query := `
		SELECT id, filename, size, content_type, s3_key, folder, created_at
		FROM files
		WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.Filename,
		&f.Size,
		&f.ContentType,
		&f.S3Key,
		&f.Folder,
		&f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FileRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM files WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
