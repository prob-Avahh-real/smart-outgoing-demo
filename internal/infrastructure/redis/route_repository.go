package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// RedisRouteRepository implements RouteRepository using Redis storage
type RedisRouteRepository struct {
	client RedisClientInterface
}

// NewRedisRouteRepository creates a new Redis route repository
func NewRedisRouteRepository(client RedisClientInterface) repositories.RouteRepository {
	return &RedisRouteRepository{
		client: client,
	}
}

// Save saves a route aggregate
func (r *RedisRouteRepository) Save(route *entities.Route) error {
	if route == nil {
		return errors.New("route cannot be nil")
	}
	ctx := context.Background()
	key := fmt.Sprintf("route:%s", route.ID)
	return r.client.Set(ctx, key, route, 24*time.Hour)
}

// FindByID finds a route by ID
func (r *RedisRouteRepository) FindByID(id string) (*entities.Route, error) {
	if id == "" {
		return nil, errors.New("route ID cannot be empty")
	}
	ctx := context.Background()
	key := fmt.Sprintf("route:%s", id)
	
	var route entities.Route
	err := r.client.Get(ctx, key, &route)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, errors.New("route not found")
		}
		return nil, fmt.Errorf("failed to get route: %w", err)
	}
	
	return &route, nil
}

// FindByVehicleID finds routes by vehicle ID
func (r *RedisRouteRepository) FindByVehicleID(vehicleID string) ([]*entities.Route, error) {
	ctx := context.Background()
	pattern := "route:*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get route keys: %w", err)
	}
	
	var vehicleRoutes []*entities.Route
	for _, key := range keys {
		var route entities.Route
		err := r.client.Get(ctx, key, &route)
		if err != nil {
			continue
		}
		if route.VehicleID == vehicleID {
			vehicleRoutes = append(vehicleRoutes, &route)
		}
	}
	
	return vehicleRoutes, nil
}

// FindLatestByVehicleID finds the latest route for a vehicle
func (r *RedisRouteRepository) FindLatestByVehicleID(vehicleID string) (*entities.Route, error) {
	routes, err := r.FindByVehicleID(vehicleID)
	if err != nil {
		return nil, err
	}
	
	if len(routes) == 0 {
		return nil, errors.New("no route found for vehicle")
	}
	
	// Return the route with the latest timestamp
	var latestRoute *entities.Route
	var latestTime int64 = 0
	
	for _, route := range routes {
		if route.UpdatedAt.Unix() > latestTime {
			latestRoute = route
			latestTime = route.UpdatedAt.Unix()
		}
	}
	
	return latestRoute, nil
}

// Delete deletes a route by ID
func (r *RedisRouteRepository) Delete(id string) error {
	if id == "" {
		return errors.New("route ID cannot be empty")
	}
	ctx := context.Background()
	key := fmt.Sprintf("route:%s", id)
	return r.client.Delete(ctx, key)
}

// DeleteByVehicleID deletes all routes for a vehicle
func (r *RedisRouteRepository) DeleteByVehicleID(vehicleID string) error {
	if vehicleID == "" {
		return errors.New("vehicle ID cannot be empty")
	}
	
	routes, err := r.FindByVehicleID(vehicleID)
	if err != nil {
		return err
	}
	
	if len(routes) == 0 {
		return errors.New("no routes found for vehicle")
	}
	
	ctx := context.Background()
	for _, route := range routes {
		key := fmt.Sprintf("route:%s", route.ID)
		r.client.Delete(ctx, key)
	}
	
	return nil
}

// Count returns the total number of routes
func (r *RedisRouteRepository) Count() (int, error) {
	ctx := context.Background()
	pattern := "route:*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return 0, fmt.Errorf("failed to count routes: %w", err)
	}
	
	return len(keys), nil
}
