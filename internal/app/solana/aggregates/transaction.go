package aggregates

import "time"

type Transaction struct {
	Slot         int64
	Signature    []byte
	IsVote       bool
	MessageType  int16
	Signatures   [][]byte
	MessageHash  []byte
	WriteVersion int64
	UpdatedOn    time.Time
	Index        int64
	//Meta              TransactionStatusMeta
	//LegacyMessage     TransactionMessage
	//V0LoadedMessage   LoadedMessageV0
}
