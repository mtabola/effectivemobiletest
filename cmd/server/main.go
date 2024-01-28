package main

import (
	"context"
	"effectivemobiletest/internal/config"
	"effectivemobiletest/internal/logger"
	mwlogger "effectivemobiletest/internal/middleware/mw-logger"
	"effectivemobiletest/internal/storage/pgsql"
	"effectivemobiletest/internal/transport/handlers"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	os.Setenv("CONFIG_PATH", "../../config/.env")

	cfg, err := config.ConfigInit()
	if err != nil {
		log.Fatal("unable to initialize config", err)
	}

	log := logger.LoggerInit(cfg.Env)
	log.Debug("Logger initialized: Debug mode")

	db, err := pgsql.NewStorage(cfg.Database)
	if err != nil {
		log.Error("Database initialization error: ", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/users", func(r chi.Router) {
		r.Get("/", handlers.NewGetter(log, db))
		r.Post("/save", handlers.NewSaver(log, db))
		r.Put("/update", handlers.NewUpdater(log, db))
		r.Delete("/delete", handlers.NewDeleteter(log, db))
	})

	log.Info("Server started", slog.String("address", cfg.Server.Address))

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("start server error: ", slog.String("error", err.Error()))
		}
	}()

	log.Info("Server is running!")

	<-done
	log.Info("Server shutdown starting...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", slog.String("error", err.Error()))
		return
	}

	log.Info("Server stopped!")
}
