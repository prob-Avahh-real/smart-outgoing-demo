package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// DDDMemoryScalingDecisionRepository implements ScalingDecisionRepository using in-memory storage
type DDDMemoryScalingDecisionRepository struct {
	decisions map[string]*domain.ScalingDecision
	mu        sync.RWMutex
}

// NewDDDMemoryScalingDecisionRepository creates a new memory scaling decision repository
func NewDDDMemoryScalingDecisionRepository() domain.ScalingDecisionRepository {
	return &DDDMemoryScalingDecisionRepository{
		decisions: make(map[string]*domain.ScalingDecision),
	}
}

// Save saves a scaling decision aggregate
func (r *DDDMemoryScalingDecisionRepository) Save(ctx context.Context, decision *domain.ScalingDecision) error {
	if decision == nil {
		return errors.New("decision cannot be nil")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.decisions[decision.ID()] = decision
	return nil
}

// FindByID finds a scaling decision by ID
func (r *DDDMemoryScalingDecisionRepository) FindByID(ctx context.Context, id string) (*domain.ScalingDecision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	decision, exists := r.decisions[id]
	if !exists {
		return nil, errors.New("decision not found")
	}
	
	return decision, nil
}

// FindByTimeRange finds scaling decisions within a time range
func (r *DDDMemoryScalingDecisionRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*domain.ScalingDecision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var results []*domain.ScalingDecision
	for _, decision := range r.decisions {
		decisionTime := decision.DecisionTime()
		if decisionTime.After(start) && decisionTime.Before(end) {
			results = append(results, decision)
		}
	}
	
	// Sort by decision time
	sort.Slice(results, func(i, j int) bool {
		return results[i].DecisionTime().Before(results[j].DecisionTime())
	})
	
	return results, nil
}

// FindLatest finds the latest scaling decision
func (r *DDDMemoryScalingDecisionRepository) FindLatest(ctx context.Context) (*domain.ScalingDecision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if len(r.decisions) == 0 {
		return nil, errors.New("no decisions found")
	}
	
	var latest *domain.ScalingDecision
	for _, decision := range r.decisions {
		if latest == nil || decision.DecisionTime().After(latest.DecisionTime()) {
			latest = decision
		}
	}
	
	return latest, nil
}

// Delete deletes a scaling decision by ID
func (r *DDDMemoryScalingDecisionRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.decisions[id]; !exists {
		return errors.New("decision not found")
	}
	
	delete(r.decisions, id)
	return nil
}
