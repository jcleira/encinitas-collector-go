package sql

import (
	"context"
	"database/sql"
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

	insertTransactionDetailQuery = `
INSERT INTO encinitas_transaction_details
(program_address, updated_on, rpc_time, solana_time, total_time)
VALUES
(:program_address, :updated_on, :rpc_time, :solana_time, :total_time);
`

	getTransactionDetailAggregated = `
WITH total_time AS (
  SELECT program_address, SUM(total_time) as total_time
  FROM encinitas_transaction_details
  GROUP BY program_address
),
total_sum AS (
  SELECT
    SUM(total_time) AS overall_total_time
  FROM
   total_time
)

SELECT
  T.program_address,
  (T.total_time / TS.overall_total_time::FLOAT) * 100 AS percentage
FROM
  total_time T,
  total_sum TS;
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

// InsertTransactionDetail inserts a new transaction detail into the database.
func (r *Repository) InsertTransactionDetail(ctx context.Context,
	detail aggregates.TransactionDetail) error {
	dbTransactionDetail := dbTransactionDetailFromAggregate(detail)

	if _, err := sqlx.NamedExec(r.db, insertTransactionDetailQuery,
		dbTransactionDetail); err != nil {
		return fmt.Errorf("sqlx.NamedExec, err: %w", err)
	}

	return nil
}

func (r *Repository) GetTransactionDetailAggregated(
	ctx context.Context) ([]aggregates.TransactionDetailAggregated, error) {
	var dbTransactionDetailsAggregated dbTransactionDetailsAggregated
	if err := r.db.SelectContext(ctx,
		&dbTransactionDetailsAggregated, getTransactionDetailAggregated); err != nil {
		return nil, fmt.Errorf("r.db.SelectContext, err: %w", err)
	}

	transactionDetailsAggregated := make([]aggregates.TransactionDetailAggregated,
		len(dbTransactionDetailsAggregated))
	for i, dbTransactionDetail := range dbTransactionDetailsAggregated {
		transactionDetailsAggregated[i] = dbTransactionDetail.toAggregate()
	}

	return transactionDetailsAggregated, nil
}

type dbTransaction struct {
	Slot            int64          `db:"slot"`
	Signature       string         `db:"signature"`
	IsVote          bool           `db:"is_vote"`
	MessageType     int            `db:"message_type"`
	LegacyMessage   string         `db:"legacy_message"`
	V0LoadedMessage sql.NullString `db:"v0_loaded_message"`
	Signatures      string         `db:"signatures"`
	MessageHash     []byte         `db:"message_hash"`
	Meta            string         `db:"meta"`
	WriteVersion    int64          `db:"write_version"`
	UpdatedOn       time.Time      `db:"updated_on"`
	TxnIndex        int64          `db:"txn_index"`
	InsertedAt      time.Time      `db:"inserted_at"`
	ProcessedAt     *time.Time     `db:"processed_at,omitempty"`
}

type dbTransactions []dbTransaction

func (dbe dbTransaction) toAggregate() aggregates.Transaction {
	return aggregates.Transaction{
		Slot:            dbe.Slot,
		Signature:       dbe.Signature,
		IsVote:          dbe.IsVote,
		MessageType:     dbe.MessageType,
		LegacyMessage:   dbe.LegacyMessage,
		V0LoadedMessage: dbe.V0LoadedMessage.String,
		Signatures:      dbe.Signatures,
		MessageHash:     dbe.MessageHash,
		Meta:            dbe.Meta,
		WriteVersion:    dbe.WriteVersion,
		UpdatedOn:       dbe.UpdatedOn,
		TxnIndex:        dbe.TxnIndex,
		ProcessedAt:     dbe.ProcessedAt,
	}
}

func dbTransactionFromAggregate(e aggregates.Transaction) dbTransaction {
	return dbTransaction{
		Slot:          e.Slot,
		Signature:     e.Signature,
		IsVote:        e.IsVote,
		MessageType:   e.MessageType,
		LegacyMessage: e.LegacyMessage,
		V0LoadedMessage: sql.NullString{
			String: e.V0LoadedMessage,
			Valid:  e.V0LoadedMessage != "",
		},
		Signatures:   e.Signatures,
		MessageHash:  e.MessageHash,
		Meta:         e.Meta,
		WriteVersion: e.WriteVersion,
		UpdatedOn:    e.UpdatedOn,
		TxnIndex:     e.TxnIndex,
		ProcessedAt:  e.ProcessedAt,
	}
}

type dbTransactionDetail struct {
	ProgramAddress string    `db:"program_address"`
	UpdatedOn      time.Time `db:"updated_on"`
	RPCTime        int64     `db:"rpc_time"`
	SolanaTime     int64     `db:"solana_time"`
	TotalTime      int64     `db:"total_time"`
}

type dbTransactionDetails []dbTransactionDetail

func (dbe dbTransactionDetail) toAggregate() aggregates.TransactionDetail {
	return aggregates.TransactionDetail{
		ProgramAddress: dbe.ProgramAddress,
		UpdatedOn:      dbe.UpdatedOn,
		RPCTime:        dbe.RPCTime,
		SolanaTime:     dbe.SolanaTime,
		TotalTime:      dbe.TotalTime,
	}
}

func dbTransactionDetailFromAggregate(e aggregates.TransactionDetail) dbTransactionDetail {
	return dbTransactionDetail{
		ProgramAddress: e.ProgramAddress,
		UpdatedOn:      e.UpdatedOn,
		RPCTime:        e.RPCTime,
		SolanaTime:     e.SolanaTime,
		TotalTime:      e.TotalTime,
	}
}

type dbTransactionDetailAggregated struct {
	ProgramAddress string  `db:"program_address"`
	Percentage     float64 `db:"percentage"`
}

type dbTransactionDetailsAggregated []dbTransactionDetailAggregated

func (dbe dbTransactionDetailAggregated) toAggregate() aggregates.TransactionDetailAggregated {
	return aggregates.TransactionDetailAggregated{
		ProgramAddress: dbe.ProgramAddress,
		Percentage:     dbe.Percentage,
	}
}

func dbTransactionDetailAggregatedFromAggregate(
	e aggregates.TransactionDetailAggregated) dbTransactionDetailAggregated {
	return dbTransactionDetailAggregated{
		ProgramAddress: e.ProgramAddress,
		Percentage:     e.Percentage,
	}
}
