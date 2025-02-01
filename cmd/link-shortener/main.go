package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"link-shortener/internal/config"
	"link-shortener/internal/http-server/handlers/redirect"
	"link-shortener/internal/http-server/handlers/url/delete"
	"link-shortener/internal/http-server/handlers/url/save"
	mwLogger "link-shortener/internal/http-server/middleware/logger"
	"link-shortener/internal/lib/logger/sl"
	"link-shortener/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	_ "time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoadConfig()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting link-shortener",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error opening storage", sl.Err(err))
		os.Exit(1)
	} else {
		log.Info("storage opened successfully")
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("link-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{id}", delete.New(log, storage)) // Delete by ID
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	// add
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// -

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
	default: // Set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
