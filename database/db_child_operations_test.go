package database

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func PrepareTestObjects(pool *pgxpool.Pool, n int) (*User, *Family, *Child, *Event, error) {
	ctx := context.Background()
	if n <= 0 {
		return nil, nil, nil, nil, nil
	}
	userRepo := &UserRepository{pool}
	testUser, err := userRepo.Create(ctx, User{Username: "Ivan"})
	if n == 1 {
		return &testUser, nil, nil, nil, err
	}
	familyRepo := &FamilyRepository{pool}
	testFamily, err := familyRepo.Create(ctx, Family{Name: "petrovы", AdminID: testUser.ID})
	if n == 2 {
		return &testUser, &testFamily, nil, nil, err
	}
	childRepo := &ChildRepository{pool}
	testChild, err := childRepo.Create(ctx, Child{Name: "petya", UserID: testUser.ID, FamilyID: testFamily.ID})
	if n == 3 {
		return &testUser, &testFamily, &testChild, nil, err
	}
	eventRepo := &EventRepository{pool}
	testEvent, err := eventRepo.Create(ctx, Event{Name: "sports", Repeat: 0, ChildID: testChild.ID, StartAt: time.Date(2027, 5, 5, 5, 5, 0, 0, time.Local), EndAt: time.Date(2027, 5, 6, 5, 5, 0, 0, time.Local)})
	return &testUser, &testFamily, &testChild, &testEvent, err
}

func TestChildRepository_Create(t *testing.T) {
	pool := SetupTestDb(t)
	repo := &ChildRepository{pool}
	ctx := context.Background()
	testUser, testFamily, _, _, err := PrepareTestObjects(pool, 2)
	require.NoError(t, err)
	_, err = repo.Create(ctx, Child{Name: "Petya", UserID: testUser.ID, FamilyID: testFamily.ID})
	require.NoError(t, err)
	_, err = repo.Create(ctx, Child{Name: "Fedya", UserID: testUser.ID + 1, FamilyID: testFamily.ID})
	require.Error(t, err)
}

func TestChildRepository_Get(t *testing.T) {
	pool := SetupTestDb(t)
	repo := &ChildRepository{pool}
	ctx := context.Background()
	testUser, testFamily, _, _, err := PrepareTestObjects(pool, 2)
	require.NoError(t, err)
	_, err = repo.Create(ctx, Child{Name: "Petya", UserID: testUser.ID, FamilyID: testFamily.ID})
	require.NoError(t, err)
	children, err := repo.GetManyByParams(ctx, ChildParams{Name: new("Petya")})
	require.NoError(t, err)
	require.Equal(t, 1, len(children))
	ch := children[0]
	child, err := repo.GetById(ctx, int(ch.ID))
	require.NoError(t, err)
	require.Equal(t, ch, child)
}

func TestChildRepository_Update(t *testing.T) {
	pool := SetupTestDb(t)
	repo := &ChildRepository{pool}
	ctx := context.Background()
	testUser, testFamily, _, _, err := PrepareTestObjects(pool, 2)
	require.NoError(t, err)
	child, err := repo.Create(ctx, Child{Name: "Petya", UserID: testUser.ID, FamilyID: testFamily.ID})
	require.NoError(t, err)
	require.NoError(t, repo.Update(ctx, child.ID, ChildParams{Name: new("Fedya")}))
	newChild, err := repo.GetById(ctx, int(child.ID))
	require.NoError(t, err)
	require.Equal(t, newChild.Name, "Fedya")
}

func TestChildRepository_Delete(t *testing.T) {
	pool := SetupTestDb(t)
	repo := &ChildRepository{pool}
	ctx := context.Background()
	testUser, testFamily, _, _, err := PrepareTestObjects(pool, 2)
	require.NoError(t, err)
	child, err := repo.Create(ctx, Child{Name: "Petya", UserID: testUser.ID, FamilyID: testFamily.ID})
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, child.ID))
	_, err = repo.GetById(ctx, int(child.ID))
	require.Error(t, err)
}
