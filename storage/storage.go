package storage

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
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
	ErrDuplicateCode = errors.New("code collision")
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
	rw       io.ReadWriter
}

const (
	codeLength     = 6
	base32Alphabet = "abcdefghijklmnopqrstuvwxyz234567"
	randSeed       = ")Bo_ItkHpnwoM7PiK9\\J]QTER\\uGB#2" // sf-pwgen -l 32 -a random
)

func NewFileStorage(rw io.ReadWriter) (*FileStorage, error) {
	var seedBytes [32]byte
	if seed := os.Getenv("SEED"); seed != "" {
		seedBytes = sha256.Sum256([]byte(seed))
	} else {
		seedBytes = sha256.Sum256([]byte(randSeed))
	}
	rng := rand.New(rand.NewChaCha8(seedBytes))
	s := []byte(base32Alphabet)
	rng.Shuffle(32, func(i, j int) { s[i], s[j] = s[j], s[i] })

	fs := &FileStorage{
		list:     make(map[string]ShortURL),
		rand:     rng,
		alphabet: s,
		rw:       rw,
	}

	var errLines []string
	buf := bufio.NewScanner(rw)
	for buf.Scan() {
		b := buf.Bytes()
		if len(bytes.TrimSpace(b)) == 0 {
			continue
		}
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
	code := s.shortID()
	if _, ok := s.list[code]; ok {
		return ShortURL{}, ErrDuplicateCode
	}

	shortURL := ShortURL{
		Code: code,
		URL:  url,
	}
	data, _ := json.Marshal(shortURL)
	if _, err := s.rw.Write(append(data, '\n')); err != nil {
		return ShortURL{}, err
	}

	s.list[code] = shortURL
	return shortURL, nil
}

// Increment the visit count and append a new line
func (s *FileStorage) Increment(code string) error {
	shortURL, ok := s.list[code]
	if !ok {
		return ErrNotFound
	}

	shortURL.Visit++

	data, _ := json.Marshal(shortURL)
	if _, err := s.rw.Write(append(data, '\n')); err != nil {
		return err
	}

	s.list[code] = shortURL
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
