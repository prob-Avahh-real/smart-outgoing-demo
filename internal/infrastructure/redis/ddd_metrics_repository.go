package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// DDDRedisMetricsRepository implements MetricsRepository using Redis
type DDDRedisMetricsRepository struct {
	client RedisClientInterface
	prefix string
}

// NewDDDRedisMetricsRepository creates a new Redis metrics repository
func NewDDDRedisMetricsRepository(client RedisClientInterface) domain.MetricsRepository {
	return &DDDRedisMetricsRepository{
		client: client,
		prefix: "metrics:",
	}
}

// Save saves a metrics aggregate to Redis
func (r *DDDRedisMetricsRepository) Save(ctx context.Context, metrics *domain.Metrics) error {
	if metrics == nil {
		return fmt.Errorf("metrics cannot be nil")
	}

	key := r.prefix + metrics.ID()
	
	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	err = r.client.Set(ctx, key, string(data), 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to save metrics to Redis: %w", err)
	}

	return nil
}

// FindLatest finds the latest metrics in Redis
func (r *DDDRedisMetricsRepository) FindLatest(ctx context.Context) (*domain.Metrics, error) {
	pattern := r.prefix + "*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics keys: %w", err)
	}

	var latest *domain.Metrics
	for _, key := range keys {
		var data string
		err := r.client.Get(ctx, key, &data)
		if err != nil {
			continue
		}

		var metrics domain.Metrics
		err = json.Unmarshal([]byte(data), &metrics)
		if err != nil {
			continue
		}

		if latest == nil || metrics.Timestamp().After(latest.Timestamp()) {
			latest = &metrics
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("no metrics found")
	}

	return latest, nil
}

// FindByTimeRange finds metrics within a time range in Redis
func (r *DDDRedisMetricsRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*domain.Metrics, error) {
	pattern := r.prefix + "*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics keys: %w", err)
	}

	var results []*domain.Metrics
	for _, key := range keys {
		var data string
		err := r.client.Get(ctx, key, &data)
		if err != nil {
			continue
		}

		var metrics domain.Metrics
		err = json.Unmarshal([]byte(data), &metrics)
		if err != nil {
			continue
		}

		timestamp := metrics.Timestamp()
		if timestamp.After(start) && timestamp.Before(end) {
			results = append(results, &metrics)
		}
	}

	return results, nil
}

// Delete deletes metrics by ID from Redis
func (r *DDDRedisMetricsRepository) Delete(ctx context.Context, id string) error {
	key := r.prefix + id
	
	err := r.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete metrics from Redis: %w", err)
	}

	return nil
}
