package storage

import (
    "github.com/vasch3nko/songlibrary/internal/types"
)

// Storage is the interface that
// describes a store of a data in API
type Storage interface {
    GetSongs(types.GetSongs, int, int) ([]types.Song, error)
    GetSongText(int) (string, error)
    CreateSong(types.CreateSong) (int, error)
    UpdateSong(int, types.UpdateSong) error
    DeleteSong(int) error
}
