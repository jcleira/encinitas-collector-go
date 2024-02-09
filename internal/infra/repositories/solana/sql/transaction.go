package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
	"github.com/jmoiron/sqlx"
)

const (
	selectTransactionByUpdatedOn = `
SELECT * FROM transactions WHERE updated_on > ?
`
	insertTransaction = `
INSERT INTO transactions
VALUES
`
	getTransaction = `
SELECT * FROM transactions WHERE id = ?;
`
)

func (r *Repository) SelectTransactionByUpdatedOn(
	ctx context.Context, updatedOn int64) ([]aggregates.Transaction, error) {
	var dbTransactions dbTransactions
	err := r.db.SelectContext(ctx, &dbTransactions, selectTransactionByUpdatedOn, updatedOn)
	if err != nil {
		return nil, fmt.Errorf("r.db.SelectContext, err: %w", err)
	}

	transactions := make([]aggregates.Transaction, len(dbTransactions))
	for i, dbTransaction := range dbTransactions {
		transactions[i] = dbTransaction.toAggregate()
	}

	return transactions, nil
}

func (r *Repository) InsertTransaction(ctx context.Context, e aggregates.Transaction) error {
	dbTransaction := dbTransactionFromAggregate(e)

	_, err := sqlx.NamedExec(r.db, insertTransaction, dbTransaction)
	if err != nil {
		return fmt.Errorf("sqlx.NamedExec, err: %w", err)
	}

	return nil
}

type dbTransaction struct {
	Slot         int64     `db:"slot"`
	Signature    []byte    `db:"signature"`
	IsVote       bool      `db:"is_vote"`
	MessageType  int16     `db:"message_type"`
	Signatures   [][]byte  `db:"signatures"`
	MessageHash  []byte    `db:"message_hash"`
	WriteVersion int64     `db:"write_version"`
	UpdatedOn    time.Time `db:"updated_on"`
	Index        int64     `db:"index"`
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
		WriteVersion: dbe.WriteVersion,
		UpdatedOn:    dbe.UpdatedOn,
		Index:        dbe.Index,
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
		WriteVersion: e.WriteVersion,
		UpdatedOn:    e.UpdatedOn,
		Index:        e.Index,
	}
}
