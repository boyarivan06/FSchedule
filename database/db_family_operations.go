package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FamilyRepository struct {
	Db *pgxpool.Pool
}
type FamilyParams struct {
	Name    *string `database:"name" json:"name"`
	AdminId *int64  `database:"user_id" json:"user_id"`
}

func (repo *FamilyRepository) GetById(ctx context.Context, id int) (Family, error) {
	var e Family
	err := repo.Db.QueryRow(ctx, "SELECT * FROM families WHERE id = $1", id).Scan(&e.ID, &e.Name, &e.AdminID)
	return e, err
}

func (repo *FamilyRepository) GetManyByParams(ctx context.Context, p FamilyParams) ([]Family, error) {
	var families []Family
	var rows pgx.Rows
	var err error
	sets, args := prepareFamilyParams(p)
	var query string
	if len(sets) == 0 {
		query = "SELECT * FROM families"
	} else {
		args = append(args)
		query = fmt.Sprintf(
			"SELECT * FROM families WHERE %s",
			strings.Join(sets, " AND "),
		)
	}
	rows, err = repo.Db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f Family
		err = rows.Scan(&f.ID, &f.Name, &f.AdminID)
		if err != nil {
			return nil, err
		}
		families = append(families, f)
	}
	return families, err
}

func prepareFamilyParams(p FamilyParams) ([]string, []any) {
	var (
		sets []string
		args []any
	)

	add := func(col string, val any) {
		args = append(args, val)
		sets = append(sets, fmt.Sprintf("%s = $%d", col, len(args)))
	}

	if p.Name != nil {
		add("name", *p.Name)
	}
	if p.AdminId != nil {
		add("admin_id", *p.AdminId)
	}
	return sets, args
}

func (repo *FamilyRepository) Create(ctx context.Context, f Family) (Family, error) {
	_, err := repo.Db.Exec(ctx, "INSERT INTO families (name, admin_id) VALUES ($1, $2)", f.Name, f.AdminID)
	family := Family{}
	err = repo.Db.QueryRow(ctx, "SELECT * FROM families WHERE name = $1 AND admin_id = $2", f.Name, f.AdminID).Scan(&family.ID, &family.Name, &family.AdminID)
	if err != nil {
		return Family{}, err
	}
	return family, err
}

func (repo *FamilyRepository) Update(ctx context.Context, id int64, p FamilyParams) error {
	sets, args := prepareFamilyParams(p)

	if len(sets) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf(
		"UPDATE families SET %s WHERE id = $%d",
		strings.Join(sets, ", "),
		len(args),
	)

	tag, err := repo.Db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (repo *FamilyRepository) Delete(ctx context.Context, id int64) error {
	_, err := repo.Db.Exec(ctx, "DELETE FROM families WHERE id = $1", id)
	return err
}
