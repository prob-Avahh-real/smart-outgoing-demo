package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// RedisMetricsRepository implements MetricsRepository using Redis storage
type RedisMetricsRepository struct {
	client RedisClientInterface
}

// NewRedisMetricsRepository creates a new Redis metrics repository
func NewRedisMetricsRepository(client RedisClientInterface) repositories.MetricsRepository {
	return &RedisMetricsRepository{
		client: client,
	}
}

// Save saves a metrics aggregate
func (r *RedisMetricsRepository) Save(metrics *entities.Metrics) error {
	if metrics == nil {
		return errors.New("metrics cannot be nil")
	}
	ctx := context.Background()
	key := fmt.Sprintf("metrics:%s", metrics.ID)
	return r.client.Set(ctx, key, metrics, 24*time.Hour)
}

// FindLatest finds the latest metrics
func (r *RedisMetricsRepository) FindLatest() (*entities.Metrics, error) {
	ctx := context.Background()
	pattern := "metrics:*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics keys: %w", err)
	}
	
	if len(keys) == 0 {
		return nil, errors.New("no metrics found")
	}
	
	// Find the latest metrics by timestamp
	var latestMetrics *entities.Metrics
	var latestTime int64 = 0
	
	for _, key := range keys {
		var metrics entities.Metrics
		err := r.client.Get(ctx, key, &metrics)
		if err != nil {
			continue
		}
		
		if metrics.Timestamp.Unix() > latestTime {
			latestMetrics = &metrics
			latestTime = metrics.Timestamp.Unix()
		}
	}
	
	if latestMetrics == nil {
		return nil, errors.New("no valid metrics found")
	}
	
	return latestMetrics, nil
}

// FindByID finds metrics by ID
func (r *RedisMetricsRepository) FindByID(id string) (*entities.Metrics, error) {
	if id == "" {
		return nil, errors.New("metrics ID cannot be empty")
	}
	ctx := context.Background()
	key := fmt.Sprintf("metrics:%s", id)
	
	var metrics entities.Metrics
	err := r.client.Get(ctx, key, &metrics)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, errors.New("metrics not found")
		}
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	
	return &metrics, nil
}

// FindByTimeRange finds metrics within a time range
func (r *RedisMetricsRepository) FindByTimeRange(start, end time.Time) ([]*entities.Metrics, error) {
	ctx := context.Background()
	pattern := "metrics:*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics keys: %w", err)
	}
	
	var result []*entities.Metrics
	for _, key := range keys {
		var metrics entities.Metrics
		err := r.client.Get(ctx, key, &metrics)
		if err != nil {
			continue
		}
		
		if metrics.Timestamp.After(start) && metrics.Timestamp.Before(end) {
			result = append(result, &metrics)
		}
	}
	
	return result, nil
}

// DeleteOlderThan deletes metrics older than the specified time
func (r *RedisMetricsRepository) DeleteOlderThan(cutoff time.Time) error {
	ctx := context.Background()
	pattern := "metrics:*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return fmt.Errorf("failed to get metrics keys: %w", err)
	}
	
	var keysToDelete []string
	for _, key := range keys {
		var metrics entities.Metrics
		err := r.client.Get(ctx, key, &metrics)
		if err != nil {
			continue
		}
		
		if metrics.Timestamp.Before(cutoff) {
			keysToDelete = append(keysToDelete, key)
		}
	}
	
	if len(keysToDelete) == 0 {
		return errors.New("no metrics to delete")
	}
	
	return r.client.Delete(ctx, keysToDelete...)
}

// Count returns the total number of metrics records
func (r *RedisMetricsRepository) Count() (int, error) {
	ctx := context.Background()
	pattern := "metrics:*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return 0, fmt.Errorf("failed to count metrics: %w", err)
	}
	
	return len(keys), nil
}
