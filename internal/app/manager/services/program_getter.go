package services

import (
	"context"
	"fmt"

	"github.com/jcleira/encinitas-collector-go/internal/app/manager/aggregates"
)

type programGetterRepository interface {
	SelectAllPrograms(context.Context) ([]aggregates.Program, error)
}

// ProgramGetter defines the methods needed to get programs.
type ProgramGetter struct {
	programGetterRepository programGetterRepository
}

// NewProgramGetter initializes a new ProgramGetter.
func NewProgramGetter(
	programGetterRepository programGetterRepository) *ProgramGetter {
	return &ProgramGetter{
		programGetterRepository: programGetterRepository,
	}
}

// GetPrograms gets all programs.
func (pg *ProgramGetter) GetPrograms(
	ctx context.Context) ([]aggregates.Program, error) {
	programs, err := pg.programGetterRepository.SelectAllPrograms(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"pg.programGetterRepository.SelectAllPrograms, err: %w", err)
	}

	return programs, nil
}
