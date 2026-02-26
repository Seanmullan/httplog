package dbmocks

import (
	"context"
	"httplog/internal/db"
)

type MockClient struct {
	InsertHttpLogRequestFn         func(ctx context.Context, request db.HttpLogRequest) error
	BulkInsertHttpLogRequestFn     func(ctx context.Context, requests []db.HttpLogRequest) error
	GetHttpLogRequestsByUrlFn      func(ctx context.Context, url string) ([]db.HttpLogRequest, error)
	GetHttpLogRequestsByUsernameFn func(ctx context.Context, username string) ([]db.HttpLogRequest, error)
}

func (f *MockClient) InsertHttpLogRequest(ctx context.Context, request db.HttpLogRequest) error {
	return f.InsertHttpLogRequestFn(ctx, request)
}

func (f *MockClient) BulkInsertHttpLogRequest(ctx context.Context, requests []db.HttpLogRequest) error {
	return f.BulkInsertHttpLogRequestFn(ctx, requests)
}

func (f *MockClient) GetHttpLogRequestsByUrl(ctx context.Context, url string) ([]db.HttpLogRequest, error) {
	return f.GetHttpLogRequestsByUrlFn(ctx, url)
}

func (f *MockClient) GetHttpLogRequestsByUsername(ctx context.Context, username string) ([]db.HttpLogRequest, error) {
	return f.GetHttpLogRequestsByUsernameFn(ctx, username)
}

func (f *MockClient) Close() error { return nil }
