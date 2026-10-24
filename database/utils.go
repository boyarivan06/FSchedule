package database

import (
	"context"
	"fmt"
	"strings"
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

func buildQuery(filters map[string]string, separator string) (string, []any) {
	allowed := map[string]bool{
		"name":     true,
		"username": true,
		"user_id":  true,
		"start_at": true,
		"end_at":   true,
		"repeat":   true,
		"child_id": true,
	}

	var conditions []string
	var args []any
	i := 1

	for key, val := range filters {
		if !allowed[key] {
			continue
		}

		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, val)
		i++
	}

	query := ""
	if len(conditions) > 0 {
		query += strings.Join(conditions, separator)
	}
	return query, args
}

func SetupTestDb(t *testing.T) *pgxpool.Pool {
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
