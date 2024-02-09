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
