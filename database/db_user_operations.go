package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Db *pgxpool.Pool
}

type UserParams struct {
	UserName *string `json:"username"`
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

//func (repo *UserRepository) GetManyByParams(ctx context.Context, p UserParams) ([]User, error) {
//	var users []User
//	var rows pgx.Rows
//	var err error
//	var (
//		sets []string
//		args []any
//	)
//
//	add := func(col string, val any) {
//		args = append(args, val)
//		sets = append(sets, fmt.Sprintf("%s = $%d", col, len(args)))
//	}
//
//	if p.UserName != nil {
//		add("username", *p.UserName)
//	}
//
//	if len(sets) == 0 {
//		return nil, nil
//	}
//
//	args = append(args)
//	query := fmt.Sprintf(
//		"SELECT * FROM users WHERE %s",
//		strings.Join(sets, " AND "),
//	)
//
//	rows, err = repo.Db.Query(ctx, query, args...)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var u User
//		err = rows.Scan(&u.ID, &u.Username)
//		if err != nil {
//			return nil, err
//		}
//		users = append(users, u)
//	}
//	return users, err
//}

func (repo *UserRepository) Create(ctx context.Context, u User) (User, error) {
	_, err := repo.Db.Exec(ctx, "INSERT INTO users (username) VALUES ($1)", u.Username)
	var user User
	err = repo.Db.QueryRow(ctx, "SELECT * FROM users WHERE username = $1", u.Username).Scan(&user.ID, &user.Username)
	return user, err
}

func (repo *UserRepository) Update(ctx context.Context, id int64, p UserParams) error {
	sets, args := prepareUserParams(p)

	if len(sets) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(sets, ", "),
		len(args),
	)

	_, err := repo.Db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	//if tag.RowsAffected() == 0 {
	//	return sql.ErrNoRows
	//}
	return nil
}

func prepareUserParams(p UserParams) ([]string, []any) {
	var (
		sets []string
		args []any
	)

	add := func(col string, val any) {
		args = append(args, val)
		sets = append(sets, fmt.Sprintf("%s = $%d", col, len(args)))
	}

	if p.UserName != nil {
		add("username", *p.UserName)
	}
	return sets, args
}

func (repo *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := repo.Db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
