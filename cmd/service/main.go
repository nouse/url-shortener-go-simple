package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nouse/url-shortener-go-simple/handlers"
	"github.com/nouse/url-shortener-go-simple/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	rootCtx := context.Background()
	returnCode := 0
	defer func() {
		os.Exit(returnCode)
	}()

	storagePath := filepath.Join(os.Getenv("DATA_PATH"), "urls.txt")
	f, err := os.OpenFile(storagePath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		logger.ErrorContext(rootCtx, "failed to open storage file", "error", err, "path", storagePath, "user", os.Getenv("USER"))
		returnCode = 1
		return
	}
	defer f.Close()

	st, err := storage.NewFileStorage(f)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidFormat) {
			logger.ErrorContext(rootCtx, "invalid format", "error", err)
			logger.InfoContext(rootCtx, "continue with remaining lines", "count", st.Len())
		} else {
			logger.ErrorContext(rootCtx, "unknown error", "error", err)
			returnCode = 1
			return
		}
	}

	shortener := handlers.NewShortener(logger, st)
	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        shortener.Handler(),
	}

	// - SIGINT: Keyboard interrupt, triggered by CTRL-C.
	// - SIGTERM: The default signal to ask the process to terminate gracefully.
	ctx, stop := signal.NotifyContext(rootCtx, syscall.SIGINT, syscall.SIGTERM)
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
