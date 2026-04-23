package repositories

import (
	"smart-outgoing-demo/internal/domain/entities"
)

// VehicleRepository defines the repository interface for Vehicle aggregate
type VehicleRepository interface {
	// Save saves a vehicle aggregate
	Save(vehicle *entities.Vehicle) error
	
	// FindByID finds a vehicle by ID
	FindByID(id string) (*entities.Vehicle, error)
	
	// FindAll finds all vehicles
	FindAll() ([]*entities.Vehicle, error)
	
	// Delete deletes a vehicle by ID
	Delete(id string) error
	
	// Count returns the total number of vehicles
	Count() (int, error)
	
	// Exists checks if a vehicle exists
	Exists(id string) (bool, error)
}
