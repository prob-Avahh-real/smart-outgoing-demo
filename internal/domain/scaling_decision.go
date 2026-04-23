package domain

import (
	"time"
)

// ScalingDecision 聚合根 - 扩容策略决策
type ScalingDecision struct {
	id              string
	decisionTime    time.Time
	fromStrategy    StorageStrategy
	toStrategy      StorageStrategy
	reason          ScalingReason
	threshold       ScalingThreshold
	currentMetrics  Metrics
	executed        bool
	executedAt      *time.Time
	version         int
}

// ScalingReason 扩容原因
type ScalingReason int

const (
	ScalingReasonMemory ScalingReason = iota
	ScalingReasonAgentCount
	ScalingReasonConnectionCount
	ScalingReasonCPUUsage
	ScalingReasonManual
)

// NewScalingDecision 创建新的扩容决策
func NewScalingDecision(id string, fromStrategy, toStrategy StorageStrategy, reason ScalingReason, threshold ScalingThreshold, currentMetrics Metrics) *ScalingDecision {
	return &ScalingDecision{
		id:             id,
		decisionTime:   time.Now(),
		fromStrategy:   fromStrategy,
		toStrategy:     toStrategy,
		reason:         reason,
		threshold:      threshold,
		currentMetrics: currentMetrics,
		executed:       false,
		version:        1,
	}
}

// Getters
func (sd *ScalingDecision) ID() string               { return sd.id }
func (sd *ScalingDecision) DecisionTime() time.Time   { return sd.decisionTime }
func (sd *ScalingDecision) FromStrategy() StorageStrategy { return sd.fromStrategy }
func (sd *ScalingDecision) ToStrategy() StorageStrategy   { return sd.toStrategy }
func (sd *ScalingDecision) Reason() ScalingReason     { return sd.reason }
func (sd *ScalingDecision) Threshold() ScalingThreshold { return sd.threshold }
func (sd *ScalingDecision) CurrentMetrics() Metrics   { return sd.currentMetrics }
func (sd *ScalingDecision) Executed() bool           { return sd.executed }
func (sd *ScalingDecision) ExecutedAt() *time.Time    { return sd.executedAt }
func (sd *ScalingDecision) Version() int             { return sd.version }

// Execute 执行扩容决策
func (sd *ScalingDecision) Execute() {
	if !sd.executed {
		sd.executed = true
		now := time.Now()
		sd.executedAt = &now
		sd.version++
	}
}

// String 返回扩容原因的字符串表示
func (sr ScalingReason) String() string {
	switch sr {
	case ScalingReasonMemory:
		return "memory_usage"
	case ScalingReasonAgentCount:
		return "agent_count"
	case ScalingReasonConnectionCount:
		return "connection_count"
	case ScalingReasonCPUUsage:
		return "cpu_usage"
	case ScalingReasonManual:
		return "manual"
	default:
		return "unknown"
	}
}
