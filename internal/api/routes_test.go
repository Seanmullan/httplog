package api

import (
	"context"
	"encoding/json"
	"httplog/internal/db"
	"httplog/internal/db/dbmocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHttpLogRequestsByUrl_OK(t *testing.T) {
	dbmock := new(dbmocks.MockClient)
	now := time.Now().UTC()
	want := &db.HttpLogRequest{
		ID:         "123",
		URL:        "https://example.com",
		Method:     "GET",
		TimeIn:     now,
		TimeOut:    now,
		Duration:   100,
		ReturnCode: 200,
		Username:   "testuser",
		Userole:    "testrole",
		OrgID:      "testorg",
		UserAgent:  "testuseragent",
		ErrorMsg:   "testerror",
	}

	dbmock.GetHttpLogRequestsByUrlFn = func(ctx context.Context, url string) ([]db.HttpLogRequest, error) {
		return []db.HttpLogRequest{*want}, nil
	}

	h := Routes(dbmock)

	req := httptest.NewRequest(http.MethodGet, "/httplog/url?url=https://example.com", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got []db.HttpLogRequest
	assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
	assert.Equal(t, want.URL, got[0].URL)
	assert.Equal(t, want.Method, got[0].Method)
	assert.Equal(t, want.TimeIn, got[0].TimeIn)
	assert.Equal(t, want.TimeOut, got[0].TimeOut)
	assert.Equal(t, want.Duration, got[0].Duration)
	assert.Equal(t, want.ReturnCode, got[0].ReturnCode)
	assert.Equal(t, want.Username, got[0].Username)
	assert.Equal(t, want.Userole, got[0].Userole)
	assert.Equal(t, want.OrgID, got[0].OrgID)
	assert.Equal(t, want.UserAgent, got[0].UserAgent)
	assert.Equal(t, want.ErrorMsg, got[0].ErrorMsg)
}

func TestGetHttpLogRequestsByUrl_NotFound(t *testing.T) {
	dbmock := new(dbmocks.MockClient)
	dbmock.GetHttpLogRequestsByUrlFn = func(ctx context.Context, url string) ([]db.HttpLogRequest, error) {
		return nil, db.ErrNotFound
	}

	h := Routes(dbmock)

	req := httptest.NewRequest(http.MethodGet, "/httplog/url?url=https://example.com", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "no logs found")
}

func TestGetHttpLogRequestsByUsername_OK(t *testing.T) {
	dbmock := new(dbmocks.MockClient)
	now := time.Now().UTC()
	want := &db.HttpLogRequest{
		ID:         "456",
		URL:        "https://example.com/api",
		Method:     "POST",
		TimeIn:     now,
		TimeOut:    now.Add(100 * time.Millisecond),
		Duration:   100,
		ReturnCode: 201,
		Username:   "testuser",
		Userole:    "admin",
		OrgID:      "testorg",
		UserAgent:  "Mozilla/5.0",
		ErrorMsg:   "",
	}

	dbmock.GetHttpLogRequestsByUsernameFn = func(ctx context.Context, username string) ([]db.HttpLogRequest, error) {
		return []db.HttpLogRequest{*want}, nil
	}

	h := Routes(dbmock)

	req := httptest.NewRequest(http.MethodGet, "/httplog/username?username=testuser", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got []db.HttpLogRequest
	assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Len(t, got, 1)
	assert.Equal(t, want.ID, got[0].ID)
	assert.Equal(t, want.URL, got[0].URL)
	assert.Equal(t, want.Method, got[0].Method)
	assert.Equal(t, want.Username, got[0].Username)
}

func TestGetHttpLogRequestsByUsername_NotFound(t *testing.T) {
	dbmock := new(dbmocks.MockClient)
	dbmock.GetHttpLogRequestsByUsernameFn = func(ctx context.Context, username string) ([]db.HttpLogRequest, error) {
		return nil, db.ErrNotFound
	}

	h := Routes(dbmock)

	req := httptest.NewRequest(http.MethodGet, "/httplog/username?username=nonexistent", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "no logs found")
}

func TestInsertHttpLogRequest_OK(t *testing.T) {
	dbmock := new(dbmocks.MockClient)
	now := time.Now().UTC()
	request := db.HttpLogRequest{
		ID:         "789",
		URL:        "https://api.example.com/data",
		Method:     "PUT",
		TimeIn:     now,
		TimeOut:    now.Add(50 * time.Millisecond),
		Duration:   50,
		ReturnCode: 200,
		Username:   "apiuser",
		Userole:    "user",
		OrgID:      "myorg",
		UserAgent:  "curl/7.68.0",
		ErrorMsg:   "",
	}

	dbmock.InsertHttpLogRequestFn = func(ctx context.Context, req db.HttpLogRequest) error {
		return nil
	}

	h := Routes(dbmock)

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/httplog", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestInsertHttpLogRequest_InvalidJSON(t *testing.T) {
	dbmock := new(dbmocks.MockClient)
	h := Routes(dbmock)

	req := httptest.NewRequest(http.MethodPost, "/httplog", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "internal server error")
}

func TestBulkInsertHttpLogRequest_OK(t *testing.T) {
	dbmock := new(dbmocks.MockClient)
	now := time.Now().UTC()
	requests := []db.HttpLogRequest{
		{
			ID:         "111",
			URL:        "https://api.example.com/endpoint1",
			Method:     "GET",
			TimeIn:     now,
			TimeOut:    now.Add(10 * time.Millisecond),
			Duration:   10,
			ReturnCode: 200,
			Username:   "user1",
			Userole:    "user",
			OrgID:      "org1",
			UserAgent:  "browser",
			ErrorMsg:   "",
		},
		{
			ID:         "222",
			URL:        "https://api.example.com/endpoint2",
			Method:     "POST",
			TimeIn:     now,
			TimeOut:    now.Add(20 * time.Millisecond),
			Duration:   20,
			ReturnCode: 201,
			Username:   "user2",
			Userole:    "admin",
			OrgID:      "org2",
			UserAgent:  "api-client",
			ErrorMsg:   "",
		},
	}

	dbmock.BulkInsertHttpLogRequestFn = func(ctx context.Context, reqs []db.HttpLogRequest) error {
		return nil
	}

	h := Routes(dbmock)

	body, _ := json.Marshal(requests)
	req := httptest.NewRequest(http.MethodPost, "/httplog/bulk", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestBulkInsertHttpLogRequest_InvalidJSON(t *testing.T) {
	dbmock := new(dbmocks.MockClient)
	h := Routes(dbmock)

	req := httptest.NewRequest(http.MethodPost, "/httplog/bulk", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "internal server error")
}
