package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
	"github.com/jcleira/encinitas-collector-go/internal/infra/http/metrics/handlers/cacheprogram"
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

var metricsCache = cacheprogram.NewMetricsCache(10 * time.Minute)

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

	if cachedItem, found := metricsCache.Get(programID); found {
		c.JSON(
			http.StatusOK,
			mapToHttpMetricsResponse(cachedItem.Performance, cachedItem.Throughput))
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

	metricsCache.Set(programID, cacheprogram.CacheItem{
		Performance: performanceMetrics,
		Throughput:  througputMetrics,
		LastUpdated: time.Now(),
	})

	c.JSON(
		http.StatusOK,
		mapToHttpMetricsResponse(performanceMetrics, througputMetrics))
}

func mapToHttpMetricsResponse(
	performance aggregates.PerformanceResults,
	throughput aggregates.ThroughputResults) interface{} {
	httpMetricsResponse := struct {
		Performance struct {
			Solana [][]interface{} `json:"solana"`
		} `json:"performance"`
		Throughput [][]interface{} `json:"throughput"`
	}{}

	for _, metric := range performance {
		httpMetricsResponse.Performance.Solana = append(
			httpMetricsResponse.Performance.Solana,
			[]interface{}{metric.Time, metric.Value})
	}

	for _, metric := range throughput {
		httpMetricsResponse.Throughput = append(
			httpMetricsResponse.Throughput,
			[]interface{}{metric.Time, metric.Value})
	}

	return httpMetricsResponse

}
