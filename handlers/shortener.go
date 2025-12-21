package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nouse/url-shortener-go-simple/storage"
)

type Shortener struct {
	Storage storage.Storage
	Logger  *slog.Logger
	*http.ServeMux
}

type encodeParams struct {
	URL string
}

func (s *Shortener) Encode(w http.ResponseWriter, r *http.Request) {
	p := encodeParams{}
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to decode request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("failed to decode request body")
		return
	}
}

func (s *Shortener) Get(w http.ResponseWriter, r *http.Request) {
}

func (s *Shortener) GetInfo(w http.ResponseWriter, r *http.Request) {
}

func heartbeat(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "OK")
}

func NewShortener(logger *slog.Logger, storage storage.Storage) http.Handler {
	s := &Shortener{
		ServeMux: http.NewServeMux(),
		Logger:   logger,
		Storage:  storage,
	}
	// Encode url to short code
	s.HandleFunc("POST /x", s.Encode)
	// Return url by short code, increment visitCount
	s.HandleFunc("GET /x/{code}", s.Get)
	// Return url and visitCount of short code
	s.HandleFunc("GET /info/{code}", s.GetInfo)
	// Heartbeat
	s.HandleFunc("GET /ping", heartbeat)
	return s
}
