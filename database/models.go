package database

// При добавлении полей добавить в разрешённые в utils.go
import "time"

type User struct {
	ID       int64  `database:"id" json:"id"`
	Username string `database:"username" json:"username"`
}

type Child struct {
	ID     int64  `database:"id" json:"id"`
	Name   string `database:"name" json:"name"`
	UserID int64  `database:"user_id" json:"user_id"`
}

type Event struct {
	ID      int64     `database:"id" json:"id"`
	Name    string    `database:"name" json:"name"`
	StartAt time.Time `database:"start_at" json:"start_at"`
	EndAt   time.Time `database:"end_at" json:"end_at"`
	Repeat  int64     `database:"repeat" json:"repeat"` // 0 - no repeat, 1 - every year, 2 - every month, 3 - every week, 4 - every day, TODO: можно мигрировать на int8
	ChildID int64     `database:"child_id" json:"child_id"`
}
