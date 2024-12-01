package config

import (
    "errors"
    "os"
    "strconv"
    "time"
)

type Config struct {
    // App environment (dev, prod)
    Env               string
    SongDetailsApiUrl string

    Server struct {
        Addr string
        // In .env string for parse duration
        // (10h / 10m / 10s / 10ns / 10us or µs / 10ms)
        ReadTimeout  time.Duration
        WriteTimeout time.Duration
        IdleTimeout  time.Duration
    }

    Db struct {
        Host     string
        Port     string
        Username string
        Password string
        Database string
        SSLMode  string // (disable / require / verify-ca / verify-full)

        MigrationsPath string // Path string ("./migrations")
    }
}

func NewConfig() *Config {
    return &Config{}
}

func (cfg *Config) LoadFromEnv() error {
    cfgPtrByEnv := map[string]interface{}{
        "SL_SONG_DETAILS_API_URL": &cfg.SongDetailsApiUrl,
        "SL_ENV":                  &cfg.Env,

        "SL_SRV_ADDR":          &cfg.Server.Addr,
        "SL_SRV_READ_TIMEOUT":  &cfg.Server.ReadTimeout,
        "SL_SRV_WRITE_TIMEOUT": &cfg.Server.WriteTimeout,
        "SL_SRV_IDLE_TIMEOUT":  &cfg.Server.IdleTimeout,

        "SL_DB_HOST":            &cfg.Db.Host,
        "SL_DB_PORT":            &cfg.Db.Port,
        "SL_DB_USERNAME":        &cfg.Db.Username,
        "SL_DB_PASSWORD":        &cfg.Db.Password,
        "SL_DB_DATABASE":        &cfg.Db.Database,
        "SL_DB_SSL_MODE":        &cfg.Db.SSLMode,
        "SL_DB_MIGRATIONS_PATH": &cfg.Db.MigrationsPath,
    }

    for env, ptr := range cfgPtrByEnv {
        temp, ok := os.LookupEnv(env)
        if !ok {
            return errors.New("env variable " + env + " not set")
        }

        switch field := ptr.(type) {
        case *string:
            *field = temp
        case *int:
            n, err := strconv.Atoi(temp)
            if err != nil {
                return err
            }

            *field = n
        case *time.Duration:
            duration, err := time.ParseDuration(temp)
            if err != nil {
                return err
            }

            *field = duration
        }
    }

    return nil
}
