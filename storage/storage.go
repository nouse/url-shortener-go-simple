package storage

import (
	"context"
	"io"
	"math/rand/v2"
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
	List     map[string]ShortURL
	alphabet []byte
	*rand.Rand
}

const (
	codeLength     = 6
	base32Alphabet = "abcdefghijklmnopqrstuvwxyz234567"
)

func NewFileStorage(reader io.Reader) (*FileStorage, error) {
	rng := rand.New(rand.NewPCG(13, 37))
	s := []byte(base32Alphabet)
	rng.Shuffle(32, func(i, j int) { s[i], s[j] = s[j], s[i] })
	return &FileStorage{Rand: rng, alphabet: s}, nil
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

func (s *FileStorage) shortID() string {
	b := make([]byte, codeLength)
	for i := range codeLength {
		b[i] = (s.alphabet[s.IntN(32)])
	}
	return string(b)
}
