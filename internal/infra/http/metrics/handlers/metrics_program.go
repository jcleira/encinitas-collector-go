package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
)

// metricsProgramRetriever defines the methods needed to retrievetmetrics.
type metricsProgramRetriever interface {
	QueryProgramPerformance(
		context.Context, string) (aggregates.PerformanceResults, error)
	QueryProgramThroughput(
		context.Context, string) (aggregates.ThroughputResults, error)
}

// MetricsProgramRetrieverHandler defines the dependencies to retrieve metrics.
type MetricsProgramRetrieverHandler struct {
	metricsProgramRetriever metricsProgramRetriever
}

// NewMetricsProgramRetrieverHandler initializes a new MetricsProgramRetrieverHandler.
func NewMetricsProgramRetrieverHandler(
	metricsProgramRetriever metricsProgramRetriever) *MetricsProgramRetrieverHandler {
	return &MetricsProgramRetrieverHandler{
		metricsProgramRetriever: metricsProgramRetriever,
	}
}

// Handle is the handler function to retrieve metrics
func (ech *MetricsProgramRetrieverHandler) Handle(c *gin.Context) {
	programID := c.Query("program_id")
	if programID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "programID is required"})
		return
	}

	performanceMetrics, err := ech.metricsProgramRetriever.QueryProgramPerformance(
		c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	througputMetrics, err := ech.metricsProgramRetriever.QueryProgramThroughput(
		c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpMetricsResponse := struct {
		Performance struct {
			RPC    [][]interface{} `json:"rpc"`
			Solana [][]interface{} `json:"solana"`
		} `json:"performance"`
		Throughput [][]interface{} `json:"throughput"`
	}{}

	for _, metric := range performanceMetrics {
		if metric.Type == aggregates.TypeRPCTime {
			httpMetricsResponse.Performance.RPC = append(
				httpMetricsResponse.Performance.RPC,
				[]interface{}{metric.Time, metric.Value})
		} else {
			httpMetricsResponse.Performance.Solana = append(
				httpMetricsResponse.Performance.Solana,
				[]interface{}{metric.Time, metric.Value})
		}
	}

	for _, metric := range througputMetrics {
		httpMetricsResponse.Throughput = append(
			httpMetricsResponse.Throughput,
			[]interface{}{metric.Time, metric.Value})
	}

	c.JSON(http.StatusOK, httpMetricsResponse)
}
