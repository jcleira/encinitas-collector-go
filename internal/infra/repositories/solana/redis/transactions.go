package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
)

const (
	channel = "solana_transactions"
)

// SubscribeToAgentTransactions subscribes to the 'solana_transactions' channel and listens
// for messages.
func (r *Repository) SubscribeToTransactions(
	ctx context.Context) (chan<- aggregates.Transaction, chan<- error) {
	pubsub := r.client.Subscribe(ctx, channel)

	transactionChannel := make(chan<- aggregates.Transaction)
	errorChannel := make(chan<- error)

	redisChannel := pubsub.Channel()
	go func() {
		defer pubsub.Close()
		for {
			select {
			case transaction := <-redisChannel:
				var redisTransaction redisTransaction
				if err := json.Unmarshal([]byte(transaction.Payload), &redisTransaction); err != nil {
					errorChannel <- fmt.Errorf("json.Unmarshal: %w", err)
					continue
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return transactionChannel, errorChannel
}

// PublishTransaction publishes an transaction to the 'solana_transactions' channel.
func (r *Repository) PublishTransaction(
	ctx context.Context, transaction aggregates.Transaction) error {
	redisTransaction := redisTransactionFromAggregate(transaction)

	message, err := json.Marshal(redisTransaction)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = r.client.Publish(ctx, channel, string(message)).Err()
	if err != nil {
		return fmt.Errorf("client.Publish: %w", err)
	}

	return nil
}
