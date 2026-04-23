package domain

import (
	"time"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventID() string
	EventType() string
	OccurredOn() time.Time
	AggregateID() string
	Data() interface{}
}

// StorageStrategyChangedEvent 存储策略变更事件
type StorageStrategyChangedEvent struct {
	eventID      string
	eventType    string
	occurredOn   time.Time
	aggregateID  string
	fromStrategy StorageStrategy
	toStrategy   StorageStrategy
	reason       ScalingReason
	metrics      Metrics
}

// NewStorageStrategyChangedEvent 创建存储策略变更事件
func NewStorageStrategyChangedEvent(aggregateID string, fromStrategy, toStrategy StorageStrategy, reason ScalingReason, metrics Metrics) *StorageStrategyChangedEvent {
	return &StorageStrategyChangedEvent{
		eventID:      generateEventID(),
		eventType:    "StorageStrategyChanged",
		occurredOn:   time.Now(),
		aggregateID:  aggregateID,
		fromStrategy: fromStrategy,
		toStrategy:   toStrategy,
		reason:       reason,
		metrics:      metrics,
	}
}

// Getters
func (e *StorageStrategyChangedEvent) EventID() string               { return e.eventID }
func (e *StorageStrategyChangedEvent) EventType() string             { return e.eventType }
func (e *StorageStrategyChangedEvent) OccurredOn() time.Time          { return e.occurredOn }
func (e *StorageStrategyChangedEvent) AggregateID() string            { return e.aggregateID }
func (e *StorageStrategyChangedEvent) FromStrategy() StorageStrategy  { return e.fromStrategy }
func (e *StorageStrategyChangedEvent) ToStrategy() StorageStrategy    { return e.toStrategy }
func (e *StorageStrategyChangedEvent) Reason() ScalingReason          { return e.reason }
func (e *StorageStrategyChangedEvent) Metrics() Metrics               { return e.metrics }

// Data 返回事件数据
func (e *StorageStrategyChangedEvent) Data() interface{} {
	return map[string]interface{}{
		"from_strategy": e.fromStrategy.String(),
		"to_strategy":   e.toStrategy.String(),
		"reason":        e.reason.String(),
		"metrics":       e.metrics,
	}
}

// ThresholdBreachedEvent 阈值突破事件
type ThresholdBreachedEvent struct {
	eventID     string
	eventType   string
	occurredOn  time.Time
	aggregateID string
	threshold   ScalingThreshold
	current     Metrics
	breachedBy  []string
}

// NewThresholdBreachedEvent 创建阈值突破事件
func NewThresholdBreachedEvent(aggregateID string, threshold ScalingThreshold, current Metrics, breachedBy []string) *ThresholdBreachedEvent {
	return &ThresholdBreachedEvent{
		eventID:     generateEventID(),
		eventType:   "ThresholdBreached",
		occurredOn:  time.Now(),
		aggregateID: aggregateID,
		threshold:   threshold,
		current:     current,
		breachedBy:  breachedBy,
	}
}

// Getters
func (e *ThresholdBreachedEvent) EventID() string            { return e.eventID }
func (e *ThresholdBreachedEvent) EventType() string          { return e.eventType }
func (e *ThresholdBreachedEvent) OccurredOn() time.Time     { return e.occurredOn }
func (e *ThresholdBreachedEvent) AggregateID() string       { return e.aggregateID }
func (e *ThresholdBreachedEvent) Threshold() ScalingThreshold { return e.threshold }
func (e *ThresholdBreachedEvent) Current() Metrics           { return e.current }
func (e *ThresholdBreachedEvent) BreachedBy() []string       { return e.breachedBy }

// Data 返回事件数据
func (e *ThresholdBreachedEvent) Data() interface{} {
	return map[string]interface{}{
		"threshold":  e.threshold,
		"current":    e.current,
		"breached_by": e.breachedBy,
	}
}

// VehicleCreatedEvent 车辆创建事件
type VehicleCreatedEvent struct {
	eventID     string
	eventType   string
	occurredOn  time.Time
	aggregateID string
	vehicle     Vehicle
}

// NewVehicleCreatedEvent 创建车辆创建事件
func NewVehicleCreatedEvent(vehicle Vehicle) *VehicleCreatedEvent {
	return &VehicleCreatedEvent{
		eventID:     generateEventID(),
		eventType:   "VehicleCreated",
		occurredOn:  time.Now(),
		aggregateID: vehicle.ID(),
		vehicle:     vehicle,
	}
}

// Getters
func (e *VehicleCreatedEvent) EventID() string     { return e.eventID }
func (e *VehicleCreatedEvent) EventType() string   { return e.eventType }
func (e *VehicleCreatedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e *VehicleCreatedEvent) AggregateID() string { return e.aggregateID }
func (e *VehicleCreatedEvent) Vehicle() Vehicle    { return e.vehicle }

// Data 返回事件数据
func (e *VehicleCreatedEvent) Data() interface{} {
	return map[string]interface{}{
		"vehicle": e.vehicle,
	}
}

// RouteCompletedEvent 路线完成事件
type RouteCompletedEvent struct {
	eventID     string
	eventType   string
	occurredOn  time.Time
	aggregateID string
	route       Route
}

// NewRouteCompletedEvent 创建路线完成事件
func NewRouteCompletedEvent(route Route) *RouteCompletedEvent {
	return &RouteCompletedEvent{
		eventID:     generateEventID(),
		eventType:   "RouteCompleted",
		occurredOn:  time.Now(),
		aggregateID: route.ID(),
		route:       route,
	}
}

// Getters
func (e *RouteCompletedEvent) EventID() string     { return e.eventID }
func (e *RouteCompletedEvent) EventType() string   { return e.eventType }
func (e *RouteCompletedEvent) OccurredOn() time.Time { return e.occurredOn }
func (e *RouteCompletedEvent) AggregateID() string { return e.aggregateID }
func (e *RouteCompletedEvent) Route() Route        { return e.route }

// Data 返回事件数据
func (e *RouteCompletedEvent) Data() interface{} {
	return map[string]interface{}{
		"route": e.route,
	}
}

// generateEventID 生成事件ID（简化实现）
func generateEventID() string {
	return "evt_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// randomString 生成随机字符串（简化实现）
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}
