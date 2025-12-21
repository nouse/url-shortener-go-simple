package storage

import (
	"context"
	"io"
)

type Storage interface {
	GetURLByCode(ctx context.Context, code string) (ShortURL, error)
	StoreURL(ctx context.Context, url string) (ShortURL, error)
	Increment(ctx context.Context, code string) error
}

type ShortURL struct {
	Code       string
	URL        string
	VisitCount int
}

type FileStorage struct {
	ShortURLList []ShortURL
}

func NewFileStorage(reader io.Reader) (*FileStorage, error) {
	return &FileStorage{}, nil
}

func (s *FileStorage) GetURLByCode(ctx context.Context, code string) (ShortURL, error) {
	return ShortURL{}, nil
}

func (s *FileStorage) StoreURL(ctx context.Context, url string) (ShortURL, error) {
	return ShortURL{}, nil
}

func (s *FileStorage) Increment(ctx context.Context, code string) error {
	return nil
}
