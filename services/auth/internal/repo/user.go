package repo

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/tyasheliy/code_rooms/services/auth/internal/entity"
	"github.com/tyasheliy/code_rooms/services/auth/internal/storage"
)

const user_table = "users"

type PostgresUserRepo struct {
	storage *storage.Storage
}

func NewPostgresUserRepo(storage *storage.Storage) *PostgresUserRepo {
	return &PostgresUserRepo{
		storage: storage,
	}
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	tx, err := r.storage.Begin()
	if err != nil {
		return nil, err
	}

	query, _, err := goqu.Insert(user_table).
		Cols("login", "password").
		Vals(goqu.Vals{user.Login, user.PasswordHash}).
		ToSQL()

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(query)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.getByLogin(ctx, user.Login)
}

func (r *PostgresUserRepo) GetById(ctx context.Context, id int64) (*entity.User, error) {
	query, _, err := goqu.From(user_table).Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return nil, err
	}

	row := r.storage.QueryRowContext(ctx, query)

	var user entity.User

	err = row.Scan(&user.Id, &user.Login, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepo) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	return r.getByLogin(ctx, login)
}

func (r *PostgresUserRepo) getByLogin(ctx context.Context, login string) (*entity.User, error) {
	query, _, err := goqu.From(user_table).Where(goqu.Ex{"login": login}).ToSQL()
	if err != nil {
		return nil, err
	}

	row := r.storage.QueryRowContext(ctx, query)

	var user entity.User

	err = row.Scan(&user.Id, &user.Login, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepo) Update(ctx context.Context, user *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *PostgresUserRepo) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
