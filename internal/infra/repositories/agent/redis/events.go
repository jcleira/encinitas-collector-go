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
	ctx context.Context) (chan<- aggregates.Event, chan<- error) {
	pubsub := r.client.Subscribe(ctx, channel)

	eventChannel := make(chan<- aggregates.Event)
	errorChannel := make(chan<- error)

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
	redisEvent := redisEventFromAggregate(event)

	message, err := json.Marshal(redisEvent)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = r.client.Publish(ctx, channel, string(message)).Err()
	if err != nil {
		return fmt.Errorf("client.Publish: %w", err)
	}

	return nil
}
