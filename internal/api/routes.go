package api

import (
	"httplog/internal/db"
	"net/http"
)

func Routes(dbclient db.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /httplog", InsertHttpLogRequest(dbclient))
	mux.HandleFunc("POST /httplog/bulk", BulkInsertHttpLogRequest(dbclient))
	mux.HandleFunc("GET /httplog/url", GetHttpLogRequestsByUrl(dbclient))
	mux.HandleFunc("GET /httplog/username", GetHttpLogRequestsByUsername(dbclient))
	return mux
}
