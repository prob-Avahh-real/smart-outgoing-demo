package services

import (
	"smart-outgoing-demo/internal/domain"
	"smart-outgoing-demo/internal/infrastructure/memory"
	"smart-outgoing-demo/internal/infrastructure/redis"
)

// UnifiedRepositoryFactory 统一仓储工厂，根据策略选择合适的仓储实现
type UnifiedRepositoryFactory struct {
	memoryFactory *memory.DDDMemoryRepositoryFactory
	redisFactory  *redis.DDDRedisRepositoryFactory
}

// NewUnifiedRepositoryFactory 创建统一仓储工厂
func NewUnifiedRepositoryFactory(
	memoryFactory *memory.DDDMemoryRepositoryFactory,
	redisFactory *redis.DDDRedisRepositoryFactory,
) *UnifiedRepositoryFactory {
	return &UnifiedRepositoryFactory{
		memoryFactory: memoryFactory,
		redisFactory:  redisFactory,
	}
}

// CreateVehicleRepository 创建车辆仓储
func (f *UnifiedRepositoryFactory) CreateVehicleRepository(strategy domain.StorageStrategy) domain.VehicleRepository {
	switch strategy {
	case domain.StorageStrategyMemory:
		return f.memoryFactory.CreateVehicleRepository(strategy)
	case domain.StorageStrategyRedis:
		return f.redisFactory.CreateVehicleRepository(strategy)
	case domain.StorageStrategyHybrid:
		// 混合策略可以优先使用Redis，失败时回退到Memory
		return f.redisFactory.CreateVehicleRepository(strategy)
	default:
		return f.memoryFactory.CreateVehicleRepository(strategy)
	}
}

// CreateRouteRepository 创建路线仓储
func (f *UnifiedRepositoryFactory) CreateRouteRepository(strategy domain.StorageStrategy) domain.RouteRepository {
	switch strategy {
	case domain.StorageStrategyMemory:
		return f.memoryFactory.CreateRouteRepository(strategy)
	case domain.StorageStrategyRedis:
		return f.redisFactory.CreateRouteRepository(strategy)
	case domain.StorageStrategyHybrid:
		return f.redisFactory.CreateRouteRepository(strategy)
	default:
		return f.memoryFactory.CreateRouteRepository(strategy)
	}
}

// CreateMetricsRepository 创建指标仓储
func (f *UnifiedRepositoryFactory) CreateMetricsRepository(strategy domain.StorageStrategy) domain.MetricsRepository {
	switch strategy {
	case domain.StorageStrategyMemory:
		return f.memoryFactory.CreateMetricsRepository(strategy)
	case domain.StorageStrategyRedis:
		return f.redisFactory.CreateMetricsRepository(strategy)
	case domain.StorageStrategyHybrid:
		return f.redisFactory.CreateMetricsRepository(strategy)
	default:
		return f.memoryFactory.CreateMetricsRepository(strategy)
	}
}

// CreateScalingDecisionRepository 创建扩容决策仓储
func (f *UnifiedRepositoryFactory) CreateScalingDecisionRepository(strategy domain.StorageStrategy) domain.ScalingDecisionRepository {
	switch strategy {
	case domain.StorageStrategyMemory:
		return f.memoryFactory.CreateScalingDecisionRepository(strategy)
	case domain.StorageStrategyRedis:
		return f.redisFactory.CreateScalingDecisionRepository(strategy)
	case domain.StorageStrategyHybrid:
		return f.redisFactory.CreateScalingDecisionRepository(strategy)
	default:
		return f.memoryFactory.CreateScalingDecisionRepository(strategy)
	}
}
