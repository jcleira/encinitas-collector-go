package influx

import (
	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

const (
	organization       = "encinitas"
	transactionsBucket = "encinitas_metrics"
	programsBucket     = "encinitas_program_metrics"
)

// Repository define the dependencies needed to store events in InfluxDB.
type Repository struct {
	client       influxdb.Client
	organization string
}

// New creates a new instance of the InfluxDB repository.
func New(client influxdb.Client) *Repository {
	return &Repository{
		client:       client,
		organization: organization,
	}
}

// Close closes the connection to the InfluxDB server.
func (r *Repository) Close() {
	r.client.Close()
}
