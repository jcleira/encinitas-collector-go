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
	selectAllPrograms = `
SELECT *
FROM programs
WHERE deleted_at IS NULL
ORDER BY priority asc, created_at desc;
`

	insertProgram = `
INSERT INTO programs
(program_address, program_name, created_at)
VALUES (:program_address, :program_name, :created_at)
`
)

func (r *Repository) SelectAllPrograms(
	ctx context.Context) ([]aggregates.Program, error) {
	var dbPrograms dbPrograms
	if err := r.db.SelectContext(ctx,
		&dbPrograms, selectAllPrograms); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext, err: %w", err)
	}

	programs := make([]aggregates.Program, len(dbPrograms))
	for i, dbProgram := range dbPrograms {
		programs[i] = dbProgram.toAggregate()
	}

	return programs, nil
}

func (r *Repository) InsertProgram(
	ctx context.Context, program aggregates.Program) error {
	dbProgram := dbProgramFromAggregate(program)
	if _, err := sqlx.NamedExec(r.db, insertProgram, dbProgram); err != nil {
		return fmt.Errorf("sqlx.NamedExec, err: %w", err)
	}

	return nil
}

type dbProgram struct {
	ProgramAddress string       `db:"program_address"`
	ProgramName    string       `db:"program_name"`
	CreatedAt      time.Time    `db:"created_at"`
	UpdateAt       time.Time    `db:"updated_at"`
	DeleteAt       sql.NullTime `db:"deleted_at"`
	Priority       int          `db:"priority"`
}

type dbPrograms []dbProgram

func (dbe dbProgram) toAggregate() aggregates.Program {
	return aggregates.Program{
		ProgramAddress: dbe.ProgramAddress,
		ProgramName:    dbe.ProgramName,
	}
}

func dbProgramFromAggregate(e aggregates.Program) dbProgram {
	return dbProgram{
		ProgramAddress: e.ProgramAddress,
		ProgramName:    e.ProgramName,
	}
}
