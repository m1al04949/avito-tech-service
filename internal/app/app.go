package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/m1al04949/avito-tech-service/internal/config"
	"github.com/m1al04949/avito-tech-service/internal/http-server/middleware/mwlog"
	"github.com/m1al04949/avito-tech-service/internal/logger"
	"github.com/m1al04949/avito-tech-service/internal/storage"
	"golang.org/x/exp/slog"
)

func RunServer() error {

	// Config Initializing
	cfg := config.MustLoad()

	// Log Initializing
	log := logger.NewLog(cfg.Env)
	log.Info("starting service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Storage Initializing
	storage := storage.New(cfg.StoragePath, cfg.DatabaseURL)
	if err := storage.Open(); err != nil {
		log.Error("failed to init storage", logger.Err(err))
		return err
	}
	log.Info("storage is initialized")

	// Router Initiziling
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlog.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/", func(r chi.Router) {
		r.Use(middleware.BasicAuth("avito-tech-service", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		// r.Post("/", save.New(log, storage))
		// r.Delete("/{alias}", delete.New(log, storage))
	})

	// router.Delete("/url/{alias}", delete.New(log, storage))

	// Start HTTP Server
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		return err
	}

	log.Error("server stopped")

	return fmt.Errorf("server is stopped")
}