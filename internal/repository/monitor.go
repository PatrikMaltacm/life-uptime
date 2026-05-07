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
	query := "SELECT id, url FROM monitors WHERE id = $1"

	row := r.db.QueryRowContext(ctx, query, id)

	var m model.MonitorResponse
	err := row.Scan(&m.ID, &m.URL)
	return &m, err
}

func (r *MonitorRepository) GetAll(ctx context.Context) ([]model.MonitorResponse, error) {
	query := `
		SELECT id, url
		FROM monitors LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allData []model.MonitorResponse

	for rows.Next() {
		var m model.MonitorResponse
		if err := rows.Scan(&m.ID, &m.URL); err != nil {
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
