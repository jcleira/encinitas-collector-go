package aggregates

// Metric represents a metric event which aggregates information coming from
// both the Solana blockchain and agents (browser/mobile).
type Metric struct {
	ProgramID  string
	EventID    string
	Signature  string
	RPCTime    int64
	SolanaTime int64
}
