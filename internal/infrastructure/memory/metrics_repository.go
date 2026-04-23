package memory

import (
	"errors"
	"sync"
	"time"
	
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// MemoryMetricsRepository implements MetricsRepository using in-memory storage
type MemoryMetricsRepository struct {
	metrics []*entities.Metrics
	mu      sync.RWMutex
}

// NewMemoryMetricsRepository creates a new memory metrics repository
func NewMemoryMetricsRepository() repositories.MetricsRepository {
	return &MemoryMetricsRepository{
		metrics: make([]*entities.Metrics, 0),
	}
}

// Save saves a metrics aggregate
func (r *MemoryMetricsRepository) Save(metrics *entities.Metrics) error {
	if metrics == nil {
		return errors.New("metrics cannot be nil")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.metrics = append(r.metrics, metrics)
	return nil
}

// FindLatest finds the latest metrics
func (r *MemoryMetricsRepository) FindLatest() (*entities.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if len(r.metrics) == 0 {
		return nil, errors.New("no metrics found")
	}
	
	// Return the last metrics (most recent)
	return r.metrics[len(r.metrics)-1], nil
}

// FindByID finds metrics by ID
func (r *MemoryMetricsRepository) FindByID(id string) (*entities.Metrics, error) {
	if id == "" {
		return nil, errors.New("metrics ID cannot be empty")
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, metrics := range r.metrics {
		if metrics.ID == id {
			return metrics, nil
		}
	}
	
	return nil, errors.New("metrics not found")
}

// FindByTimeRange finds metrics within a time range
func (r *MemoryMetricsRepository) FindByTimeRange(start, end time.Time) ([]*entities.Metrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*entities.Metrics
	for _, metrics := range r.metrics {
		if metrics.Timestamp.After(start) && metrics.Timestamp.Before(end) {
			result = append(result, metrics)
		}
	}
	
	return result, nil
}

// DeleteOlderThan deletes metrics older than the specified time
func (r *MemoryMetricsRepository) DeleteOlderThan(cutoff time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	var newMetrics []*entities.Metrics
	deleted := 0
	
	for _, metrics := range r.metrics {
		if metrics.Timestamp.After(cutoff) {
			newMetrics = append(newMetrics, metrics)
		} else {
			deleted++
		}
	}
	
	r.metrics = newMetrics
	
	if deleted == 0 {
		return errors.New("no metrics to delete")
	}
	
	return nil
}

// Count returns the total number of metrics records
func (r *MemoryMetricsRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return len(r.metrics), nil
}
