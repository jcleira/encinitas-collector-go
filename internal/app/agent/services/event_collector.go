package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
)

type eventsRedisRepository interface {
	SubscribeToEvents(context.Context) (chan aggregates.Event, chan error)
	SetEvent(context.Context, string, aggregates.Event) error
}

// EventCollector define the dependencies to collect events.
type EventCollector struct {
	repository eventsRedisRepository
}

// NewEventCollector creates a new EventCollector.
func NewEventCollector(repository eventsRedisRepository) *EventCollector {
	return &EventCollector{
		repository: repository,
	}
}

// Collect starts collecting and processing events.
func (ec *EventCollector) Collect(ctx context.Context) {
	eventChan, errChan := ec.repository.SubscribeToEvents(ctx)

	for {
		select {
		case event := <-eventChan:
			// For the moment we are only interested in successful responses
			// so we can ignore the rest.
			//
			// But failure responses are event more important than successful
			// ones, so we should handle them as well.
			if event.Response == nil || event.Response.Status != 200 || event.Response.Body == nil {
				continue
			}

			if event.Request == nil || event.Request.Body == nil {
				continue
			}

			solanaRequestBody := struct {
				Method string `json:"method"`
			}{}

			if err := json.Unmarshal([]byte(*event.Request.Body), &solanaRequestBody); err != nil {
				// Collect doesn't return an error, so we should log it.
				slog.Error("can't unmarshal solana request body: ", err)
				continue
			}

			if solanaRequestBody.Method != "sendTransaction" {
				continue
			}

			solanaResponseBody := struct {
				Result string `json:"result"`
			}{}

			if err := json.Unmarshal([]byte(*event.Response.Body), &solanaResponseBody); err != nil {
				slog.Error("can't unmarshal solana response body: ", err)
				continue
			}

			slog.Info("Set event: ",
				solanaRequestBody.Method,
				solanaResponseBody.Result,
				slog.Any("event", event),
			)

			ec.repository.SetEvent(ctx,
				fmt.Sprintf("%s.%s", solanaRequestBody.Method, solanaResponseBody.Result),
				event,
			)

		case err := <-errChan:
			slog.Error("error in the agent events redis repository: ", err)

		case <-ctx.Done():
			return
		}
	}
}
