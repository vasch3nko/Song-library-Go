package api

import (
    "encoding/json"
    "net/http"
)

// WriteJson is the helper function that encodes
// a value of any type to json, sets a response status code
// and writes it to http.ResponseWriter. Returns the error.
func WriteJson(w http.ResponseWriter, status int, v any) error {
    w.Header().Add("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(v)
}
