package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/events"
	"smart-outgoing-demo/internal/domain/repositories"
	"smart-outgoing-demo/internal/domain/services"
	"smart-outgoing-demo/internal/infrastructure/memory"
)

// ScalingOrchestrator coordinates the entire scaling process
type ScalingOrchestrator struct {
	decisionService *services.ScalingDecisionService
	vehicleService  *services.VehicleManagementService
	memoryFactory   repositories.RepositoryFactory
	eventBus        chan events.DomainEvent
	ctx             context.Context
	cancel          context.CancelFunc
}

// NewScalingOrchestrator creates a new scaling orchestrator
func NewScalingOrchestrator(
	decisionRepo repositories.ScalingDecisionRepository,
	metricsRepo repositories.MetricsRepository,
	vehicleRepo repositories.VehicleRepository,
	routeRepo repositories.RouteRepository,
) *ScalingOrchestrator {
	eventBus := make(chan events.DomainEvent, 100)
	ctx, cancel := context.WithCancel(context.Background())

	memoryFactory := memory.NewMemoryRepositoryFactory()

	decisionService := services.NewScalingDecisionService(decisionRepo, metricsRepo, eventBus)
	vehicleService := services.NewVehicleManagementService(vehicleRepo, routeRepo, metricsRepo, eventBus)

	return &ScalingOrchestrator{
		decisionService: decisionService,
		vehicleService:  vehicleService,
		memoryFactory:   memoryFactory,
		eventBus:        eventBus,
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start starts the scaling orchestrator
func (o *ScalingOrchestrator) Start() error {
	log.Println("Starting scaling orchestrator...")

	// Start event processor
	go o.processEvents()

	// Start metrics collection
	go o.collectMetrics()

	// Start scaling evaluation
	go o.evaluateScaling()

	log.Println("Scaling orchestrator started successfully")
	return nil
}

// Stop stops the scaling orchestrator
func (o *ScalingOrchestrator) Stop() error {
	log.Println("Stopping scaling orchestrator...")
	o.cancel()
	close(o.eventBus)
	log.Println("Scaling orchestrator stopped")
	return nil
}

// collectMetrics periodically collects system metrics
func (o *ScalingOrchestrator) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second) // Collect metrics every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-o.ctx.Done():
			return
		case <-ticker.C:
			// Get current storage strategy
			decision, err := o.decisionService.GetDecisionRepository().FindActive()
			if err != nil {
				log.Printf("Error getting current decision: %v\n", err)
				continue
			}

			// Collect current metrics
			_, err = o.decisionService.GetCurrentMetrics(string(decision.CurrentStrategy))
			if err != nil {
				log.Printf("Error collecting metrics: %v\n", err)
			}
		}
	}
}

// evaluateScaling periodically evaluates if scaling is needed
func (o *ScalingOrchestrator) evaluateScaling() {
	ticker := time.NewTicker(60 * time.Second) // Evaluate every minute
	defer ticker.Stop()

	for {
		select {
		case <-o.ctx.Done():
			return
		case <-ticker.C:
			decision, err := o.decisionService.EvaluateScaling()
			if err != nil {
				log.Printf("Error evaluating scaling: %v\n", err)
				continue
			}

			log.Printf("Current storage strategy: %s, Reason: %s\n",
				decision.CurrentStrategy, decision.Reason)
		}
	}
}

// processEvents processes domain events
func (o *ScalingOrchestrator) processEvents() {
	for {
		select {
		case <-o.ctx.Done():
			return
		case event, ok := <-o.eventBus:
			if !ok {
				return // Event bus closed
			}

			switch e := event.(type) {
			case *events.StorageStrategyChanged:
				o.handleStorageStrategyChanged(e)
			case *events.ThresholdBreached:
				o.handleThresholdBreached(e)
			default:
				log.Printf("Unknown event type: %T\n", event)
			}
		}
	}
}

// handleStorageStrategyChanged handles storage strategy change events
func (o *ScalingOrchestrator) handleStorageStrategyChanged(event *events.StorageStrategyChanged) {
	log.Printf("Storage strategy changed from %s to %s. Reason: %s\n",
		event.FromStrategy, event.ToStrategy, event.Reason)

	// In a real implementation, this would trigger data migration
	// For now, we just log the change
	err := o.migrateData(event.FromStrategy, event.ToStrategy)
	if err != nil {
		log.Printf("Error migrating data: %v\n", err)
	}
}

// handleThresholdBreached handles threshold breach events
func (o *ScalingOrchestrator) handleThresholdBreached(event *events.ThresholdBreached) {
	log.Printf("Threshold breached: %s. Current values: %+v, Thresholds: %+v\n",
		event.BreachedType, event.CurrentValues, event.Thresholds)

	// Trigger immediate scaling evaluation
	decision, err := o.decisionService.EvaluateScaling()
	if err != nil {
		log.Printf("Error evaluating scaling after threshold breach: %v\n", err)
		return
	}

	log.Printf("Post-breach evaluation: Strategy: %s, Reason: %s\n",
		decision.CurrentStrategy, decision.Reason)
}

// migrateData migrates data between storage strategies
func (o *ScalingOrchestrator) migrateData(fromStrategy, toStrategy string) error {
	log.Printf("Starting data migration from %s to %s...\n", fromStrategy, toStrategy)

	// In a real implementation, this would:
	// 1. Read all data from the source storage
	// 2. Write all data to the target storage
	// 3. Verify data integrity
	// 4. Switch active repositories
	// 5. Clean up old storage if needed

	// For now, we simulate the migration
	time.Sleep(1 * time.Second) // Simulate migration time
	log.Printf("Data migration from %s to %s completed\n", fromStrategy, toStrategy)

	return nil
}

// GetCurrentScalingStatus returns the current scaling status
func (o *ScalingOrchestrator) GetCurrentScalingStatus() (*entities.ScalingDecision, error) {
	return o.decisionService.GetDecisionRepository().FindActive()
}

// ForceScaleToMemory forces scaling to memory storage
func (o *ScalingOrchestrator) ForceScaleToMemory(reason string) error {
	decision, err := o.decisionService.GetDecisionRepository().FindActive()
	if err != nil {
		return fmt.Errorf("failed to get current decision: %w", err)
	}

	decision.SwitchToMemory(reason)
	return o.decisionService.GetDecisionRepository().Save(decision)
}

// ForceScaleToRedis forces scaling to Redis storage
func (o *ScalingOrchestrator) ForceScaleToRedis(reason string) error {
	decision, err := o.decisionService.GetDecisionRepository().FindActive()
	if err != nil {
		return fmt.Errorf("failed to get current decision: %w", err)
	}

	decision.SwitchToRedis(reason)
	return o.decisionService.GetDecisionRepository().Save(decision)
}
