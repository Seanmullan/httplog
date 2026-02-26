package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"httplog/internal/api"
	"httplog/internal/db"
)

func main() {
	addr := getenv("HTTP_ADDR", ":8080")
	dsn := getenv("POSTGRES_DSN", "postgres://user:password@localhost:5432/httplog?sslmode=disable")

	dbclient, err := db.NewPostgresClient(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbclient.Close()

	srv := &http.Server{
		Addr:    addr,
		Handler: api.Routes(dbclient),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("listening on %s", addr)
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
