package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// TestDDDIntegration 测试DDD集成
func TestDDDIntegration(t *testing.T) {
	// 创建集成服务
	integrationService := NewDDDIntegrationService()

	// 启动上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 启动系统
	err := integrationService.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start DDD integration: %v", err)
	}

	// 测试创建车辆
	vehicle, err := integrationService.CreateVehicle(ctx, "test-vehicle-1", "Test Vehicle", "transport", 0.0, 0.0, 0.0)
	if err != nil {
		t.Fatalf("Failed to create vehicle: %v", err)
	}

	if vehicle.ID() != "test-vehicle-1" {
		t.Errorf("Expected vehicle ID 'test-vehicle-1', got '%s'", vehicle.ID())
	}

	if vehicle.Status() != domain.VehicleStatusIdle {
		t.Errorf("Expected vehicle status 'idle', got '%s'", vehicle.Status().String())
	}

	// 测试更新位置
	err = integrationService.UpdateVehiclePosition(ctx, "test-vehicle-1", 1.0, 1.0, 0.0)
	if err != nil {
		t.Fatalf("Failed to update vehicle position: %v", err)
	}

	// 测试分配路线
	waypoints := []domain.Coordinates{
		{Longitude: 0.0, Latitude: 0.0},
		{Longitude: 1.0, Latitude: 1.0},
		{Longitude: 2.0, Latitude: 2.0},
	}

	route, err := integrationService.AssignRoute(ctx, "test-vehicle-1", waypoints)
	if err != nil {
		t.Fatalf("Failed to assign route: %v", err)
	}

	if route.VehicleID() != "test-vehicle-1" {
		t.Errorf("Expected route vehicle ID 'test-vehicle-1', got '%s'", route.VehicleID())
	}

	// 测试获取指标
	metrics, err := integrationService.GetCurrentMetrics(ctx)
	if err != nil {
		t.Fatalf("Failed to get current metrics: %v", err)
	}

	if metrics.StorageStrategy() != domain.StorageStrategyMemory {
		t.Errorf("Expected storage strategy 'memory', got '%s'", metrics.StorageStrategy().String())
	}

	// 测试获取扩容状态
	status := integrationService.GetScalingStatus()
	if status.CurrentStrategy != domain.StorageStrategyMemory.String() {
		t.Errorf("Expected current strategy 'memory', got '%s'", status.CurrentStrategy)
	}

	t.Log("DDD Integration test completed successfully")
}

// TestScalingDecisionService 测试扩容决策服务
func TestScalingDecisionService(t *testing.T) {
	// 创建模拟组件
	eventBus := domain.NewInMemoryEventBus()
	memoryFactory := NewMockMemoryRepositoryFactory()

	metricsRepo := memoryFactory.CreateMetricsRepository(domain.StorageStrategyMemory)
	scalingDecisionRepo := memoryFactory.CreateScalingDecisionRepository(domain.StorageStrategyMemory)

	// 创建扩容决策服务
	scalingService := domain.NewScalingDecisionService(
		metricsRepo,
		scalingDecisionRepo,
		eventBus,
		domain.DefaultThresholds(),
	)

	// 创建测试指标
	memoryUsage := domain.MemoryMetrics{
		TotalMB:      1024,
		UsedMB:       900, // 超过80%阈值
		AvailableMB:  124,
		UsagePercent: 87.5,
	}

	metrics := domain.NewMetrics(
		"test-metrics-1",
		memoryUsage,
		150,  // 超过100阈值
		1200, // 超过1000阈值
		80.0, // 超过75%阈值
		domain.StorageStrategyMemory,
	)

	// 保存指标
	ctx := context.Background()
	err := metricsRepo.Save(ctx, metrics)
	if err != nil {
		t.Fatalf("Failed to save metrics: %v", err)
	}

	// 评估扩容决策
	decision, err := scalingService.EvaluateScaling(ctx)
	if err != nil {
		t.Fatalf("Failed to evaluate scaling: %v", err)
	}

	if decision == nil {
		t.Error("Expected scaling decision, got nil")
	} else {
		if decision.ToStrategy() != domain.StorageStrategyRedis {
			t.Errorf("Expected target strategy 'redis', got '%s'", decision.ToStrategy().String())
		}

		if decision.Reason() != domain.ScalingReasonMemory {
			t.Errorf("Expected scaling reason 'memory', got '%s'", decision.Reason().String())
		}
	}

	t.Log("Scaling decision service test completed successfully")
}

// MockMemoryRepositoryFactory 模拟内存仓储工厂
type MockMemoryRepositoryFactory struct{}

func NewMockMemoryRepositoryFactory() *MockMemoryRepositoryFactory {
	return &MockMemoryRepositoryFactory{}
}

func (f *MockMemoryRepositoryFactory) CreateVehicleRepository(strategy domain.StorageStrategy) domain.VehicleRepository {
	return &MockVehicleRepository{vehicles: make(map[string]*domain.Vehicle)}
}

func (f *MockMemoryRepositoryFactory) CreateRouteRepository(strategy domain.StorageStrategy) domain.RouteRepository {
	return &MockRouteRepository{routes: make(map[string]*domain.Route)}
}

func (f *MockMemoryRepositoryFactory) CreateMetricsRepository(strategy domain.StorageStrategy) domain.MetricsRepository {
	return &MockMetricsRepository{metrics: make(map[string]*domain.Metrics)}
}

func (f *MockMemoryRepositoryFactory) CreateScalingDecisionRepository(strategy domain.StorageStrategy) domain.ScalingDecisionRepository {
	return &MockScalingDecisionRepository{decisions: make(map[string]*domain.ScalingDecision)}
}

// Mock implementations for testing
type MockVehicleRepository struct {
	vehicles map[string]*domain.Vehicle
}

func (r *MockVehicleRepository) Save(ctx context.Context, vehicle *domain.Vehicle) error {
	r.vehicles[vehicle.ID()] = vehicle
	return nil
}

func (r *MockVehicleRepository) FindByID(ctx context.Context, id string) (*domain.Vehicle, error) {
	if v, exists := r.vehicles[id]; exists {
		return v, nil
	}
	return nil, fmt.Errorf("vehicle not found")
}

func (r *MockVehicleRepository) FindAll(ctx context.Context) ([]*domain.Vehicle, error) {
	var vehicles []*domain.Vehicle
	for _, v := range r.vehicles {
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (r *MockVehicleRepository) Delete(ctx context.Context, id string) error {
	delete(r.vehicles, id)
	return nil
}

func (r *MockVehicleRepository) FindByStatus(ctx context.Context, status domain.VehicleStatus) ([]*domain.Vehicle, error) {
	var vehicles []*domain.Vehicle
	for _, v := range r.vehicles {
		if v.Status() == status {
			vehicles = append(vehicles, v)
		}
	}
	return vehicles, nil
}

func (r *MockVehicleRepository) Count(ctx context.Context) (int, error) {
	return len(r.vehicles), nil
}

// Mock implementations for other repositories (simplified)
type MockRouteRepository struct {
	routes map[string]*domain.Route
}

func (r *MockRouteRepository) Save(ctx context.Context, route *domain.Route) error {
	r.routes[route.ID()] = route
	return nil
}

func (r *MockRouteRepository) FindByID(ctx context.Context, id string) (*domain.Route, error) {
	if route, exists := r.routes[id]; exists {
		return route, nil
	}
	return nil, fmt.Errorf("route not found")
}

func (r *MockRouteRepository) FindByVehicleID(ctx context.Context, vehicleID string) (*domain.Route, error) {
	for _, route := range r.routes {
		if route.VehicleID() == vehicleID {
			return route, nil
		}
	}
	return nil, fmt.Errorf("route not found for vehicle")
}

func (r *MockRouteRepository) FindAll(ctx context.Context) ([]*domain.Route, error) {
	var routes []*domain.Route
	for _, route := range r.routes {
		routes = append(routes, route)
	}
	return routes, nil
}

func (r *MockRouteRepository) Delete(ctx context.Context, id string) error {
	delete(r.routes, id)
	return nil
}

func (r *MockRouteRepository) FindByStatus(ctx context.Context, status domain.RouteStatus) ([]*domain.Route, error) {
	var routes []*domain.Route
	for _, route := range r.routes {
		if route.Status() == status {
			routes = append(routes, route)
		}
	}
	return routes, nil
}

type MockMetricsRepository struct {
	metrics map[string]*domain.Metrics
}

func (r *MockMetricsRepository) Save(ctx context.Context, metrics *domain.Metrics) error {
	r.metrics[metrics.ID()] = metrics
	return nil
}

func (r *MockMetricsRepository) FindLatest(ctx context.Context) (*domain.Metrics, error) {
	var latest *domain.Metrics
	for _, metrics := range r.metrics {
		if latest == nil || metrics.Timestamp().After(latest.Timestamp()) {
			latest = metrics
		}
	}
	if latest == nil {
		return nil, fmt.Errorf("no metrics found")
	}
	return latest, nil
}

func (r *MockMetricsRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*domain.Metrics, error) {
	var results []*domain.Metrics
	for _, metrics := range r.metrics {
		if metrics.Timestamp().After(start) && metrics.Timestamp().Before(end) {
			results = append(results, metrics)
		}
	}
	return results, nil
}

func (r *MockMetricsRepository) Delete(ctx context.Context, id string) error {
	delete(r.metrics, id)
	return nil
}

type MockScalingDecisionRepository struct {
	decisions map[string]*domain.ScalingDecision
}

func (r *MockScalingDecisionRepository) Save(ctx context.Context, decision *domain.ScalingDecision) error {
	r.decisions[decision.ID()] = decision
	return nil
}

func (r *MockScalingDecisionRepository) FindByID(ctx context.Context, id string) (*domain.ScalingDecision, error) {
	if decision, exists := r.decisions[id]; exists {
		return decision, nil
	}
	return nil, fmt.Errorf("decision not found")
}

func (r *MockScalingDecisionRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*domain.ScalingDecision, error) {
	var results []*domain.ScalingDecision
	for _, decision := range r.decisions {
		if decision.DecisionTime().After(start) && decision.DecisionTime().Before(end) {
			results = append(results, decision)
		}
	}
	return results, nil
}

func (r *MockScalingDecisionRepository) FindLatest(ctx context.Context) (*domain.ScalingDecision, error) {
	var latest *domain.ScalingDecision
	for _, decision := range r.decisions {
		if latest == nil || decision.DecisionTime().After(latest.DecisionTime()) {
			latest = decision
		}
	}
	if latest == nil {
		return nil, fmt.Errorf("no decisions found")
	}
	return latest, nil
}

func (r *MockScalingDecisionRepository) Delete(ctx context.Context, id string) error {
	delete(r.decisions, id)
	return nil
}
