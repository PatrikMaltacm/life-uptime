package repository

import (
	"context"
	"database/sql"

	"github.com/PatrikMaltacm/life-uptime/internal/model"
)

type MonitorRepository struct {
	db *sql.DB
}

func NewMonitorRepository(db *sql.DB) *MonitorRepository {
	return &MonitorRepository{
		db: db,
	}
}

func (r *MonitorRepository) GetByID(ctx context.Context, id string) (*model.MonitorResponse, error) {
	query := "SELECT id, url, interval, active FROM monitors WHERE id = $1"

	row := r.db.QueryRowContext(ctx, query, id)

	var m model.MonitorResponse
	err := row.Scan(&m.ID, &m.URL, &m.Interval, &m.Active)
	return &m, err
}

func (r *MonitorRepository) GetAll(ctx context.Context) ([]model.MonitorResponse, error) {
	query := `
		SELECT 
			m.id, 
			m.url, 
			m.interval, 
			m.active,
			pl.status_code,
			pl.latency_ms,
			pl.timestamp,
			pl.error
		FROM monitors m
		LEFT JOIN LATERAL (
			SELECT status_code, latency_ms, timestamp, error
			FROM ping_logs
			WHERE monitor_id = m.id
			ORDER BY timestamp DESC
			LIMIT 1
		) pl ON true
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allData []model.MonitorResponse

	for rows.Next() {
		var m model.MonitorResponse
		if err := rows.Scan(
			&m.ID,
			&m.URL,
			&m.Interval,
			&m.Active,
			&m.StatusCode,
			&m.LatencyMs,
			&m.LastPingAt,
			&m.Error,
		); err != nil {
			return nil, err
		}
		allData = append(allData, m)
	}

	return allData, nil
}

func (r *MonitorRepository) Create(ctx context.Context, m model.MonitorRequest) (*model.MonitorResponse, error) {
	query := `
		INSERT INTO monitors (url, interval, active)
		VALUES ($1, $2, $3)
		RETURNING id, url, interval, active
	`

	var res model.MonitorResponse

	err := r.db.QueryRowContext(ctx, query, m.URL, m.Interval, m.Active).Scan(
		&res.ID,
		&res.URL,
		&res.Interval,
		&res.Active,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *MonitorRepository) Delete(ctx context.Context, id string) (*model.MonitorResponse, error) {
	query := "DELETE FROM monitors WHERE id = $1 RETURNING id, url"

	var res model.MonitorResponse

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID,
		&res.URL,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *MonitorRepository) Update(ctx context.Context, m model.MonitorRequest, id string) (*model.MonitorResponse, error) {
	query := `
		UPDATE monitors 
		SET url = $1, interval = $2, active = $3 
		WHERE id = $4 
		RETURNING id, url, interval, active;
	`

	var res model.MonitorResponse

	err := r.db.QueryRowContext(ctx, query, m.URL, m.Interval, m.Active, id).Scan(
		&res.ID,
		&res.URL,
		&res.Interval,
		&res.Active,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *MonitorRepository) GetAllActive(ctx context.Context) ([]model.MonitorResponse, error) {
	query := "SELECT id, url, interval FROM monitors WHERE active = true"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []model.MonitorResponse
	for rows.Next() {
		var m model.MonitorResponse
		if err := rows.Scan(&m.ID, &m.URL, &m.Interval); err != nil {
			return nil, err
		}
		monitors = append(monitors, m)
	}
	return monitors, nil
}

func (r *MonitorRepository) GetPingHistory(ctx context.Context, monitorID string) ([]model.PingLogResponse, error) {
	query := `
		SELECT monitor_id, status_code, latency_ms, timestamp, error
		FROM ping_logs
		WHERE monitor_id = $1
		ORDER BY timestamp DESC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query, monitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.PingLogResponse

	for rows.Next() {
		var l model.PingLogResponse
		if err := rows.Scan(
			&l.MonitorID,
			&l.StatusCode,
			&l.Latency,
			&l.Timestamp,
			&l.Error,
		); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, nil
}
