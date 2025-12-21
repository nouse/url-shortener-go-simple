package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"strings"
)

type Storage interface {
	GetURLByCode(code string) (ShortURL, error)
	StoreURL(url string) (ShortURL, error)
	Increment(code string) error
}

var (
	ErrNotFound      = errors.New("code not found")
	ErrInvalidFormat = errors.New("invalid format")
)

type ShortURL struct {
	Code  string
	URL   string
	Visit int
}

type FileStorage struct {
	list     map[string]ShortURL
	alphabet []byte
	rand     *rand.Rand
	rwCloser io.ReadWriter
}

const (
	codeLength     = 6
	base32Alphabet = "abcdefghijklmnopqrstuvwxyz234567"
)

func NewFileStorage(rwCloser io.ReadWriter) (*FileStorage, error) {
	rng := rand.New(rand.NewPCG(13, 37))
	s := []byte(base32Alphabet)
	rng.Shuffle(32, func(i, j int) { s[i], s[j] = s[j], s[i] })

	fs := &FileStorage{
		list:     make(map[string]ShortURL),
		rand:     rng,
		alphabet: s,
		rwCloser: rwCloser,
	}

	var errLines []string
	buf := bufio.NewScanner(rwCloser)
	for buf.Scan() {
		var s ShortURL
		if err := json.Unmarshal(buf.Bytes(), &s); err != nil {
			errLines = append(errLines, buf.Text())
			continue
		}
		fs.list[s.Code] = s
	}

	if len(errLines) > 0 {
		return fs, fmt.Errorf("lines: %s, err: %w", strings.Join(errLines, "\n"), ErrInvalidFormat)
	}
	return fs, nil
}

func (s *FileStorage) GetURLByCode(code string) (ShortURL, error) {
	url, ok := s.list[code]
	if !ok {
		return ShortURL{}, ErrNotFound
	}
	return url, nil
}

// StoreURL save url and code by appending a new line to writer
func (s *FileStorage) StoreURL(url string) (ShortURL, error) {
	return ShortURL{}, nil
}

// Increment the visit count and append a new line
func (s *FileStorage) Increment(code string) error {
	return nil
}

func (s *FileStorage) Len() int {
	return len(s.list)
}

func (s *FileStorage) shortID() string {
	b := make([]byte, codeLength)
	for i := range codeLength {
		b[i] = s.alphabet[s.rand.IntN(32)]
	}
	return string(b)
}
