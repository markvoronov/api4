package main

import (
	"fmt"
	_ "github.com/go-chi/chi/v5/middleware"
	"github.com/markvoronov/shortener/internal/api"
	"github.com/markvoronov/shortener/internal/config"
	"github.com/markvoronov/shortener/internal/logger/slogpretty"
	"github.com/markvoronov/shortener/internal/repository"
	"github.com/markvoronov/shortener/internal/repository/memory"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Info("Конфигурация сервера", slog.Any("config", cfg))
	log.Debug("debug massages are enable")

	// объявляем storage здесь, в пределах функции main
	var storage repository.Storage

	if cfg.Storage == "memory" {
		storage = memory.NewStorage()
	} else {
		// например, ошибка или default
		log.Info("unknown storage type", "type", cfg.Storage)
		os.Exit(1)
	}

	server := api.New(cfg, log, storage)
	server.Start()

	log.Error("server stopeed")
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
		//log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
