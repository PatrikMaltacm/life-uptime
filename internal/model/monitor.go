package model

import "time"

type MonitorRequest struct {
	ID       string        `json:"id"`
	URL      string        `json:"url" binding:"url"`
	Interval time.Duration `json:"interval" binding:"gt=0"`
	Active   bool          `json:"active"`
}

type MonitorResponse struct {
	ID         string        `json:"id"`
	URL        string        `json:"url"`
	Interval   time.Duration `json:"interval"`
	Active     bool          `json:"active"`
	StatusCode *int          `json:"status_code"`
	LatencyMs  *int64        `json:"latency_ms"`
	LastPingAt *time.Time    `json:"last_ping_at"`
	Error      *string       `json:"error"`
}

type MonitorDeletedResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
