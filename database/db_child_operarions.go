package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChildRepository struct {
	Db *pgxpool.Pool
}
type ChildParams struct {
	Name   *string `database:"name" json:"name"`
	UserID *int64  `database:"user_id" json:"user_id"`
}

func (repo *ChildRepository) GetById(ctx context.Context, id int) (Child, error) {
	var e Child
	err := repo.Db.QueryRow(ctx, "SELECT * FROM children WHERE id = $1", id).Scan(&e.ID, &e.Name, &e.UserID)
	return e, err
}
func (repo *ChildRepository) GetManyByParams(ctx context.Context, p ChildParams) ([]Child, error) {
	var children []Child
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
	if p.UserID != nil {
		add("user_id", *p.UserID)
	}
	var query string
	if len(sets) == 0 {
		query = "SELECT * FROM children"
	} else {
		args = append(args)
		query = fmt.Sprintf(
			"SELECT * FROM children WHERE ",
			strings.Join(sets, " AND "),
			len(args),
		)
	}

	rows, err = repo.Db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Child
		err = rows.Scan(&c.ID, &c.Name, &c.UserID)
		if err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, err
}

func (repo *ChildRepository) Create(ctx context.Context, c Child) error {
	_, err := repo.Db.Exec(ctx, "INSERT INTO children (name, user_id) VALUES ($1, $2)", c.Name, c.UserID)
	return err
}

func (repo *ChildRepository) Update(ctx context.Context, id int64, p ChildParams) error {
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
	if p.UserID != nil {
		add("user_id", *p.UserID)
	}

	if len(sets) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE children SET %s WHERE id = $%d",
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

func (repo *ChildRepository) Delete(ctx context.Context, id int64) error {
	_, err := repo.Db.Exec(ctx, "DELETE FROM children WHERE id = $1", id)
	return err
}
