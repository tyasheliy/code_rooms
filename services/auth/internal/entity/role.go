package entity

import "context"

type Role struct {
	Id   int
	Name string
}

type RoleRepo interface {
	GetById(ctx context.Context, id int) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
}
