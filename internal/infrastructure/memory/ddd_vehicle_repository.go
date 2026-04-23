package memory

import (
	"context"
	"errors"
	"sync"

	"smart-outgoing-demo/internal/domain"
)

// DDDMemoryVehicleRepository implements VehicleRepository using in-memory storage
type DDDMemoryVehicleRepository struct {
	vehicles map[string]*domain.Vehicle
	mu       sync.RWMutex
}

// NewDDDMemoryVehicleRepository creates a new memory vehicle repository
func NewDDDMemoryVehicleRepository() domain.VehicleRepository {
	return &DDDMemoryVehicleRepository{
		vehicles: make(map[string]*domain.Vehicle),
	}
}

// Save saves a vehicle aggregate
func (r *DDDMemoryVehicleRepository) Save(ctx context.Context, vehicle *domain.Vehicle) error {
	if vehicle == nil {
		return errors.New("vehicle cannot be nil")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.vehicles[vehicle.ID()] = vehicle
	return nil
}

// FindByID finds a vehicle by ID
func (r *DDDMemoryVehicleRepository) FindByID(ctx context.Context, id string) (*domain.Vehicle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	vehicle, exists := r.vehicles[id]
	if !exists {
		return nil, errors.New("vehicle not found")
	}
	
	return vehicle, nil
}

// FindAll finds all vehicles
func (r *DDDMemoryVehicleRepository) FindAll(ctx context.Context) ([]*domain.Vehicle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	vehicles := make([]*domain.Vehicle, 0, len(r.vehicles))
	for _, vehicle := range r.vehicles {
		vehicles = append(vehicles, vehicle)
	}
	
	return vehicles, nil
}

// Delete deletes a vehicle by ID
func (r *DDDMemoryVehicleRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.vehicles[id]; !exists {
		return errors.New("vehicle not found")
	}
	
	delete(r.vehicles, id)
	return nil
}

// FindByStatus finds vehicles by status
func (r *DDDMemoryVehicleRepository) FindByStatus(ctx context.Context, status domain.VehicleStatus) ([]*domain.Vehicle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var vehicles []*domain.Vehicle
	for _, vehicle := range r.vehicles {
		if vehicle.Status() == status {
			vehicles = append(vehicles, vehicle)
		}
	}
	
	return vehicles, nil
}

// Count returns the total number of vehicles
func (r *DDDMemoryVehicleRepository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return len(r.vehicles), nil
}
