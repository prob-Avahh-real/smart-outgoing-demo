package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"smart-outgoing-demo/internal/domain"
	"smart-outgoing-demo/internal/infrastructure/memory"
	"smart-outgoing-demo/internal/infrastructure/redis"
)

// DDDIntegrationService DDD架构集成服务
type DDDIntegrationService struct {
	scalingOrchestrator *DDDScalingOrchestrator
	vehicleService      *domain.VehicleManagementService
	scalingService      *domain.ScalingDecisionService
	eventBus            domain.EventBus
	currentStrategy     domain.StorageStrategy
}

// NewDDDIntegrationService 创建DDD集成服务
func NewDDDIntegrationService() *DDDIntegrationService {
	// 创建事件总线
	eventBus := domain.NewInMemoryEventBus()

	// 创建内存仓储工厂
	memoryFactory := memory.NewDDDMemoryRepositoryFactory()

	// 创建Redis客户端（使用Mock实现）
	redisClient := redis.NewMockRedisClient()
	redisFactory := redis.NewDDDRedisRepositoryFactory(redisClient)

	// 创建统一的仓储工厂
	unifiedFactory := NewUnifiedRepositoryFactory(memoryFactory.(*memory.DDDMemoryRepositoryFactory), redisFactory)

	// 创建仓储实例
	metricsRepo := unifiedFactory.CreateMetricsRepository(domain.StorageStrategyMemory)
	scalingDecisionRepo := unifiedFactory.CreateScalingDecisionRepository(domain.StorageStrategyMemory)
	vehicleRepo := unifiedFactory.CreateVehicleRepository(domain.StorageStrategyMemory)
	routeRepo := unifiedFactory.CreateRouteRepository(domain.StorageStrategyMemory)

	// 创建领域服务
	scalingService := domain.NewScalingDecisionService(
		metricsRepo,
		scalingDecisionRepo,
		eventBus,
		domain.DefaultThresholds(),
	)

	// 创建一致性检查器
	consistencyChecker := NewDefaultConsistencyChecker(vehicleRepo, routeRepo)

	vehicleService := domain.NewVehicleManagementService(
		vehicleRepo,
		routeRepo,
		eventBus,
		consistencyChecker,
	)

	// 创建迁移服务
	migrationService := NewFixedMigrationService(memoryFactory.(*memory.DDDMemoryRepositoryFactory), redisFactory)

	// 创建扩容协调器
	scalingOrchestrator := NewDDDMScalingOrchestrator(
		scalingService,
		vehicleService,
		unifiedFactory,
		migrationService,
		eventBus,
		domain.StorageStrategyMemory,
	)

	return &DDDIntegrationService{
		scalingOrchestrator: scalingOrchestrator,
		vehicleService:      vehicleService,
		scalingService:      scalingService,
		eventBus:            eventBus,
		currentStrategy:     domain.StorageStrategyMemory,
	}
}

// Start 启动DDD系统
func (s *DDDIntegrationService) Start(ctx context.Context) error {
	log.Printf("Starting DDD Auto-Scaling System with strategy: %s", s.currentStrategy.String())

	// 启动自动扩容
	go s.scalingOrchestrator.StartAutoScaling(ctx)

	// 注册事件处理器
	s.registerEventHandlers()

	log.Printf("DDD Auto-Scaling System started successfully")
	return nil
}

// registerEventHandlers 注册事件处理器
func (s *DDDIntegrationService) registerEventHandlers() {
	// 注册存储策略变更事件处理器
	s.eventBus.Subscribe("StorageStrategyChanged", &StorageStrategyChangedHandler{
		orchestrator: s.scalingOrchestrator,
	})

	// 注册阈值突破事件处理器
	s.eventBus.Subscribe("ThresholdBreached", &ThresholdBreachedHandler{
		scalingService: s.scalingService,
	})

	// 注册车辆创建事件处理器
	s.eventBus.Subscribe("VehicleCreated", &VehicleCreatedHandler{
		vehicleService: s.vehicleService,
	})

	// 注册路线完成事件处理器
	s.eventBus.Subscribe("RouteCompleted", &RouteCompletedHandler{
		vehicleService: s.vehicleService,
	})
}

// CreateVehicle 创建车辆（对外接口）
func (s *DDDIntegrationService) CreateVehicle(ctx context.Context, id, name, role string, lng, lat, alt float64) (*domain.Vehicle, error) {
	coords := domain.Coordinates{
		Longitude: lng,
		Latitude:  lat,
		Altitude:  alt,
	}

	return s.vehicleService.CreateVehicle(ctx, id, name, role, coords)
}

// AssignRoute 分配路线（对外接口）
func (s *DDDIntegrationService) AssignRoute(ctx context.Context, vehicleID string, waypoints []domain.Coordinates) (*domain.Route, error) {
	return s.vehicleService.AssignRoute(ctx, vehicleID, waypoints)
}

// UpdateVehiclePosition 更新车辆位置（对外接口）
func (s *DDDIntegrationService) UpdateVehiclePosition(ctx context.Context, vehicleID string, lng, lat, alt float64) error {
	coords := domain.Coordinates{
		Longitude: lng,
		Latitude:  lat,
		Altitude:  alt,
	}

	return s.vehicleService.UpdateVehiclePosition(ctx, vehicleID, coords)
}

// GetCurrentMetrics 获取当前指标（对外接口）
func (s *DDDIntegrationService) GetCurrentMetrics(ctx context.Context) (*domain.Metrics, error) {
	// 创建示例指标（实际实现中应该从系统监控获取）
	memoryUsage := domain.MemoryMetrics{
		TotalMB:      1024,
		UsedMB:       512,
		AvailableMB:  512,
		UsagePercent: 50.0,
	}

	metrics := domain.NewMetrics(
		fmt.Sprintf("metrics_%d", time.Now().Unix()),
		memoryUsage,
		10,   // agent count
		100,  // connection count
		25.0, // CPU usage
		s.currentStrategy,
	)

	return metrics, nil
}

// ForceScaling 强制扩容（对外接口）
func (s *DDDIntegrationService) ForceScaling(ctx context.Context, targetStrategy domain.StorageStrategy) error {
	switch targetStrategy {
	case domain.StorageStrategyMemory:
		return s.scalingOrchestrator.ForceScaleToMemory(ctx)
	case domain.StorageStrategyRedis:
		return s.scalingOrchestrator.ForceScaleToRedis(ctx)
	default:
		return fmt.Errorf("unsupported target strategy: %s", targetStrategy.String())
	}
}

// GetScalingStatus 获取扩容状态（对外接口）
func (s *DDDIntegrationService) GetScalingStatus() ScalingStatusInfo {
	return ScalingStatusInfo{
		CurrentStrategy:       s.currentStrategy.String(),
		IsMigrationInProgress: s.scalingOrchestrator.IsMigrationInProgress(),
		LastScalingTime:       time.Now(), // 简化实现
	}
}

// ScalingStatusInfo 扩容状态信息
type ScalingStatusInfo struct {
	CurrentStrategy       string    `json:"current_strategy"`
	IsMigrationInProgress bool      `json:"is_migration_in_progress"`
	LastScalingTime       time.Time `json:"last_scaling_time"`
}
