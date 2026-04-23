package entities

import (
	"fmt"
	"time"
)

// ScalingDecision represents the scaling decision aggregate
type ScalingDecision struct {
	ID                   string           `json:"id"`
	CurrentStrategy      StorageStrategy `json:"current_strategy"`
	PreviousStrategy     StorageStrategy `json:"previous_strategy"`
	MemoryThreshold      uint64           `json:"memory_threshold"`      // bytes
	CPUThreshold         float64          `json:"cpu_threshold"`         // percentage
	AgentThreshold       int              `json:"agent_threshold"`        // count
	ConnectionThreshold  int              `json:"connection_threshold"`   // count
	DecisionAt           time.Time        `json:"decision_at"`
	Reason               string           `json:"reason"`
	AutoSwitchEnabled    bool             `json:"auto_switch_enabled"`
}

// StorageStrategy represents the storage strategy value object
type StorageStrategy string

const (
	StorageStrategyMemory StorageStrategy = "memory"
	StorageStrategyRedis  StorageStrategy = "redis"
)

// ScalingThreshold represents scaling threshold configuration
type ScalingThreshold struct {
	MemoryUsage     uint64  `json:"memory_usage"`     // bytes
	CPUUsage        float64 `json:"cpu_usage"`        // percentage
	ActiveAgents    int     `json:"active_agents"`
	WebSocketConns  int     `json:"websocket_conns"`
}

// NewScalingDecision creates a new scaling decision aggregate
func NewScalingDecision(strategy StorageStrategy, thresholds ScalingThreshold) *ScalingDecision {
	return &ScalingDecision{
		ID:                  generateDecisionID(),
		CurrentStrategy:     strategy,
		PreviousStrategy:    strategy, // initially same as current
		MemoryThreshold:     thresholds.MemoryUsage,
		CPUThreshold:        thresholds.CPUUsage,
		AgentThreshold:      thresholds.ActiveAgents,
		ConnectionThreshold: thresholds.WebSocketConns,
		DecisionAt:          time.Now(),
		Reason:              "Initial configuration",
		AutoSwitchEnabled:   true,
	}
}

// ShouldScaleToMemory determines if scaling to memory is appropriate
func (sd *ScalingDecision) ShouldScaleToMemory(metrics *Metrics) bool {
	if !sd.AutoSwitchEnabled || sd.CurrentStrategy == StorageStrategyMemory {
		return false
	}
	
	return metrics.MemoryUsage < sd.MemoryThreshold*8/10 && // 80% of threshold
		metrics.CPUUsage < sd.CPUThreshold*8/10 &&
		metrics.ActiveAgents < sd.AgentThreshold*8/10 &&
		metrics.WebSocketConns < sd.ConnectionThreshold*8/10
}

// ShouldScaleToRedis determines if scaling to Redis is appropriate
func (sd *ScalingDecision) ShouldScaleToRedis(metrics *Metrics) bool {
	if !sd.AutoSwitchEnabled || sd.CurrentStrategy == StorageStrategyRedis {
		return false
	}
	
	return metrics.MemoryUsage > sd.MemoryThreshold ||
		metrics.CPUUsage > sd.CPUThreshold ||
		metrics.ActiveAgents > sd.AgentThreshold ||
		metrics.WebSocketConns > sd.ConnectionThreshold
}

// SwitchToMemory switches the storage strategy to memory
func (sd *ScalingDecision) SwitchToMemory(reason string) {
	if sd.CurrentStrategy == StorageStrategyMemory {
		return
	}
	
	sd.PreviousStrategy = sd.CurrentStrategy
	sd.CurrentStrategy = StorageStrategyMemory
	sd.DecisionAt = time.Now()
	sd.Reason = reason
}

// SwitchToRedis switches the storage strategy to Redis
func (sd *ScalingDecision) SwitchToRedis(reason string) {
	if sd.CurrentStrategy == StorageStrategyRedis {
		return
	}
	
	sd.PreviousStrategy = sd.CurrentStrategy
	sd.CurrentStrategy = StorageStrategyRedis
	sd.DecisionAt = time.Now()
	sd.Reason = reason
}

// EnableAutoSwitch enables automatic switching
func (sd *ScalingDecision) EnableAutoSwitch() {
	sd.AutoSwitchEnabled = true
	sd.DecisionAt = time.Now()
	sd.Reason = "Auto-switch enabled"
}

// DisableAutoSwitch disables automatic switching
func (sd *ScalingDecision) DisableAutoSwitch() {
	sd.AutoSwitchEnabled = false
	sd.DecisionAt = time.Now()
	sd.Reason = "Auto-switch disabled"
}

// generateDecisionID creates a unique decision ID
func generateDecisionID() string {
	return fmt.Sprintf("decision_%d", time.Now().UnixNano())
}
