package api

import (
    "errors"
    "log/slog"
    "net/http"
    "time"
)

type errorHandlerFunc func(http.ResponseWriter, *http.Request) error

// LoggingMux is the wrapper over http.ServeMux that logs and handle errors
type LoggingMux struct {
    mux *http.ServeMux
    log *slog.Logger
}

// NewLoggingMux is the constructor for LoggingMux that returns pointer
func NewLoggingMux(logger *slog.Logger) *LoggingMux {
    log := logger.With("component", "logging mux")

    return &LoggingMux{
        mux: http.NewServeMux(),
        log: log,
    }
}

// HandleFunc handles error, delegates handling pattern
// to internal ServeMux and logs it
func (m *LoggingMux) HandleFunc(pattern string, handlerFunc errorHandlerFunc) {
    m.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        // Starting request logging
        start := time.Now()
        entry := m.log.With(
            slog.String("method", r.Method),
            slog.String("path", r.URL.Path),
            slog.String("remote_addr", r.RemoteAddr),
            slog.String("user_agent", r.UserAgent()),
        )

        err := handlerFunc(w, r)
        if err == nil {
            // Completing request logging
            duration := time.Since(start)
            entry.Info(
                "Request completed",
                slog.Duration("duration", duration),
            )
            return
        }

        // Logging error
        m.log.Error(
            "Request error",
            slog.Any("error", err),
        )

        // If error is the httpError
        var httpErr *HttpError
        if errors.As(err, &httpErr) {
            // Writing json response with status and message from httpError
            if err := WriteJson(w, httpErr.StatusCode, struct{}{}); err != nil {
                m.log.Error("Failed to write response")
                return
            }
            return
        }

        // Else responding internal server error
        if err := WriteJson(w, http.StatusInternalServerError, struct{}{}); err != nil {
            m.log.Error("Failed to write response")
            return
        }
    })
}

// ServeHTTP delegates request processing to internal ServeMux
func (m *LoggingMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    m.mux.ServeHTTP(w, r)
}
