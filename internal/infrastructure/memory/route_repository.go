package memory

import (
	"errors"
	"sync"
	
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// MemoryRouteRepository implements RouteRepository using in-memory storage
type MemoryRouteRepository struct {
	routes map[string]*entities.Route
	mu     sync.RWMutex
}

// NewMemoryRouteRepository creates a new memory route repository
func NewMemoryRouteRepository() repositories.RouteRepository {
	return &MemoryRouteRepository{
		routes: make(map[string]*entities.Route),
	}
}

// Save saves a route aggregate
func (r *MemoryRouteRepository) Save(route *entities.Route) error {
	if route == nil {
		return errors.New("route cannot be nil")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.routes[route.ID] = route
	return nil
}

// FindByID finds a route by ID
func (r *MemoryRouteRepository) FindByID(id string) (*entities.Route, error) {
	if id == "" {
		return nil, errors.New("route ID cannot be empty")
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	route, exists := r.routes[id]
	if !exists {
		return nil, errors.New("route not found")
	}
	
	return route, nil
}

// FindByVehicleID finds routes by vehicle ID
func (r *MemoryRouteRepository) FindByVehicleID(vehicleID string) ([]*entities.Route, error) {
	if vehicleID == "" {
		return nil, errors.New("vehicle ID cannot be empty")
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var vehicleRoutes []*entities.Route
	for _, route := range r.routes {
		if route.VehicleID == vehicleID {
			vehicleRoutes = append(vehicleRoutes, route)
		}
	}
	
	return vehicleRoutes, nil
}

// FindLatestByVehicleID finds the latest route for a vehicle
func (r *MemoryRouteRepository) FindLatestByVehicleID(vehicleID string) (*entities.Route, error) {
	if vehicleID == "" {
		return nil, errors.New("vehicle ID cannot be empty")
	}
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var latestRoute *entities.Route
	var latestTime int64 = 0
	
	for _, route := range r.routes {
		if route.VehicleID == vehicleID && route.UpdatedAt.Unix() > latestTime {
			latestRoute = route
			latestTime = route.UpdatedAt.Unix()
		}
	}
	
	if latestRoute == nil {
		return nil, errors.New("no route found for vehicle")
	}
	
	return latestRoute, nil
}

// Delete deletes a route by ID
func (r *MemoryRouteRepository) Delete(id string) error {
	if id == "" {
		return errors.New("route ID cannot be empty")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.routes[id]; !exists {
		return errors.New("route not found")
	}
	
	delete(r.routes, id)
	return nil
}

// DeleteByVehicleID deletes all routes for a vehicle
func (r *MemoryRouteRepository) DeleteByVehicleID(vehicleID string) error {
	if vehicleID == "" {
		return errors.New("vehicle ID cannot be empty")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	deleted := 0
	for id, route := range r.routes {
		if route.VehicleID == vehicleID {
			delete(r.routes, id)
			deleted++
		}
	}
	
	if deleted == 0 {
		return errors.New("no routes found for vehicle")
	}
	
	return nil
}

// Count returns the total number of routes
func (r *MemoryRouteRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return len(r.routes), nil
}
