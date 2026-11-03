package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFamilyRepository_Create(t *testing.T) {
	pool := SetupTestDb(t)
	ctx := context.Background()
	repo := FamilyRepository{pool}
	testUser, _, _, _, err := PrepareTestObjects(pool, 1)
	require.NoError(t, err)
	family, err := repo.Create(ctx, Family{Name: "family", AdminID: testUser.ID})
	require.NoError(t, err)
	require.Equal(t, "family", family.Name)
}

func TestFamilyRepository_Get(t *testing.T) {
	pool := SetupTestDb(t)
	ctx := context.Background()
	repo := FamilyRepository{pool}
	testUser, _, _, _, err := PrepareTestObjects(pool, 1)
	require.NoError(t, err)
	family, err := repo.Create(ctx, Family{Name: "family", AdminID: testUser.ID})
	require.NoError(t, err)
	f2, err := repo.GetById(ctx, int(family.ID))
	require.NoError(t, err)
	require.Equal(t, family, f2)
	f3s, err := repo.GetManyByParams(ctx, FamilyParams{Name: new("family")})
	require.NoError(t, err)
	require.Equal(t, 1, len(f3s))
	require.Equal(t, f2, f3s[0])
}

func TestFamilyRepository_Update(t *testing.T) {
	pool := SetupTestDb(t)
	ctx := context.Background()
	repo := FamilyRepository{pool}
	testUser, _, _, _, err := PrepareTestObjects(pool, 1)
	require.NoError(t, err)
	family, err := repo.Create(ctx, Family{Name: "family", AdminID: testUser.ID})
	require.NoError(t, err)
	require.NoError(t, repo.Update(ctx, family.ID, FamilyParams{Name: new("new_name")}))
	f2, err := repo.GetById(ctx, int(family.ID))
	require.NoError(t, err)
	require.Equal(t, "new_name", f2.Name)
}

func TestFamilyRepository_Delete(t *testing.T) {
	pool := SetupTestDb(t)
	ctx := context.Background()
	repo := FamilyRepository{pool}
	testUser, _, _, _, err := PrepareTestObjects(pool, 1)
	require.NoError(t, err)
	family, err := repo.Create(ctx, Family{Name: "family", AdminID: testUser.ID})
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, family.ID))
	_, err = repo.GetById(ctx, int(family.ID))
	require.Error(t, err)
}
