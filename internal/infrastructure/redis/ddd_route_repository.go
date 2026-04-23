package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// DDDRedisRouteRepository implements RouteRepository using Redis
type DDDRedisRouteRepository struct {
	client RedisClientInterface
	prefix string
}

// NewDDDRedisRouteRepository creates a new Redis route repository
func NewDDDRedisRouteRepository(client RedisClientInterface) domain.RouteRepository {
	return &DDDRedisRouteRepository{
		client: client,
		prefix: "route:",
	}
}

// Save saves a route aggregate to Redis
func (r *DDDRedisRouteRepository) Save(ctx context.Context, route *domain.Route) error {
	if route == nil {
		return fmt.Errorf("route cannot be nil")
	}

	key := r.prefix + route.ID()
	
	data, err := json.Marshal(route)
	if err != nil {
		return fmt.Errorf("failed to marshal route: %w", err)
	}

	err = r.client.Set(ctx, key, string(data), 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to save route to Redis: %w", err)
	}

	return nil
}

// FindByID finds a route by ID in Redis
func (r *DDDRedisRouteRepository) FindByID(ctx context.Context, id string) (*domain.Route, error) {
	key := r.prefix + id
	
	var data string
	err := r.client.Get(ctx, key, &data)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	var route domain.Route
	err = json.Unmarshal([]byte(data), &route)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal route: %w", err)
	}

	return &route, nil
}

// FindByVehicleID finds a route by vehicle ID in Redis
func (r *DDDRedisRouteRepository) FindByVehicleID(ctx context.Context, vehicleID string) (*domain.Route, error) {
	pattern := r.prefix + "*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get route keys: %w", err)
	}

	for _, key := range keys {
		var data string
		err := r.client.Get(ctx, key, &data)
		if err != nil {
			continue
		}

		var route domain.Route
		err = json.Unmarshal([]byte(data), &route)
		if err != nil {
			continue
		}

		if route.VehicleID() == vehicleID {
			return &route, nil
		}
	}

	return nil, fmt.Errorf("route not found for vehicle")
}

// FindAll finds all routes in Redis
func (r *DDDRedisRouteRepository) FindAll(ctx context.Context) ([]*domain.Route, error) {
	pattern := r.prefix + "*"
	
	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get route keys: %w", err)
	}

	var routes []*domain.Route
	for _, key := range keys {
		var data string
		err := r.client.Get(ctx, key, &data)
		if err != nil {
			continue
		}

		var route domain.Route
		err = json.Unmarshal([]byte(data), &route)
		if err != nil {
			continue
		}

		routes = append(routes, &route)
	}

	return routes, nil
}

// Delete deletes a route by ID from Redis
func (r *DDDRedisRouteRepository) Delete(ctx context.Context, id string) error {
	key := r.prefix + id
	
	err := r.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete route from Redis: %w", err)
	}

	return nil
}

// FindByStatus finds routes by status in Redis
func (r *DDDRedisRouteRepository) FindByStatus(ctx context.Context, status domain.RouteStatus) ([]*domain.Route, error) {
	allRoutes, err := r.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var routes []*domain.Route
	for _, route := range allRoutes {
		if route.Status() == status {
			routes = append(routes, route)
		}
	}

	return routes, nil
}
