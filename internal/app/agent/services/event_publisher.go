package services

import (
	"context"
	"fmt"

	"github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
)

type eventPublisher interface {
	PublishEvent(context.Context, aggregates.Event) error
}

// EventPublisher defines the dependencies to publish events.
type EventPublisher struct {
	eventPublisher eventPublisher
}

// NewEventPublisher initializes a new EventPublisher.
func NewEventPublisher(
	eventPublisher eventPublisher) *EventPublisher {
	return &EventPublisher{
		eventPublisher: eventPublisher,
	}
}

// Publish publishes an event.
func (ep *EventPublisher) Publish(
	ctx context.Context, event aggregates.Event) error {
	if err := ep.eventPublisher.PublishEvent(ctx, event); err != nil {
		return fmt.Errorf("eventPublisher.PublishEvent: %w", err)
	}

	return nil
}
