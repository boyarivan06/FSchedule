package database

import (
	"context"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDb(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	//const connStr = "postgres://go_user:machine_banana@localhost:5432/f_schedule"
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2). // Postgres логирует это дважды
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = pgContainer.Terminate(ctx)
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	m, err := migrate.New("file://../migrations", connStr)
	require.NoError(t, err)
	require.NoError(t, m.Up())
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return pool
}

func TestUserRepository_Get(t *testing.T) {
	pool := setupTestDb(t)
	repo := &UserRepository{pool}
	ctx := context.Background()
	require.NoError(t, repo.Create(ctx, User{Username: "alice"}))
	user, err := repo.GetByUsername(ctx, "alice")
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
	pool := setupTestDb(t)
	repo := &UserRepository{pool}
	ctx := context.Background()
	require.NoError(t, repo.Create(ctx, User{Username: "Biba"}))
	require.Error(t, repo.Create(ctx, User{Username: "Biba"}))
}

func TestUserRepository_Update(t *testing.T) {
	pool := setupTestDb(t)
	repo := &UserRepository{pool}
	ctx := context.Background()
	require.NoError(t, repo.Create(ctx, User{Username: "Biba"}))
	biba, err := repo.GetByUsername(ctx, "Biba")
	require.NoError(t, err)
	require.NoError(t, repo.Update(ctx, biba.ID, UserParams{UserName: new("Johnson")}))
	require.NoError(t, repo.Update(ctx, biba.ID, UserParams{}))
	new_biba, err := repo.GetById(ctx, int(biba.ID))
	require.NoError(t, err)
	require.Equal(t, new_biba.ID, biba.ID)
	require.Equal(t, "Johnson", new_biba.Username)
}

func TestUserRepository_Delete(t *testing.T) {
	pool := setupTestDb(t)
	repo := &UserRepository{pool}
	ctx := context.Background()
	require.NoError(t, repo.Create(ctx, User{Username: "Biba"}))
	biba, err := repo.GetByUsername(ctx, "Biba")
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, biba.ID))
	_, err = repo.GetByUsername(ctx, "biba")
	require.Error(t, err)
	_, err = repo.GetById(ctx, int(biba.ID))
	require.Error(t, err)
}
