package entity

import "github.com/google/uuid"

type Permission struct {
	ID        uuid.UUID `db:"id"`
	Operation string    `db:"operation"`
	Resource  string    `db:"resource"`
}

func (p *Permission) Code() string {
	return p.Operation + "_" + p.Resource
}
