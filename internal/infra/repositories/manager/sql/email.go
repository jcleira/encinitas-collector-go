package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/jcleira/encinitas-collector-go/internal/app/manager/aggregates"
)

const (
	selectEmailsByEmail = `
SELECT *
FROM emails
WHERE email = :email
deleted_at IS NULL
`

	insertEmail = `
INSERT INTO emails
(email, email_name, created_at)
VALUES (:email, :email_name, :created_at)
`
)

func (r *Repository) CheckEmailExists(ctx context.Context, email string) error {
	dbEmails := []dbEmails{}
	if err := r.db.SelectContext(ctx,
		&dbEmails, selectEmailsByEmail, email); err != nil {
		return fmt.Errorf("r.db.SelectContext, err: %w", err)
	}

	if len(dbEmails) > 0 {
		return aggregates.ErrEmailAlreadyExists
	}

	return nil
}

func (r *Repository) InsertEmail(
	ctx context.Context, email string) error {
	dbEmail := dbEmailFromAggregate(email)
	if _, err := sqlx.NamedExec(r.db, insertEmail, dbEmail); err != nil {
		return fmt.Errorf("sqlx.NamedExec, err: %w", err)
	}

	return nil
}

type dbEmail struct {
	Email     string       `db:"email"`
	CreatedAt time.Time    `db:"created_at"`
	DeleteAt  sql.NullTime `db:"deleted_at"`
}

type dbEmails []dbEmail

func (dbe dbEmail) toAggregate() string {
	return dbe.Email
}

func dbEmailFromAggregate(email string) dbEmail {
	return dbEmail{
		Email: email,
	}
}
