package services

import (
	"context"
	"fmt"

	"github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
)

type eventSubscriber interface {
	SubscribeToEvents(context.Context) (chan aggregates.Event, chan error)
}

// EventCollector define the dependencies to collect events.
type EventCollector struct {
	subscriber eventSubscriber
}

// NewEventCollector creates a new EventCollector.
func NewEventCollector(subscriber eventSubscriber) *EventCollector {
	return &EventCollector{
		subscriber: subscriber,
	}
}

// Collect starts collecting and processing events.
func (ec *EventCollector) Collect(ctx context.Context) {
	eventChan, errChan := ec.subscriber.SubscribeToEvents(ctx)

	for {
		select {
		case event := <-eventChan:
			if event.Response == nil {
				continue
			}

			// For the moment we are only interested in successful responses
			// so we can ignore the rest.
			//
			// But failure responses are event more important than successful
			// ones, so we should handle them as well.
			if event.Response.Status != 200 {
				continue
			}

			if event.Request.Body == nil {
				continue
			}

			if event.Response.Body == nil {
				continue
			}

			struct  solanaBody {
				Method string `json:"method"`
				Params string `json:"params"`
			}{
			}


			jsonBody, err := event.Response.Body.MarshalJSON()


		case err := <-errChan:
			fmt.Println(err)
			// Handle error
		case <-ctx.Done():
			return
		}
	}
}
