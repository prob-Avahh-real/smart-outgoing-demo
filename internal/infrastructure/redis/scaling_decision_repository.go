package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// RedisScalingDecisionRepository implements ScalingDecisionRepository using Redis storage
type RedisScalingDecisionRepository struct {
	client RedisClientInterface
}

// NewRedisScalingDecisionRepository creates a new Redis scaling decision repository
func NewRedisScalingDecisionRepository(client RedisClientInterface) repositories.ScalingDecisionRepository {
	return &RedisScalingDecisionRepository{
		client: client,
	}
}

// Save saves a scaling decision aggregate
func (r *RedisScalingDecisionRepository) Save(decision *entities.ScalingDecision) error {
	if decision == nil {
		return errors.New("scaling decision cannot be nil")
	}
	ctx := context.Background()
	
	// Save both by ID and as "latest"
	idKey := fmt.Sprintf("decision:%s", decision.ID)
	latestKey := fmt.Sprintf("decision:latest")
	
	// Save individual decision
	err := r.client.Set(ctx, idKey, decision, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to save decision: %w", err)
	}
	
	// Save as latest
	return r.client.Set(ctx, latestKey, decision, 24*time.Hour)
}

// FindLatest finds the latest scaling decision
func (r *RedisScalingDecisionRepository) FindLatest() (*entities.ScalingDecision, error) {
	ctx := context.Background()
	key := "decision:latest"
	
	var decision entities.ScalingDecision
	err := r.client.Get(ctx, key, &decision)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, errors.New("no scaling decisions found")
		}
		return nil, fmt.Errorf("failed to get latest decision: %w", err)
	}
	
	return &decision, nil
}

// FindByID finds a scaling decision by ID
func (r *RedisScalingDecisionRepository) FindByID(id string) (*entities.ScalingDecision, error) {
	if id == "" {
		return nil, errors.New("scaling decision ID cannot be empty")
	}
	ctx := context.Background()
	key := fmt.Sprintf("decision:%s", id)
	
	var decision entities.ScalingDecision
	err := r.client.Get(ctx, key, &decision)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, errors.New("scaling decision not found")
		}
		return nil, fmt.Errorf("failed to get scaling decision: %w", err)
	}
	
	return &decision, nil
}

// FindActive finds the currently active scaling decision
func (r *RedisScalingDecisionRepository) FindActive() (*entities.ScalingDecision, error) {
	// In this implementation, the latest decision is the active one
	return r.FindLatest()
}

// DeleteOlderThan deletes scaling decisions older than the specified time
func (r *RedisScalingDecisionRepository) DeleteOlderThan(cutoff time.Time) error {
	ctx := context.Background()
	pattern := "decision:*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return fmt.Errorf("failed to get decision keys: %w", err)
	}
	
	var keysToDelete []string
	for _, key := range keys {
		// Skip the latest key
		if key == "decision:latest" {
			continue
		}
		
		var decision entities.ScalingDecision
		err := r.client.Get(ctx, key, &decision)
		if err != nil {
			continue
		}
		
		if decision.DecisionAt.Before(cutoff) {
			keysToDelete = append(keysToDelete, key)
		}
	}
	
	if len(keysToDelete) == 0 {
		return errors.New("no scaling decisions to delete")
	}
	
	return r.client.Delete(ctx, keysToDelete...)
}

// Count returns the total number of scaling decisions
func (r *RedisScalingDecisionRepository) Count() (int, error) {
	ctx := context.Background()
	pattern := "decision:*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return 0, fmt.Errorf("failed to count scaling decisions: %w", err)
	}
	
	// Exclude the latest key from count
	count := 0
	for _, key := range keys {
		if key != "decision:latest" {
			count++
		}
	}
	
	return count, nil
}
