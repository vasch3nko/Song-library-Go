package app

import (
    "context"
    "errors"
    "github.com/joho/godotenv"
    "github.com/vasch3nko/songlibrary/internal/api"
    "github.com/vasch3nko/songlibrary/internal/config"
    "github.com/vasch3nko/songlibrary/internal/services"
    "github.com/vasch3nko/songlibrary/internal/storage"
    "log/slog"
    "net/http"
    "os"
)

const (
    envDev  = "dev"
    envProd = "prod"
)

func Run(ctx context.Context) error {
    if err := godotenv.Load(); err != nil {
        return err
    }

    cfg := config.NewConfig()
    if err := cfg.LoadFromEnv(); err != nil {
        return err
    }

    log, err := setupLogger(cfg.Env)
    if err != nil {
        return err
    }

    log.Info("Starting song library app", slog.String("env", cfg.Env))

    // Postgres storage initialization and connecting
    store, err := storage.NewPostgresStore(
        cfg.Db.Host,
        cfg.Db.Port,
        cfg.Db.Username,
        cfg.Db.Password,
        cfg.Db.Database,
        cfg.Db.SSLMode,
        log,
    )
    if err != nil {
        log.Error("Failed to init storage", slog.String("error", err.Error()))
        return err
    }

    // Postgres storage migrating
    if err = store.Migrate(cfg.Db.MigrationsPath); err != nil {
        log.Error("Failed to migrate storage", slog.String("error", err.Error()))
        return err
    }

    songService := services.NewSongService(store, cfg.SongDetailsApiUrl, log)

    mux := api.NewLoggingMux(log)
    api.NewSongHandler(songService, mux).RegisterSongRoutes()

    srv := &http.Server{
        Addr:         cfg.Server.Addr,
        Handler:      mux,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
        IdleTimeout:  cfg.Server.IdleTimeout,
        ErrorLog:     slog.NewLogLogger(log.Handler(), slog.LevelInfo),
    }

    // Goroutine that handles an interrupt
    go func() {
        <-ctx.Done()
        log.Info("Shutting down server")
        if err = srv.Shutdown(ctx); err != nil {
            log.Error("error", err)
        }
    }()

    log.Info("Starting server", slog.String("addr", cfg.Server.Addr))
    err = srv.ListenAndServe()
    if err != nil && !errors.Is(err, http.ErrServerClosed) {
        log.Error("Unexpected server shutdown", slog.String("error", err.Error()))
        return err
    }

    return nil
}

func setupLogger(env string) (*slog.Logger, error) {
    var logger *slog.Logger

    switch env {
    case envDev:
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        }))
    case envProd:
        logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        }))
    default:
        return nil, errors.New("invalid environment provided")
    }

    return logger, nil
}
