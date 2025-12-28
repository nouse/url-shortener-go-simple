package handlers

import (
	"log/slog"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
	"github.com/nouse/url-shortener-go-simple/storage"
)

type Shortener struct {
	storage storage.Storage
	logger  *slog.Logger
	mux     *http.ServeMux
}

type encodeParams struct {
	URL string
}

func NewShortener(logger *slog.Logger, storage storage.Storage) *Shortener {
	s := &Shortener{
		logger:  logger,
		storage: storage,
		mux:     http.NewServeMux(),
	}
	s.mux.HandleFunc("POST /x", s.Encode)
	s.mux.HandleFunc("GET /x/{code}", s.Get)
	s.mux.HandleFunc("GET /info/{code}", s.GetInfo)
	s.mux.HandleFunc("GET /ping", heartbeat)

	s.mux.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	return s
}

func (s *Shortener) Handler() http.Handler {
	return s.mux
}
