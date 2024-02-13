package services

import (
	"context"
	"fmt"

	"github.com/jcleira/encinitas-collector-go/internal/app/manager/aggregates"
)

type programCreatorRepository interface {
	InsertProgram(context.Context, aggregates.Program) error
}

// ProgramCreator defines the methods needed to create programs.
type ProgramCreator struct {
	programCreatorRepository
}

// NewProgramCreator initializes a new ProgramCreator.
func NewProgramCreator(
	programCreatorRepository programCreatorRepository) *ProgramCreator {
	return &ProgramCreator{
		programCreatorRepository: programCreatorRepository,
	}
}

// Create creates a new program.
func (pc *ProgramCreator) Create(ctx context.Context, program aggregates.Program) error {
	program.ProgramName = program.ProgramAddress
	if err := pc.programCreatorRepository.InsertProgram(ctx, program); err != nil {
		return fmt.Errorf("pc.programCreatorRepository.InsertProgram, err: %w", err)
	}

	return nil
}
