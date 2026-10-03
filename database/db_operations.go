package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Db *pgxpool.Pool
}

func (repo *UserRepository) GetOneUserById(ctx context.Context, id int) (User, error) {
	var user User
	err := repo.Db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username)
	return user, err
}
