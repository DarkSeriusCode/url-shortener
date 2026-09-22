package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/DarkSeriusCode/url-shortener/internal/config"
	"github.com/DarkSeriusCode/url-shortener/internal/handler"
	"github.com/DarkSeriusCode/url-shortener/internal/service"
	"github.com/DarkSeriusCode/url-shortener/internal/storage"
	"github.com/lmittmann/tint"
)

func setupLogger(env config.ServerEnv) {
	var logger *slog.Logger

	switch env {
	case config.EnvLocal:
		opt := tint.Options{
			Level: slog.LevelDebug,
		}
		logger = slog.New(tint.NewTextHandler(os.Stdout, &opt))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)
}

func initStorage(cfg *config.ServerConfig) storage.Storage {
	var s storage.Storage
	var err error

	switch cfg.Server.Storage {
	case config.StorageMemory:
		s = storage.NewMemoryStorage()
		err = nil
	case config.StoragePostgres:
		s, err = storage.NewPostgresStorage(cfg)
	}

	if err != nil {
		panic(err)
	}
	slog.Info(fmt.Sprintf("Using %s as a storage", cfg.Server.Storage))
	return s
}

func main() {
	// CLI
	configPath := flag.String("config", "config/config.toml", "Path to the server config")
	flag.Parse()

	// Config
	cfg, err := config.LoadConfigFromFile(*configPath)
	if err != nil {
		panic(err)
	}

	// Logger
	setupLogger(cfg.Server.Env)

	storage := initStorage(cfg)
	service := service.NewShortener(storage)

	// Handler
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler.NewHandler(service),
	}

	slog.Info(fmt.Sprintf("Running %s server on %s", cfg.Server.Env, srv.Addr))

	if err = srv.ListenAndServe(); err != nil {
		slog.Error(fmt.Sprintf("Server error: %s", err))
	}
}
