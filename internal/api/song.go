package api

import (
    "encoding/json"
    "github.com/vasch3nko/songlibrary/internal/services"
    "github.com/vasch3nko/songlibrary/internal/types"
    "math"
    "net/http"
    "strconv"
)

type SongHandler struct {
    service services.SongService
    mux     *LoggingMux
}

func NewSongHandler(service services.SongService, mux *LoggingMux) *SongHandler {
    return &SongHandler{
        service: service,
        mux:     mux,
    }
}

func (s SongHandler) RegisterSongRoutes() {
    s.mux.HandleFunc("GET /songs", s.handleGetSongs)
    s.mux.HandleFunc("GET /songs/{id}", s.handleGetSongText)
    s.mux.HandleFunc("POST /songs", s.handleCreateSong)
    s.mux.HandleFunc("PATCH /songs/{id}", s.handleUpdateSong)
    s.mux.HandleFunc("DELETE /songs/{id}", s.handleDeleteSong)
}

func (s SongHandler) handleGetSongs(w http.ResponseWriter, r *http.Request) error {
    var req types.GetSongs
    ptrByParam := map[string]interface{}{
        "id":           &req.Id,
        "song":         &req.Song,
        "group":        &req.Group,
        "text":         &req.Text,
        "link":         &req.Link,
        "release_date": &req.ReleaseDate,
    }

    // Filling get songs request struct from query params
    for param, ptr := range ptrByParam {
        if !r.URL.Query().Has(param) {
            continue
        }
        value := r.URL.Query().Get(param)

        switch field := ptr.(type) {
        case **int:
            n, err := strconv.Atoi(value)
            if err != nil {
                return NewHttpError(http.StatusBadRequest)
            }
            *field = &n
        case **string:
            *field = &value
        case **types.Date:
            var date types.Date
            if err := date.Scan(value); err != nil {
                return NewHttpError(http.StatusBadRequest)
            }
            *field = &date
        default:
            return NewHttpError(http.StatusBadRequest)
        }
    }

    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    // Validating the page parameter
    if page < 1 || page > math.MaxInt32 {
        return NewHttpError(http.StatusBadRequest)
    }

    limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    // Validating the limit parameter
    if limit < 1 || limit > math.MaxInt32 {
        return NewHttpError(http.StatusBadRequest)
    }

    songs, err := s.service.GetSongs(req, page, limit)
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    if songs == nil {
        return WriteJson(w, http.StatusOK, []interface{}{})
    }

    return WriteJson(w, http.StatusOK, songs)
}

func (s SongHandler) handleGetSongText(w http.ResponseWriter, r *http.Request) error {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    verse, err := s.service.GetSongText(id, page)
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    return WriteJson(w, http.StatusOK, map[string]interface{}{
        "id":    id,
        "page":  page,
        "verse": verse,
    })
}

func (s SongHandler) handleCreateSong(w http.ResponseWriter, r *http.Request) error {
    // Decoding the request in CreateSong struct
    var req types.CreateSong
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return NewHttpError(http.StatusBadRequest)
    }
    defer r.Body.Close()

    id, err := s.service.CreateSong(req)
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    return WriteJson(w, http.StatusCreated, map[string]interface{}{
        "id": id,
    })
}

func (s SongHandler) handleUpdateSong(w http.ResponseWriter, r *http.Request) error {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    // Decoding the request in UpdateSong struct
    var req types.UpdateSong
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return NewHttpError(http.StatusBadRequest)
    }
    defer r.Body.Close()

    if err := s.service.UpdateSong(id, req); err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    return WriteJson(w, http.StatusNoContent, struct{}{})
}

func (s SongHandler) handleDeleteSong(w http.ResponseWriter, r *http.Request) error {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    if err := s.service.DeleteSong(id); err != nil {
        return NewHttpError(http.StatusBadRequest)
    }

    return WriteJson(w, http.StatusNoContent, struct{}{})
}
