package entity

import "context"

type User struct {
	Id           int64
	RoleId       int
	Login        string
	PasswordHash string
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
