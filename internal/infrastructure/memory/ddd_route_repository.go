package memory

import (
	"context"
	"errors"
	"sync"

	"smart-outgoing-demo/internal/domain"
)

// DDDMemoryRouteRepository implements RouteRepository using in-memory storage
type DDDMemoryRouteRepository struct {
	routes map[string]*domain.Route
	mu     sync.RWMutex
}

// NewDDDMemoryRouteRepository creates a new memory route repository
func NewDDDMemoryRouteRepository() domain.RouteRepository {
	return &DDDMemoryRouteRepository{
		routes: make(map[string]*domain.Route),
	}
}

// Save saves a route aggregate
func (r *DDDMemoryRouteRepository) Save(ctx context.Context, route *domain.Route) error {
	if route == nil {
		return errors.New("route cannot be nil")
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.routes[route.ID()] = route
	return nil
}

// FindByID finds a route by ID
func (r *DDDMemoryRouteRepository) FindByID(ctx context.Context, id string) (*domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	route, exists := r.routes[id]
	if !exists {
		return nil, errors.New("route not found")
	}
	
	return route, nil
}

// FindByVehicleID finds a route by vehicle ID
func (r *DDDMemoryRouteRepository) FindByVehicleID(ctx context.Context, vehicleID string) (*domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, route := range r.routes {
		if route.VehicleID() == vehicleID {
			return route, nil
		}
	}
	
	return nil, errors.New("route not found for vehicle")
}

// FindAll finds all routes
func (r *DDDMemoryRouteRepository) FindAll(ctx context.Context) ([]*domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	routes := make([]*domain.Route, 0, len(r.routes))
	for _, route := range r.routes {
		routes = append(routes, route)
	}
	
	return routes, nil
}

// Delete deletes a route by ID
func (r *DDDMemoryRouteRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.routes[id]; !exists {
		return errors.New("route not found")
	}
	
	delete(r.routes, id)
	return nil
}

// FindByStatus finds routes by status
func (r *DDDMemoryRouteRepository) FindByStatus(ctx context.Context, status domain.RouteStatus) ([]*domain.Route, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var routes []*domain.Route
	for _, route := range r.routes {
		if route.Status() == status {
			routes = append(routes, route)
		}
	}
	
	return routes, nil
}
