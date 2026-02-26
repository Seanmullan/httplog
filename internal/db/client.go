package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var ErrNotFound = errors.New("not found")

type Client interface {
	InsertHttpLogRequest(ctx context.Context, request HttpLogRequest) error
	BulkInsertHttpLogRequest(ctx context.Context, requests []HttpLogRequest) error
	GetHttpLogRequestsByUrl(ctx context.Context, url string) ([]HttpLogRequest, error)
	GetHttpLogRequestsByUsername(ctx context.Context, username string) ([]HttpLogRequest, error)
	Close() error
}

type PostgresClient struct {
	db *sql.DB
}

func NewPostgresClient(dsn string) (Client, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}

	client := &PostgresClient{db: db}
	if err := client.runSchemaMigration(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("schema migration failed: %w", err)
	}
	return client, nil
}

func (c *PostgresClient) Close() error {
	return c.db.Close()
}

func (c *PostgresClient) InsertHttpLogRequest(ctx context.Context, request HttpLogRequest) error {
	const requestInsert = `
		INSERT INTO httplog (id, url, method, time_in, time_out, duration, return_code, username, userole, org_id, user_agent, error_msg)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO NOTHING`

	_, err := c.db.ExecContext(
		ctx,
		requestInsert,
		request.ID,
		request.URL,
		request.Method,
		request.TimeIn,
		request.TimeOut,
		request.Duration,
		request.ReturnCode,
		request.Username,
		request.Userole,
		request.OrgID,
		request.UserAgent,
		request.ErrorMsg,
	)
	if err != nil {
		return fmt.Errorf("failed to insert http log request: %w", err)
	}
	return nil
}

func (c *PostgresClient) BulkInsertHttpLogRequest(ctx context.Context, requests []HttpLogRequest) error {
	// TODO: Use a transaction or batch insert for better performance
	for _, request := range requests {
		if err := c.InsertHttpLogRequest(ctx, request); err != nil {
			return err
		}
	}
	return nil
}

func (c *PostgresClient) GetHttpLogRequestsByUrl(ctx context.Context, url string) ([]HttpLogRequest, error) {
	const query = `
		SELECT id, url, method, time_in, time_out, duration, return_code, username, userole, org_id, user_agent, error_msg
		FROM httplog
		WHERE url = $1
		ORDER BY time_in DESC`

	rows, err := c.db.QueryContext(ctx, query, url)
	if err != nil {
		return nil, fmt.Errorf("failed to query http log requests by url: %w", err)
	}
	defer rows.Close()

	requests := make([]HttpLogRequest, 0)

	for rows.Next() {
		var r HttpLogRequest
		if err := rows.Scan(
			&r.ID,
			&r.URL,
			&r.Method,
			&r.TimeIn,
			&r.TimeOut,
			&r.Duration,
			&r.ReturnCode,
			&r.Username,
			&r.Userole,
			&r.OrgID,
			&r.UserAgent,
			&r.ErrorMsg,
		); err != nil {
			return nil, fmt.Errorf("failed to scan http log request: %w", err)
		}
		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if len(requests) == 0 {
		return nil, ErrNotFound
	}

	return requests, nil
}

func (c *PostgresClient) GetHttpLogRequestsByUsername(ctx context.Context, username string) ([]HttpLogRequest, error) {
	const query = `
		SELECT id, url, method, time_in, time_out, duration, return_code, username, userole, org_id, user_agent, error_msg
		FROM httplog
		WHERE username = $1
		ORDER BY time_in DESC`

	rows, err := c.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query http log requests by username: %w", err)
	}
	defer rows.Close()

	requests := make([]HttpLogRequest, 0)

	for rows.Next() {
		var r HttpLogRequest
		if err := rows.Scan(
			&r.ID,
			&r.URL,
			&r.Method,
			&r.TimeIn,
			&r.TimeOut,
			&r.Duration,
			&r.ReturnCode,
			&r.Username,
			&r.Userole,
			&r.OrgID,
			&r.UserAgent,
			&r.ErrorMsg,
		); err != nil {
			return nil, fmt.Errorf("failed to scan http log request: %w", err)
		}
		requests = append(requests, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if len(requests) == 0 {
		return nil, ErrNotFound
	}

	return requests, nil
}

func (c *PostgresClient) runSchemaMigration(ctx context.Context) error {
	const ddl = `
	CREATE TABLE IF NOT EXISTS httplog (
		id          TEXT PRIMARY KEY,
		url         TEXT         NOT NULL,
		method      TEXT         NOT NULL,
		time_in     TIMESTAMPTZ  NOT NULL,
		time_out    TIMESTAMPTZ  NOT NULL,
		duration    INT          NOT NULL,
		return_code INT          NOT NULL,
		username    TEXT         NOT NULL,
		userole     TEXT         NOT NULL,
		org_id      TEXT         NOT NULL,
		user_agent  TEXT         NOT NULL,
		error_msg   TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_httplog_url_time_in ON httplog (url, time_in DESC);
	CREATE INDEX IF NOT EXISTS idx_httplog_username_time_in ON httplog (username, time_in DESC);
	CREATE INDEX IF NOT EXISTS idx_httplog_time_in ON httplog (time_in DESC);
	`
	_, err := c.db.ExecContext(ctx, ddl)
	if err != nil {
		return fmt.Errorf("schema migration failed: %w", err)
	}
	return nil
}
