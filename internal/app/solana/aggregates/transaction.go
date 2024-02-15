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
	ProcessedAt     *time.Time
}
