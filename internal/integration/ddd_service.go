package integration

import (
	"context"
	"fmt"
	"log"

	"smart-outgoing-demo/internal/application/services"
	"smart-outgoing-demo/internal/domain"
)

// DDDService integrates the DDD architecture with the existing application
type DDDService struct {
	integrationService *services.DDDIntegrationService
	currentStrategy    domain.StorageStrategy
	ctx                context.Context
	cancel             context.CancelFunc
}

// NewDDDService creates a new DDD service
func NewDDDService() *DDDService {
	// Create the new DDD integration service
	integrationService := services.NewDDDIntegrationService()

	ctx, cancel := context.WithCancel(context.Background())

	return &DDDService{
		integrationService: integrationService,
		currentStrategy:    domain.StorageStrategyMemory,
		ctx:                ctx,
		cancel:             cancel,
	}
}

// Start starts the DDD service
func (s *DDDService) Start() error {
	log.Println("Starting DDD service...")

	// Start the integration service
	err := s.integrationService.Start(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to start integration service: %w", err)
	}

	log.Println("DDD service started successfully")
	return nil
}

// Stop stops the DDD service
func (s *DDDService) Stop() error {
	log.Println("Stopping DDD service...")

	s.cancel()

	log.Println("DDD service stopped")
	return nil
}

// GetCurrentStrategy returns the current storage strategy
func (s *DDDService) GetCurrentStrategy() domain.StorageStrategy {
	return s.currentStrategy
}

// SwitchToRedis switches to Redis storage strategy
func (s *DDDService) SwitchToRedis() error {
	if s.currentStrategy == domain.StorageStrategyRedis {
		return fmt.Errorf("already using Redis storage")
	}

	log.Println("Switching to Redis storage...")

	err := s.integrationService.ForceScaling(s.ctx, domain.StorageStrategyRedis)
	if err != nil {
		return fmt.Errorf("failed to switch to Redis: %w", err)
	}

	s.currentStrategy = domain.StorageStrategyRedis

	log.Println("Switched to Redis storage successfully")
	return nil
}

// SwitchToMemory switches to memory storage strategy
func (s *DDDService) SwitchToMemory() error {
	if s.currentStrategy == domain.StorageStrategyMemory {
		return fmt.Errorf("already using memory storage")
	}

	log.Println("Switching to memory storage...")

	err := s.integrationService.ForceScaling(s.ctx, domain.StorageStrategyMemory)
	if err != nil {
		return fmt.Errorf("failed to switch to memory: %w", err)
	}

	s.currentStrategy = domain.StorageStrategyMemory

	log.Println("Switched to memory storage successfully")
	return nil
}

// GetIntegrationService returns the integration service
func (s *DDDService) GetIntegrationService() *services.DDDIntegrationService {
	return s.integrationService
}

// GetScalingStatus returns the current scaling status
func (s *DDDService) GetScalingStatus() services.ScalingStatusInfo {
	return s.integrationService.GetScalingStatus()
}

// ForceScaleToMemory forces scaling to memory
func (s *DDDService) ForceScaleToMemory(reason string) error {
	return s.SwitchToMemory()
}

// ForceScaleToRedis forces scaling to Redis
func (s *DDDService) ForceScaleToRedis(reason string) error {
	return s.SwitchToRedis()
}

// CreateVehicle creates a vehicle through the integration service
func (s *DDDService) CreateVehicle(id, name, role string, lng, lat, alt float64) error {
	_, err := s.integrationService.CreateVehicle(s.ctx, id, name, role, lng, lat, alt)
	return err
}

// AssignRoute assigns a route through the integration service
func (s *DDDService) AssignRoute(vehicleID string, waypoints []domain.Coordinates) error {
	_, err := s.integrationService.AssignRoute(s.ctx, vehicleID, waypoints)
	return err
}

// UpdateVehiclePosition updates vehicle position through the integration service
func (s *DDDService) UpdateVehiclePosition(vehicleID string, lng, lat, alt float64) error {
	return s.integrationService.UpdateVehiclePosition(s.ctx, vehicleID, lng, lat, alt)
}
