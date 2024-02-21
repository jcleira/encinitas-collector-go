package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"math/rand"
	"strings"
	"time"

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

type influxTelegrafRepository interface {
	WriteTransaction(context.Context, aggregates.TransactionMetric) error
	WriteProgram(context.Context, aggregates.ProgramMetric) error
}

type solanaSQLRepository interface {
	InsertTransactionDetail(context.Context, solanaAggregates.TransactionDetail) error
	GetBlockTimeByBlockHash(context.Context, string) (time.Time, error)
}

// Ingester is a service that ingests information coming from both the Solana
// blockchain and agents events (browser/mobile).
type Ingester struct {
	solanaRedisRepository    solanaRedisRepository
	agentRedisRepository     agentRedisRepository
	influxTelegrafRepository influxTelegrafRepository
	solanaSQLRepository      solanaSQLRepository
}

// NewIngester creates a new instance of the Ingester service.
func NewIngester(
	solanaRedisRepository solanaRedisRepository,
	agentRedisRepository agentRedisRepository,
	influxTelegrafRepository influxTelegrafRepository,
	solanaSQLRepository solanaSQLRepository,
) *Ingester {
	return &Ingester{
		solanaRedisRepository:    solanaRedisRepository,
		agentRedisRepository:     agentRedisRepository,
		influxTelegrafRepository: influxTelegrafRepository,
		solanaSQLRepository:      solanaSQLRepository,
	}
}

// Ingest starts the ingestion process.
func (i *Ingester) Ingest(ctx context.Context) {
	transactions, errors := i.solanaRedisRepository.SubscribeToTransactions(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case transaction := <-transactions:
			if rand.Intn(100) < 5 {
				continue
			}

			bytes, err := hex.DecodeString(transaction.Signature[2:])
			if err != nil {
				slog.Error("error while decoding transaction signature", err)
				continue
			}

			metric := aggregates.TransactionMetric{
				UpdatedOn: transaction.UpdatedOn,
				Signature: transaction.Signature,
				// TODO: We are setting the error rate to keep some randomness around
				// 0.18% for now,but we should be using the error rate from the
				// transactions.
				Error: rand.Float64() < 0.018,
			}

			metric.EventID = transaction.Signature

			transactionLegacyMessage := struct {
				RecentBlockhash string `json:"recent_blockhash"`
			}{}

			if json.Unmarshal(
				[]byte(transaction.LegacyMessage), &transactionLegacyMessage); err != nil {
				slog.Error("error while unmarshalling transaction legacy message", err)
				continue
			}

			if len(transactionLegacyMessage.RecentBlockhash) < 2 {
				slog.Error("transactionLegacyMessage.RecentBlockhash is empty")
				continue
			}

			bytes, err = hex.DecodeString(
				transactionLegacyMessage.RecentBlockhash[2:])
			if err != nil {
				slog.Error("error while decoding program account", err)
				continue
			}

			blockTime, err := i.solanaSQLRepository.GetBlockTimeByBlockHash(
				ctx, base58.Encode(bytes))
			if err != nil {
				slog.Error("error while getting block time by block hash", err)
				continue
			}

			metric.SolanaTime = transaction.UpdatedOn.Sub(blockTime).Milliseconds()

			if err := json.Unmarshal(
				[]byte(transaction.LegacyMessage), &transactionLegacyMessage); err != nil {
				slog.Error("error while unmarshalling transaction legacy message", err)
				continue
			}

			if err := i.influxTelegrafRepository.WriteTransaction(
				ctx, metric); err != nil {
				slog.Error("error while writing transaction metric", err)
				continue
			}

			if transaction.LegacyMessage == "" {
				slog.Error("transaction.LegacyMessage is empty")
				continue
			}

			transactionData := solanaAggregates.TransactionData{}
			if err := json.Unmarshal([]byte(transaction.LegacyMessage), &transactionData); err != nil {
				slog.Error("error while unmarshalling transaction data", err)
				continue
			}

			for _, instruction := range transactionData.Instructions {
				programAccountKey := transactionData.AccountKeys[instruction.ProgramIDIndex]
				bytes, err := hex.DecodeString(programAccountKey[2:])
				if err != nil {
					slog.Error("error while decoding program account", err)
					continue
				}

				transactionDetail := solanaAggregates.TransactionDetail{
					ProgramAddress: base58.Encode(bytes),
					UpdatedOn:      transaction.UpdatedOn,
					SolanaTime:     metric.SolanaTime,
				}

				if err := i.solanaSQLRepository.InsertTransactionDetail(
					ctx, transactionDetail); err != nil {
					slog.Error("error while inserting transaction detail", err)
					continue
				}

				if err := i.influxTelegrafRepository.WriteProgram(ctx,
					aggregates.ProgramMetric{
						ProgramAddress: base58.Encode(bytes),
						SolanaTime:     metric.SolanaTime,
					}); err != nil {
					slog.Error("error while writing program metric", err)
					continue
				}
			}

		case err := <-errors:
			slog.Error("error while ingesting a transaction", err)

		}
	}
}

func isSolanaProgramDemoID(meta string) bool {
	for _, demoID := range solanaProgramDemoIDs() {
		if strings.Contains(meta, demoID) {
			return true
		}
	}

	return false
}

// solanaProgramDemoIDs returns the list of Solana program IDs that we are
// using as a demo for encinitas, whenever we capture information for any of
// these programs, we will provide demo data.
func solanaProgramDemoIDs() []string {
	return []string{
		"8tfDNiaEyrV6Q1U4DEXrEigs9DoDtkugzFbybENEbCDz",
	}
}
