package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// DDDRedisScalingDecisionRepository implements ScalingDecisionRepository using Redis
type DDDRedisScalingDecisionRepository struct {
	client RedisClientInterface
	prefix string
}

// NewDDDRedisScalingDecisionRepository creates a new Redis scaling decision repository
func NewDDDRedisScalingDecisionRepository(client RedisClientInterface) domain.ScalingDecisionRepository {
	return &DDDRedisScalingDecisionRepository{
		client: client,
		prefix: "scaling_decision:",
	}
}

// Save saves a scaling decision aggregate to Redis
func (r *DDDRedisScalingDecisionRepository) Save(ctx context.Context, decision *domain.ScalingDecision) error {
	if decision == nil {
		return fmt.Errorf("decision cannot be nil")
	}

	key := r.prefix + decision.ID()
	
	data, err := json.Marshal(decision)
	if err != nil {
		return fmt.Errorf("failed to marshal decision: %w", err)
	}

	err = r.client.Set(ctx, key, string(data), 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to save decision to Redis: %w", err)
	}

	return nil
}

// FindByID finds a scaling decision by ID in Redis
func (r *DDDRedisScalingDecisionRepository) FindByID(ctx context.Context, id string) (*domain.ScalingDecision, error) {
	key := r.prefix + id
	
	var data string
	err := r.client.Get(ctx, key, &data)
	if err != nil {
		return nil, fmt.Errorf("decision not found: %w", err)
	}

	var decision domain.ScalingDecision
	err = json.Unmarshal([]byte(data), &decision)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal decision: %w", err)
	}

	return &decision, nil
}

// FindByTimeRange finds scaling decisions within a time range in Redis
func (r *DDDRedisScalingDecisionRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*domain.ScalingDecision, error) {
	pattern := r.prefix + "*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get decision keys: %w", err)
	}

	var results []*domain.ScalingDecision
	for _, key := range keys {
		var data string
		err := r.client.Get(ctx, key, &data)
		if err != nil {
			continue
		}

		var decision domain.ScalingDecision
		err = json.Unmarshal([]byte(data), &decision)
		if err != nil {
			continue
		}

		decisionTime := decision.DecisionTime()
		if decisionTime.After(start) && decisionTime.Before(end) {
			results = append(results, &decision)
		}
	}

	return results, nil
}

// FindLatest finds the latest scaling decision in Redis
func (r *DDDRedisScalingDecisionRepository) FindLatest(ctx context.Context) (*domain.ScalingDecision, error) {
	pattern := r.prefix + "*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get decision keys: %w", err)
	}

	var latest *domain.ScalingDecision
	for _, key := range keys {
		var data string
		err := r.client.Get(ctx, key, &data)
		if err != nil {
			continue
		}

		var decision domain.ScalingDecision
		err = json.Unmarshal([]byte(data), &decision)
		if err != nil {
			continue
		}

		if latest == nil || decision.DecisionTime().After(latest.DecisionTime()) {
			latest = &decision
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("no decisions found")
	}

	return latest, nil
}

// Delete deletes a scaling decision by ID from Redis
func (r *DDDRedisScalingDecisionRepository) Delete(ctx context.Context, id string) error {
	key := r.prefix + id
	
	err := r.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete decision from Redis: %w", err)
	}

	return nil
}
