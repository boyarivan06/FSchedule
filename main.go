package main

import (
	"FSchedule/database"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	const connStr = "postgres://go_user:machine_banana@localhost:5432/f_schedule"

	ctx := context.Background()
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось подключиться к базе данных: %v\n", err)
		os.Exit(1)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			fmt.Println("aaaa not closing")
		}
	}(conn, ctx)
	pool, err := pgxpool.New(ctx, connStr)
	userRepo := &database.UserRepository{Db: pool}
	user, err := userRepo.GetOneUserById(ctx, 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка %v\n", err)
	}
	fmt.Println("got user", user.Username)
}
