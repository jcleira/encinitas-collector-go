package influxdb

import (
	"context"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"

	"github.com/jcleira/encinitas-collector-go/internal/app/collector/aggregates"
)

const (
	organization = "encinitas"
	bucket       = "encinitas_metrics"
)

// Repository define the dependencies needed to store events in InfluxDB.
type Repository struct {
	client       influxdb.Client
	organization string
	bucket       string
}

// New creates a new instance of the InfluxDB repository.
func New(client influxdb.Client) *Repository {
	return &Repository{
		client:       client,
		organization: organization,
		bucket:       bucket,
	}
}

// Close closes the connection to the InfluxDB server.
func (r *Repository) Close() {
	r.client.Close()
}

// WriteMetric writes a metric event to the InfluxDB server.
func (r *Repository) WriteMetric(ctx context.Context,
	event aggregates.Metric) error {
	writeAPI := r.client.WriteAPI(r.organization, r.bucket)

	p := influxdb.NewPointWithMeasurement("events").
		AddTag("event_ID", event.EventID).
		AddTag("signature", event.Signature).
		AddField("rpc_time", event.RPCTime).
		AddField("solana_time", event.SolanaTime).
		SetTime(time.Now())

	writeAPI.WritePoint(p)

	writeAPI.Flush()

	return nil
}
