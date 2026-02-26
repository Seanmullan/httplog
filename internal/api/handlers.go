package api

import (
	"encoding/json"
	"errors"
	"httplog/internal/db"
	"net/http"
)

func InsertHttpLogRequest(dbclient db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request db.HttpLogRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := dbclient.InsertHttpLogRequest(r.Context(), request); err != nil {
			http.Error(w, "failed to insert http log request: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func BulkInsertHttpLogRequest(dbclient db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requests []db.HttpLogRequest
		if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
			http.Error(w, "internal server error: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := dbclient.BulkInsertHttpLogRequest(r.Context(), requests); err != nil {
			http.Error(w, "failed to insert http log requests: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetHttpLogRequestsByUrl(dbclient db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlParam := r.URL.Query().Get("url")
		if urlParam == "" {
			http.Error(w, "missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		requests, err := dbclient.GetHttpLogRequestsByUrl(r.Context(), urlParam)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				http.Error(w, "no logs found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to get http log requests: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(requests)
	}
}

func GetHttpLogRequestsByUsername(dbclient db.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "missing 'username' query parameter", http.StatusBadRequest)
			return
		}

		requests, err := dbclient.GetHttpLogRequestsByUsername(r.Context(), username)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				http.Error(w, "no logs found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to get http log requests: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(requests)
	}
}
