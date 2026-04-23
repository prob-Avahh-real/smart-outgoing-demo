package entities

import (
	"fmt"
	"time"
)

// Metrics represents the metrics aggregate
type Metrics struct {
	ID              string    `json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	MemoryUsage     uint64    `json:"memory_usage"` // in bytes
	CPUUsage        float64   `json:"cpu_usage"`    // percentage
	ActiveAgents    int       `json:"active_agents"`
	WebSocketConns  int       `json:"websocket_conns"`
	StorageStrategy string    `json:"storage_strategy"` // "memory" or "redis"
}

// NewMetrics creates a new metrics aggregate
func NewMetrics(memoryUsage uint64, cpuUsage float64, activeAgents, webSocketConns int, storageStrategy string) *Metrics {
	return &Metrics{
		ID:              generateMetricsID(),
		Timestamp:       time.Now(),
		MemoryUsage:     memoryUsage,
		CPUUsage:        cpuUsage,
		ActiveAgents:    activeAgents,
		WebSocketConns:  webSocketConns,
		StorageStrategy: storageStrategy,
	}
}

// Update updates the metrics values
func (m *Metrics) Update(memoryUsage uint64, cpuUsage float64, activeAgents, webSocketConns int, storageStrategy string) {
	m.Timestamp = time.Now()
	m.MemoryUsage = memoryUsage
	m.CPUUsage = cpuUsage
	m.ActiveAgents = activeAgents
	m.WebSocketConns = webSocketConns
	m.StorageStrategy = storageStrategy
}

// generateMetricsID creates a unique metrics ID
func generateMetricsID() string {
	return fmt.Sprintf("metrics_%d", time.Now().UnixNano())
}
