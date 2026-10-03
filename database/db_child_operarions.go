package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChildRepository struct {
	Db *pgxpool.Pool
}

func (repo *ChildRepository) GetById(ctx context.Context, id int) (Child, error) {
	var e Child
	err := repo.Db.QueryRow(ctx, "SELECT * FROM children WHERE id = $1", id).Scan(&e.ID, &e.Name, &e.UserID)
	return e, err
}
func (repo *ChildRepository) GetManyByParams(ctx context.Context, params map[string]string) ([]Child, error) {
	whereQuery, args := buildQuery(params, " AND ")
	var children []Child
	rows, err := repo.Db.Query(ctx, "SELECT * FROM children "+whereQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Child
		err := rows.Scan(&c.ID, &c.Name, &c.UserID)
		if err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, err
}

func (repo *ChildRepository) Create(ctx context.Context, name string, userId int64) error {
	_, err := repo.Db.Exec(ctx, "INSERT INTO children (name, user_id) VALUES ($1, $2)", name, userId)
	return err
}

func (repo *ChildRepository) Update(ctx context.Context, id int64, params map[string]string) error {
	preQuery, args := buildQuery(params, ", ")
	query := "UPDATE children SET " + preQuery + " WHERE id = " + strconv.FormatInt(id, 10)
	fmt.Println(query)
	_, err := repo.Db.Exec(ctx, query, args...)
	return err
}

func (repo *ChildRepository) Delete(ctx context.Context, id int64) error {
	_, err := repo.Db.Exec(ctx, "DELETE FROM children WHERE id = $1", id)
	return err
}
