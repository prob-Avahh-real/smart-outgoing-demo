package memory

import (
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// MemoryRepositoryFactory implements RepositoryFactory for memory storage
type MemoryRepositoryFactory struct{}

// NewMemoryRepositoryFactory creates a new memory repository factory
func NewMemoryRepositoryFactory() repositories.RepositoryFactory {
	return &MemoryRepositoryFactory{}
}

// CreateVehicleRepository creates a memory vehicle repository
func (f *MemoryRepositoryFactory) CreateVehicleRepository(strategy entities.StorageStrategy) repositories.VehicleRepository {
	return NewMemoryVehicleRepository()
}

// CreateRouteRepository creates a memory route repository
func (f *MemoryRepositoryFactory) CreateRouteRepository(strategy entities.StorageStrategy) repositories.RouteRepository {
	return NewMemoryRouteRepository()
}

// CreateMetricsRepository creates a memory metrics repository
func (f *MemoryRepositoryFactory) CreateMetricsRepository(strategy entities.StorageStrategy) repositories.MetricsRepository {
	return NewMemoryMetricsRepository()
}

// CreateScalingDecisionRepository creates a memory scaling decision repository
func (f *MemoryRepositoryFactory) CreateScalingDecisionRepository(strategy entities.StorageStrategy) repositories.ScalingDecisionRepository {
	return NewMemoryScalingDecisionRepository()
}
