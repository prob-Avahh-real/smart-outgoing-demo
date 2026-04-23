package domain

import (
	"context"
	"fmt"
	"log"
)

// VehicleManagementService 车辆管理领域服务
type VehicleManagementService struct {
	vehicleRepo        VehicleRepository
	routeRepo          RouteRepository
	eventBus           EventBus
	consistencyChecker ConsistencyChecker
}

// ConsistencyChecker 一致性检查器接口
type ConsistencyChecker interface {
	CheckVehicleRouteConsistency(ctx context.Context, vehicleID string) error
	ReconcileVehicleData(ctx context.Context, vehicleID string) error
}

// NewVehicleManagementService 创建车辆管理服务
func NewVehicleManagementService(
	vehicleRepo VehicleRepository,
	routeRepo RouteRepository,
	eventBus EventBus,
	consistencyChecker ConsistencyChecker,
) *VehicleManagementService {
	return &VehicleManagementService{
		vehicleRepo:        vehicleRepo,
		routeRepo:          routeRepo,
		eventBus:           eventBus,
		consistencyChecker: consistencyChecker,
	}
}

// CreateVehicle 创建新车辆
func (s *VehicleManagementService) CreateVehicle(ctx context.Context, id, name, role string, coords Coordinates) (*Vehicle, error) {
	// 检查车辆是否已存在
	existing, err := s.vehicleRepo.FindByID(ctx, id)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("vehicle with ID %s already exists", id)
	}

	// 创建车辆聚合
	vehicle := NewVehicle(id, name, role, coords)

	// 保存车辆
	if err := s.vehicleRepo.Save(ctx, vehicle); err != nil {
		return nil, fmt.Errorf("failed to save vehicle: %w", err)
	}

	// 发布车辆创建事件
	event := NewVehicleCreatedEvent(*vehicle)
	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("Warning: failed to publish vehicle created event: %v", err)
	}

	log.Printf("Vehicle created: %s (%s)", id, name)
	return vehicle, nil
}

// AssignRoute 为车辆分配路线
func (s *VehicleManagementService) AssignRoute(ctx context.Context, vehicleID string, waypoints []Coordinates) (*Route, error) {
	// 获取车辆
	vehicle, err := s.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("vehicle not found: %w", err)
	}

	// 检查车辆状态
	if vehicle.Status() != VehicleStatusIdle {
		return nil, fmt.Errorf("vehicle %s is not idle, current status: %s", 
			vehicleID, vehicle.Status().String())
	}

	// 创建路线聚合
	routeID := fmt.Sprintf("route_%s_%d", vehicleID, len(waypoints))
	route := NewRoute(routeID, vehicleID, waypoints)

	// 保存路线
	if err := s.routeRepo.Save(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to save route: %w", err)
	}

	// 更新车辆状态
	destination := waypoints[len(waypoints)-1]
	vehicle.SetDestination(destination)
	if err := s.vehicleRepo.Save(ctx, vehicle); err != nil {
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	// 激活路线
	route.Activate()
	if err := s.routeRepo.Save(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to activate route: %w", err)
	}

	log.Printf("Route assigned to vehicle %s: %d waypoints", vehicleID, len(waypoints))
	return route, nil
}

// CompleteRoute 完成路线
func (s *VehicleManagementService) CompleteRoute(ctx context.Context, vehicleID string) error {
	// 获取车辆
	vehicle, err := s.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("vehicle not found: %w", err)
	}

	// 获取当前路线
	route, err := s.routeRepo.FindByVehicleID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("route not found for vehicle: %w", err)
	}

	// 完成路线
	route.Complete()
	if err := s.routeRepo.Save(ctx, route); err != nil {
		return fmt.Errorf("failed to complete route: %w", err)
	}

	// 更新车辆状态
	vehicle.Arrive()
	if err := s.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	// 发布路线完成事件
	event := NewRouteCompletedEvent(*route)
	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("Warning: failed to publish route completed event: %v", err)
	}

	// 检查一致性
	if err := s.consistencyChecker.CheckVehicleRouteConsistency(ctx, vehicleID); err != nil {
		log.Printf("Consistency check failed for vehicle %s: %v", vehicleID, err)
		// 尝试修复
		if reconcileErr := s.consistencyChecker.ReconcileVehicleData(ctx, vehicleID); reconcileErr != nil {
			log.Printf("Failed to reconcile vehicle %s: %v", vehicleID, reconcileErr)
		}
	}

	log.Printf("Route completed for vehicle %s", vehicleID)
	return nil
}

// UpdateVehiclePosition 更新车辆位置
func (s *VehicleManagementService) UpdateVehiclePosition(ctx context.Context, vehicleID string, coords Coordinates) error {
	// 获取车辆
	vehicle, err := s.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("vehicle not found: %w", err)
	}

	// 更新位置
	vehicle.UpdatePosition(coords)
	if err := s.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to update vehicle position: %w", err)
	}

	// 如果车辆有活跃路线，更新路线进度
	if vehicle.Status() == VehicleStatusMoving {
		route, err := s.routeRepo.FindByVehicleID(ctx, vehicleID)
		if err == nil && route != nil {
			route.UpdateProgress(0) // 简化实现，实际应该计算进度
			if err := s.routeRepo.Save(ctx, route); err != nil {
				log.Printf("Warning: failed to update route progress: %v", err)
			}
		}
	}

	return nil
}

// SetVehicleMaintenance 设置车辆维护状态
func (s *VehicleManagementService) SetVehicleMaintenance(ctx context.Context, vehicleID string) error {
	// 获取车辆
	vehicle, err := s.vehicleRepo.FindByID(ctx, vehicleID)
	if err != nil {
		return fmt.Errorf("vehicle not found: %w", err)
	}

	// 如果车辆有活跃路线，先取消路线
	if vehicle.Status() == VehicleStatusMoving {
		route, err := s.routeRepo.FindByVehicleID(ctx, vehicleID)
		if err == nil && route != nil {
			route.Cancel()
			if err := s.routeRepo.Save(ctx, route); err != nil {
				log.Printf("Warning: failed to cancel route: %v", err)
			}
		}
	}

	// 设置维护状态
	vehicle.SetMaintenance()
	if err := s.vehicleRepo.Save(ctx, vehicle); err != nil {
		return fmt.Errorf("failed to set vehicle maintenance: %w", err)
	}

	log.Printf("Vehicle %s set to maintenance", vehicleID)
	return nil
}
