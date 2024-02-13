package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
	"github.com/jmoiron/sqlx"
)

const (
	selectTransactionsByProcessedAt = `
SELECT * FROM encinitas_transactions WHERE processed_at is NULL LIMIT 1000;
`
	updateTransactionProcessedAt = `
UPDATE encinitas_transactions
SET processed_at = :processed_at
WHERE signature = :signature;
`
)

func (r *Repository) SelectTransactionsByProcessedAt(
	ctx context.Context) ([]aggregates.Transaction, error) {
	var dbTransactions dbTransactions
	if err := r.db.SelectContext(ctx,
		&dbTransactions, selectTransactionsByProcessedAt); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext, err: %w", err)
	}

	transactions := make([]aggregates.Transaction, len(dbTransactions))
	for i, dbTransaction := range dbTransactions {
		transactions[i] = dbTransaction.toAggregate()
	}

	return transactions, nil
}

func (r *Repository) UpdateTransactionProcessedAt(
	ctx context.Context, signature string, processedAt time.Time) error {
	if _, err := sqlx.NamedExec(r.db, updateTransactionProcessedAt,
		map[string]interface{}{
			"processed_at": processedAt,
			"signature":    signature,
		}); err != nil {
		return fmt.Errorf("sqlx.NamedExec, err: %w", err)
	}

	return nil
}

type dbTransaction struct {
	Slot            int64      `db:"slot"`
	Signature       string     `db:"signature"`
	IsVote          bool       `db:"is_vote"`
	MessageType     int        `db:"message_type"`
	LegacyMessage   string     `db:"legacy_message"`
	V0LoadedMessage *string    `db:"v0_loaded_message"`
	Signatures      string     `db:"signatures"`
	MessageHash     []byte     `db:"message_hash"`
	Meta            string     `db:"meta"`
	WriteVersion    int64      `db:"write_version"`
	UpdatedOn       time.Time  `db:"updated_on"`
	TxnIndex        int64      `db:"txn_index"`
	InsertedAt      time.Time  `db:"inserted_at"`
	ProcessedAt     *time.Time `db:"processed_at,omitempty"`
	ErrorInfo       *string    `db:"error_info,omitempty"`
}

type dbTransactions []dbTransaction

func (dbe dbTransaction) toAggregate() aggregates.Transaction {
	return aggregates.Transaction{
		Slot:         dbe.Slot,
		Signature:    dbe.Signature,
		IsVote:       dbe.IsVote,
		MessageType:  dbe.MessageType,
		Signatures:   dbe.Signatures,
		MessageHash:  dbe.MessageHash,
		Meta:         dbe.Meta,
		WriteVersion: dbe.WriteVersion,
		UpdatedOn:    dbe.UpdatedOn,
		TxnIndex:     dbe.TxnIndex,
		ErrorInfo:    dbe.ErrorInfo,
		ProcessedAt:  dbe.ProcessedAt,
	}
}

func dbTransactionFromAggregate(e aggregates.Transaction) dbTransaction {
	return dbTransaction{
		Slot:         e.Slot,
		Signature:    e.Signature,
		IsVote:       e.IsVote,
		MessageType:  e.MessageType,
		Signatures:   e.Signatures,
		MessageHash:  e.MessageHash,
		Meta:         e.Meta,
		WriteVersion: e.WriteVersion,
		UpdatedOn:    e.UpdatedOn,
		TxnIndex:     e.TxnIndex,
		ErrorInfo:    e.ErrorInfo,
		ProcessedAt:  e.ProcessedAt,
	}
}
