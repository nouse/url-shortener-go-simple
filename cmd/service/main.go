package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nouse/url-shortener-go-simple/handlers"
	"github.com/nouse/url-shortener-go-simple/storage"
)

func main() {
	returnCode := 0
	defer func() {
		os.Exit(returnCode)
	}()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	shortener, err := initService(logger)
	if err != nil {
		logger.ErrorContext(ctx, "failed to init service", "error", err)
		returnCode = 1
		return
	}

	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        shortener,
	}

	go func() {
		logger.InfoContext(ctx, "Listening on", "port", 8080)
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.ErrorContext(ctx, "Server closed unexpectedly", "error", err)
			}
		}
	}()

	<-ctx.Done()
	if err := server.Shutdown(ctx); err != nil {
		logger.ErrorContext(ctx, "server shutdown failed", "error", err)
		returnCode = 1
	}
	logger.InfoContext(ctx, "server exited")
}

func initService(logger *slog.Logger) (http.Handler, error) {
	f, err := os.Open("urls.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err = os.Create("urls.txt")
		}
		if err != nil { // reevaluate error after os.Create
			return nil, fmt.Errorf("failed to open urls.txt: %w", err)
		}
	}
	s, err := storage.NewFileStorage(f)
	if err != nil { // reevaluate error after os.Create
		return nil, fmt.Errorf("failed to parse urls.txt: %w", err)
	}
	return handlers.NewShortener(logger, s), nil
}
