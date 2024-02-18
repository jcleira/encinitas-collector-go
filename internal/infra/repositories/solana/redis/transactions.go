package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
)

const (
	channel      = "solana_transactions"
	updatedOnKey = "solana_transactions_updated_on"
)

// SubscribeToAgentTransactions subscribes to the 'solana_transactions' channel and listens
// for messages.
func (r *Repository) SubscribeToTransactions(
	ctx context.Context) (chan aggregates.Transaction, chan error) {
	pubsub := r.client.Subscribe(ctx, channel)

	transactionChannel := make(chan aggregates.Transaction)
	errorChannel := make(chan error)

	redisChannel := pubsub.Channel()
	go func() {
		defer pubsub.Close()
		for {
			select {
			case <-ctx.Done():
				return

			case transaction := <-redisChannel:
				var redisTransaction redisTransaction
				if err := json.Unmarshal([]byte(transaction.Payload), &redisTransaction); err != nil {
					errorChannel <- fmt.Errorf("json.Unmarshal: %w", err)
					continue
				}

				transactionChannel <- redisTransaction.toAggregate()
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

// SetUpdatedOn sets the updatedOn value in the Redis database.
func (r *Repository) SetUpdatedOn(ctx context.Context, updatedOn int64) error {
	if err := r.client.Set(ctx, updatedOnKey, updatedOn, 0).Err(); err != nil {
		return fmt.Errorf("client.Set: %w", err)
	}

	return nil
}

// GetUpdatedOn gets the updatedOn value from the Redis database.
func (r *Repository) GetUpdatedOn(ctx context.Context) (int64, error) {
	updatedOn, err := r.client.Get(ctx, updatedOnKey).Int64()
	if err != nil {
		return 0, fmt.Errorf("client.Get: %w", err)
	}

	return updatedOn, nil
}
