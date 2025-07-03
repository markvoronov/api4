package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5/middleware"
	"github.com/markvoronov/shortener/internal/api"
	"github.com/markvoronov/shortener/internal/api/health"
	"github.com/markvoronov/shortener/internal/api/shortener"
	"github.com/markvoronov/shortener/internal/config"
	"github.com/markvoronov/shortener/internal/logger/slogpretty"
	"github.com/markvoronov/shortener/internal/repository"
	"github.com/markvoronov/shortener/internal/repository/memory"
	"github.com/markvoronov/shortener/internal/repository/postgres"
	"github.com/markvoronov/shortener/internal/service"
	"github.com/markvoronov/shortener/migrations"
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

	// 1️⃣ ➜ вызываем миграции
	if err := migrations.RunMigrations(cfg, log); err != nil {
		log.Error("migrations failed", slog.Any("error", err.Error()))
		os.Exit(1)
	}

	log.Info("migrations applyied")

	var (
		repo repository.Repo
		err  error
	)
	switch cfg.Database.Driver {
	case "memory":
		repo = memory.NewStorage() // map[string]URL
	case "postgres":
		repo, err = postgres.NewPostgresDB(cfg)
		if err != nil {
			log.Info("Can`t start postgres db", slog.Any("error", err.Error()))
			os.Exit(1)
		}
		log.Info("postrgres db connected")
	default:
		log.Info("unknown database driver: %s", slog.String("Driver", cfg.Database.Driver))
		os.Exit(1)
	}

	// сервис и HTTP-хендлер для health
	pingSvc := service.NewHealthService(repo)
	healthH := health.NewHandler(pingSvc, log)

	shortenSvc := service.NewShortenerService(repo)
	shortenH := shortener.NewHandler(shortenSvc, log)

	apiRouter := api.NewAPIRouter(chi.NewRouter(), log, healthH, shortenH)

	server := api.New(cfg, log, apiRouter)
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
