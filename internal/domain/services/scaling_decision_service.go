package services

import (
	"fmt"
	"runtime"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/events"
	"smart-outgoing-demo/internal/domain/repositories"
)

// ScalingDecisionService handles scaling decisions based on system metrics
type ScalingDecisionService struct {
	decisionRepo repositories.ScalingDecisionRepository
	metricsRepo  repositories.MetricsRepository
	eventBus     chan events.DomainEvent
}

// NewScalingDecisionService creates a new scaling decision service
func NewScalingDecisionService(decisionRepo repositories.ScalingDecisionRepository, metricsRepo repositories.MetricsRepository, eventBus chan events.DomainEvent) *ScalingDecisionService {
	return &ScalingDecisionService{
		decisionRepo: decisionRepo,
		metricsRepo:  metricsRepo,
		eventBus:     eventBus,
	}
}

// EvaluateScaling evaluates current metrics and decides if scaling is needed
func (s *ScalingDecisionService) EvaluateScaling() (*entities.ScalingDecision, error) {
	// Get latest metrics
	latestMetrics, err := s.metricsRepo.FindLatest()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest metrics: %w", err)
	}

	// Get current scaling decision
	currentDecision, err := s.decisionRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("failed to get current scaling decision: %w", err)
	}

	// Check if scaling to Redis is needed
	if currentDecision.ShouldScaleToRedis(latestMetrics) {
		return s.scaleToRedis(currentDecision, latestMetrics)
	}

	// Check if scaling to memory is needed
	if currentDecision.ShouldScaleToMemory(latestMetrics) {
		return s.scaleToMemory(currentDecision, latestMetrics)
	}

	return currentDecision, nil
}

// GetCurrentMetrics collects current system metrics
func (s *ScalingDecisionService) GetCurrentMetrics(storageStrategy string) (*entities.Metrics, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get current agent count and WebSocket connections from repositories
	// For now, we'll use placeholder values that would be collected from actual sources
	activeAgents := 0
	webSocketConns := 0

	// Try to get vehicle count as agent count
	// This would be injected or passed in a real implementation
	// For now, we'll use a simple approach

	metrics := entities.NewMetrics(
		m.Alloc,        // Memory usage in bytes
		0,              // CPU usage (would need actual CPU monitoring)
		activeAgents,   // Active agents
		webSocketConns, // WebSocket connections
		storageStrategy,
	)

	// Save metrics
	err := s.metricsRepo.Save(metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to save metrics: %w", err)
	}

	return metrics, nil
}

// scaleToRedis switches storage strategy to Redis
func (s *ScalingDecisionService) scaleToRedis(decision *entities.ScalingDecision, metrics *entities.Metrics) (*entities.ScalingDecision, error) {
	oldStrategy := string(decision.CurrentStrategy)

	decision.SwitchToRedis(fmt.Sprintf("Thresholds breached: Memory: %.2fMB, Agents: %d, Connections: %d",
		float64(metrics.MemoryUsage)/1024/1024, metrics.ActiveAgents, metrics.WebSocketConns))

	// Save the updated decision
	err := s.decisionRepo.Save(decision)
	if err != nil {
		return nil, fmt.Errorf("failed to save scaling decision: %w", err)
	}

	// Publish domain event
	event := events.NewStorageStrategyChanged(
		decision.ID,
		oldStrategy,
		string(entities.StorageStrategyRedis),
		decision.Reason,
		metrics,
	)

	if s.eventBus != nil {
		s.eventBus <- event
	}

	return decision, nil
}

// scaleToMemory switches storage strategy to memory
func (s *ScalingDecisionService) scaleToMemory(decision *entities.ScalingDecision, metrics *entities.Metrics) (*entities.ScalingDecision, error) {
	oldStrategy := string(decision.CurrentStrategy)

	decision.SwitchToMemory(fmt.Sprintf("Load reduced: Memory: %.2fMB, Agents: %d, Connections: %d",
		float64(metrics.MemoryUsage)/1024/1024, metrics.ActiveAgents, metrics.WebSocketConns))

	// Save the updated decision
	err := s.decisionRepo.Save(decision)
	if err != nil {
		return nil, fmt.Errorf("failed to save scaling decision: %w", err)
	}

	// Publish domain event
	event := events.NewStorageStrategyChanged(
		decision.ID,
		oldStrategy,
		string(entities.StorageStrategyMemory),
		decision.Reason,
		metrics,
	)

	if s.eventBus != nil {
		s.eventBus <- event
	}

	return decision, nil
}

// UpdateThresholds updates the scaling thresholds
func (s *ScalingDecisionService) UpdateThresholds(decision *entities.ScalingDecision, thresholds entities.ScalingThreshold) error {
	decision.MemoryThreshold = thresholds.MemoryUsage
	decision.CPUThreshold = thresholds.CPUUsage
	decision.AgentThreshold = thresholds.ActiveAgents
	decision.ConnectionThreshold = thresholds.WebSocketConns

	return s.decisionRepo.Save(decision)
}

// EnableAutoSwitch enables automatic switching
func (s *ScalingDecisionService) EnableAutoSwitch(decision *entities.ScalingDecision) error {
	decision.EnableAutoSwitch()
	return s.decisionRepo.Save(decision)
}

// DisableAutoSwitch disables automatic switching
func (s *ScalingDecisionService) DisableAutoSwitch(decision *entities.ScalingDecision) error {
	decision.DisableAutoSwitch()
	return s.decisionRepo.Save(decision)
}

// GetDecisionRepository returns the decision repository (for orchestrator access)
func (s *ScalingDecisionService) GetDecisionRepository() repositories.ScalingDecisionRepository {
	return s.decisionRepo
}
