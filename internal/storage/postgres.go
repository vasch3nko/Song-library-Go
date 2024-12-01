package storage

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/pressly/goose/v3"
    "github.com/vasch3nko/songlibrary/internal/types"
    "log/slog"
    "strings"
)

// PostgresStore is the struct that
// implements Storage interface
// for Postgres database
type PostgresStore struct {
    db  *sql.DB
    log *slog.Logger
}

// NewPostgresStore is a constructor function
// that creates a PostgresStore struct
// and connects to postgres DB
func NewPostgresStore(host, port, user, password, dbname, sslmode string, logger *slog.Logger) (*PostgresStore, error) {
    log := logger.With("component", "storage/postgres")

    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, dbname, sslmode,
    )

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Debug(
            "Failed to connect to Postgres",
            slog.String("host", host),
            slog.String("port", port),
            slog.String("user", user),
            slog.String("password", password),
            slog.String("dbname", dbname),
            slog.String("sslmode", sslmode),
            slog.Any("error", err),
        )
        return nil, err
    }

    if err := db.Ping(); err != nil {
        log.Debug(
            "Failed to ping to Postgres",
            slog.Any("error", err),
        )
        return nil, err
    }

    log.Info("Connected to Postgres successfully")

    return &PostgresStore{db: db, log: log}, nil
}

func (s *PostgresStore) Migrate(path string) error {
    const dialect = "postgres"
    entry := s.log.With(slog.String("method", "migrate"))

    goose.SetLogger(slog.NewLogLogger(entry.Handler(), slog.LevelInfo))

    if err := goose.SetDialect(dialect); err != nil {
        entry.Error("Failed to set goose dialect",
            slog.String("error", err.Error()),
        )
        return err
    }

    entry.Debug("Goose dialect sets successfully", slog.String("dialect", dialect))

    if err := goose.Up(s.db, path); err != nil {
        entry.Error("Failed to up migrations",
            slog.String("path", path),
            slog.String("error", err.Error()),
        )
        return err
    }

    entry.Info("Migrations completed successfully")

    return nil
}

func (s *PostgresStore) GetSongs(filter types.GetSongs, offset, limit int) ([]types.Song, error) {
    entry := s.log.With(slog.String("method", "get songs"))

    query := `SELECT "id", "name", "group", "text", "link", "release_date" FROM song`

    var whereClauses []string
    var args []interface{}
    i := 1

    if filter.Id != nil {
        whereClauses = append(whereClauses, fmt.Sprintf(`"id" = $%d`, i))
        args = append(args, *filter.Id)
        i++
    }
    if filter.Song != nil {
        whereClauses = append(whereClauses, fmt.Sprintf(`"song" = $%d`, i))
        args = append(args, *filter.Song)
        i++
    }
    if filter.Group != nil {
        whereClauses = append(whereClauses, fmt.Sprintf(`"group" = $%d`, i))
        args = append(args, *filter.Group)
        i++
    }
    if filter.Text != nil {
        whereClauses = append(whereClauses, fmt.Sprintf(`"text" = $%d`, i))
        args = append(args, *filter.Text)
        i++
    }
    if filter.Link != nil {
        whereClauses = append(whereClauses, fmt.Sprintf(`"link" = $%d`, i))
        args = append(args, *filter.Link)
        i++
    }
    if filter.ReleaseDate != nil {
        whereClauses = append(whereClauses, fmt.Sprintf(`"release_date" = $%d`, i))
        args = append(args, *filter.ReleaseDate)
        i++
    }

    if len(whereClauses) > 0 {
        query += " WHERE " + strings.Join(whereClauses, " AND ")
    }

    query += fmt.Sprintf(" OFFSET $%d LIMIT $%d", i, i+1)
    args = append(args, offset, limit)

    rows, err := s.db.Query(query, args...)
    if err != nil {
        entry.Error("Get songs query failed",
            slog.String("query", query),
            slog.Any("error", err),
        )
        return nil, err
    }

    entry.Debug("Get songs query completed successfully")

    var songs []types.Song
    for rows.Next() {
        var song types.Song
        if err := rows.Scan(
            &song.Id,
            &song.Song,
            &song.Group,
            &song.Text,
            &song.Link,
            &song.ReleaseDate,
        ); err != nil {
            entry.Error("Failed to scan song",
                slog.Group("song",
                    slog.Int("id", song.Id),
                    slog.String("song", song.Song),
                    slog.String("group", song.Group),
                    slog.String("text", song.Text),
                    slog.String("link", song.Link),
                    slog.String("release_date", song.ReleaseDate.String()),
                ),
                slog.Any("error", err),
            )
            return nil, err
        }
        songs = append(songs, song)
    }
    entry.Debug("Songs scanned successfully")
    entry.Info("Got songs successfully")

    return songs, nil
}

func (s *PostgresStore) GetSongText(id int) (string, error) {
    entry := s.log.With(slog.String("method", "get song text"))

    query := `SELECT text FROM song WHERE id = $1;`
    row := s.db.QueryRow(query, id)

    var text string
    if err := row.Scan(&text); err != nil {
        entry.Error("Failed to get song text",
            slog.String("query", query),
            slog.Any("error", err),
        )
        return "", err
    }
    entry.Info("Got song text successfully")

    return text, nil
}

func (s *PostgresStore) CreateSong(song types.CreateSong) (int, error) {
    entry := s.log.With(slog.String("method", "create song"))
    var id int

    query := `
            INSERT INTO song ("name", "group", "text", "link", "release_date")
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id;
        `

    err := s.db.QueryRow(
        query,
        song.Song,
        song.Group,
        song.Text,
        song.Link,
        song.ReleaseDate,
    ).Scan(&id)

    if err != nil {
        entry.Error("Failed to create song",
            slog.String("query", query),
            slog.Any("error", err),
        )
        return -1, err
    }

    entry.Info("Song successfully created")

    return id, nil
}

func (s *PostgresStore) UpdateSong(id int, song types.UpdateSong) error {
    entry := s.log.With(slog.String("method", "update song"))
    updates := make(map[string]interface{})

    if song.Song != nil {
        updates["name"] = *song.Song
    }
    if song.Group != nil {
        updates["group"] = *song.Group
    }
    if song.Text != nil {
        updates["text"] = *song.Text
    }
    if song.Link != nil {
        updates["link"] = *song.Link
    }
    if song.ReleaseDate != nil {
        updates["release_date"] = *song.ReleaseDate
    }

    if len(updates) == 0 {
        err := fmt.Errorf("no fields to update")
        entry.Error("Failed to update song",
            slog.Any("error", err),
        )
        return err
    }
    entry.Debug("Updates len greater than 0")

    var fields []string
    var args []interface{}
    counter := 1

    for field, value := range updates {
        fields = append(fields, fmt.Sprintf(`"%s" = $%d`, field, counter))
        args = append(args, value)
        counter++
    }

    args = append(args, id)
    query := fmt.Sprintf("UPDATE song SET %s WHERE id = $%d", strings.Join(fields, ", "), counter)

    if _, err := s.db.Exec(query, args...); err != nil {
        entry.Error("Failed to update song",
            slog.String("query", query),
            slog.Any("error", err),
        )
        return err
    }

    entry.Info("Song updated successfully", slog.Int("id", id))

    return nil
}

func (s *PostgresStore) DeleteSong(id int) error {
    entry := s.log.With(slog.String("method", "delete song"))

    query := `DELETE FROM song WHERE id = $1;`

    if _, err := s.db.Exec(
        query,
        id,
    ); err != nil {
        entry.Error("Failed to delete song",
            slog.String("query", query),
            slog.Any("error", err),
        )
        return err
    }

    slog.Debug("Song deleted successfully", slog.Int("id", id))

    return nil
}
