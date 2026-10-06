package main

import (
	"FSchedule/api"
	"FSchedule/database"
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// database
	const connStr = "postgres://go_user:machine_banana@localhost:5432/f_schedule"
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	initCtx, cancelInit := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelInit()
	pool, err := pgxpool.New(initCtx, connStr)
	if err != nil {
		logger.Error("pgxpool.New", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(initCtx); err != nil {
		logger.Error("db ping", "err", err)
		os.Exit(1)
	}
	storage := database.NewStorage(pool)
	handler := api.NewHandler(storage, logger)
	router := api.NewRouter(handler, logger)
	// app

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return context.Background()
		},
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		logger.Error("server failed", "err", err)
		os.Exit(1)
	case sig := <-stop:
		logger.Info("shutdown signal received", "signal", sig.String())
	}

	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(initCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
		_ = srv.Close()
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}
