package memory

import (
	"smart-outgoing-demo/internal/domain"
)

// DDDMemoryRepositoryFactory implements RepositoryFactory for memory storage
type DDDMemoryRepositoryFactory struct{}

// NewDDDMemoryRepositoryFactory creates a new memory repository factory
func NewDDDMemoryRepositoryFactory() domain.RepositoryFactory {
	return &DDDMemoryRepositoryFactory{}
}

// CreateVehicleRepository creates a vehicle repository using memory storage
func (f *DDDMemoryRepositoryFactory) CreateVehicleRepository(strategy domain.StorageStrategy) domain.VehicleRepository {
	if strategy == domain.StorageStrategyMemory || strategy == domain.StorageStrategyHybrid {
		return NewDDDMemoryVehicleRepository()
	}
	// For Redis strategy, we would return Redis implementation
	// For now, fall back to memory
	return NewDDDMemoryVehicleRepository()
}

// CreateRouteRepository creates a route repository using memory storage
func (f *DDDMemoryRepositoryFactory) CreateRouteRepository(strategy domain.StorageStrategy) domain.RouteRepository {
	if strategy == domain.StorageStrategyMemory || strategy == domain.StorageStrategyHybrid {
		return NewDDDMemoryRouteRepository()
	}
	return NewDDDMemoryRouteRepository()
}

// CreateMetricsRepository creates a metrics repository using memory storage
func (f *DDDMemoryRepositoryFactory) CreateMetricsRepository(strategy domain.StorageStrategy) domain.MetricsRepository {
	if strategy == domain.StorageStrategyMemory || strategy == domain.StorageStrategyHybrid {
		return NewDDDMemoryMetricsRepository()
	}
	return NewDDDMemoryMetricsRepository()
}

// CreateScalingDecisionRepository creates a scaling decision repository using memory storage
func (f *DDDMemoryRepositoryFactory) CreateScalingDecisionRepository(strategy domain.StorageStrategy) domain.ScalingDecisionRepository {
	if strategy == domain.StorageStrategyMemory || strategy == domain.StorageStrategyHybrid {
		return NewDDDMemoryScalingDecisionRepository()
	}
	return NewDDDMemoryScalingDecisionRepository()
}
