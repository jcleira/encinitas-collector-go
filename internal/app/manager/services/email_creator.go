package services

import (
	"context"
	"fmt"
)

type emailCreatorRepository interface {
	CheckEmailExists(context.Context, string) error
	InsertEmail(context.Context, string) error
}

// EmailCreator defines the methods needed to create emails.
type EmailCreator struct {
	emailCreatorRepository
}

// NewEmailCreator initializes a new EmailCreator.
func NewEmailCreator(
	emailCreatorRepository emailCreatorRepository) *EmailCreator {
	return &EmailCreator{
		emailCreatorRepository: emailCreatorRepository,
	}
}

// Create creates a new email.
func (pc *EmailCreator) Create(ctx context.Context, email string) error {
	if err := pc.emailCreatorRepository.CheckEmailExists(ctx, email); err != nil {
		return fmt.Errorf("pc.emailCreatorRepository.CheckEmailExists, err: %w", err)
	}

	if err := pc.emailCreatorRepository.InsertEmail(ctx, email); err != nil {
		return fmt.Errorf("pc.emailCreatorRepository.InsertEmail, err: %w", err)
	}

	return nil
}
