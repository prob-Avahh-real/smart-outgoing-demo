package repositories

import (
	"smart-outgoing-demo/internal/domain/entities"
)

// RouteRepository defines the repository interface for Route aggregate
type RouteRepository interface {
	// Save saves a route aggregate
	Save(route *entities.Route) error
	
	// FindByID finds a route by ID
	FindByID(id string) (*entities.Route, error)
	
	// FindByVehicleID finds routes by vehicle ID
	FindByVehicleID(vehicleID string) ([]*entities.Route, error)
	
	// FindLatestByVehicleID finds the latest route for a vehicle
	FindLatestByVehicleID(vehicleID string) (*entities.Route, error)
	
	// Delete deletes a route by ID
	Delete(id string) error
	
	// DeleteByVehicleID deletes all routes for a vehicle
	DeleteByVehicleID(vehicleID string) error
	
	// Count returns the total number of routes
	Count() (int, error)
}
