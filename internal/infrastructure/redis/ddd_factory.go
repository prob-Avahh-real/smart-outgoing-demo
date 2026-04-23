package redis

import (
	"smart-outgoing-demo/internal/domain"
)

// DDDRedisRepositoryFactory implements RepositoryFactory for Redis storage
type DDDRedisRepositoryFactory struct {
	client RedisClientInterface
}

// NewDDDRedisRepositoryFactory creates a new Redis repository factory
func NewDDDRedisRepositoryFactory(client RedisClientInterface) *DDDRedisRepositoryFactory {
	return &DDDRedisRepositoryFactory{
		client: client,
	}
}

// CreateVehicleRepository creates a vehicle repository using Redis storage
func (f *DDDRedisRepositoryFactory) CreateVehicleRepository(strategy domain.StorageStrategy) domain.VehicleRepository {
	if strategy == domain.StorageStrategyRedis || strategy == domain.StorageStrategyHybrid {
		return NewDDDRedisVehicleRepository(f.client)
	}
	// For non-Redis strategies, return nil or fallback
	return nil
}

// CreateRouteRepository creates a route repository using Redis storage
func (f *DDDRedisRepositoryFactory) CreateRouteRepository(strategy domain.StorageStrategy) domain.RouteRepository {
	if strategy == domain.StorageStrategyRedis || strategy == domain.StorageStrategyHybrid {
		return NewDDDRedisRouteRepository(f.client)
	}
	return nil
}

// CreateMetricsRepository creates a metrics repository using Redis storage
func (f *DDDRedisRepositoryFactory) CreateMetricsRepository(strategy domain.StorageStrategy) domain.MetricsRepository {
	if strategy == domain.StorageStrategyRedis || strategy == domain.StorageStrategyHybrid {
		return NewDDDRedisMetricsRepository(f.client)
	}
	return nil
}

// CreateScalingDecisionRepository creates a scaling decision repository using Redis storage
func (f *DDDRedisRepositoryFactory) CreateScalingDecisionRepository(strategy domain.StorageStrategy) domain.ScalingDecisionRepository {
	if strategy == domain.StorageStrategyRedis || strategy == domain.StorageStrategyHybrid {
		return NewDDDRedisScalingDecisionRepository(f.client)
	}
	return nil
}
