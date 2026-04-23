package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// DDDMemoryMetricsRepository implements MetricsRepository using in-memory storage
type DDDMemoryMetricsRepository struct {
	metrics map[string]*domain.Metrics
	mu      sync.RWMutex
}

// NewDDDMemoryMetricsRepository creates a new memory metrics repository
func NewDDDMemoryMetricsRepository() domain.MetricsRepository {
	return &DDDMemoryMetricsRepository{
		metrics: make(map[string]*domain.Metrics),
	}
}

// Save saves a metrics aggregate
func (r *DDDMemoryMetricsRepository) Save(ctx context.Context, metrics *domain.Metrics) error {
	if metrics == nil {
		return errors.New("metrics cannot be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics[metrics.ID()] = metrics
	return nil
}

// FindLatest finds the latest metrics
func (r *DDDMemoryMetricsRepository) FindLatest(ctx context.Context) (*domain.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.metrics) == 0 {
		return nil, errors.New("no metrics found")
	}

	var latest *domain.Metrics
	for _, metrics := range r.metrics {
		if latest == nil || metrics.Timestamp().After(latest.Timestamp()) {
			latest = metrics
		}
	}

	return latest, nil
}

// FindByTimeRange finds metrics within a time range
func (r *DDDMemoryMetricsRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*domain.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*domain.Metrics
	for _, metrics := range r.metrics {
		timestamp := metrics.Timestamp()
		if timestamp.After(start) && timestamp.Before(end) {
			results = append(results, metrics)
		}
	}

	// Sort by timestamp
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp().Before(results[j].Timestamp())
	})

	return results, nil
}

// Delete deletes metrics by ID
func (r *DDDMemoryMetricsRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.metrics[id]; !exists {
		return errors.New("metrics not found")
	}

	delete(r.metrics, id)
	return nil
}
