package database

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	Users    *UserRepository
	Children *ChildRepository
	Events   *EventRepository
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		Users:    &UserRepository{pool},
		Children: &ChildRepository{pool},
		Events:   &EventRepository{pool},
	}
}
