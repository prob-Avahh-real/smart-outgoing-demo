package services

import (
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/events"
	"smart-outgoing-demo/internal/domain/repositories"
)

// VehicleManagementService coordinates vehicle operations across aggregates
type VehicleManagementService struct {
	vehicleRepo repositories.VehicleRepository
	routeRepo   repositories.RouteRepository
	metricsRepo repositories.MetricsRepository
	eventBus    chan events.DomainEvent
}

// NewVehicleManagementService creates a new vehicle management service
func NewVehicleManagementService(
	vehicleRepo repositories.VehicleRepository,
	routeRepo repositories.RouteRepository,
	metricsRepo repositories.MetricsRepository,
	eventBus chan events.DomainEvent,
) *VehicleManagementService {
	return &VehicleManagementService{
		vehicleRepo: vehicleRepo,
		routeRepo:   routeRepo,
		metricsRepo: metricsRepo,
		eventBus:    eventBus,
	}
}

// CreateVehicle creates a new vehicle and associated aggregates
func (s *VehicleManagementService) CreateVehicle(name string, start entities.Coordinates) (*entities.Vehicle, error) {
	// Generate vehicle ID
	vehicleID := generateVehicleID()

	// Create vehicle aggregate
	vehicle := entities.NewVehicle(vehicleID, name, start)

	// Save vehicle
	err := s.vehicleRepo.Save(vehicle)
	if err != nil {
		return nil, fmt.Errorf("failed to save vehicle: %w", err)
	}

	// Update metrics to reflect new vehicle count
	err = s.updateAgentCount()
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to update agent count: %v\n", err)
	}

	return vehicle, nil
}

// UpdateVehicleDestination updates a vehicle's destination and creates a route
func (s *VehicleManagementService) UpdateVehicleDestination(vehicleID string, destination entities.Coordinates) error {
	// Find vehicle
	vehicle, err := s.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return fmt.Errorf("vehicle not found: %w", err)
	}

	// Update vehicle destination
	vehicle.SetDestination(destination)

	// Save updated vehicle
	err = s.vehicleRepo.Save(vehicle)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}

	// Create a simple route (in a real implementation, this would use routing service)
	path := [][]float64{
		{vehicle.Start.Lng, vehicle.Start.Lat},
		{destination.Lng, destination.Lat},
	}

	route := entities.NewRoute(vehicleID, path)

	// Save route
	err = s.routeRepo.Save(route)
	if err != nil {
		return fmt.Errorf("failed to save route: %w", err)
	}

	return nil
}

// DeleteVehicle deletes a vehicle and all associated data
func (s *VehicleManagementService) DeleteVehicle(vehicleID string) error {
	// Check if vehicle exists
	_, err := s.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return fmt.Errorf("vehicle not found: %w", err)
	}

	// Delete vehicle
	err = s.vehicleRepo.Delete(vehicleID)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	// Delete associated routes
	err = s.routeRepo.DeleteByVehicleID(vehicleID)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to delete routes for vehicle %s: %v\n", vehicleID, err)
	}

	// Update metrics to reflect vehicle count change
	err = s.updateAgentCount()
	if err != nil {
		fmt.Printf("Warning: failed to update agent count: %v\n", err)
	}

	return nil
}

// GetVehicleWithRoute gets a vehicle along with its latest route
func (s *VehicleManagementService) GetVehicleWithRoute(vehicleID string) (*entities.Vehicle, *entities.Route, error) {
	// Find vehicle
	vehicle, err := s.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return nil, nil, fmt.Errorf("vehicle not found: %w", err)
	}

	// Find latest route
	route, err := s.routeRepo.FindLatestByVehicleID(vehicleID)
	if err != nil {
		// Route not found is not an error - vehicle might not have a route yet
		return vehicle, nil, nil
	}

	return vehicle, route, nil
}

// GetAllVehiclesWithRoutes gets all vehicles with their latest routes
func (s *VehicleManagementService) GetAllVehiclesWithRoutes() (map[string]*entities.Vehicle, map[string]*entities.Route, error) {
	// Get all vehicles
	vehicles, err := s.vehicleRepo.FindAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get vehicles: %w", err)
	}

	vehicleMap := make(map[string]*entities.Vehicle)
	routeMap := make(map[string]*entities.Route)

	// Populate vehicle map
	for _, vehicle := range vehicles {
		vehicleMap[vehicle.ID] = vehicle
	}

	// Get routes for all vehicles
	for _, vehicle := range vehicles {
		route, err := s.routeRepo.FindLatestByVehicleID(vehicle.ID)
		if err == nil {
			routeMap[vehicle.ID] = route
		}
		// If route not found, we just skip it
	}

	return vehicleMap, routeMap, nil
}

// updateAgentCount updates the agent count in metrics
func (s *VehicleManagementService) updateAgentCount() error {
	// Get current vehicle count
	count, err := s.vehicleRepo.Count()
	if err != nil {
		return fmt.Errorf("failed to count vehicles: %w", err)
	}

	// Get latest metrics
	latestMetrics, err := s.metricsRepo.FindLatest()
	if err != nil {
		// If no metrics exist, create new ones
		latestMetrics = entities.NewMetrics(0, 0, count, 0, "memory")
	} else {
		// Update existing metrics
		latestMetrics.Update(latestMetrics.MemoryUsage, latestMetrics.CPUUsage, count, latestMetrics.WebSocketConns, latestMetrics.StorageStrategy)
	}

	// Save updated metrics
	return s.metricsRepo.Save(latestMetrics)
}

// generateVehicleID generates a unique vehicle ID
func generateVehicleID() string {
	return fmt.Sprintf("vehicle_%d", time.Now().UnixNano())
}
