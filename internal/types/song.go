package types

// Song is the model that represents storing of song
type Song struct {
    Id          int    `json:"id"`
    Song        string `json:"song"`
    Group       string `json:"group"`
    Text        string `json:"text"`
    Link        string `json:"link"`
    ReleaseDate Date   `json:"releaseDate"`
}

// GetSongs represents data that uses
// for getting songs from storage.
type GetSongs struct {
    Id          *int    `json:"id"`
    Song        *string `json:"song"`
    Group       *string `json:"group"`
    Text        *string `json:"text"`
    Link        *string `json:"link"`
    ReleaseDate *Date   `json:"releaseDate"`
}

// CreateSong represents data that uses
// for creating song in the storage.
type CreateSong struct {
    Song  string `json:"song"`
    Group string `json:"group"`
    SongDetail
}

// SongDetail represents data that gets
// from external api for enrichment data from request.
type SongDetail struct {
    Text        string `json:"text"`
    Link        string `json:"link"`
    ReleaseDate Date   `json:"releaseDate"`
}

// UpdateSong represents data that uses
// for updating song in the storage.
type UpdateSong struct {
    Song        *string `json:"song"`
    Group       *string `json:"group"`
    Text        *string `json:"text"`
    Link        *string `json:"link"`
    ReleaseDate *Date   `json:"releaseDate"`
}
