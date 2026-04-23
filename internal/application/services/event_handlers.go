package services

import (
	"context"
	"log"

	"smart-outgoing-demo/internal/domain"
)

// StorageStrategyChangedHandler 存储策略变更事件处理器
type StorageStrategyChangedHandler struct {
	orchestrator *DDDScalingOrchestrator
}

// Handle 处理存储策略变更事件
func (h *StorageStrategyChangedHandler) Handle(event domain.DomainEvent) error {
	log.Printf("Handling StorageStrategyChanged event: %s", event.EventType())
	// 可以在这里添加额外的处理逻辑，如通知、日志记录等
	return nil
}

// ThresholdBreachedHandler 阈值突破事件处理器
type ThresholdBreachedHandler struct {
	scalingService *domain.ScalingDecisionService
}

// Handle 处理阈值突破事件
func (h *ThresholdBreachedHandler) Handle(event domain.DomainEvent) error {
	log.Printf("Handling ThresholdBreached event: %s", event.EventType())
	// 可以在这里添加额外的处理逻辑，如告警、通知等
	return nil
}

// VehicleCreatedHandler 车辆创建事件处理器
type VehicleCreatedHandler struct {
	vehicleService *domain.VehicleManagementService
}

// Handle 处理车辆创建事件
func (h *VehicleCreatedHandler) Handle(event domain.DomainEvent) error {
	log.Printf("Handling VehicleCreated event: %s", event.EventType())
	// 可以在这里添加额外的处理逻辑，如初始化路线、统计更新等
	return nil
}

// RouteCompletedHandler 路线完成事件处理器
type RouteCompletedHandler struct {
	vehicleService *domain.VehicleManagementService
}

// Handle 处理路线完成事件
func (h *RouteCompletedHandler) Handle(event domain.DomainEvent) error {
	log.Printf("Handling RouteCompleted event: %s", event.EventType())
	// 可以在这里添加额外的处理逻辑，如清理资源、更新统计等
	return nil
}

// DefaultConsistencyChecker 默认一致性检查器
type DefaultConsistencyChecker struct {
	vehicleRepo domain.VehicleRepository
	routeRepo   domain.RouteRepository
}

// NewDefaultConsistencyChecker 创建默认一致性检查器
func NewDefaultConsistencyChecker(vehicleRepo domain.VehicleRepository, routeRepo domain.RouteRepository) *DefaultConsistencyChecker {
	return &DefaultConsistencyChecker{
		vehicleRepo: vehicleRepo,
		routeRepo:   routeRepo,
	}
}

// CheckVehicleRouteConsistency 检查车辆路线一致性
func (c *DefaultConsistencyChecker) CheckVehicleRouteConsistency(ctx context.Context, vehicleID string) error {
	// 获取车辆
	vehicle, err := c.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return err
	}

	// 获取车辆路线
	route, err := c.routeRepo.FindByVehicleID(ctx, vehicleID)
	if err != nil {
		// 车辆没有路线是正常情况
		return nil
	}

	// 检查一致性：如果车辆状态为moving，应该有活跃路线
	if vehicle.Status() == domain.VehicleStatusMoving && route.Status() != domain.RouteStatusActive {
		log.Printf("Consistency issue: vehicle %s is moving but route %s is not active", vehicleID, route.ID())
	}

	// 检查一致性：如果车辆状态为idle，不应该有活跃路线
	if vehicle.Status() == domain.VehicleStatusIdle && route.Status() == domain.RouteStatusActive {
		log.Printf("Consistency issue: vehicle %s is idle but route %s is still active", vehicleID, route.ID())
	}

	return nil
}

// ReconcileVehicleData 修复车辆数据一致性
func (c *DefaultConsistencyChecker) ReconcileVehicleData(ctx context.Context, vehicleID string) error {
	// 获取车辆
	vehicle, err := c.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return err
	}

	// 获取车辆路线
	route, err := c.routeRepo.FindByVehicleID(ctx, vehicleID)
	if err != nil {
		// 如果车辆状态为moving但没有路线，重置为idle
		if vehicle.Status() == domain.VehicleStatusMoving {
			// 车辆到达目的地
			vehicle.Arrive()
			return c.vehicleRepo.Save(ctx, vehicle)
		}
		return nil
	}

	// 修复逻辑：确保车辆状态和路线状态一致
	if vehicle.Status() == domain.VehicleStatusIdle && route.Status() == domain.RouteStatusActive {
		// 取消路线
		route.Cancel()
		if err := c.routeRepo.Save(ctx, route); err != nil {
			return err
		}
	}

	return nil
}
