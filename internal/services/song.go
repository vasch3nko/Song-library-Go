package services

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/vasch3nko/songlibrary/internal/storage"
    "github.com/vasch3nko/songlibrary/internal/types"
    "io"
    "log/slog"
    "net/http"
    "net/url"
    "strings"
)

type SongService struct {
    store            storage.Storage
    songDetailApiUrl string
    log              *slog.Logger
}

func NewSongService(store storage.Storage, songDetailApiUrl string, logger *slog.Logger) SongService {
    log := logger.With("component", "services/song")

    return SongService{store: store, songDetailApiUrl: songDetailApiUrl, log: log}
}

func (s SongService) GetSongs(req types.GetSongs, page int, limit int) ([]types.Song, error) {
    entry := s.log.With(slog.String("method", "get songs"))

    // Getting songs from storage
    songs, err := s.store.GetSongs(req, (page-1)*limit, limit)
    if err != nil {
        return nil, err
    }

    entry.Info("Songs received successfully")

    return songs, nil
}

func (s SongService) GetSongText(id int, page int) (string, error) {
    entry := s.log.With(slog.String("method", "get song text"))

    // Getting song's text by id from storage
    text, err := s.store.GetSongText(id)
    if err != nil {
        return "", err
    }

    // Splitting text by verses
    verses := strings.Split(text, "\n\n")

    // Validating the page parameter
    if len(verses) < page || page < 1 {
        err := errors.New("page out of range")
        entry.Error("Invalid parameter page",
            slog.Any("error", err),
        )
        return "", err
    }

    entry.Debug("Page validated successfully", slog.Int("page", page))
    entry.Info("Song text received successfully")

    return verses[page-1], nil
}

func (s SongService) CreateSong(req types.CreateSong) (int, error) {
    entry := s.log.With(slog.String("method", "create song"))

    // Adding request params
    params := url.Values{}
    params.Add("song", req.Song)
    params.Add("group", req.Group)

    // Requesting external API
    fullURL := fmt.Sprintf("%s/info?%s", s.songDetailApiUrl, params.Encode())
    resp, err := http.Get(fullURL)
    if err != nil {
        entry.Error("Failed to get response from external API",
            slog.String("url", fullURL),
            slog.String("error", err.Error()),
        )
        return -1, err
    }
    defer resp.Body.Close()

    entry.Debug("Got response from external API successfully")

    if resp.StatusCode != http.StatusOK {
        entry.Error("Response status from external API is not OK",
            slog.Int("status_code", resp.StatusCode),
        )
        return -1, err
    }

    entry.Debug("Response status from external API is OK")

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        entry.Error("Failed to read response from external API",
            slog.String("error", err.Error()),
        )
        return -1, err
    }

    entry.Debug("Response body read from external API successfully")

    // Adding song details (text, link, release date)
    // from external API response
    songDetail := types.SongDetail{}
    if err = json.Unmarshal(body, &songDetail); err != nil {
        entry.Error("Failed to unmarshal response from external API",
            slog.String("error", err.Error()),
        )
        return -1, err
    }
    req.SongDetail = songDetail

    entry.Debug("Response body unmarshalled successfully")

    // Creating song in the storage
    id, err := s.store.CreateSong(req)
    if err != nil {
        return -1, err
    }

    entry.Debug("Song created successfully", slog.Int("id", id))

    return id, nil
}

func (s SongService) UpdateSong(id int, req types.UpdateSong) error {
    entry := s.log.With(slog.String("method", "update song"))

    if err := s.store.UpdateSong(id, req); err != nil {
        return err
    }

    entry.Info("Song updated successfully", slog.Int("id", id))

    return nil
}

func (s SongService) DeleteSong(id int) error {
    entry := s.log.With(slog.String("method", "delete song"))

    if err := s.store.DeleteSong(id); err != nil {
        return err
    }

    entry.Info("Song deleted successfully", slog.Int("id", id))

    return nil
}
