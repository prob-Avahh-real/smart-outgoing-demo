package repositories

import (
	"smart-outgoing-demo/internal/domain/entities"
	"time"
)

// ScalingDecisionRepository defines the repository interface for ScalingDecision aggregate
type ScalingDecisionRepository interface {
	// Save saves a scaling decision aggregate
	Save(decision *entities.ScalingDecision) error

	// FindLatest finds the latest scaling decision
	FindLatest() (*entities.ScalingDecision, error)

	// FindByID finds a scaling decision by ID
	FindByID(id string) (*entities.ScalingDecision, error)

	// FindActive finds the currently active scaling decision
	FindActive() (*entities.ScalingDecision, error)

	// DeleteOlderThan deletes scaling decisions older than the specified time
	DeleteOlderThan(time.Time) error

	// Count returns the total number of scaling decisions
	Count() (int, error)
}
