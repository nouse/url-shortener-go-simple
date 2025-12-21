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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	rootCtx := context.Background()

	shortener, err := initService(logger)
	if err != nil {
		logger.ErrorContext(rootCtx, "failed to init service", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        shortener,
	}

	returnCode := 0
	defer func() {
		os.Exit(returnCode)
	}()

	ctx, stop := signal.NotifyContext(rootCtx, os.Interrupt, os.Kill)
	go func() {
		logger.InfoContext(ctx, "Listening on", "port", 8080)
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.ErrorContext(ctx, "Server closed unexpectedly", "error", err)
				returnCode = 1
				stop()
			}
		}
	}()
	<-ctx.Done()

	// Use rootCtx as ctx is already cancelled
	if err := server.Shutdown(rootCtx); err != nil {
		logger.ErrorContext(rootCtx, "server shutdown failed", "error", err)
		returnCode = 1
	}
	logger.InfoContext(rootCtx, "server exited")
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
