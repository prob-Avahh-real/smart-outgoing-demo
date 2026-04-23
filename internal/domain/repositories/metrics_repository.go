package repositories

import (
	"smart-outgoing-demo/internal/domain/entities"
	"time"
)

// MetricsRepository defines the repository interface for Metrics aggregate
type MetricsRepository interface {
	// Save saves a metrics aggregate
	Save(metrics *entities.Metrics) error

	// FindLatest finds the latest metrics
	FindLatest() (*entities.Metrics, error)

	// FindByID finds metrics by ID
	FindByID(id string) (*entities.Metrics, error)

	// FindByTimeRange finds metrics within a time range
	FindByTimeRange(start, end time.Time) ([]*entities.Metrics, error)

	// DeleteOlderThan deletes metrics older than the specified time
	DeleteOlderThan(time.Time) error

	// Count returns the total number of metrics records
	Count() (int, error)
}
