package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
)

const (
	channel = "agent_events"
)

// SubscribeToAgentEvents subscribes to the 'agent_events' channel and listens
// for messages.
func (r *Repository) SubscribeToEvents(
	ctx context.Context) (chan aggregates.Event, chan error) {
	pubsub := r.client.Subscribe(ctx, channel)

	eventChannel := make(chan aggregates.Event)
	errorChannel := make(chan error)

	redisChannel := pubsub.Channel()
	go func() {
		defer pubsub.Close()
		for {
			select {
			case event := <-redisChannel:
				var redisEvent redisEvent
				if err := json.Unmarshal([]byte(event.Payload), &redisEvent); err != nil {
					errorChannel <- fmt.Errorf("json.Unmarshal: %w", err)
					continue
				}

				eventChannel <- redisEvent.toAggregate()

			case <-ctx.Done():
				return
			}
		}
	}()

	return eventChannel, errorChannel
}

// PublishEvent publishes an event to the 'agent_events' channel.
func (r *Repository) PublishEvent(
	ctx context.Context, event aggregates.Event) error {
	message, err := json.Marshal(redisEventFromAggregate(event))
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = r.client.Publish(ctx, channel, string(message)).Err()
	if err != nil {
		return fmt.Errorf("client.Publish: %w", err)
	}

	return nil
}

// SetEvent sets an event in the redis repository.
func (r *Repository) SetEvent(
	ctx context.Context, key string, event aggregates.Event) error {
	message, err := json.Marshal(redisEventFromAggregate(event))
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	if err = r.client.Set(ctx, key, string(message), 0).Err(); err != nil {
		return fmt.Errorf("client.Set: %w", err)
	}

	return nil
}
