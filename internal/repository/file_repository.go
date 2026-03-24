package repository

import (
	"context"

	"github.com/drobyshevv/object-storage/internal/model"
	"github.com/jackc/pgx/v5"
)

type FileRepository struct {
	db *pgx.Conn
}

func NewFileRepository(db *pgx.Conn) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, f *model.File) error {
	query := `
	INSERT INTO files (filename, size, content_type, s3_key)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		f.Filename,
		f.Size,
		f.ContentType,
		f.S3Key,
	).Scan(&f.ID, &f.CreatedAt)
}

func (r *FileRepository) GetAll(ctx context.Context) ([]model.File, error) {
	rows, err := r.db.Query(ctx, `SELECT id, filename, size, created_at FROM files`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []model.File

	for rows.Next() {
		var f model.File
		err := rows.Scan(&f.ID, &f.Filename, &f.Size, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (r *FileRepository) GetByID(ctx context.Context, id int) (*model.File, error) {
	var f model.File

	err := r.db.QueryRow(ctx,
		`SELECT id, filename, s3_key FROM files WHERE id=$1`, id,
	).Scan(&f.ID, &f.Filename, &f.S3Key)

	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (r *FileRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM files WHERE id=$1`, id)
	return err
}
