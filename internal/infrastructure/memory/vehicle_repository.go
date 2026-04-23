package memory

import (
	"errors"
	"sync"
	
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// MemoryVehicleRepository implements VehicleRepository using in-memory storage
type MemoryVehicleRepository struct {
	vehicles map[string]*entities.Vehicle
	mu       sync.RWMutex
}

// NewMemoryVehicleRepository creates a new memory vehicle repository
func NewMemoryVehicleRepository() repositories.VehicleRepository {
	return &MemoryVehicleRepository{
		vehicles: make(map[string]*entities.Vehicle),
	}
}

// Save saves a vehicle aggregate
func (r *MemoryVehicleRepository) Save(vehicle *entities.Vehicle) error {
	if vehicle == nil {
		return errors.New("vehicle cannot be nil")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.vehicles[vehicle.ID] = vehicle
	return nil
}

// FindByID finds a vehicle by ID
func (r *MemoryVehicleRepository) FindByID(id string) (*entities.Vehicle, error) {
	if id == "" {
		return nil, errors.New("vehicle ID cannot be empty")
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	vehicle, exists := r.vehicles[id]
	if !exists {
		return nil, errors.New("vehicle not found")
	}
	
	return vehicle, nil
}

// FindAll finds all vehicles
func (r *MemoryVehicleRepository) FindAll() ([]*entities.Vehicle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	vehicles := make([]*entities.Vehicle, 0, len(r.vehicles))
	for _, vehicle := range r.vehicles {
		vehicles = append(vehicles, vehicle)
	}
	
	return vehicles, nil
}

// Delete deletes a vehicle by ID
func (r *MemoryVehicleRepository) Delete(id string) error {
	if id == "" {
		return errors.New("vehicle ID cannot be empty")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.vehicles[id]; !exists {
		return errors.New("vehicle not found")
	}
	
	delete(r.vehicles, id)
	return nil
}

// Count returns the total number of vehicles
func (r *MemoryVehicleRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return len(r.vehicles), nil
}

// Exists checks if a vehicle exists
func (r *MemoryVehicleRepository) Exists(id string) (bool, error) {
	if id == "" {
		return false, errors.New("vehicle ID cannot be empty")
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	_, exists := r.vehicles[id]
	return exists, nil
}
