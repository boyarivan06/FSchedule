package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Db *pgxpool.Pool
}

func (repo *UserRepository) GetById(ctx context.Context, id int) (User, error) {
	var user User
	err := repo.Db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username)
	return user, err
}
func (repo *UserRepository) GetByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := repo.Db.QueryRow(ctx, "SELECT * FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username)
	return user, err
}
func (repo *UserRepository) GetManyByParams(ctx context.Context, params map[string]string) ([]User, error) {
	whereQuery, args := buildQuery(params, " AND ")
	var users []User
	rows, err := repo.Db.Query(ctx, "SELECT * FROM users "+whereQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, err
}

func (repo *UserRepository) Create(ctx context.Context, username string) error {
	_, err := repo.Db.Exec(ctx, "INSERT INTO users (username) VALUES ($1)", username)
	return err
}

func (repo *UserRepository) Update(ctx context.Context, id int64, params map[string]string) error {
	preQuery, args := buildQuery(params, ", ")
	query := "UPDATE users SET " + preQuery + " WHERE id = " + strconv.FormatInt(id, 10)
	fmt.Println(query)
	_, err := repo.Db.Exec(ctx, query, args...)
	return err
}

func (repo *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := repo.Db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
