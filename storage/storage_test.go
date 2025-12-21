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
		if len(fs.ShortURLList) != 0 {
			t.Errorf("NewFileStorage() got = %v, want %v", len(fs.ShortURLList), 0)
		}
	})
}
