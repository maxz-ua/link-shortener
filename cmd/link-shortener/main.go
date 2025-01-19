package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"link-shortener/internal/config"
	"link-shortener/internal/http-server/handlers/url/save"
	"link-shortener/internal/lib/logger/sl"
	"link-shortener/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
	_ "time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoadConfig()

	// Validate environment variable
	validEnvs := map[string]bool{
		envLocal: true,
		envDev:   true,
		envProd:  true,
	}

	if !validEnvs[cfg.Env] {
		println("Invalid environment configuration:", cfg.Env)
		println("Expected one of: local, dev, or prod.")
		os.Exit(1)
	}

	log := setupLogger(cfg.Env)
	log.Info("starting link-shortener")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("error opening storage", sl.Err(err))
		os.Exit(1)
	} else {
		log.Info("storage opened successfully")
	}

	_ = storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error starting server", sl.Err(err))
	}

	log.Error("shutting down")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
