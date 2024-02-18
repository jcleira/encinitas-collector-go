package influx

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"

	"github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
)

const (
	apdexSatisfactory = 500
	apdexTolerable    = 1500
)

// WriteTransaction writes a transaction metric to the InfluxDB server.
func (r *Repository) WriteTransaction(ctx context.Context,
	metric aggregates.TransactionMetric) {
	r.influxdbWriter.WritePoint(
		influxdb.NewPointWithMeasurement("events").
			AddTag("event_ID", metric.EventID).
			AddTag("signature", metric.Signature).
			AddField("rpc_time", metric.RPCTime).
			AddField("solana_time", metric.SolanaTime).
			AddField("error", metric.Error).
			SetTime(metric.UpdatedOn),
	)

	r.influxdbWriter.Flush()
}

// WriteProgram writes a program transaction metric to the InfluxDB server.
func (r *Repository) WriteProgram(ctx context.Context,
	metric aggregates.ProgramMetric) {
	r.influxdbWriter.WritePoint(
		influxdb.NewPointWithMeasurement(metric.ProgramAddress).
			AddTag("program_address", metric.ProgramAddress).
			AddField("rpc_time", metric.RPCTime).
			AddField("solana_time", metric.SolanaTime).
			SetTime(metric.UpdatedOn),
	)

	r.influxdbWriter.Flush()
}

// QueryPerformance queries the InfluxDB server for performance metrics.
//
// Currently is agreggating the metrics by the mean of a fixed value which is
// not ideal; but it's a good starting point.
func (r *Repository) QueryPerformance(ctx context.Context) (aggregates.PerformanceResults, error) {
	result, err := r.client.QueryAPI(organization).Query(ctx,
		fmt.Sprintf(
			`from(bucket:"%s")
    |> range(start: -2d)
    |> filter(fn: (r) => r._measurement == "events")
    |> filter(fn: (r) => r._field == "rpc_time"
			or r._field == "solana_time")
		|> group(columns: ["_field"])
    |> aggregateWindow(every: 30m, fn: mean)`, r.bucket),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	performanceResults := make([]aggregates.PerformanceResult, 0)

	for result.Next() {
		performanceResult := aggregates.PerformanceResult{
			Time: result.Record().Time(),
		}

		switch result.Record().Field() {
		case "rpc_time":
			performanceResult.Type = aggregates.TypeRPCTime
		case "solana_time":
			performanceResult.Type = aggregates.TypeSolanaTime
		}

		if result.Record().Value() != nil {
			if value, ok := result.Record().Value().(float64); ok {
				performanceResult.Value = value
			} else {
				slog.Error("result.Record().Value() is not an float64", result.Record().Value())
			}
		}

		performanceResults = append(performanceResults, performanceResult)
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("result.Err: %w", result.Err())
	}

	return performanceResults, nil
}

// QueryProgramPerformance queries the InfluxDB server for performance metrics.
//
// Currently is agreggating the metrics by the mean of a fixed value which is
// not ideal; but it's a good starting point.
func (r *Repository) QueryProgramPerformance(
	ctx context.Context, program string) (aggregates.PerformanceResults, error) {
	result, err := r.client.QueryAPI(organization).Query(ctx,
		fmt.Sprintf(`
			from(bucket:"%s")
    |> range(start: -2d)
    |> filter(fn: (r) => r._measurement == "%s")
    |> filter(fn: (r) => r._field == "rpc_time"
			or r._field == "solana_time")
		|> group(columns: ["_field"])
    |> aggregateWindow(every: 30m, fn: mean)`, r.bucket, program),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	performanceResults := make([]aggregates.PerformanceResult, 0)

	for result.Next() {
		performanceResult := aggregates.PerformanceResult{
			Time: result.Record().Time(),
		}

		switch result.Record().Field() {
		case "rpc_time":
			performanceResult.Type = aggregates.TypeRPCTime
		case "solana_time":
			performanceResult.Type = aggregates.TypeSolanaTime
		}

		if result.Record().Value() != nil {
			if value, ok := result.Record().Value().(float64); ok {
				performanceResult.Value = value
			} else {
				slog.Error("result.Record().Value() is not an float64", result.Record().Value())
			}
		}

		performanceResults = append(performanceResults, performanceResult)
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("result.Err: %w", result.Err())
	}

	return performanceResults, nil
}

// QueryThroughput queries the InfluxDB server for performance metrics.
//
// Currently is agreggating the metrics by the mean of a fixed value which is
// not ideal; but it's a good starting point.
func (r *Repository) QueryThroughput(ctx context.Context) (aggregates.ThroughputResults, error) {
	result, err := r.client.QueryAPI(organization).Query(ctx,
		fmt.Sprintf(
			`from(bucket:"%s")
    |> range(start: -2d)
    |> filter(fn: (r) => r._measurement == "events")
    |> filter(fn: (r) => r._field == "rpc_time")
		|> group()
    |> aggregateWindow(every: 30m, fn: count)`, r.bucket),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	throughputResults := make([]aggregates.ThroughputResult, 0)

	for result.Next() {
		throughputResult := aggregates.ThroughputResult{
			Time: result.Record().Time(),
		}

		if result.Record().Value() != nil {
			if value, ok := result.Record().Value().(int64); ok {
				throughputResult.Value = value
			} else {
				slog.Error("result.Record().Value() is not a int64", result.Record().Value())
			}
		}

		throughputResults = append(throughputResults, throughputResult)
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("result.Err: %w", result.Err())
	}

	return throughputResults, nil
}

// QueryProgramThroughPut queries the InfluxDB server for performance metrics.
//
// Currently is agreggating the metrics by the mean of a fixed value which is
// not ideal; but it's a good starting point.
func (r *Repository) QueryProgramThroughput(
	ctx context.Context, programAddress string) (aggregates.ThroughputResults, error) {
	result, err := r.client.QueryAPI(organization).Query(ctx,
		fmt.Sprintf(
			`from(bucket:"%s")
    |> range(start: -2d)
    |> filter(fn: (r) => r._measurement == "%s")
		|> group()
    |> aggregateWindow(every: 30m, fn: count)`, r.bucket, programAddress),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	throughputResults := make([]aggregates.ThroughputResult, 0)

	for result.Next() {
		throughputResult := aggregates.ThroughputResult{
			Time: result.Record().Time(),
		}

		if result.Record().Value() != nil {
			if value, ok := result.Record().Value().(int64); ok {
				throughputResult.Value = value
			} else {
				slog.Error("result.Record().Value() is not a int64", result.Record().Value())
			}
		}

		throughputResults = append(throughputResults, throughputResult)
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("result.Err: %w", result.Err())
	}

	return throughputResults, nil
}

// QueryApdex queries the InfluxDB server for performance metrics, then it doesn
// some calculations to get the Apdex score.
//
// Currently is agreggating the metrics by the mean of a fixed value which is
// not ideal; but it's a good starting point.
func (r *Repository) QueryApdex(ctx context.Context) (aggregates.ApdexResults, error) {
	influxMetrics, err := r.client.QueryAPI(organization).Query(ctx,
		fmt.Sprintf(`
			from(bucket:"%s")
		|> range(start: -2d)
		|> filter(fn: (r) => r._measurement == "events")
		|> filter(fn: (r) => r._field == "rpc_time" or r._field == "solana_time")
		|> group(columns: ["_time"])
		|> aggregateWindow(every: 30m, fn: sum, createEmpty: false)`, r.bucket),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	apdexMetricMap := make(map[string][]aggregates.ApdexMetric)

	startTime := time.Now().Add(-48 * time.Hour).Truncate(30 * time.Minute)
	for t := startTime; t.Before(time.Now()); t = t.Add(30 * time.Minute) {
		timeStr := t.Format(time.RFC3339)
		apdexMetricMap[timeStr] = []aggregates.ApdexMetric{}
	}

	for influxMetrics.Next() {
		influxMetric := influxMetrics.Record()

		apdexMetric := aggregates.ApdexMetric{
			Time: influxMetric.Time(),
		}

		if influxMetric.Value() != nil {
			if value, ok := influxMetric.Value().(int64); ok {
				apdexMetric.Value = value
			} else {
				slog.Error("result.Record().Value() is not a int64", influxMetric.Value())
			}
		}

		apdexValueTimeStr := apdexMetric.Time.Format(time.RFC3339)
		apdexValues, ok := apdexMetricMap[apdexValueTimeStr]
		if !ok {
			slog.Error(
				"result.Time not found in apdexMetricMap",
				slog.Time("apdexMetric", apdexMetric.Time))
			continue
		}

		apdexMetricMap[apdexValueTimeStr] = append(apdexValues, apdexMetric)
	}

	if influxMetrics.Err() != nil {
		return nil, fmt.Errorf("result.Err: %w", influxMetrics.Err())
	}

	apdexResults := make([]aggregates.ApdexResult, 0)
	for apdexMetricsTimeStr, apdexMetrics := range apdexMetricMap {
		apdexMetricsTime, err := time.Parse(time.RFC3339, apdexMetricsTimeStr)
		if err != nil {
			// This should never happen, but if it does, we should log it and
			// continue to the next iteration.
			slog.Error("time.Parse", err)
			continue
		}

		if len(apdexMetrics) == 0 {
			apdexResults = append(apdexResults, aggregates.ApdexResult{
				Time:  apdexMetricsTime,
				Value: 1,
			})

			continue
		}

		var (
			satisfactory = 0
			tolerating   = 0
		)

		for _, apdexMetric := range apdexMetrics {
			if apdexMetric.Value <= apdexSatisfactory {
				satisfactory++
			} else if apdexMetric.Value <= apdexTolerable {
				tolerating++
			}
		}

		apdex := (float64(satisfactory) + float64(tolerating)/2) / float64(len(apdexMetrics))
		apdexResults = append(apdexResults, aggregates.ApdexResult{
			Time:  apdexMetricsTime,
			Value: apdex,
		})
	}

	sort.Slice(apdexResults, func(i, j int) bool {
		return apdexResults[i].Time.Before(apdexResults[j].Time)
	})

	return apdexResults, nil
}

func (r *Repository) QueryErrors(
	ctx context.Context) (aggregates.ErrorResults, error) {
	influxErrors, err := r.client.QueryAPI(organization).Query(ctx,
		fmt.Sprintf(`
			from(bucket: "%s")
			|> range(start: -2d)
			|> filter(fn: (r) => r._field == "error" and r._value == true)
			|> group(columns: ["_time"])
			|> group()
			|> aggregateWindow(every: 30m, fn: count, createEmpty: false)
			|> yield(name: "errors")`, r.bucket),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	influxTotals, err := r.client.QueryAPI(organization).Query(ctx,
		fmt.Sprintf(`
			from(bucket: "%s")
			|> range(start: -2d)
			|> filter(fn: (r) => r._measurement == "events")
			|> filter(fn: (r) => r._field == "error")
			|> group(columns: ["_time"])
			|> group()
			|> aggregateWindow(every: 30m, fn: count, createEmpty: false)
			|> yield(name: "total")`, r.bucket),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	errorMetricMap := make(map[string]aggregates.ErrorResult)

	startTime := time.Now().Add(-48 * time.Hour).Truncate(30 * time.Minute)
	for t := startTime; t.Before(time.Now()); t = t.Add(30 * time.Minute) {
		timeStr := t.Format(time.RFC3339)
		errorMetricMap[timeStr] = aggregates.ErrorResult{
			Time: t,
		}
	}

	for influxErrors.Next() {
		influxError := influxErrors.Record()
		influxErrorTimeStr := influxError.Time().Format(time.RFC3339)
		errorResult, ok := errorMetricMap[influxErrorTimeStr]
		if !ok {
			slog.Error(
				"result.Time not found in apdexMetricMap",
				slog.Time("apdexMetric", errorResult.Time))
			continue
		}

		errorResult.Time = influxError.Time()
		if influxError.Value() != nil {
			if value, ok := influxError.Value().(int64); ok {
				errorResult.TotalErrors = value
			} else {
				slog.Error("result.Record().Value() is not a int64", influxError.Value())
			}
		}

		errorMetricMap[influxErrorTimeStr] = errorResult
	}

	if influxErrors.Err() != nil {
		return nil, fmt.Errorf("result.Err: %w", influxErrors.Err())
	}

	for influxTotals.Next() {
		influxTotal := influxTotals.Record()

		influxTotalTimeStr := influxTotal.Time().Format(time.RFC3339)
		errorResult, ok := errorMetricMap[influxTotalTimeStr]
		if !ok {
			slog.Error(
				"result.Time not found in apdexMetricMap",
				slog.Time("apdexMetric", errorResult.Time))
			continue
		}

		if influxTotal.Value() != nil {
			if value, ok := influxTotal.Value().(int64); ok {
				errorResult.TotalCount = value
			} else {
				slog.Error("result.Record().Value() is not a int64", influxTotal.Value())
			}
		}

		errorResult.Value = float64(errorResult.TotalErrors) / float64(errorResult.TotalCount)
		errorMetricMap[influxTotalTimeStr] = errorResult
	}

	errorsResults := make([]aggregates.ErrorResult, 0)
	for _, errorResult := range errorMetricMap {
		errorsResults = append(errorsResults, errorResult)
	}

	sort.Slice(errorsResults, func(i, j int) bool {
		return errorsResults[i].Time.Before(errorsResults[j].Time)
	})

	return errorsResults, nil
}
