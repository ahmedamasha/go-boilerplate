package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/refda/backend/internal/pkg/config"
)

type Storage interface {
	Save(file *multipart.FileHeader, folder string) (string, error)
}

type LocalStorage struct {
	basePath string
	baseURL  string
}

func New(cfg *config.Config) (Storage, error) {
	if cfg.Storage.Type != "local" {
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}
	if err := os.MkdirAll(cfg.Storage.LocalPath, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{
		basePath: cfg.Storage.LocalPath,
		baseURL:  strings.TrimRight(cfg.Storage.BaseURL, "/"),
	}, nil
}

func (s *LocalStorage) Save(file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	name := uuid.New().String() + ext
	dir := filepath.Join(s.basePath, folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	destPath := filepath.Join(dir, name)
	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.baseURL, folder, name), nil
}
