package redis

import (
	"time"

	"github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
)

type redisTransaction struct {
	Slot         int64     `json:"slot"`
	Signature    []byte    `json:"signature"`
	IsVote       bool      `json:"is_vote"`
	MessageType  int16     `json:"message_type"`
	Signatures   [][]byte  `json:"signatures"`
	MessageHash  []byte    `json:"message_hash"`
	WriteVersion int64     `json:"write_version"`
	UpdatedOn    time.Time `json:"updated_on"`
	Index        int64     `json:"index"`
	//Meta              TransactionStatusMeta
	//LegacyMessage     TransactionMessage
	//V0LoadedMessage   LoadedMessageV0
}

func redisTransactionFromAggregate(
	transaction aggregates.Transaction) redisTransaction {
	return redisTransaction{
		Slot:         transaction.Slot,
		Signature:    transaction.Signature,
		IsVote:       transaction.IsVote,
		MessageType:  transaction.MessageType,
		Signatures:   transaction.Signatures,
		MessageHash:  transaction.MessageHash,
		WriteVersion: transaction.WriteVersion,
		UpdatedOn:    transaction.UpdatedOn,
		Index:        transaction.Index,
	}
}
