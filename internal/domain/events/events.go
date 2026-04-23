package events

import (
	"fmt"
	"time"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	ID() string
	AggregateID() string
	EventType() string
	OccurredAt() time.Time
	Data() interface{}
}

// StorageStrategyChanged represents a domain event when storage strategy changes
type StorageStrategyChanged struct {
	id           string
	aggregateID  string
	eventType    string
	occurredAt   time.Time
	FromStrategy string      `json:"from_strategy"`
	ToStrategy   string      `json:"to_strategy"`
	Reason       string      `json:"reason"`
	Metrics      interface{} `json:"metrics"`
}

// NewStorageStrategyChanged creates a new StorageStrategyChanged event
func NewStorageStrategyChanged(aggregateID, fromStrategy, toStrategy, reason string, metrics interface{}) *StorageStrategyChanged {
	return &StorageStrategyChanged{
		id:           generateEventID(),
		aggregateID:  aggregateID,
		eventType:    "StorageStrategyChanged",
		occurredAt:   time.Now(),
		FromStrategy: fromStrategy,
		ToStrategy:   toStrategy,
		Reason:       reason,
		Metrics:      metrics,
	}
}

// ID returns the event ID
func (e *StorageStrategyChanged) ID() string {
	return e.id
}

// AggregateID returns the aggregate ID
func (e *StorageStrategyChanged) AggregateID() string {
	return e.aggregateID
}

// EventType returns the event type
func (e *StorageStrategyChanged) EventType() string {
	return e.eventType
}

// OccurredAt returns when the event occurred
func (e *StorageStrategyChanged) OccurredAt() time.Time {
	return e.occurredAt
}

// Data returns the event data
func (e *StorageStrategyChanged) Data() interface{} {
	return e
}

// ThresholdBreached represents a domain event when system thresholds are breached
type ThresholdBreached struct {
	id            string
	aggregateID   string
	eventType     string
	occurredAt    time.Time
	Thresholds    map[string]interface{} `json:"thresholds"`
	CurrentValues map[string]interface{} `json:"current_values"`
	BreachedType  string                 `json:"breached_type"` // "memory", "cpu", "agents", "connections"
}

// NewThresholdBreached creates a new ThresholdBreached event
func NewThresholdBreached(aggregateID string, thresholds, currentValues map[string]interface{}, breachedType string) *ThresholdBreached {
	return &ThresholdBreached{
		id:            generateEventID(),
		aggregateID:   aggregateID,
		eventType:     "ThresholdBreached",
		occurredAt:    time.Now(),
		Thresholds:    thresholds,
		CurrentValues: currentValues,
		BreachedType:  breachedType,
	}
}

// ID returns the event ID
func (e *ThresholdBreached) ID() string {
	return e.id
}

// AggregateID returns the aggregate ID
func (e *ThresholdBreached) AggregateID() string {
	return e.aggregateID
}

// EventType returns the event type
func (e *ThresholdBreached) EventType() string {
	return e.eventType
}

// OccurredAt returns when the event occurred
func (e *ThresholdBreached) OccurredAt() time.Time {
	return e.occurredAt
}

// Data returns the event data
func (e *ThresholdBreached) Data() interface{} {
	return e
}

// generateEventID creates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("event_%d", time.Now().UnixNano())
}
