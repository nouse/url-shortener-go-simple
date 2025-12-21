package storage

import (
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
