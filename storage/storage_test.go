package storage

import (
	"errors"
	"strings"
	"testing"
)

func TestNewFileStorage(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		s := strings.NewReader("")
		fs, err := NewFileStorage(s)
		if err != nil {
			t.Fatalf("NewFileStorage() error = %v", err)
		}
		if len(fs.List) != 0 {
			t.Errorf("NewFileStorage() got = %v, want %v", len(fs.List), 0)
		}
	})
}

func TestFileStorageShortCode(t *testing.T) {
	fs, err := NewFileStorage(strings.NewReader(""))
	if err != nil {
		t.Fatalf("NewFileStorage() error = %v", err)
	}

	// generate 10 codes and check uniqueness
	codes := make(map[string]bool, 1000)
	for i := range 1000 {
		c := fs.shortID()
		if _, ok := codes[c]; ok {
			t.Errorf("shortID() collides after %d, code: %s", i, c)
			break
		}
		codes[c] = true
	}
}

func TestFileStorage_GetURLByCode(t *testing.T) {
	j := `{
  "aaa": { "url" : "http://example.com/v1" },
  "bbb": { "url" : "http://example.com/v2", "visit": 2 }
}`
	fs, err := NewFileStorage(strings.NewReader(j))
	if err != nil {
		t.Fatalf("NewFileStorage() error = %v", err)
	}

	testCases := []struct {
		code,
		url string
		visit int
		err   error
	}{
		{
			code: "aaa",
			url:  "http://example.com/v1",
		},
		{
			code:  "bbb",
			url:   "http://example.com/v2",
			visit: 2,
		},
		{
			code: "ccc",
			url:  "",
			err:  ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run("fetching "+tc.code, func(t *testing.T) {
			s, err := fs.GetURLByCode(t.Context(), tc.code)
			if !errors.Is(err, tc.err) {
				t.Fatalf("GetURLByCode() error = %v, want %v", err, tc.err)
			}
			if s.URL != tc.url {
				t.Errorf("GetURLByCode() got = %v, want %v", s.URL, tc.url)
			}
			if s.Visit != tc.visit {
				t.Errorf("GetURLByCode() got = %v, want %v", s.Visit, tc.visit)
			}
		})
	}
}
