package repo

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/auth/internal/storage"
)

const roles_table = "roles"

type PostgresRoleRepo struct {
	storage *storage.Storage
}

func NewPostgresRoleRepo(storage *storage.Storage) *PostgresRoleRepo {
	return &PostgresRoleRepo{storage: storage}
}

func (r *PostgresRoleRepo) GetById(ctx context.Context, id int) (*entity.Role, error) {
	query, _, err := goqu.From(roles_table).
		Select("id", "name").
		Where(goqu.Ex{"id": id}).
		ToSQL()
	if err != nil {
		return nil, err
	}

	row := r.storage.QueryRowContext(ctx, query)

	var role entity.Role

	err = row.Scan(&role.Id, &role.Name)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *PostgresRoleRepo) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	query, _, err := goqu.From(roles_table).
		Select("id", "name").
		Where(goqu.Ex{"name": name}).
		ToSQL()
	if err != nil {
		return nil, err
	}

	row := r.storage.QueryRowContext(ctx, query)

	var role entity.Role

	if err = row.Scan(&role.Id, &role.Name); err != nil {
		return nil, err
	}

	return &role, nil
}
