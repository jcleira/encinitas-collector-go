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

// TransactionDetail is the domain representation of a solana transaction detail.
type TransactionDetail struct {
	ProgramAddress string
	UpdatedOn      time.Time
	RPCTime        int64
	SolanaTime     int64
	TotalTime      int64
}

// TransactionDetailAggregated is the domain representation of a solana
// transaction detail aggregated.
type TransactionDetailAggregated struct {
	ProgramAddress string
	Percentage     float64
}

type TransactionMeta struct {
	Error             *string  `json:"error"`
	Fee               int      `json:"fee"`
	PreBalances       []int64  `json:"pre_balances"`
	PostBalances      []int64  `json:"post_balances"`
	InnerInstructions []string `json:"inner_instructions"`
	LogMessages       []string `json:"log_messages"`
	PreTokenBalances  []string `json:"pre_token_balances"`
	PostTokenBalances []string `json:"post_token_balances"`
	Rewards           []string `json:"rewards"`
}

type TransactionData struct {
	Header          Header        `json:"header"`
	AccountKeys     []string      `json:"account_keys"`
	RecentBlockhash string        `json:"recent_blockhash"`
	Instructions    []Instruction `json:"instructions"`
}

type Header struct {
	NumRequiredSignatures       int `json:"num_required_signatures"`
	NumReadonlySignedAccounts   int `json:"num_readonly_signed_accounts"`
	NumReadonlyUnsignedAccounts int `json:"num_readonly_unsigned_accounts"`
}

type Instruction struct {
	ProgramIDIndex int    `json:"program_id_index"`
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
}
