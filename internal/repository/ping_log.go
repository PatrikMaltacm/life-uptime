package repository

import (
	"context"
	"database/sql"

	"github.com/PatrikMaltacm/life-uptime/internal/model"
)

type PingLogRepository struct {
	db *sql.DB
}

func NewPingLogRepository(db *sql.DB) *PingLogRepository {
	return &PingLogRepository{db: db}
}

func (r *PingLogRepository) Create(ctx context.Context, l model.PingLogRequest) error {
	query := `
		INSERT INTO ping_logs (monitor_id, status_code, latency_ms, timestamp, error)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		l.MonitorID,
		l.StatusCode,
		l.Latency,
		l.Timestamp,
		l.Error,
	)
	return err
}

func (r *PingLogRepository) Get(ctx context.Context, l model.PingLogRequest) error {
	query := `
		INSERT INTO ping_logs (monitor_id, status_code, latency_ms, timestamp, error)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		l.MonitorID,
		l.StatusCode,
		l.Latency,
		l.Timestamp,
		l.Error,
	)
	return err
}




