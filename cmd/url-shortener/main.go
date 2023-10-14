package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"
	"url-shortener/internal/config"
	slogpkg "url-shortener/internal/packages/logger/slog"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slogpkg.Err(err))
		os.Exit(1)
	}

	err = storage.Save("https://google.com", "google")
	if err != nil {
		log.Error("failed to save url", slogpkg.Err(err))
		os.Exit(1)
	}

	url, err := storage.GetUrlByAlias("google")
	fmt.Println(url)
	if err != nil {
		log.Error("failed to get url by alias", slogpkg.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
