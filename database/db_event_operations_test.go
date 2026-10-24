package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEventRepository_Create(t *testing.T) {
	pool := setupTestDb(t)
	repo := &EventRepository{pool}
	ctx := context.Background()
	userRepo := &UserRepository{pool}
	childRepo := &ChildRepository{pool}
	require.NoError(t, userRepo.Create(ctx, User{Username: "Ivan"}))
	petya, err := childRepo.Create(ctx, Child{Name: "Petya", UserID: 1})
	require.NoError(t, err)
	_, err = repo.Create(ctx, Event{Name: "theatre", StartAt: time.Date(2021, 5, 4, 1, 1, 1, 0, time.Local),
		EndAt:  time.Date(2021, 5, 5, 1, 1, 1, 0, time.Local),
		Repeat: 0, ChildID: petya.ID})
	require.NoError(t, err)
	_, err = repo.Create(ctx, Event{Name: "theatre", StartAt: time.Date(2021, 5, 4, 1, 1, 1, 0, time.Local),
		EndAt:  time.Date(2021, 5, 5, 1, 1, 1, 0, time.Local),
		Repeat: 0, ChildID: petya.ID + 1})
	require.Error(t, err)
}

func TestEventRepository_Get(t *testing.T) {
	pool := setupTestDb(t)
	repo := &EventRepository{pool}
	ctx := context.Background()
	userRepo := &UserRepository{pool}
	childRepo := &ChildRepository{pool}
	require.NoError(t, userRepo.Create(ctx, User{Username: "Ivan"}))
	petya, err := childRepo.Create(ctx, Child{Name: "Petya", UserID: 1})
	require.NoError(t, err)
	e, err := repo.Create(ctx, Event{Name: "theatre", StartAt: time.Date(2021, 5, 4, 1, 1, 1, 0, time.Local),
		EndAt:  time.Date(2021, 5, 5, 1, 1, 1, 0, time.Local),
		Repeat: 0, ChildID: petya.ID})
	require.NoError(t, err)
	theatre, err := repo.GetById(ctx, int(e.ID))
	require.NoError(t, err)
	another, err := repo.GetManyByParams(ctx, EventParams{Name: new("theatre"), StartAt: new(time.Date(2021, 5, 4, 1, 1, 1, 0, time.Local))})
	require.Equal(t, 1, len(another))
	require.Equal(t, theatre, another[0])
}

func TestEventRepository_Update(t *testing.T) {
	pool := setupTestDb(t)
	repo := &EventRepository{pool}
	ctx := context.Background()
	userRepo := &UserRepository{pool}
	childRepo := &ChildRepository{pool}
	require.NoError(t, userRepo.Create(ctx, User{Username: "Ivan"}))
	petya, err := childRepo.Create(ctx, Child{Name: "Petya", UserID: 1})
	require.NoError(t, err)
	e, err := repo.Create(ctx, Event{Name: "theatre", StartAt: time.Date(2021, 5, 4, 1, 1, 1, 0, time.Local),
		EndAt:  time.Date(2021, 5, 5, 1, 1, 1, 0, time.Local),
		Repeat: 0, ChildID: petya.ID})
	require.NoError(t, err)
	err = repo.Update(ctx, e.ID, EventParams{Name: new("sport"), Repeat: new(int64(1))})
	require.NoError(t, err)
	event, err := repo.GetById(ctx, int(e.ID))
	require.NoError(t, err)
	require.NotEqual(t, e, event)
	require.Equal(t, event.Name, "sport")
}

func TestEventRepository_Delete(t *testing.T) {
	pool := setupTestDb(t)
	repo := &EventRepository{pool}
	ctx := context.Background()
	userRepo := &UserRepository{pool}
	childRepo := &ChildRepository{pool}
	require.NoError(t, userRepo.Create(ctx, User{Username: "Ivan"}))
	petya, err := childRepo.Create(ctx, Child{Name: "Petya", UserID: 1})
	require.NoError(t, err)
	e, err := repo.Create(ctx, Event{Name: "theatre", StartAt: time.Date(2021, 5, 4, 1, 1, 1, 0, time.Local),
		EndAt:  time.Date(2021, 5, 5, 1, 1, 1, 0, time.Local),
		Repeat: 0, ChildID: petya.ID})
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, e.ID))
	_, err = repo.GetById(ctx, int(e.ID))
	require.Error(t, err)
}
