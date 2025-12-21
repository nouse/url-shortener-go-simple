package storage

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

const fileContent = `{"code": "aaa", "url": "http://example.com/v1"}
{"code": bbb", "url": "http://example.com/v2"}
{"code": "bbb", "url": "http://example.com/v2", "visit": 2}`

func TestNewFileStorage(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		fs, err := NewFileStorage(bytes.NewBufferString(""))
		if err != nil {
			t.Fatalf("NewFileStorage() error = %v", err)
		}
		if len(fs.list) != 0 {
			t.Errorf("NewFileStorage() got = %v, want %v", len(fs.list), 0)
		}
	})
	t.Run("test with fileContent", func(t *testing.T) {
		fs, err := NewFileStorage(bytes.NewBufferString(fileContent))
		if !errors.Is(err, ErrInvalidFormat) {
			t.Fatalf("NewFileStorage() error = %v", err)
		}
		if len(fs.list) != 2 {
			t.Errorf("NewFileStorage() got = %v, want %v", len(fs.list), 2)
		}
	})
}

func TestFileStorageShortCodeCollision(t *testing.T) {
	fs, err := NewFileStorage(bytes.NewBufferString(""))
	if err != nil {
		t.Fatalf("NewFileStorage() error = %v", err)
	}

	// generate 1000 codes and check uniqueness
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
	fs, err := NewFileStorage(bytes.NewBufferString(fileContent))
	if !errors.Is(err, ErrInvalidFormat) {
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
			s, err := fs.GetURLByCode(tc.code)
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

func TestFileStorage_StoreURL(t *testing.T) {
	buf := bytes.NewBufferString(fileContent)
	fs, err := NewFileStorage(buf)
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("NewFileStorage() error = %v", err)
	}
	url := "http://example.com/v3"
	r, err := fs.StoreURL(url)
	if err != nil {
		t.Fatalf("StoreURL() error = %v", err)
	}
	if r.URL != url {
		t.Errorf("StoreURL() got = %v, want %v", r.URL, url)
	}
	if r.Visit != 0 {
		t.Errorf("StoreURL() got = %v, want %v", r.Visit, 0)
	}
	if strings.Count(buf.String(), "\n") != 1 {
		t.Errorf("StoreURL() should increase a new line, current content %s", buf.String())
	}
}

func TestFileStorage_Increment(t *testing.T) {
	buf := bytes.NewBufferString(fileContent)
	fs, err := NewFileStorage(buf)
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("NewFileStorage() error = %v", err)
	}

	// TODO change to testcases
	err = fs.Increment("aaa")
	if err != nil {
		t.Fatalf("Increment() error = %v", err)
	}
	r, _ := fs.GetURLByCode("aaa")
	if r.URL != "http://example.com/v1" {
		t.Errorf("Increment() got = %v, want %v", r.URL, "http://example.com/v1")
	}
	if r.Visit != 1 {
		t.Errorf("Increment() got = %v, want %v", r.Visit, 1)
	}
	if strings.Count(buf.String(), "\n") != 1 {
		t.Errorf("Increment() should increase a new line, current content %s", buf.String())
	}

	err = fs.Increment("bbb")
	if err != nil {
		t.Fatalf("Increment() error = %v", err)
	}
	r, _ = fs.GetURLByCode("bbb")
	if r.URL != "http://example.com/v2" {
		t.Errorf("Increment() got = %v, want %v", r.URL, "http://example.com/v2")
	}
	if r.Visit != 3 {
		t.Errorf("Increment() got = %v, want %v", r.Visit, 3)
	}
	if strings.Count(buf.String(), "\n") != 2 {
		t.Errorf("Increment() should increase a new line, current content %s", buf.String())
	}

	err = fs.Increment("ccc")
	if err == nil {
		t.Error("Incremenet() should return an error")
	}
}
