package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"

	"github.com/btcsuite/btcutil/base58"

	agentAggregates "github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
	aggregates "github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
	solanaAggregates "github.com/jcleira/encinitas-collector-go/internal/app/solana/aggregates"
)

type solanaRedisRepository interface {
	SubscribeToTransactions(
		context.Context) (chan solanaAggregates.Transaction, chan error)
}

type agentRedisRepository interface {
	GetEvent(context.Context, string) (agentAggregates.Event, error)
}

type influxDBRepository interface {
	WriteMetric(context.Context, aggregates.Metric) error
}

// Ingester is a service that ingests information coming from both the Solana
// blockchain and agents events (browser/mobile).
type Ingester struct {
	solanaRedisRepository solanaRedisRepository
	agentRedisRepository  agentRedisRepository
	influxDBRepository    influxDBRepository
}

// NewIngester creates a new instance of the Ingester service.
func NewIngester(
	solanaRedisRepository solanaRedisRepository,
	agentRedisRepository agentRedisRepository,
	influxDBRepository influxDBRepository,
) *Ingester {
	return &Ingester{
		solanaRedisRepository: solanaRedisRepository,
		agentRedisRepository:  agentRedisRepository,
		influxDBRepository:    influxDBRepository,
	}
}

// Ingest starts the ingestion process.
func (i *Ingester) Ingest(ctx context.Context) {
	transactions, errors := i.solanaRedisRepository.SubscribeToTransactions(ctx)

	for {
		select {
		case transaction := <-transactions:
			bytes, err := hex.DecodeString(transaction.Signature[2:])
			if err != nil {
				slog.Error("error while decoding transaction signature", err)
				continue
			}

			event, err := i.agentRedisRepository.GetEvent(ctx,
				fmt.Sprintf("sendTransaction.%s", base58.Encode(bytes)),
			)
			if err != nil {
				slog.Error("error while getting event from redis", err)
				continue
			}

			metric := aggregates.Metric{
				EventID:   event.ID,
				Signature: transaction.Signature,
				RPCTime: event.Response.ResponseTime.Sub(
					event.Request.RequestTime).Milliseconds(),
				SolanaTime: transaction.UpdatedOn.Sub(
					event.Response.ResponseTime).Milliseconds(),
			}

			if err = i.influxDBRepository.WriteMetric(ctx, metric); err != nil {
				slog.Error("error while writing metric to influxdb", err)
				continue
			}

		case err := <-errors:
			slog.Error("error while ingesting a transaction", err)

		case <-ctx.Done():
			return
		}
	}
}
