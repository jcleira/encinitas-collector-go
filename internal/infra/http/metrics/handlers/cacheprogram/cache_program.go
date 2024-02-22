package cacheprogram

import (
	"sync"
	"time"

	"github.com/jcleira/encinitas-collector-go/internal/app/metrics/aggregates"
)

type CacheItem struct {
	Performance aggregates.PerformanceResults
	Throughput  aggregates.ThroughputResults
	LastUpdated time.Time
}

type MetricsCache struct {
	mu       sync.RWMutex
	cache    map[string]CacheItem
	lifetime time.Duration
}

func NewMetricsCache(lifetime time.Duration) *MetricsCache {
	return &MetricsCache{
		cache:    make(map[string]CacheItem),
		lifetime: lifetime,
	}
}

func (m *MetricsCache) Get(key string) (CacheItem, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, found := m.cache[key]
	if !found || time.Since(item.LastUpdated) > m.lifetime {
		return CacheItem{}, false
	}
	return item, true
}

func (m *MetricsCache) Set(key string, item CacheItem) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache[key] = item
}
