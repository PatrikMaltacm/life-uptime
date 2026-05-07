package model

import "time"

type MonitorRequest struct {
	ID       string        `json:"id"`
	URL      string        `json:"url" binding:"url"`
	Interval time.Duration `json:"interval" binding:"gt=0"`
	Active   bool          `json:"active" binding:"required"`
}

type MonitorResponse struct {
	ID       string        `json:"id"`
	URL      string        `json:"url"`
	Interval time.Duration `json:"interval"`
	Active   bool          `json:"active"`
}
