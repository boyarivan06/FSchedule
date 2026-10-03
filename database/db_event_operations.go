package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	Db *pgxpool.Pool
}

func (repo *EventRepository) GetById(ctx context.Context, id int) (Event, error) {
	var e Event
	err := repo.Db.QueryRow(ctx, "SELECT * FROM events WHERE id = $1", id).Scan(&e.ID, &e.Name, &e.ChildID, &e.StartAt, &e.EndAt, &e.Repeat)
	return e, err
}
func (repo *EventRepository) GetManyByParams(ctx context.Context, params map[string]string) ([]Event, error) {
	whereQuery, args := buildQuery(params, " AND ")
	var events []Event
	rows, err := repo.Db.Query(ctx, "SELECT * FROM events "+whereQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Name, &e.ChildID, &e.StartAt, &e.EndAt, &e.Repeat)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, err
}

func (repo *EventRepository) Create(ctx context.Context, name string, startAt time.Time, endAt time.Time,
	repeat int64, childId int64) error {
	_, err := repo.Db.Exec(ctx, "INSERT INTO events (name, start_at, end_at, repeat, child_id) VALUES ($1, $2, $3, $4, $5)", name, startAt, endAt, repeat, childId)
	return err
}

func (repo *EventRepository) Update(ctx context.Context, id int64, params map[string]string) error {
	preQuery, args := buildQuery(params, ", ")
	query := "UPDATE events SET " + preQuery + " WHERE id = " + strconv.FormatInt(id, 10)
	fmt.Println(query)
	_, err := repo.Db.Exec(ctx, query, args...)
	return err
}

func (repo *EventRepository) Delete(ctx context.Context, id int64) error {
	_, err := repo.Db.Exec(ctx, "DELETE FROM events WHERE id = $1", id)
	return err
}
