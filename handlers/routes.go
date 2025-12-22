package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nouse/url-shortener-go-simple/storage"
)

func (s *Shortener) Encode(w http.ResponseWriter, r *http.Request) {
	p := encodeParams{}
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to decode request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "failed to decode request body")
		return
	}
	_, err = url.Parse(p.URL)
	if err != nil {
		s.logger.ErrorContext(ctx, "invalid URL", "url", p.URL, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "invalid URL")
		return
	}
	shortURL, err := s.storage.StoreURL(p.URL)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to store URL", "url", p.URL, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "failed to encode URL")
		return
	}
	if err := json.NewEncoder(w).Encode(shortURL); err != nil {
		s.logger.ErrorContext(ctx, "failed to encode response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Get returns url of short code, and increment visit by 1
func (s *Shortener) Get(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	ctx := r.Context()

	shortURL, err := s.storage.GetURLByCode(code)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			s.logger.WarnContext(ctx, "code not found", "code", code)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.logger.ErrorContext(ctx, "failed to get URL", "error", err, "code", code)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.storage.Increment(code); err != nil {
		s.logger.ErrorContext(ctx, "failed to increment visit count", "error", err, "code", code)
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.Redirect(w, r, shortURL.URL, http.StatusPermanentRedirect)
}

// GetInfo returns url and visit of code
func (s *Shortener) GetInfo(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	ctx := r.Context()

	shortURL, err := s.storage.GetURLByCode(code)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			s.logger.WarnContext(ctx, "code not found", "code", code)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.logger.ErrorContext(ctx, "failed to get URL", "error", err, "code", code)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shortURL); err != nil {
		s.logger.ErrorContext(ctx, "failed to encode response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func heartbeat(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "OK")
}
