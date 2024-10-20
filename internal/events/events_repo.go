package events

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/google/uuid"

	"github.com/forest-shadow/calendar/internal/database"
)

type Repository struct {
	db database.DBConnection
}

func NewEventsRepository(db database.DBConnection) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, event Event) error {
	query := `
		INSERT INTO events (id, title, start_time, end_time, description, user_id, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserID,
		event.NotifyBefore,
	)
	if err != nil {
		return fmt.Errorf("event creation error: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, id uuid.UUID) (Event, error) {
	query := `
		SELECT * FROM events WHERE id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return Event{}, fmt.Errorf("exec query: %w", err)
	}

	var event Event
	err = dbscan.ScanOne(&event, rows)
	if err != nil {
		return Event{}, fmt.Errorf("scan rows: %w", err)
	}

	return event, nil
}

func (r *Repository) GetList(ctx context.Context, from time.Time, to time.Time) ([]Event, error) {
	query := `
		SELECT * FROM events WHERE start_time BETWEEN $1 AND $2
	`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	var events []Event
	err = dbscan.ScanAll(&events, rows)
	if err != nil {
		return nil, fmt.Errorf("scan rows: %w", err)
	}

	return events, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, eventUpdateDTO EventUpdateDTO) error {
	qb := squirrel.Update("events").Where(squirrel.Eq{"id": id})
	qb = qb.SetMap(changesBuilder(eventUpdateDTO).ToMap())
	query, args, err := qb.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("event update error: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM events WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("event deletion error: %w", err)
	}

	return nil
}
