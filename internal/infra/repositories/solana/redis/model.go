package redis

import (
	"time"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
)

type redisTransaction struct {
	Slot            int64      `json:"slot"`
	Signature       string     `json:"signature"`
	IsVote          bool       `json:"is_vote"`
	MessageType     int        `json:"message_type"`
	LegacyMessage   string     `json:"legacy_message"`
	V0LoadedMessage string     `json:"v0_loaded_message"`
	Signatures      string     `json:"signatures"`
	MessageHash     []byte     `json:"message_hash"`
	Meta            string     `json:"meta"`
	WriteVersion    int64      `json:"write_version"`
	UpdatedOn       time.Time  `json:"updated_on"`
	TxnIndex        int64      `json:"txn_index"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
}

func (r redisTransaction) toAggregate() aggregates.Transaction {
	return aggregates.Transaction{
		Slot:            r.Slot,
		Signature:       r.Signature,
		IsVote:          r.IsVote,
		MessageType:     r.MessageType,
		LegacyMessage:   r.LegacyMessage,
		V0LoadedMessage: r.V0LoadedMessage,
		Signatures:      r.Signatures,
		MessageHash:     r.MessageHash,
		Meta:            r.Meta,
		WriteVersion:    r.WriteVersion,
		UpdatedOn:       r.UpdatedOn,
		TxnIndex:        r.TxnIndex,
		ProcessedAt:     r.ProcessedAt,
	}
}

func redisTransactionFromAggregate(
	transaction aggregates.Transaction) redisTransaction {
	return redisTransaction{
		Slot:            transaction.Slot,
		Signature:       transaction.Signature,
		IsVote:          transaction.IsVote,
		MessageType:     transaction.MessageType,
		LegacyMessage:   transaction.LegacyMessage,
		V0LoadedMessage: transaction.V0LoadedMessage,
		Signatures:      transaction.Signatures,
		MessageHash:     transaction.MessageHash,
		Meta:            transaction.Meta,
		WriteVersion:    transaction.WriteVersion,
		UpdatedOn:       transaction.UpdatedOn,
		TxnIndex:        transaction.TxnIndex,
		ProcessedAt:     transaction.ProcessedAt,
	}
}
