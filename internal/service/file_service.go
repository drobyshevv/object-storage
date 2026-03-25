package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/drobyshevv/object-storage/internal/model"
	"github.com/drobyshevv/object-storage/internal/repository"
	"github.com/drobyshevv/object-storage/internal/storage"
	"github.com/google/uuid"
)

type FileService struct {
	repo    *repository.FileRepository
	storage *storage.Storage
}

func NewFileService(r *repository.FileRepository, s *storage.Storage) *FileService {
	return &FileService{repo: r, storage: s}
}

func (s *FileService) Upload(ctx context.Context, file *multipart.FileHeader, folder string) (*model.File, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	key := uuid.New().String() + "_" + file.Filename

	if folder != "" {
		key = folder + "/" + key
	}

	err = s.storage.Upload(ctx, key, src)
	if err != nil {
		return nil, err
	}

	f := &model.File{
		Filename:    file.Filename,
		Size:        file.Size,
		ContentType: file.Header.Get("Content-Type"),
		S3Key:       key,
		Folder:      folder,
	}

	err = s.repo.Create(ctx, f)
	return f, err
}

func (s *FileService) GetAll(ctx context.Context, folder string) ([]model.File, error) {
	files, err := s.repo.GetAll(ctx, folder)
	if files == nil {
		files = []model.File{}
	}
	return files, err
}

func (s *FileService) GetFile(ctx context.Context, id int) (io.ReadCloser, string, error) {
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, "", err
	}

	body, err := s.storage.Download(ctx, f.S3Key)
	return body, f.Filename, err
}

func (s *FileService) Delete(ctx context.Context, id int) error {
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.storage.Delete(ctx, f.S3Key)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
