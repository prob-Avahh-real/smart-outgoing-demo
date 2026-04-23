package redis

import (
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// RedisRepositoryFactory implements RepositoryFactory for Redis storage
type RedisRepositoryFactory struct {
	client RedisClientInterface
}

// NewRedisRepositoryFactory creates a new Redis repository factory
func NewRedisRepositoryFactory(client RedisClientInterface) repositories.RepositoryFactory {
	return &RedisRepositoryFactory{
		client: client,
	}
}

// CreateVehicleRepository creates a Redis vehicle repository
func (f *RedisRepositoryFactory) CreateVehicleRepository(strategy entities.StorageStrategy) repositories.VehicleRepository {
	return NewRedisVehicleRepository(f.client)
}

// CreateRouteRepository creates a Redis route repository
func (f *RedisRepositoryFactory) CreateRouteRepository(strategy entities.StorageStrategy) repositories.RouteRepository {
	return NewRedisRouteRepository(f.client)
}

// CreateMetricsRepository creates a Redis metrics repository
func (f *RedisRepositoryFactory) CreateMetricsRepository(strategy entities.StorageStrategy) repositories.MetricsRepository {
	return NewRedisMetricsRepository(f.client)
}

// CreateScalingDecisionRepository creates a Redis scaling decision repository
func (f *RedisRepositoryFactory) CreateScalingDecisionRepository(strategy entities.StorageStrategy) repositories.ScalingDecisionRepository {
	return NewRedisScalingDecisionRepository(f.client)
}
