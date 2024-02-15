package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
)

// TransactionsSQLRepository define the methods to become a transaction's
// SQL sqlRepository.
type TransactionsSQLRepository interface {
	SelectTransactionsByProcessedAt(
		context.Context) ([]aggregates.Transaction, error)
	UpdateTransactionProcessedAt(
		context.Context, string, time.Time) error
}

type TransactionsRedisRepository interface {
	PublishTransaction(context.Context, aggregates.Transaction) error
}

// TransactionsCollector  define the dependencies needed to perform the
// transaction collection.
type TransactionsCollector struct {
	sqlRepository   TransactionsSQLRepository
	redisRepository TransactionsRedisRepository
}

// NewTransactionsCollector creates a new transaction collector.
func NewTransactionsCollector(
	sqlRepository TransactionsSQLRepository,
	redisRepository TransactionsRedisRepository,
) *TransactionsCollector {
	return &TransactionsCollector{
		sqlRepository:   sqlRepository,
		redisRepository: redisRepository,
	}
}

// CollectTransactions collects transactions from the database.
func (tc *TransactionsCollector) Collect(ctx context.Context) {
	for {
		// Sleep for 1 second before collecting transactions.
		// We should use a better approach to avoid busy waiting.
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			return

		default:
			transactions, err := tc.sqlRepository.SelectTransactionsByProcessedAt(ctx)
			if err != nil {
				slog.Error("tc.sqlRepository.SelectTransactionsByProcessedAt", err)
				continue
			}

			for _, transaction := range transactions {
				if err := tc.redisRepository.PublishTransaction(ctx, transaction); err != nil {
					slog.Error("tc.sqlRepository.PublishTransaction", err)
				}

				if err := tc.sqlRepository.UpdateTransactionProcessedAt(
					ctx, transaction.Signature, time.Now()); err != nil {
					slog.Error("tc.sqlRepository.SetUpdatedOn", err)
				}
			}
		}
	}
}
