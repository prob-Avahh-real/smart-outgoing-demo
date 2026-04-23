package services

import (
	"context"
	"fmt"
	"log"

	"smart-outgoing-demo/internal/domain"
	"smart-outgoing-demo/internal/infrastructure/memory"
	"smart-outgoing-demo/internal/infrastructure/redis"
)

// FixedMigrationService 简化的数据迁移服务实现
type FixedMigrationService struct {
	memoryFactory *memory.DDDMemoryRepositoryFactory
	redisFactory  *redis.DDDRedisRepositoryFactory
}

// NewFixedMigrationService 创建简化迁移服务
func NewFixedMigrationService(
	memoryFactory *memory.DDDMemoryRepositoryFactory,
	redisFactory *redis.DDDRedisRepositoryFactory,
) *FixedMigrationService {
	return &FixedMigrationService{
		memoryFactory: memoryFactory,
		redisFactory:  redisFactory,
	}
}

// MigrateVehicles 迁移车辆数据
func (s *FixedMigrationService) MigrateVehicles(ctx context.Context, from, to domain.StorageStrategy) error {
	log.Printf("Migrating vehicles from %s to %s", from.String(), to.String())

	// 获取源仓储
	sourceRepo := s.getVehicleRepository(from)
	if sourceRepo == nil {
		return fmt.Errorf("failed to create source repository for strategy %s", from.String())
	}

	// 获取目标仓储
	targetRepo := s.getVehicleRepository(to)
	if targetRepo == nil {
		return fmt.Errorf("failed to create target repository for strategy %s", to.String())
	}

	// 读取所有车辆
	vehicles, err := sourceRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to read source vehicles: %w", err)
	}

	// 迁移到目标存储
	for _, vehicle := range vehicles {
		if err := targetRepo.Save(ctx, vehicle); err != nil {
			log.Printf("Warning: failed to migrate vehicle %s: %v", vehicle.ID(), err)
			continue
		}
	}

	log.Printf("Successfully migrated %d vehicles", len(vehicles))
	return nil
}

// MigrateRoutes 迁移路线数据
func (s *FixedMigrationService) MigrateRoutes(ctx context.Context, from, to domain.StorageStrategy) error {
	log.Printf("Migrating routes from %s to %s", from.String(), to.String())

	sourceRepo := s.getRouteRepository(from)
	if sourceRepo == nil {
		return fmt.Errorf("failed to create source route repository for strategy %s", from.String())
	}

	targetRepo := s.getRouteRepository(to)
	if targetRepo == nil {
		return fmt.Errorf("failed to create target route repository for strategy %s", to.String())
	}

	routes, err := sourceRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to read source routes: %w", err)
	}

	for _, route := range routes {
		if err := targetRepo.Save(ctx, route); err != nil {
			log.Printf("Warning: failed to migrate route %s: %v", route.ID(), err)
			continue
		}
	}

	log.Printf("Successfully migrated %d routes", len(routes))
	return nil
}

// MigrateMetrics 迁移指标数据
func (s *FixedMigrationService) MigrateMetrics(ctx context.Context, from, to domain.StorageStrategy) error {
	log.Printf("Migrating metrics from %s to %s", from.String(), to.String())

	sourceRepo := s.getMetricsRepository(from)
	if sourceRepo == nil {
		return fmt.Errorf("failed to create source metrics repository for strategy %s", from.String())
	}

	targetRepo := s.getMetricsRepository(to)
	if targetRepo == nil {
		return fmt.Errorf("failed to create target metrics repository for strategy %s", to.String())
	}

	metrics, err := sourceRepo.FindLatest(ctx)
	if err != nil {
		return fmt.Errorf("failed to read source metrics: %w", err)
	}

	if err := targetRepo.Save(ctx, metrics); err != nil {
		return fmt.Errorf("failed to migrate metrics: %w", err)
	}

	log.Printf("Successfully migrated latest metrics")
	return nil
}

// MigrateScalingDecisions 迁移扩容决策数据
func (s *FixedMigrationService) MigrateScalingDecisions(ctx context.Context, from, to domain.StorageStrategy) error {
	log.Printf("Migrating scaling decisions from %s to %s", from.String(), to.String())

	sourceRepo := s.getScalingDecisionRepository(from)
	if sourceRepo == nil {
		return fmt.Errorf("failed to create source scaling decision repository for strategy %s", from.String())
	}

	targetRepo := s.getScalingDecisionRepository(to)
	if targetRepo == nil {
		return fmt.Errorf("failed to create target scaling decision repository for strategy %s", to.String())
	}

	decision, err := sourceRepo.FindLatest(ctx)
	if err != nil {
		return fmt.Errorf("failed to read source scaling decision: %w", err)
	}

	if err := targetRepo.Save(ctx, decision); err != nil {
		return fmt.Errorf("failed to migrate scaling decision: %w", err)
	}

	log.Printf("Successfully migrated latest scaling decision")
	return nil
}

// Helper methods to get repositories
func (s *FixedMigrationService) getVehicleRepository(strategy domain.StorageStrategy) domain.VehicleRepository {
	switch strategy {
	case domain.StorageStrategyMemory:
		return s.memoryFactory.CreateVehicleRepository(strategy)
	case domain.StorageStrategyRedis:
		return s.redisFactory.CreateVehicleRepository(strategy)
	default:
		return nil
	}
}

func (s *FixedMigrationService) getRouteRepository(strategy domain.StorageStrategy) domain.RouteRepository {
	switch strategy {
	case domain.StorageStrategyMemory:
		return s.memoryFactory.CreateRouteRepository(strategy)
	case domain.StorageStrategyRedis:
		return s.redisFactory.CreateRouteRepository(strategy)
	default:
		return nil
	}
}

func (s *FixedMigrationService) getMetricsRepository(strategy domain.StorageStrategy) domain.MetricsRepository {
	switch strategy {
	case domain.StorageStrategyMemory:
		return s.memoryFactory.CreateMetricsRepository(strategy)
	case domain.StorageStrategyRedis:
		return s.redisFactory.CreateMetricsRepository(strategy)
	default:
		return nil
	}
}

func (s *FixedMigrationService) getScalingDecisionRepository(strategy domain.StorageStrategy) domain.ScalingDecisionRepository {
	switch strategy {
	case domain.StorageStrategyMemory:
		return s.memoryFactory.CreateScalingDecisionRepository(strategy)
	case domain.StorageStrategyRedis:
		return s.redisFactory.CreateScalingDecisionRepository(strategy)
	default:
		return nil
	}
}
