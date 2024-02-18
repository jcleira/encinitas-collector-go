package influx

import (
	"log/slog"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	influxdbapi "github.com/influxdata/influxdb-client-go/v2/api"
)

const (
	organization       = "encinitas"
	TransactionsBucket = "encinitas_metrics"
	ProgramsBucket     = "encinitas_program_metrics"
)

// Repository define the dependencies needed to store events in InfluxDB.
type Repository struct {
	client         influxdb.Client
	influxdbWriter influxdbapi.WriteAPI
	organization   string
	bucket         string
}

// New creates a new instance of the InfluxDB repository.
func New(client influxdb.Client, bucket string) *Repository {
	influxdbWriter := client.WriteAPI(organization, bucket)
	errorsCh := influxdbWriter.Errors()
	go func() {
		for err := range errorsCh {
			slog.Error("error writing to InfluxDB: ", err)
		}
	}()

	return &Repository{
		client:         client,
		organization:   organization,
		bucket:         bucket,
		influxdbWriter: influxdbWriter,
	}
}

// Close closes the connection to the InfluxDB server.
func (r *Repository) Close() {
	r.client.Close()
}
