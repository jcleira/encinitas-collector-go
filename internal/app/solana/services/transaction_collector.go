package services

import (
	"fmt"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
)

// TransactionRepository define the methods to become a transaction repository.
type TransactionRepository interface {
	SelectTransactionByUpdatedOn(
		updatedOn int64) ([]aggregates.Transaction, error)
}

// TransactionsCollector define the dependencies needed to perform the
// transaction collection.
type TransactionCollector struct {
	repository TransactionRepository
}

func NewTransactionCollector(
	repository TransactionRepository) *TransactionCollector {
	return &TransactionCollector{
		repository: repository,
	}
}

// CollectTransactions collects transactions from the database.
func (tc *TransactionCollector) CollectTransactions(
	updatedOn int64) ([]aggregates.Transaction, error) {
	transactions, err := tc.repository.SelectTransactionByUpdatedOn(updatedOn)
	if err != nil {
		return nil, fmt.Errorf("tc.repository.SelectTransactionByUpdatedOn, err: %w", err)
	}

	return transactions, nil
}
