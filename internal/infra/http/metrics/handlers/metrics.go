package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
)

// metricsRetriever defines the methods needed to retrievetmetrics.
type metricsRetriever interface {
	QueryPerformance(context.Context) (aggregates.PerformanceResults, error)
	QueryThroughput(context.Context) (aggregates.ThroughputResults, error)
	QueryApdex(context.Context) (aggregates.ApdexResults, error)
	QueryErrors(context.Context) (aggregates.ErrorResults, error)
}

// MetricsRetrieverHandler defines the dependencies to retrieve metrics.
type MetricsRetrieverHandler struct {
	metricsRetriever metricsRetriever
}

// NewMetricsRetriever initializes a new MetricsRetrieverHandler.
func NewMetricsRetriever(
	metricsRetriever metricsRetriever) *MetricsRetrieverHandler {
	return &MetricsRetrieverHandler{
		metricsRetriever: metricsRetriever,
	}
}

// Handle is the handler function to retrieve metrics
func (ech *MetricsRetrieverHandler) Handle(c *gin.Context) {
	performanceMetrics, err := ech.metricsRetriever.QueryPerformance(
		c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	througputMetrics, err := ech.metricsRetriever.QueryThroughput(
		c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apdexMetrics, err := ech.metricsRetriever.QueryApdex(
		c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	errorMetrics, err := ech.metricsRetriever.QueryErrors(
		c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	httpMetricsResponse := struct {
		Performance struct {
			Solana [][]interface{} `json:"solana"`
		} `json:"performance"`
		Throughput [][]interface{} `json:"throughput"`
		Apdex      [][]interface{} `json:"apdex"`
		Errors     [][]interface{} `json:"errors"`
	}{}

	for _, metric := range performanceMetrics {
		httpMetricsResponse.Performance.Solana = append(
			httpMetricsResponse.Performance.Solana,
			[]interface{}{metric.Time, int64(metric.Value)})
	}

	for _, metric := range througputMetrics {
		httpMetricsResponse.Throughput = append(
			httpMetricsResponse.Throughput,
			[]interface{}{metric.Time, metric.Value})
	}

	for _, metric := range apdexMetrics {
		httpMetricsResponse.Apdex = append(
			httpMetricsResponse.Apdex,
			[]interface{}{metric.Time, metric.Value})
	}

	for _, metric := range errorMetrics {
		httpMetricsResponse.Errors = append(
			httpMetricsResponse.Errors,
			[]interface{}{metric.Time, metric.Value})
	}

	c.JSON(http.StatusOK, httpMetricsResponse)
}
