package repositories

import (
	"smart-outgoing-demo/internal/domain/entities"
)

// RepositoryFactory creates repository instances based on storage strategy
type RepositoryFactory interface {
	// CreateVehicleRepository creates a vehicle repository
	CreateVehicleRepository(strategy entities.StorageStrategy) VehicleRepository
	
	// CreateRouteRepository creates a route repository
	CreateRouteRepository(strategy entities.StorageStrategy) RouteRepository
	
	// CreateMetricsRepository creates a metrics repository
	CreateMetricsRepository(strategy entities.StorageStrategy) MetricsRepository
	
	// CreateScalingDecisionRepository creates a scaling decision repository
	CreateScalingDecisionRepository(strategy entities.StorageStrategy) ScalingDecisionRepository
}
