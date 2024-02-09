package sql

import (
	"github.com/jmoiron/sqlx"
)

// Repository is a SQL repository for interactions.
type Repository struct {
	db *sqlx.DB
}

// New returns a new SQL repository for interactions.
func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
