package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"

	bin "github.com/gagliardetto/binary"
	solana "github.com/gagliardetto/solana-go"

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
				Method  string        `json:"method"`
				JsonRPC string        `json:"jsonrpc"`
				Params  []interface{} `json:"params"`
			}{}

			if err := json.Unmarshal([]byte(*event.Request.Body), &solanaRequestBody); err != nil {
				// Collect doesn't return an error, so we should log it.
				slog.Error("can't unmarshal solana request body: ", err)
				continue
			}

			if solanaRequestBody.Method != "sendTransaction" {
				continue
			}

			if len(solanaRequestBody.Params) > 0 {
				if transactionData, ok := solanaRequestBody.Params[0].(string); ok {
					// We don't need to decode the transaction signature, we just
					// need to store it.

					data, err := base64.StdEncoding.DecodeString(transactionData)
					if err != nil {
						slog.Error(
							"can't decode transaction base64: ",
							slog.String("transactionData", transactionData),
							err)
					}

					decodedTx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(data))
					if err != nil {
						slog.Error("can't decode transaction: ", err)
					}

					programIDs := make([]string, 0)

					for _, instruction := range decodedTx.Message.Instructions {
						programID, err := decodedTx.ResolveProgramIDIndex(instruction.ProgramIDIndex)
						if err != nil {
							slog.Error("can't resolve program ID index: ", err)
						}

						programIDs = append(programIDs, programID.String())
					}

					event.ProgramIDs = programIDs
				}
			}

			solanaResponseBody := struct {
				Result string `json:"result"`
			}{}

			if err := json.Unmarshal([]byte(*event.Response.Body), &solanaResponseBody); err != nil {
				slog.Error("can't unmarshal solana response body: ", err)
				continue
			}

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
