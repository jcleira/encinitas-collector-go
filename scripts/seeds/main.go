package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/jcleira/encinitas-collector-go/config"
	"github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
	"github.com/kelseyhightower/envconfig"
)

const (
	influxOrg    = "encinitas"
	influxBucket = "encinitas_metrics"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var config config.Config

	if err := envconfig.Process("", &config); err != nil {
		slog.Error("can't process envconfig: ", err)
		os.Exit(1)
	}

	client := influxdb2.NewClient(config.InfluxDB.URL, config.InfluxDB.Token)
	defer client.Close()

	writeAPI := client.WriteAPI(influxOrg, influxBucket)

	for i := 0; i < 1000; i++ {
		metric := generateRandomMetric()

		point := influxdb2.NewPoint(
			"solana_metrics",
			map[string]string{
				"programID": metric.ProgramID,
				"eventID":   metric.EventID,
				"signature": metric.Signature,
			},
			map[string]interface{}{
				"rpcTime":    metric.RPCTime,
				"solanaTime": metric.SolanaTime,
			},
			time.Now().Add(-24*time.Hour+time.Duration(i)*time.Minute),
		)

		writeAPI.WritePoint(point)
	}

	writeAPI.Flush()
	fmt.Println("Finished writing 1000 metrics.")
}

// generateRandomMetric generates a mock Metric instance with random data.
func generateRandomMetric() aggregates.Metric {
	return aggregates.Metric{
		ProgramID:  fmt.Sprintf("Program%d", rand.Intn(100)),
		EventID:    fmt.Sprintf("Event%d", rand.Intn(100)),
		Signature:  fmt.Sprintf("Signature%d", rand.Intn(100)),
		RPCTime:    rand.Int63n(1000),
		SolanaTime: rand.Int63n(1000),
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
