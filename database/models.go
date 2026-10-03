package database

// При добавлении полей добавить в разрешённые в utils.go
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
	ID      int64     `db:"id"`
	Name    string    `db:"name"`
	StartAt time.Time `db:"start_at"`
	EndAt   time.Time `db:"end_at"`
	Repeat  int64     `db:"repeat"` // 0 - no repeat, 1 - every year, 2 - every month, 3 - every week, 4 - every day, TODO: можно мигрировать на int8
	ChildID int64     `db:"child_id"`
}
