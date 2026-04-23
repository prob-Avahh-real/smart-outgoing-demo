package domain

import (
	"time"
)

// Metrics 聚合根 - 系统性能指标
type Metrics struct {
	id              string
	timestamp       time.Time
	memoryUsage     MemoryMetrics
	agentCount      int
	connectionCount int
	cpuUsage        float64
	storageStrategy StorageStrategy
	version         int
}

// MemoryMetrics 内存使用指标
type MemoryMetrics struct {
	TotalMB      int
	UsedMB       int
	AvailableMB  int
	UsagePercent float64
}

// StorageStrategy 存储策略值对象
type StorageStrategy int

const (
	StorageStrategyMemory StorageStrategy = iota
	StorageStrategyRedis
	StorageStrategyHybrid
)

// ScalingThreshold 扩容阈值值对象
type ScalingThreshold struct {
	MemoryUsagePercent float64
	AgentCount         int
	ConnectionCount    int
	CPUUsagePercent    float64
}

// DefaultThresholds 默认阈值配置
func DefaultThresholds() ScalingThreshold {
	return ScalingThreshold{
		MemoryUsagePercent: 80.0,
		AgentCount:         100,
		ConnectionCount:    1000,
		CPUUsagePercent:    75.0,
	}
}

// NewMetrics 创建新的指标聚合
func NewMetrics(id string, memoryUsage MemoryMetrics, agentCount, connectionCount int, cpuUsage float64, strategy StorageStrategy) *Metrics {
	return &Metrics{
		id:              id,
		timestamp:       time.Now(),
		memoryUsage:     memoryUsage,
		agentCount:      agentCount,
		connectionCount: connectionCount,
		cpuUsage:        cpuUsage,
		storageStrategy: strategy,
		version:         1,
	}
}

// Getters
func (m *Metrics) ID() string                       { return m.id }
func (m *Metrics) Timestamp() time.Time             { return m.timestamp }
func (m *Metrics) MemoryUsage() MemoryMetrics       { return m.memoryUsage }
func (m *Metrics) AgentCount() int                  { return m.agentCount }
func (m *Metrics) ConnectionCount() int             { return m.connectionCount }
func (m *Metrics) CPUUsage() float64                { return m.cpuUsage }
func (m *Metrics) StorageStrategy() StorageStrategy { return m.storageStrategy }
func (m *Metrics) Version() int                     { return m.version }

// ShouldScaleUp 检查是否应该扩容
func (m *Metrics) ShouldScaleUp(threshold ScalingThreshold) bool {
	return m.memoryUsage.UsagePercent > threshold.MemoryUsagePercent ||
		m.agentCount > threshold.AgentCount ||
		m.connectionCount > threshold.ConnectionCount ||
		m.cpuUsage > threshold.CPUUsagePercent
}

// ShouldScaleDown 检查是否应该缩容
func (m *Metrics) ShouldScaleDown(threshold ScalingThreshold) bool {
	// 缩容阈值更保守，避免频繁切换
	downThreshold := ScalingThreshold{
		MemoryUsagePercent: threshold.MemoryUsagePercent * 0.6,
		AgentCount:         int(float64(threshold.AgentCount) * 0.6),
		ConnectionCount:    int(float64(threshold.ConnectionCount) * 0.6),
		CPUUsagePercent:    threshold.CPUUsagePercent * 0.6,
	}

	return m.memoryUsage.UsagePercent < downThreshold.MemoryUsagePercent &&
		m.agentCount < downThreshold.AgentCount &&
		m.connectionCount < downThreshold.ConnectionCount &&
		m.cpuUsage < downThreshold.CPUUsagePercent
}

// UpdateStorageStrategy 更新存储策略
func (m *Metrics) UpdateStorageStrategy(strategy StorageStrategy) {
	m.storageStrategy = strategy
	m.timestamp = time.Now()
	m.version++
}

// UpdateMetrics 更新指标数据
func (m *Metrics) UpdateMetrics(memoryUsage MemoryMetrics, agentCount, connectionCount int, cpuUsage float64) {
	m.memoryUsage = memoryUsage
	m.agentCount = agentCount
	m.connectionCount = connectionCount
	m.cpuUsage = cpuUsage
	m.timestamp = time.Now()
	m.version++
}

// String 返回存储策略的字符串表示
func (s StorageStrategy) String() string {
	switch s {
	case StorageStrategyMemory:
		return "memory"
	case StorageStrategyRedis:
		return "redis"
	case StorageStrategyHybrid:
		return "hybrid"
	default:
		return "unknown"
	}
}
