package cache

import (
	"sync"
	"time"

	"github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
)

type MetricsCache struct {
	LastUpdate  time.Time
	Performance aggregates.PerformanceResults
	Throughput  aggregates.ThroughputResults
	Apdex       aggregates.ApdexResults
	Errors      aggregates.ErrorResults
	mu          sync.RWMutex // Ensure thread-safe access to the cache
}

// Global instance of the cache
var metricsCache = MetricsCache{}

// UpdateCache updates the cache with new data
func UpdateCache(
	performance aggregates.PerformanceResults,
	throughput aggregates.ThroughputResults,
	apdex aggregates.ApdexResults,
	errors aggregates.ErrorResults) {
	metricsCache.mu.Lock()
	defer metricsCache.mu.Unlock()
	metricsCache.Performance = performance
	metricsCache.Throughput = throughput
	metricsCache.Apdex = apdex
	metricsCache.Errors = errors
	metricsCache.LastUpdate = time.Now()
}

// GetCache returns the current cache
func GetCache() (
	aggregates.PerformanceResults,
	aggregates.ThroughputResults,
	aggregates.ApdexResults,
	aggregates.ErrorResults,
	bool) {
	metricsCache.mu.RLock()
	defer metricsCache.mu.RUnlock()
	if time.Since(metricsCache.LastUpdate) > 10*time.Minute {
		return aggregates.PerformanceResults{},
			aggregates.ThroughputResults{},
			aggregates.ApdexResults{},
			aggregates.ErrorResults{}, false
	}
	return metricsCache.Performance,
		metricsCache.Throughput,
		metricsCache.Apdex,
		metricsCache.Errors, true
}
