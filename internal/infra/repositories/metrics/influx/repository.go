package influx

import (
	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

const (
	organization       = "encinitas"
	TransactionsBucket = "encinitas_metrics"
	ProgramsBucket     = "encinitas_program_metrics"
)

// Repository define the dependencies needed to store events in InfluxDB.
type Repository struct {
	client       influxdb.Client
	telegrafURL  string
	organization string
	bucket       string
}

// New creates a new instance of the InfluxDB repository.
func New(client influxdb.Client, telegrafURL, bucket string) *Repository {
	return &Repository{
		client:       client,
		organization: organization,
		bucket:       bucket,
		telegrafURL:  telegrafURL,
	}
}

// Close closes the connection to the InfluxDB server.
func (r *Repository) Close() {
	r.client.Close()
}
