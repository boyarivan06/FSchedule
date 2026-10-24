package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	Db *pgxpool.Pool
}
type EventParams struct {
	Name    *string    `database:"name" json:"name"`
	StartAt *time.Time `database:"start_at" json:"start_at"`
	EndAt   *time.Time `database:"end_at" json:"end_at"`
	Repeat  *int64     `database:"repeat" json:"repeat"`
	ChildID *int64     `database:"child_id" json:"child_id"`
}

func (repo *EventRepository) GetById(ctx context.Context, id int) (Event, error) {
	var e Event
	err := repo.Db.QueryRow(ctx, "SELECT * FROM events WHERE id = $1", id).Scan(&e.ID, &e.Name, &e.StartAt, &e.EndAt, &e.Repeat, &e.ChildID)
	return e, err
}
func (repo *EventRepository) GetManyByParams(ctx context.Context, p EventParams) ([]Event, error) {
	var events []Event
	var rows pgx.Rows
	var err error
	var (
		sets []string
		args []any
	)

	add := func(col string, val any) {
		args = append(args, val)
		sets = append(sets, fmt.Sprintf("%s = $%d", col, len(args)))
	}

	if p.Name != nil {
		add("name", *p.Name)
	}
	if p.ChildID != nil {
		add("child_id", *p.ChildID)
	}
	if p.StartAt != nil {
		add("start_at", *p.StartAt)
	}
	if p.EndAt != nil {
		add("end_at", *p.EndAt)
	}
	if p.Repeat != nil {
		add("repeat", *p.Repeat)
	}
	var query string
	if len(sets) == 0 {
		query = "SELECT * FROM events"
	} else {
		args = append(args)
		query = fmt.Sprintf(
			"SELECT * FROM events WHERE %s",
			strings.Join(sets, " AND "),
		)
	}

	rows, err = repo.Db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e Event
		err = rows.Scan(&e.ID, &e.Name, &e.StartAt, &e.EndAt, &e.Repeat, &e.ChildID)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, err
}

func (repo *EventRepository) Create(ctx context.Context, e Event) (Event, error) {
	_, err := repo.Db.Exec(ctx, "INSERT INTO events (name, start_at, end_at, repeat, child_id) VALUES ($1, $2, $3, $4, $5)", e.Name, e.StartAt, e.EndAt, e.Repeat, e.ChildID)
	event := Event{}
	err = repo.Db.QueryRow(ctx, "SELECT * FROM events WHERE name = $1 AND child_id = $2 AND start_at = $3", e.Name, e.ChildID, e.StartAt).Scan(&event.ID, &event.Name, &event.StartAt, &event.EndAt, &event.Repeat, &event.ChildID)
	if err != nil {
		return Event{}, err
	}
	return event, err
}

func (repo *EventRepository) Update(ctx context.Context, id int64, p EventParams) error {
	var (
		sets []string
		args []any
	)

	add := func(col string, val any) {
		args = append(args, val)
		sets = append(sets, fmt.Sprintf("%s = $%d", col, len(args)))
	}

	if p.Name != nil {
		add("name", *p.Name)
	}
	if p.ChildID != nil {
		add("child_id", *p.ChildID)
	}
	if p.StartAt != nil {
		add("start_at", *p.StartAt)
	}
	if p.EndAt != nil {
		add("end_at", *p.EndAt)
	}
	if p.Repeat != nil {
		add("repeat", *p.Repeat)
	}

	if len(sets) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE events SET %s WHERE id = $%d",
		strings.Join(sets, ", "),
		len(args),
	)

	tag, err := repo.Db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *EventRepository) Delete(ctx context.Context, id int64) error {
	_, err := repo.Db.Exec(ctx, "DELETE FROM events WHERE id = $1", id)
	return err
}
