package influx

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
)

const (
	apdexSatisfactory = 21000
	apdexTolerable    = 24000
)

// WriteTransaction writes a transaction metric to Telegraf using HTTP.
func (r *Repository) WriteTransaction(ctx context.Context,
	metric aggregates.TransactionMetric) error {
	// TODO I'm writing the metric with time.Now().UnixNano() as the timestamp
	// but I should be using the timestamp from the Solana block.
	data := fmt.Sprintf(
		"transactions,event_id=%s,signature=%s,error=%t solana_time=%d %d",
		metric.EventID, metric.Signature, metric.Error,
		metric.SolanaTime, time.Now().UTC().UnixNano())

	req, err := http.NewRequestWithContext(ctx,
		"POST", r.telegrafURL, bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("http.DefaultClient.Do: %d", resp.StatusCode)
	}

	return nil
}

// WriteProgram writes a program transaction metric to Telegraf using HTTP.
func (r *Repository) WriteProgram(ctx context.Context,
	metric aggregates.ProgramMetric) error {
	// TODO I'm writing the metric with time.Now().UnixNano() as the timestamp
	// but I should be using the timestamp from the Solana block.
	data := fmt.Sprintf(
		"%s,program_address=%s solana_time=%d %d",
		metric.ProgramAddress, metric.ProgramAddress,
		metric.SolanaTime, time.Now().UTC().UnixNano())

	req, err := http.NewRequestWithContext(ctx,
		"POST", r.telegrafURL, bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to write metric, status code: %d", resp.StatusCode)
	}

	return nil
}

// QueryPerformance queries the InfluxDB server for performance metrics.
//
// Currently is agreggating the metrics by the mean of a fixed value which is
// not ideal; but it's a good starting point.
func (r *Repository) QueryPerformance(ctx context.Context) (aggregates.PerformanceResults, error) {
	result, err := r.client.QueryAPI(organization).Query(ctx,
		fmt.Sprintf(
			`from(bucket:"%s")
    |> range(start: -8h)
    |> filter(fn: (r) => r._measurement == "transactions")
    |> filter(fn: (r) => r._field == "solana_time_mean")
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
		case "solana_time_mean":
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
    |> range(start: -8h)
    |> filter(fn: (r) => r._measurement == "%s")
    |> filter(fn: (r) => r._field == "solana_time_mean")
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
		case "solana_time_mean":
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

	// remove the last element from the slice, as it's always empty
	performanceResults = performanceResults[:len(performanceResults)-1]

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
    |> range(start: -8h)
    |> filter(fn: (r) => r._measurement == "transactions")
    |> filter(fn: (r) => r._field == "solana_time_count")
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

	// remove the last element from the slice, as it has not the proper date
	throughputResults = throughputResults[:len(throughputResults)-1]

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
    |> range(start: -8h)
    |> filter(fn: (r) => r._measurement == "%s")
    |> filter(fn: (r) => r._field == "solana_time_count")
		|> group()
    |> aggregateWindow(every: 30m, fn: sum)`, r.bucket, programAddress),
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
			if value, ok := result.Record().Value().(float64); ok {
				throughputResult.Value = int64(value)
			} else {
				slog.Error("result.Record().Value() is not a int64", result.Record().Value())
			}
		}

		throughputResults = append(throughputResults, throughputResult)
	}

	// remove the last element from the slice, as it has not the proper date
	throughputResults = throughputResults[:len(throughputResults)-1]

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
	satisfactoryMetrics, err := r.client.QueryAPI(organization).Query(ctx,
		// TODO: This query is using only the solana time for the apdex calculation,
		// but it should be using both the solana time.
		fmt.Sprintf(`
			from(bucket:"%s")
		|> range(start: -8h)
		|> filter(fn: (r) => r._measurement == "transactions")
		|> filter(fn: (r) => r._field == "solana_time_mean")
		|> filter(fn: (r) => r._value < %d)
		|> group()
		|> aggregateWindow(every: 30m, fn: mean, createEmpty: false)`, r.bucket, apdexSatisfactory),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	tolerableMetrics, err := r.client.QueryAPI(organization).Query(ctx,
		// TODO: This query is using only the solana time for the apdex calculation,
		// but it should be using both the solana time.
		fmt.Sprintf(`
			from(bucket:"%s")
		|> range(start: -8h)
		|> filter(fn: (r) => r._measurement == "transactions")
		|> filter(fn: (r) => r._field == "solana_time_mean")
		|> filter(fn: (r) => r._value > %d)
		|> filter(fn: (r) => r._value < %d)
		|> group()
		|> aggregateWindow(every: 30m, fn: mean, createEmpty: false)`,
			r.bucket, apdexSatisfactory, apdexTolerable),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	frustratingMetrics, err := r.client.QueryAPI(organization).Query(ctx,
		// TODO: This query is using only the solana time for the apdex calculation,
		// but it should be using both the solana time.
		fmt.Sprintf(`
			from(bucket:"%s")
		|> range(start: -8h)
		|> filter(fn: (r) => r._measurement == "transactions")
		|> filter(fn: (r) => r._field == "solana_time_mean")
		|> filter(fn: (r) => r._value > %d)
		|> group()
		|> aggregateWindow(every: 30m, fn: mean, createEmpty: false)`,
			r.bucket, apdexTolerable),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	apdexMetricMap := make(map[string]aggregates.ApdexMetric)

	startTime := time.Now().UTC().Add(-8 * time.Hour).Truncate(30 * time.Minute)
	for t := startTime; t.Before(time.Now().UTC()); t = t.Add(30 * time.Minute) {
		timeStr := t.Format(time.RFC3339)
		apdexMetricMap[timeStr] = aggregates.ApdexMetric{
			Time: t,
		}
	}

	for satisfactoryMetrics.Next() {
		influxMetric := satisfactoryMetrics.Record()

		apdexValueTimeStr := influxMetric.Time().UTC().Format(time.RFC3339)
		apdexMetric, ok := apdexMetricMap[apdexValueTimeStr]
		if !ok {
			continue
		}

		if influxMetric.Value() != nil {
			if value, ok := influxMetric.Value().(float64); ok {
				apdexMetric.SatisfactoryCount = int64(value)
			} else {
				slog.Error("result.Record().Value() is not a int64", influxMetric.Value())
			}
		}

		apdexMetricMap[apdexValueTimeStr] = apdexMetric
	}

	for tolerableMetrics.Next() {
		influxMetric := tolerableMetrics.Record()

		apdexValueTimeStr := influxMetric.Time().UTC().Format(time.RFC3339)
		apdexMetric, ok := apdexMetricMap[apdexValueTimeStr]
		if !ok {
			continue
		}

		if influxMetric.Value() != nil {
			if value, ok := influxMetric.Value().(float64); ok {
				apdexMetric.TolerableCount = int64(value)
			} else {
				slog.Error("result.Record().Value() is not a int64", influxMetric.Value())
			}
		}

		apdexMetricMap[apdexValueTimeStr] = apdexMetric
	}

	for frustratingMetrics.Next() {
		influxMetric := frustratingMetrics.Record()

		apdexValueTimeStr := influxMetric.Time().UTC().Format(time.RFC3339)
		apdexMetric, ok := apdexMetricMap[apdexValueTimeStr]
		if !ok {
			continue
		}

		if influxMetric.Value() != nil {
			if value, ok := influxMetric.Value().(float64); ok {
				apdexMetric.FrustratingCount = int64(value)
			} else {
				slog.Error("result.Record().Value() is not a int64", influxMetric.Value())
			}
		}

		apdexMetricMap[apdexValueTimeStr] = apdexMetric
	}

	apdexResults := make([]aggregates.ApdexResult, 0)
	for apdexMetricsTimeStr, apdexMetric := range apdexMetricMap {
		apdexMetricsTime, err := time.Parse(time.RFC3339, apdexMetricsTimeStr)
		if err != nil {
			// This should never happen, but if it does, we should log it and
			// continue to the next iteration.
			slog.Error("time.Parse", err)
			continue
		}

		var (
			satisfactory = apdexMetric.SatisfactoryCount
			tolerating   = apdexMetric.TolerableCount
			total        = satisfactory + tolerating + apdexMetric.FrustratingCount
		)

		if total > 0 {
			apdex := (float64(satisfactory) + float64(tolerating)/2) / float64(total)
			apdexResults = append(apdexResults, aggregates.ApdexResult{
				Time:  apdexMetricsTime,
				Value: apdex,
			})
		}
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
			|> range(start: -8h)
			|> filter(fn: (r) => r._measurement == "transactions")
			|> filter(fn: (r) => r.error == "true")
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
			|> range(start: -8h)
			|> filter(fn: (r) => r._measurement == "transactions")
			|> filter(fn: (r) => r.error == "false" or r.error == "true")
			|> group(columns: ["_time"])
			|> group()
			|> aggregateWindow(every: 30m, fn: count, createEmpty: false)
			|> yield(name: "total")`, r.bucket),
	)
	if err != nil {
		return nil, fmt.Errorf("r.client.QueryAPI(organization).Query: %w", err)
	}

	errorMetricMap := make(map[string]aggregates.ErrorResult)

	startTime := time.Now().UTC().Add(-8 * time.Hour).Truncate(30 * time.Minute)
	for t := startTime; t.Before(time.Now().UTC()); t = t.Add(30 * time.Minute) {
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
