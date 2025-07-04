package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
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

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	// log.Debug("Конфигурация сервера", slog.Any("config", cfg))
	log.Debug("debug massages are enable")

	storage := initRepository(cfg, log)
	runMigrations(cfg, log)

	ulrs, err2 := storage.GetAllUrls(context.Background())
	if err2 != nil {
		log.Error(err2.Error())
		return
	}
	fmt.Println(ulrs)
	fmt.Println("до старта сервера")

	err := buildServer(cfg, log, storage)
	if err != nil {
		log.Error("Can`t start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	//log.Error("server stoped")
}

func buildServer(cfg *config.Config, log *slog.Logger, repo repository.Storage) error {

	// Сборка service → handler → router
	pingSvc := service.NewHealthService(repo)
	healthH := health.NewHandler(pingSvc, log)
	shortSvc := service.NewShortenerService(repo, log)
	shortenH := shortener.NewHandler(shortSvc, log)
	apiRouter := api.NewAPIRouter(chi.NewRouter(), log, healthH, shortenH)
	apiRouter.ConfigureRouterField()

	serverAPI := api.New(cfg, log, apiRouter)
	err := serverAPI.Start()
	return err

}

func initRepository(cfg *config.Config, log *slog.Logger) repository.Storage {
	switch cfg.Database.Driver {
	case "memory":
		return memory.NewStorage()
	case "postgres":
		repo, err := postgres.NewPostgresDB(cfg)
		if err != nil {
			log.Error("db connect failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
		return repo
	default:
		log.Error("unknown driver", slog.String("driver", cfg.Database.Driver))
		os.Exit(1)
		return nil
	}
}

func runMigrations(cfg *config.Config, log *slog.Logger) {

	switch cfg.Database.Driver {
	case "postgres":
		if err := migrations.RunMigrations(cfg, log); err != nil {
			log.Error("migrations failed", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}
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
