package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChildRepository_Create(t *testing.T) {
	pool := setupTestDb(t)
	repo := &ChildRepository{pool}
	ctx := context.Background()
	userRepo := &UserRepository{pool}
	require.NoError(t, userRepo.Create(ctx, User{Username: "Ivan"}))
	_, err := repo.Create(ctx, Child{Name: "Petya", UserID: 1})
	require.NoError(t, err)
	_, err = repo.Create(ctx, Child{Name: "Fedya", UserID: 2})
	require.Error(t, err)
}

func TestChildRepository_Get(t *testing.T) {
	pool := setupTestDb(t)
	repo := &ChildRepository{pool}
	ctx := context.Background()
	userRepo := &UserRepository{pool}
	require.NoError(t, userRepo.Create(ctx, User{Username: "Ivan"}))
	_, err := repo.Create(ctx, Child{Name: "Petya", UserID: 1})
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
	pool := setupTestDb(t)
	repo := &ChildRepository{pool}
	ctx := context.Background()
	userRepo := &UserRepository{pool}
	require.NoError(t, userRepo.Create(ctx, User{Username: "Ivan"}))
	child, err := repo.Create(ctx, Child{Name: "Petya", UserID: 1})
	require.NoError(t, err)
	require.NoError(t, repo.Update(ctx, child.ID, ChildParams{Name: new("Fedya")}))
	newChild, err := repo.GetById(ctx, int(child.ID))
	require.NoError(t, err)
	require.Equal(t, newChild.Name, "Fedya")
}

func TestChildRepository_Delete(t *testing.T) {
	pool := setupTestDb(t)
	repo := &ChildRepository{pool}
	ctx := context.Background()
	userRepo := &UserRepository{pool}
	require.NoError(t, userRepo.Create(ctx, User{Username: "Ivan"}))
	child, err := repo.Create(ctx, Child{Name: "Petya", UserID: 1})
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, child.ID))
	_, err = repo.GetById(ctx, int(child.ID))
	require.Error(t, err)
}
