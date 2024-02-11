package influx

import (
	influxdb "github.com/influxdata/influxdb-client-go/v2"
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
