package memory

import (
	"errors"
	"sync"
	"time"
	
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// MemoryScalingDecisionRepository implements ScalingDecisionRepository using in-memory storage
type MemoryScalingDecisionRepository struct {
	decisions map[string]*entities.ScalingDecision
	mu        sync.RWMutex
}

// NewMemoryScalingDecisionRepository creates a new memory scaling decision repository
func NewMemoryScalingDecisionRepository() repositories.ScalingDecisionRepository {
	return &MemoryScalingDecisionRepository{
		decisions: make(map[string]*entities.ScalingDecision),
	}
}

// Save saves a scaling decision aggregate
func (r *MemoryScalingDecisionRepository) Save(decision *entities.ScalingDecision) error {
	if decision == nil {
		return errors.New("scaling decision cannot be nil")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.decisions[decision.ID] = decision
	return nil
}

// FindLatest finds the latest scaling decision
func (r *MemoryScalingDecisionRepository) FindLatest() (*entities.ScalingDecision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if len(r.decisions) == 0 {
		return nil, errors.New("no scaling decisions found")
	}
	
	var latestDecision *entities.ScalingDecision
	var latestTime int64 = 0
	
	for _, decision := range r.decisions {
		if decision.DecisionAt.Unix() > latestTime {
			latestDecision = decision
			latestTime = decision.DecisionAt.Unix()
		}
	}
	
	return latestDecision, nil
}

// FindByID finds a scaling decision by ID
func (r *MemoryScalingDecisionRepository) FindByID(id string) (*entities.ScalingDecision, error) {
	if id == "" {
		return nil, errors.New("scaling decision ID cannot be empty")
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	decision, exists := r.decisions[id]
	if !exists {
		return nil, errors.New("scaling decision not found")
	}
	
	return decision, nil
}

// FindActive finds the currently active scaling decision
func (r *MemoryScalingDecisionRepository) FindActive() (*entities.ScalingDecision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if len(r.decisions) == 0 {
		return nil, errors.New("no scaling decisions found")
	}
	
	// Return the latest decision as the active one
	return r.FindLatest()
}

// DeleteOlderThan deletes scaling decisions older than the specified time
func (r *MemoryScalingDecisionRepository) DeleteOlderThan(cutoff time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	deleted := 0
	for id, decision := range r.decisions {
		if decision.DecisionAt.Before(cutoff) {
			delete(r.decisions, id)
			deleted++
		}
	}
	
	if deleted == 0 {
		return errors.New("no scaling decisions to delete")
	}
	
	return nil
}

// Count returns the total number of scaling decisions
func (r *MemoryScalingDecisionRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return len(r.decisions), nil
}
