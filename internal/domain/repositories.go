package domain

import (
	"context"
	"time"
)

// VehicleRepository 车辆仓储接口
type VehicleRepository interface {
	Save(ctx context.Context, vehicle *Vehicle) error
	FindByID(ctx context.Context, id string) (*Vehicle, error)
	FindAll(ctx context.Context) ([]*Vehicle, error)
	Delete(ctx context.Context, id string) error
	FindByStatus(ctx context.Context, status VehicleStatus) ([]*Vehicle, error)
	Count(ctx context.Context) (int, error)
}

// RouteRepository 路线仓储接口
type RouteRepository interface {
	Save(ctx context.Context, route *Route) error
	FindByID(ctx context.Context, id string) (*Route, error)
	FindByVehicleID(ctx context.Context, vehicleID string) (*Route, error)
	FindAll(ctx context.Context) ([]*Route, error)
	Delete(ctx context.Context, id string) error
	FindByStatus(ctx context.Context, status RouteStatus) ([]*Route, error)
}

// MetricsRepository 指标仓储接口
type MetricsRepository interface {
	Save(ctx context.Context, metrics *Metrics) error
	FindLatest(ctx context.Context) (*Metrics, error)
	FindByTimeRange(ctx context.Context, start, end time.Time) ([]*Metrics, error)
	Delete(ctx context.Context, id string) error
}

// ScalingDecisionRepository 扩容决策仓储接口
type ScalingDecisionRepository interface {
	Save(ctx context.Context, decision *ScalingDecision) error
	FindByID(ctx context.Context, id string) (*ScalingDecision, error)
	FindByTimeRange(ctx context.Context, start, end time.Time) ([]*ScalingDecision, error)
	FindLatest(ctx context.Context) (*ScalingDecision, error)
	Delete(ctx context.Context, id string) error
}

// RepositoryFactory 仓储工厂接口
type RepositoryFactory interface {
	CreateVehicleRepository(strategy StorageStrategy) VehicleRepository
	CreateRouteRepository(strategy StorageStrategy) RouteRepository
	CreateMetricsRepository(strategy StorageStrategy) MetricsRepository
	CreateScalingDecisionRepository(strategy StorageStrategy) ScalingDecisionRepository
}

// EventRepository 事件仓储接口（用于事件溯源）
type EventRepository interface {
	Save(ctx context.Context, event DomainEvent) error
	FindByAggregateID(ctx context.Context, aggregateID string) ([]DomainEvent, error)
	FindByEventType(ctx context.Context, eventType string) ([]DomainEvent, error)
	FindByTimeRange(ctx context.Context, start, end time.Time) ([]DomainEvent, error)
}
