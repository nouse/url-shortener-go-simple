package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
)

type Storage interface {
	GetURLByCode(ctx context.Context, code string) (ShortURL, error)
	StoreURL(ctx context.Context, url string) (ShortURL, error)
	Increment(ctx context.Context, code string) error
}

var ErrNotFound = errors.New("code not found")

type ShortURL struct {
	Code  string
	URL   string
	Visit int
}

type ShortURLStorage map[string]struct {
	URL   string
	Visit int
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

	fs := &FileStorage{Rand: rng, alphabet: s, List: make(map[string]ShortURL)}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return fs, nil
	}

	var storage ShortURLStorage
	if err := json.Unmarshal(data, &storage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal storage data: %w", err)
	}

	for code, entry := range storage {
		fs.List[code] = ShortURL{
			Code:  code,
			URL:   entry.URL,
			Visit: entry.Visit,
		}
	}

	return fs, nil
}

func (s *FileStorage) GetURLByCode(ctx context.Context, code string) (ShortURL, error) {
	url, ok := s.List[code]
	if !ok {
		return ShortURL{}, ErrNotFound
	}
	return url, nil
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
		b[i] = s.alphabet[s.IntN(32)]
	}
	return string(b)
}
