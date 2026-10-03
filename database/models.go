package database

import "time"

type User struct {
	ID       int64  `database:"id"`
	Username string `database:"username"`
}

type Child struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	UserID int64  `db:"user_id"`
}

type Event struct {
	ID      int64     `database:"id"`
	Name    string    `database:"name"`
	StartAt time.Time `database:"start_at"`
	EndAt   time.Time `database:"end_at"`
	Repeat  int64     `database:"repeat"`
	ChildID int64     `database:"child_id"`
}
