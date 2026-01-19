package db

import (
	"context"
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) EnsureTables(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS plans (
			id SERIAL PRIMARY KEY,
			name VARCHAR(80) NOT NULL,
			description VARCHAR(255),
			price NUMERIC(10,2) NOT NULL,
			active BOOLEAN DEFAULT TRUE
		);`,
		`CREATE TABLE IF NOT EXISTS class_schedules (
			id SERIAL PRIMARY KEY,
			weekday VARCHAR(3) NOT NULL,
			start_time TIME NOT NULL,
			end_time TIME,
			modality VARCHAR(80),
			coach VARCHAR(80),
			active BOOLEAN DEFAULT TRUE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_class_schedules_weekday ON class_schedules(weekday);`,
	}

	for _, q := range queries {
		if _, err := r.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) ListPlans(ctx context.Context) ([]Plan, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, price::text, active
		FROM plans
		WHERE active = true
		ORDER BY price ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Active); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *Repo) CreatePlan(ctx context.Context, name string, description *string, price float64, active bool) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO plans (name, description, price, active)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, name, description, price, active).Scan(&id)
	return id, err
}

func (r *Repo) ListSchedules(ctx context.Context) ([]Schedule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, weekday, start_time, end_time, modality, coach, active
		FROM class_schedules
		WHERE active = true
		ORDER BY weekday ASC, start_time ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var s Schedule
		if err := rows.Scan(
			&s.ID, &s.Weekday, &s.StartTime, &s.EndTime,
			&s.Modality, &s.Coach, &s.Active,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}

	return schedules, nil
}

func (r *Repo) CreateSchedule(
	ctx context.Context,
	weekday string,
	start string,
	end *string,
	modality *string,
	coach *string,
	active bool,
) (int64, error) {
	var id int64

	var endTime any = nil
	if end != nil && *end != "" {
		endTime = *end
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO class_schedules (weekday, start_time, end_time, modality, coach, active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, weekday, start, endTime, modality, coach, active).Scan(&id)

	return id, err
}
