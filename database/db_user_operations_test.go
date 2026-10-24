package database

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Get(t *testing.T) {
	pool := SetupTestDb(t)
	repo := &UserRepository{pool}
	ctx := context.Background()
	user, err := repo.Create(ctx, User{Username: "ivan"})
	require.NoError(t, err)
	require.Equal(t, "alice", user.Username)
	id := user.ID
	newUser, err := repo.GetById(ctx, int(id))
	require.NoError(t, err)
	require.Equal(t, "alice", newUser.Username)
}

//func TestUserRepository_GetMany(t *testing.T) {
//	pool := setupTestDb(t)
//	repo := &UserRepository{pool}
//	ctx := context.Background()
//	require.NoError(t, repo.Create(ctx, User{Username: "Biba"}))
//	require.NoError(t, repo.Create(ctx, User{Username: "Boba"}))
//	params := UserParams{new("Boba")}
//	users1, err := repo.GetManyByParams(ctx, params)
//	require.NoError(t, err)
//	require.Equal(t, len(users1), 1)
//	require.Equal(t, users1[0].Username, "Boba")
//	users2, err := repo.GetManyByParams(ctx, UserParams{})
//	require.NoError(t, err)
//	require.GreaterOrEqual(t, len(users2), 2)
//}

func TestUserRepository_Create(t *testing.T) {
	pool := SetupTestDb(t)
	repo := &UserRepository{pool}
	ctx := context.Background()
	_, err := repo.Create(ctx, User{Username: "Biba"})
	require.NoError(t, err)
	_, err = repo.Create(ctx, User{Username: "Biba"})
	require.Error(t, err)
}

func TestUserRepository_Update(t *testing.T) {
	pool := SetupTestDb(t)
	repo := &UserRepository{pool}
	ctx := context.Background()
	biba, err := repo.Create(ctx, User{Username: "Biba"})
	require.NoError(t, err)
	require.NoError(t, repo.Update(ctx, biba.ID, UserParams{UserName: new("Johnson")}))
	require.NoError(t, repo.Update(ctx, biba.ID, UserParams{}))
	newBiba, err := repo.GetById(ctx, int(biba.ID))
	require.NoError(t, err)
	require.Equal(t, newBiba.ID, biba.ID)
	require.Equal(t, "Johnson", newBiba.Username)
}

func TestUserRepository_Delete(t *testing.T) {
	pool := SetupTestDb(t)
	repo := &UserRepository{pool}
	ctx := context.Background()
	biba, err := repo.Create(ctx, User{Username: "Biba"})
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, biba.ID))
	_, err = repo.GetByUsername(ctx, "biba")
	require.Error(t, err)
	_, err = repo.GetById(ctx, int(biba.ID))
	require.Error(t, err)
}
