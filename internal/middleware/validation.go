package middleware

import (
    "encoding/json"
    "net/http"
)

func ValidateJSON(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Content-Type") != "application/json" && r.Method != "GET" {
            http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func ParseJSONBody(r *http.Request, v interface{}) error {
    defer r.Body.Close()
    return json.NewDecoder(r.Body).Decode(v)
}