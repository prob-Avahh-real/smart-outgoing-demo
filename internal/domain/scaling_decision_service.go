package domain

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ScalingDecisionService 扩容决策领域服务
type ScalingDecisionService struct {
	metricsRepo         MetricsRepository
	scalingDecisionRepo ScalingDecisionRepository
	eventBus            EventBus
	threshold           ScalingThreshold
	evaluationInterval  time.Duration
}

// NewScalingDecisionService 创建扩容决策服务
func NewScalingDecisionService(
	metricsRepo MetricsRepository,
	scalingDecisionRepo ScalingDecisionRepository,
	eventBus EventBus,
	threshold ScalingThreshold,
) *ScalingDecisionService {
	return &ScalingDecisionService{
		metricsRepo:         metricsRepo,
		scalingDecisionRepo: scalingDecisionRepo,
		eventBus:            eventBus,
		threshold:           threshold,
		evaluationInterval:  30 * time.Second,
	}
}

// EvaluateScaling 评估是否需要扩容
func (s *ScalingDecisionService) EvaluateScaling(ctx context.Context) (*ScalingDecision, error) {
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
func (s *ScalingDecisionService) createScaleUpDecision(ctx context.Context, currentMetrics *Metrics) (*ScalingDecision, error) {
	currentStrategy := currentMetrics.StorageStrategy()
	targetStrategy := s.determineTargetStrategy(currentMetrics, true)
	
	if currentStrategy == targetStrategy {
		return nil, nil // 已经是目标策略
	}

	// 确定扩容原因
	reason := s.determineScalingReason(currentMetrics, s.threshold)

	// 创建扩容决策
	decisionID := fmt.Sprintf("scale_up_%d", time.Now().Unix())
	decision := NewScalingDecision(
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
	event := NewStorageStrategyChangedEvent(
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
func (s *ScalingDecisionService) createScaleDownDecision(ctx context.Context, currentMetrics *Metrics) (*ScalingDecision, error) {
	currentStrategy := currentMetrics.StorageStrategy()
	targetStrategy := s.determineTargetStrategy(currentMetrics, false)
	
	if currentStrategy == targetStrategy {
		return nil, nil // 已经是目标策略
	}

	// 确定缩容原因
	reason := ScalingReasonManual // 缩容通常是手动或自动优化

	// 创建缩容决策
	decisionID := fmt.Sprintf("scale_down_%d", time.Now().Unix())
	decision := NewScalingDecision(
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
	event := NewStorageStrategyChangedEvent(
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
func (s *ScalingDecisionService) determineTargetStrategy(metrics *Metrics, scaleUp bool) StorageStrategy {
	currentStrategy := metrics.StorageStrategy()

	if scaleUp {
		// 扩容：从内存转向Redis
		if currentStrategy == StorageStrategyMemory {
			return StorageStrategyRedis
		} else if currentStrategy == StorageStrategyRedis {
			return StorageStrategyHybrid
		}
		return StorageStrategyRedis
	} else {
		// 缩容：从Redis转向内存
		if currentStrategy == StorageStrategyHybrid {
			return StorageStrategyRedis
		} else if currentStrategy == StorageStrategyRedis {
			return StorageStrategyMemory
		}
		return StorageStrategyMemory
	}
}

// determineScalingReason 确定扩容原因
func (s *ScalingDecisionService) determineScalingReason(metrics *Metrics, threshold ScalingThreshold) ScalingReason {
	memoryUsage := metrics.MemoryUsage()
	
	if memoryUsage.UsagePercent > threshold.MemoryUsagePercent {
		return ScalingReasonMemory
	}
	
	if metrics.AgentCount() > threshold.AgentCount {
		return ScalingReasonAgentCount
	}
	
	if metrics.ConnectionCount() > threshold.ConnectionCount {
		return ScalingReasonConnectionCount
	}
	
	if metrics.CPUUsage() > threshold.CPUUsagePercent {
		return ScalingReasonCPUUsage
	}
	
	return ScalingReasonManual
}

// StartMonitoring 开始监控指标
func (s *ScalingDecisionService) StartMonitoring(ctx context.Context) {
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
func (s *ScalingDecisionService) evaluateAndPublishThresholdEvents(ctx context.Context) error {
	currentMetrics, err := s.metricsRepo.FindLatest(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current metrics: %w", err)
	}

	if currentMetrics.ShouldScaleUp(s.threshold) {
		breachedBy := s.getBreachedThresholds(currentMetrics, s.threshold)
		
		event := NewThresholdBreachedEvent(
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
func (s *ScalingDecisionService) getBreachedThresholds(metrics *Metrics, threshold ScalingThreshold) []string {
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
func (s *ScalingDecisionService) UpdateThreshold(newThreshold ScalingThreshold) {
	s.threshold = newThreshold
	log.Printf("Scaling threshold updated: %+v", newThreshold)
}
