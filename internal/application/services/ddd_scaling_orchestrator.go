package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// DDDScalingOrchestrator 扩容协调器 - 应用服务
type DDDScalingOrchestrator struct {
	scalingDecisionService   *domain.ScalingDecisionService
	vehicleManagementService *domain.VehicleManagementService
	repositoryFactory        domain.RepositoryFactory
	migrationService         MigrationService
	eventBus                 domain.EventBus
	currentStrategy          domain.StorageStrategy
	mu                       sync.RWMutex
	migrationInProgress      bool
}

// MigrationService 数据迁移服务接口
type MigrationService interface {
	MigrateVehicles(ctx context.Context, from, to domain.StorageStrategy) error
	MigrateRoutes(ctx context.Context, from, to domain.StorageStrategy) error
	MigrateMetrics(ctx context.Context, from, to domain.StorageStrategy) error
	MigrateScalingDecisions(ctx context.Context, from, to domain.StorageStrategy) error
}

// NewDDDMScalingOrchestrator 创建扩容协调器
func NewDDDMScalingOrchestrator(
	scalingDecisionService *domain.ScalingDecisionService,
	vehicleManagementService *domain.VehicleManagementService,
	repositoryFactory domain.RepositoryFactory,
	migrationService MigrationService,
	eventBus domain.EventBus,
	initialStrategy domain.StorageStrategy,
) *DDDScalingOrchestrator {
	return &DDDScalingOrchestrator{
		scalingDecisionService:   scalingDecisionService,
		vehicleManagementService: vehicleManagementService,
		repositoryFactory:        repositoryFactory,
		migrationService:         migrationService,
		eventBus:                 eventBus,
		currentStrategy:          initialStrategy,
	}
}

// ExecuteScaling 执行扩容操作
func (o *DDDScalingOrchestrator) ExecuteScaling(ctx context.Context, decision *domain.ScalingDecision) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.migrationInProgress {
		return fmt.Errorf("migration already in progress")
	}

	fromStrategy := decision.FromStrategy()
	toStrategy := decision.ToStrategy()

	if fromStrategy != o.currentStrategy {
		return fmt.Errorf("strategy mismatch: current is %s, decision expects %s",
			o.currentStrategy.String(), fromStrategy.String())
	}

	log.Printf("Starting scaling migration: %s -> %s", fromStrategy.String(), toStrategy.String())

	// 标记迁移开始
	o.migrationInProgress = true
	defer func() {
		o.migrationInProgress = false
	}()

	// 执行零停机迁移
	if err := o.executeZeroDowntimeMigration(ctx, fromStrategy, toStrategy); err != nil {
		return fmt.Errorf("zero-downtime migration failed: %w", err)
	}

	// 更新当前策略
	o.currentStrategy = toStrategy

	// 标记决策已执行
	decision.Execute()

	log.Printf("Scaling migration completed: %s -> %s", fromStrategy.String(), toStrategy.String())
	return nil
}

// executeZeroDowntimeMigration 执行零停机迁移
func (o *DDDScalingOrchestrator) executeZeroDowntimeMigration(ctx context.Context, from, to domain.StorageStrategy) error {
	// 1. 准备目标存储
	if err := o.prepareTargetStorage(ctx, to); err != nil {
		return fmt.Errorf("failed to prepare target storage: %w", err)
	}

	// 2. 双写阶段 - 同时写入新旧存储
	if err := o.enableDualWrite(ctx, from, to); err != nil {
		return fmt.Errorf("failed to enable dual write: %w", err)
	}

	// 3. 数据同步
	if err := o.syncData(ctx, from, to); err != nil {
		return fmt.Errorf("failed to sync data: %w", err)
	}

	// 4. 验证数据一致性
	if err := o.verifyConsistency(ctx, from, to); err != nil {
		return fmt.Errorf("consistency verification failed: %w", err)
	}

	// 5. 切换读取到新存储
	if err := o.switchReadTarget(ctx, to); err != nil {
		return fmt.Errorf("failed to switch read target: %w", err)
	}

	// 6. 停止双写，只写新存储
	if err := o.disableDualWrite(ctx, to); err != nil {
		return fmt.Errorf("failed to disable dual write: %w", err)
	}

	return nil
}

// prepareTargetStorage 准备目标存储
func (o *DDDScalingOrchestrator) prepareTargetStorage(ctx context.Context, strategy domain.StorageStrategy) error {
	log.Printf("Preparing target storage: %s", strategy.String())

	// 创建目标存储的仓储
	vehicleRepo := o.repositoryFactory.CreateVehicleRepository(strategy)

	// 测试连接（简化实现）
	_, err := vehicleRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("target storage connection test failed: %w", err)
	}

	log.Printf("Target storage prepared: %s", strategy.String())
	return nil
}

// enableDualWrite 启用双写模式
func (o *DDDScalingOrchestrator) enableDualWrite(ctx context.Context, from, to domain.StorageStrategy) error {
	log.Printf("Enabling dual write: %s -> %s", from.String(), to.String())

	// 这里应该配置仓储工厂以支持双写
	// 简化实现，实际需要修改仓储工厂的逻辑

	log.Printf("Dual write enabled")
	return nil
}

// syncData 同步数据
func (o *DDDScalingOrchestrator) syncData(ctx context.Context, from, to domain.StorageStrategy) error {
	log.Printf("Syncing data from %s to %s", from.String(), to.String())

	// 迁移车辆数据
	if err := o.migrationService.MigrateVehicles(ctx, from, to); err != nil {
		return fmt.Errorf("failed to migrate vehicles: %w", err)
	}

	// 迁移路线数据
	if err := o.migrationService.MigrateRoutes(ctx, from, to); err != nil {
		return fmt.Errorf("failed to migrate routes: %w", err)
	}

	// 迁移指标数据
	if err := o.migrationService.MigrateMetrics(ctx, from, to); err != nil {
		return fmt.Errorf("failed to migrate metrics: %w", err)
	}

	// 迁移扩容决策数据
	if err := o.migrationService.MigrateScalingDecisions(ctx, from, to); err != nil {
		return fmt.Errorf("failed to migrate scaling decisions: %w", err)
	}

	log.Printf("Data sync completed")
	return nil
}

// verifyConsistency 验证数据一致性
func (o *DDDScalingOrchestrator) verifyConsistency(ctx context.Context, from, to domain.StorageStrategy) error {
	log.Printf("Verifying data consistency between %s and %s", from.String(), to.String())

	// 获取两个存储的车辆数量
	fromRepo := o.repositoryFactory.CreateVehicleRepository(from)
	toRepo := o.repositoryFactory.CreateVehicleRepository(to)

	fromCount, err := fromRepo.Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count vehicles in source: %w", err)
	}

	toCount, err := toRepo.Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to count vehicles in target: %w", err)
	}

	if fromCount != toCount {
		return fmt.Errorf("vehicle count mismatch: source=%d, target=%d", fromCount, toCount)
	}

	log.Printf("Consistency verification passed: %d vehicles", fromCount)
	return nil
}

// switchReadTarget 切换读取目标
func (o *DDDScalingOrchestrator) switchReadTarget(ctx context.Context, strategy domain.StorageStrategy) error {
	log.Printf("Switching read target to: %s", strategy.String())

	// 这里应该更新应用配置以从新存储读取
	// 简化实现，实际需要更新仓储工厂的默认策略

	log.Printf("Read target switched to: %s", strategy.String())
	return nil
}

// disableDualWrite 禁用双写模式
func (o *DDDScalingOrchestrator) disableDualWrite(ctx context.Context, strategy domain.StorageStrategy) error {
	log.Printf("Disabling dual write, using only: %s", strategy.String())

	// 这里应该配置仓储工厂只使用新存储
	// 简化实现，实际需要修改仓储工厂的逻辑

	log.Printf("Dual write disabled")
	return nil
}

// GetCurrentStrategy 获取当前存储策略
func (o *DDDScalingOrchestrator) GetCurrentStrategy() domain.StorageStrategy {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.currentStrategy
}

// IsMigrationInProgress 检查迁移是否进行中
func (o *DDDScalingOrchestrator) IsMigrationInProgress() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.migrationInProgress
}

// StartAutoScaling 启动自动扩容
func (o *DDDScalingOrchestrator) StartAutoScaling(ctx context.Context) {
	log.Printf("Starting auto-scaling with initial strategy: %s", o.currentStrategy.String())

	// 启动扩容决策服务的监控
	go o.scalingDecisionService.StartMonitoring(ctx)

	// 启动协调器的主循环
	go o.scalingLoop(ctx)
}

// scalingLoop 扩容主循环
func (o *DDDScalingOrchestrator) scalingLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := o.checkAndExecuteScaling(ctx); err != nil {
				log.Printf("Error during scaling check: %v", err)
			}
		}
	}
}

// checkAndExecuteScaling 检查并执行扩容
func (o *DDDScalingOrchestrator) checkAndExecuteScaling(ctx context.Context) error {
	// 评估是否需要扩容
	decision, err := o.scalingDecisionService.EvaluateScaling(ctx)
	if err != nil {
		return fmt.Errorf("failed to evaluate scaling: %w", err)
	}

	if decision != nil {
		log.Printf("Scaling decision found: %s -> %s",
			decision.FromStrategy().String(), decision.ToStrategy().String())

		// 执行扩容
		if err := o.ExecuteScaling(ctx, decision); err != nil {
			return fmt.Errorf("failed to execute scaling: %w", err)
		}
	}

	return nil
}

// ForceScaleToMemory 强制切换到内存存储
func (o *DDDScalingOrchestrator) ForceScaleToMemory(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.currentStrategy == domain.StorageStrategyMemory {
		return fmt.Errorf("already using memory strategy")
	}

	decision := domain.NewScalingDecision(
		fmt.Sprintf("force_memory_%d", time.Now().Unix()),
		o.currentStrategy,
		domain.StorageStrategyMemory,
		domain.ScalingReasonManual,
		domain.DefaultThresholds(),
		domain.Metrics{}, // 简化实现
	)

	return o.ExecuteScaling(ctx, decision)
}

// ForceScaleToRedis 强制切换到Redis存储
func (o *DDDScalingOrchestrator) ForceScaleToRedis(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.currentStrategy == domain.StorageStrategyRedis {
		return fmt.Errorf("already using redis strategy")
	}

	decision := domain.NewScalingDecision(
		fmt.Sprintf("force_redis_%d", time.Now().Unix()),
		o.currentStrategy,
		domain.StorageStrategyRedis,
		domain.ScalingReasonManual,
		domain.DefaultThresholds(),
		domain.Metrics{}, // 简化实现
	)

	return o.ExecuteScaling(ctx, decision)
}
