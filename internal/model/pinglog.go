package model

import "time"

type PingLogRequest struct {
	MonitorID  string    `json:"monitor_id" binding:"required"`
	StatusCode int       `json:"status_code" binding:"required"`
	Latency    int64     `json:"latency_ms" binding:"required"`
	Timestamp  time.Time `json:"timestamp" binding:"required"`
	Error      string    `json:"error,omitempty" `
}

type PingLogResponse struct {
	ID         string    `json:"id"`
	MonitorID  string    `json:"monitor_id"`
	StatusCode *int      `json:"status_code"`
	LatencyMs  *int64    `json:"latency_ms"`
	Timestamp  time.Time `json:"timestamp"`
	Error      *string   `json:"error"`
}
