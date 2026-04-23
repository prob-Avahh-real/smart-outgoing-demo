package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// DDDScalingDecisionService 扩容决策领域服务
type DDDScalingDecisionService struct {
	metricsRepo         domain.MetricsRepository
	scalingDecisionRepo domain.ScalingDecisionRepository
	eventBus            EventBus
	threshold           domain.ScalingThreshold
	evaluationInterval  time.Duration
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(event domain.DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event domain.DomainEvent) error
}

// NewDDDMScalingDecisionService 创建扩容决策服务
func NewDDDMScalingDecisionService(
	metricsRepo domain.MetricsRepository,
	scalingDecisionRepo domain.ScalingDecisionRepository,
	eventBus EventBus,
	threshold domain.ScalingThreshold,
) *DDDScalingDecisionService {
	return &DDDScalingDecisionService{
		metricsRepo:         metricsRepo,
		scalingDecisionRepo: scalingDecisionRepo,
		eventBus:            eventBus,
		threshold:           threshold,
		evaluationInterval:  30 * time.Second,
	}
}

// EvaluateScaling 评估是否需要扩容
func (s *DDDScalingDecisionService) EvaluateScaling(ctx context.Context) (*domain.ScalingDecision, error) {
	// 获取当前指标
	currentMetrics, err := s.metricsRepo.FindLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current metrics: %w", err)
	}

	// 检查是否需要扩容
	if currentMetrics.ShouldScaleUp(s.threshold) {
		return s.createScaleUpDecision(ctx, currentMetrics)
	}

	// 检查是否需要缩容
	if currentMetrics.ShouldScaleDown(s.threshold) {
		return s.createScaleDownDecision(ctx, currentMetrics)
	}

	return nil, nil // 无需扩容
}

// createScaleUpDecision 创建扩容决策
func (s *DDDScalingDecisionService) createScaleUpDecision(ctx context.Context, currentMetrics *domain.Metrics) (*domain.ScalingDecision, error) {
	currentStrategy := currentMetrics.StorageStrategy()
	targetStrategy := s.determineTargetStrategy(currentMetrics, true)

	if currentStrategy == targetStrategy {
		return nil, nil // 已经是目标策略
	}

	// 确定扩容原因
	reason := s.determineScalingReason(currentMetrics, s.threshold)

	// 创建扩容决策
	decisionID := fmt.Sprintf("scale_up_%d", time.Now().Unix())
	decision := domain.NewScalingDecision(
		decisionID,
		currentStrategy,
		targetStrategy,
		reason,
		s.threshold,
		*currentMetrics,
	)

	// 保存决策
	if err := s.scalingDecisionRepo.Save(ctx, decision); err != nil {
		return nil, fmt.Errorf("failed to save scaling decision: %w", err)
	}

	// 发布扩容事件
	event := domain.NewStorageStrategyChangedEvent(
		decisionID,
		currentStrategy,
		targetStrategy,
		reason,
		*currentMetrics,
	)

	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("Warning: failed to publish scaling event: %v", err)
	}

	log.Printf("Scale-up decision created: %s -> %s, reason: %s",
		currentStrategy.String(), targetStrategy.String(), reason.String())

	return decision, nil
}

// createScaleDownDecision 创建缩容决策
func (s *DDDScalingDecisionService) createScaleDownDecision(ctx context.Context, currentMetrics *domain.Metrics) (*domain.ScalingDecision, error) {
	currentStrategy := currentMetrics.StorageStrategy()
	targetStrategy := s.determineTargetStrategy(currentMetrics, false)

	if currentStrategy == targetStrategy {
		return nil, nil // 已经是目标策略
	}

	// 确定缩容原因
	reason := domain.ScalingReasonManual // 缩容通常是手动或自动优化

	// 创建缩容决策
	decisionID := fmt.Sprintf("scale_down_%d", time.Now().Unix())
	decision := domain.NewScalingDecision(
		decisionID,
		currentStrategy,
		targetStrategy,
		reason,
		s.threshold,
		*currentMetrics,
	)

	// 保存决策
	if err := s.scalingDecisionRepo.Save(ctx, decision); err != nil {
		return nil, fmt.Errorf("failed to save scaling decision: %w", err)
	}

	// 发布缩容事件
	event := domain.NewStorageStrategyChangedEvent(
		decisionID,
		currentStrategy,
		targetStrategy,
		reason,
		*currentMetrics,
	)

	if err := s.eventBus.Publish(event); err != nil {
		log.Printf("Warning: failed to publish scaling event: %v", err)
	}

	log.Printf("Scale-down decision created: %s -> %s",
		currentStrategy.String(), targetStrategy.String())

	return decision, nil
}

// determineTargetStrategy 确定目标存储策略
func (s *DDDScalingDecisionService) determineTargetStrategy(metrics *domain.Metrics, scaleUp bool) domain.StorageStrategy {
	currentStrategy := metrics.StorageStrategy()

	if scaleUp {
		// 扩容：从内存转向Redis
		if currentStrategy == domain.StorageStrategyMemory {
			return domain.StorageStrategyRedis
		} else if currentStrategy == domain.StorageStrategyRedis {
			return domain.StorageStrategyHybrid
		}
		return domain.StorageStrategyRedis
	} else {
		// 缩容：从Redis转向内存
		if currentStrategy == domain.StorageStrategyHybrid {
			return domain.StorageStrategyRedis
		} else if currentStrategy == domain.StorageStrategyRedis {
			return domain.StorageStrategyMemory
		}
		return domain.StorageStrategyMemory
	}
}

// determineScalingReason 确定扩容原因
func (s *DDDScalingDecisionService) determineScalingReason(metrics *domain.Metrics, threshold domain.ScalingThreshold) domain.ScalingReason {
	memoryUsage := metrics.MemoryUsage()

	if memoryUsage.UsagePercent > threshold.MemoryUsagePercent {
		return domain.ScalingReasonMemory
	}

	if metrics.AgentCount() > threshold.AgentCount {
		return domain.ScalingReasonAgentCount
	}

	if metrics.ConnectionCount() > threshold.ConnectionCount {
		return domain.ScalingReasonConnectionCount
	}

	if metrics.CPUUsage() > threshold.CPUUsagePercent {
		return domain.ScalingReasonCPUUsage
	}

	return domain.ScalingReasonManual
}

// StartMonitoring 开始监控指标
func (s *DDDScalingDecisionService) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(s.evaluationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.evaluateAndPublishThresholdEvents(ctx); err != nil {
				log.Printf("Error during scaling evaluation: %v", err)
			}
		}
	}
}

// evaluateAndPublishThresholdEvents 评估并发布阈值突破事件
func (s *DDDScalingDecisionService) evaluateAndPublishThresholdEvents(ctx context.Context) error {
	currentMetrics, err := s.metricsRepo.FindLatest(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current metrics: %w", err)
	}

	if currentMetrics.ShouldScaleUp(s.threshold) {
		breachedBy := s.getBreachedThresholds(currentMetrics, s.threshold)

		event := domain.NewThresholdBreachedEvent(
			currentMetrics.ID(),
			s.threshold,
			*currentMetrics,
			breachedBy,
		)

		if err := s.eventBus.Publish(event); err != nil {
			return fmt.Errorf("failed to publish threshold breached event: %w", err)
		}
	}

	return nil
}

// getBreachedThresholds 获取突破的阈值列表
func (s *DDDScalingDecisionService) getBreachedThresholds(metrics *domain.Metrics, threshold domain.ScalingThreshold) []string {
	var breached []string

	if metrics.MemoryUsage().UsagePercent > threshold.MemoryUsagePercent {
		breached = append(breached, "memory_usage")
	}

	if metrics.AgentCount() > threshold.AgentCount {
		breached = append(breached, "agent_count")
	}

	if metrics.ConnectionCount() > threshold.ConnectionCount {
		breached = append(breached, "connection_count")
	}

	if metrics.CPUUsage() > threshold.CPUUsagePercent {
		breached = append(breached, "cpu_usage")
	}

	return breached
}

// UpdateThreshold 更新阈值配置
func (s *DDDScalingDecisionService) UpdateThreshold(newThreshold domain.ScalingThreshold) {
	s.threshold = newThreshold
	log.Printf("Scaling threshold updated: %+v", newThreshold)
}
