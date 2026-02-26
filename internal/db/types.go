package db

import "time"

type HttpLogRequest struct {
	ID         string    `json:"id"`
	URL        string    `json:"url"`
	Method     string    `json:"method"`
	TimeIn     time.Time `json:"time_in"`
	TimeOut    time.Time `json:"time_out"`
	Duration   int       `json:"duration"`
	ReturnCode int       `json:"return_code"`
	Username   string    `json:"username"`
	Userole    string    `json:"userole"`
	OrgID      string    `json:"org_id"`
	UserAgent  string    `json:"user_agent"`
	ErrorMsg   string    `json:"error_msg"`
}
