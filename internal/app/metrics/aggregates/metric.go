package aggregates

import "time"

// Metric represents a metric event which aggregates information coming from
// both the Solana blockchain and agents (browser/mobile).
type Metric struct {
	ProgramID  string
	EventID    string
	Signature  string
	RPCTime    int64
	SolanaTime int64
	Error      bool
}

// Type represents the type of metric.
type Type string

const (
	// TypeRPCTime represents the type of metric for the RPC time.
	TypeRPCTime Type = "rpc_time"

	// TypeSolanaTime represents the type of metric for the Solana time.
	TypeSolanaTime Type = "solana_time"
)

// PerformanceResult represents a metric result.
type PerformanceResult struct {
	Time  time.Time
	Type  Type
	Value float64
}

// PerformanceResults represents a slice of PerformanceResult.
type PerformanceResults []PerformanceResult

// ThroughputResult represents a throughput result.
type ThroughputResult struct {
	Time  time.Time
	Value int64
}

// ThroughputResults represents a slice of ThroughputResult.
type ThroughputResults []ThroughputResult

type ApdexMetric struct {
	Time  time.Time
	Value int64
}

// ApdexResult represents an Apdex result.
type ApdexResult struct {
	Time  time.Time
	Value float64
}

// ApdexResults represents a slice of ApdexResult.
type ApdexResults []ApdexResult

// ErrorResult represents an Error result.
type ErrorResult struct {
	Time        time.Time
	TotalErrors int64
	TotalCount  int64
	Value       float64
}

// ErrorResults represents a slice of ErrorResult.
type ErrorResults []ErrorResult
