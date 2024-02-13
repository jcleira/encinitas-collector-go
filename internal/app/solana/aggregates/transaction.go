package aggregates

import "time"

// Transaction is the domain representation of a solana transaction.
type Transaction struct {
	Slot            int64
	Signature       string
	IsVote          bool
	MessageType     int
	LegacyMessage   string
	V0LoadedMessage string
	Signatures      string
	MessageHash     []byte
	Meta            string
	WriteVersion    int64
	UpdatedOn       time.Time
	TxnIndex        int64
	// I'm going to keep the error info as a pointer to a string, for the
	// moment, till we get a transaction with an error.
	ErrorInfo   *string
	ProcessedAt *time.Time
}
