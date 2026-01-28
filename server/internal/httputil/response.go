package httputil

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError writes an error response with logging
func WriteError(w http.ResponseWriter, r *http.Request, status int, message string, err error) {
	if err != nil {
		log.Printf("ERROR: %s %s - %v", r.Method, r.URL.Path, err)
	} else {
		log.Printf("WARN: %s %s - %s", r.Method, r.URL.Path, message)
	}
	WriteJSON(w, status, ErrorResponse{Error: message})
}

// WriteClientError writes a 400 Bad Request error
func WriteClientError(w http.ResponseWriter, r *http.Request, message string, err error) {
	WriteError(w, r, http.StatusBadRequest, message, err)
}

// WriteInternalError writes a 500 Internal Server Error
func WriteInternalError(w http.ResponseWriter, r *http.Request, err error) {
	WriteError(w, r, http.StatusInternalServerError, "An internal error occurred", err)
}

// WriteUnauthorized writes a 401 Unauthorized error
func WriteUnauthorized(w http.ResponseWriter, r *http.Request) {
	log.Printf("WARN: %s %s - unauthorized", r.Method, r.URL.Path)
	WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
