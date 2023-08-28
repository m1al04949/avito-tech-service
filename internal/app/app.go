package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/m1al04949/avito-tech-service/internal/config"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/addtouser"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/adduser"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/createsegment"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/deletefromuser"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/deletesegment"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/deleteuser"
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
	store := storage.New(cfg.StoragePath, cfg.DatabaseURL)
	if err := store.Open(); err != nil {
		log.Error("failed to init storage", logger.Err(err))
		return err
	}
	defer store.Close()
	if err := store.CreateTabs(); err != nil {
		log.Error("failed to init tabs", logger.Err(err))
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

		r.Post("/segment_manage", createsegment.NewSegment(log, store))          // Add Segment
		r.Post("/users/id={id}", adduser.AddUser(log, store))                    // Add User
		r.Post("/user_manage/{id}", addtouser.AddToUser(log, store))             // Add Segment To User
		r.Delete("/segment_manage", deletesegment.DelSegment(log, store))        // Delete Segment
		r.Delete("/users/id={id}", deleteuser.DeleteUser(log, store))            // Delete User
		r.Delete("/user_manage/{id}", deletefromuser.DeleteFromUser(log, store)) // Delete Segment From User
	})

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
