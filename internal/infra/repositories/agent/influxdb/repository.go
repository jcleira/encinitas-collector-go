package influxdb

import (
	"context"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jcleira/encinitas-collector-go/internal/app/agent/aggregates"
)

type Repository struct {
	client       influxdb.Client
	organization string
	bucket       string
}

func New(url, token, organization, bucket string) (*Repository, error) {
	client := influxdb.NewClient(url, token)
	return &Repository{
		client:       client,
		organization: organization,
		bucket:       bucket,
	}, nil
}

func (r *Repository) Close() {
	r.client.Close()
}

func (r *Repository) WriteEvent(
	ctx context.Context, event aggregates.Event) error {
	writeAPI := r.client.WriteAPI(r.organization, r.bucket)

	p := influxdb.NewPointWithMeasurement("events").
		AddTag("id", event.ID).
		AddField("browser_id", event.BrowserID).
		AddField("client_id", event.ClientID).
		SetTime(time.Now())

	writeAPI.WritePoint(p)

	writeAPI.Flush()

	return nil
}
