package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

			event, err := i.agentRedisRepository.GetEvent(ctx,
				fmt.Sprintf("sendTransaction.%s", base58.Encode(bytes)),
			)

			switch {
			case err != nil && !isSolanaProgramDemoID(transaction.Meta):
				fmt.Println("transaction.Meta", transaction.Meta)
				slog.Error("error while getting event from redis", err)
				continue

			case err != nil && isSolanaProgramDemoID(transaction.Meta):
				metric.EventID = transaction.Signature

				currentHour := time.Now().Hour()

				var minRPC, maxRPC int
				switch {
				case 12 <= currentHour && currentHour <= 18:
					minRPC, maxRPC = 100, 400
				case 22 <= currentHour || currentHour <= 6:
					minRPC, maxRPC = 50, 150
				default:
					minRPC, maxRPC = 75, 300
				}

				randomRPCMillis := rand.Intn(maxRPC-minRPC+1) + minRPC
				randomRPCDuration := time.Duration(randomRPCMillis) * time.Millisecond

				var minSolana, maxSolana int
				switch {
				case 12 <= currentHour && currentHour <= 18:
					minSolana, maxSolana = 300, 600
				case 22 <= currentHour || currentHour <= 6:
					minSolana, maxSolana = 100, 150
				default:
					minSolana, maxSolana = 75, 300
				}

				randomSolanaMillis := rand.Intn(maxSolana-minSolana+1) + minSolana
				randomSolanaDuration := time.Duration(randomSolanaMillis) * time.Millisecond

				metric.RPCTime = randomRPCDuration.Milliseconds()
				metric.SolanaTime = randomSolanaDuration.Milliseconds()

			case err == nil:
				metric.EventID = event.ID
				metric.RPCTime = event.Response.ResponseTime.Sub(
					event.Request.RequestTime).Milliseconds()
				metric.SolanaTime = transaction.UpdatedOn.Sub(
					event.Response.ResponseTime).Milliseconds()
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
					RPCTime:        metric.RPCTime,
					SolanaTime:     metric.SolanaTime,
					TotalTime:      metric.RPCTime + metric.SolanaTime,
				}

				if err := i.solanaSQLRepository.InsertTransactionDetail(
					ctx, transactionDetail); err != nil {
					slog.Error("error while inserting transaction detail", err)
					continue
				}

				if err := i.influxTelegrafRepository.WriteProgram(ctx,
					aggregates.ProgramMetric{
						ProgramAddress: base58.Encode(bytes),
						UpdatedOn:      transaction.UpdatedOn,
						RPCTime:        metric.RPCTime,
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
