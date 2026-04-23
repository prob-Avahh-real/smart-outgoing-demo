package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain"
)

// DDDRedisVehicleRepository implements VehicleRepository using Redis
type DDDRedisVehicleRepository struct {
	client RedisClientInterface
	prefix string
}

// NewDDDRedisVehicleRepository creates a new Redis vehicle repository
func NewDDDRedisVehicleRepository(client RedisClientInterface) domain.VehicleRepository {
	return &DDDRedisVehicleRepository{
		client: client,
		prefix: "vehicle:",
	}
}

// Save saves a vehicle aggregate to Redis
func (r *DDDRedisVehicleRepository) Save(ctx context.Context, vehicle *domain.Vehicle) error {
	if vehicle == nil {
		return fmt.Errorf("vehicle cannot be nil")
	}

	key := r.prefix + vehicle.ID()

	// Convert vehicle to JSON
	data, err := json.Marshal(vehicle)
	if err != nil {
		return fmt.Errorf("failed to marshal vehicle: %w", err)
	}

	// Save to Redis with TTL
	err = r.client.Set(ctx, key, string(data), 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to save vehicle to Redis: %w", err)
	}

	return nil
}

// FindByID finds a vehicle by ID in Redis
func (r *DDDRedisVehicleRepository) FindByID(ctx context.Context, id string) (*domain.Vehicle, error) {
	key := r.prefix + id

	var data string
	err := r.client.Get(ctx, key, &data)
	if err != nil {
		return nil, fmt.Errorf("vehicle not found: %w", err)
	}

	// Unmarshal vehicle
	var vehicle domain.Vehicle
	err = json.Unmarshal([]byte(data), &vehicle)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vehicle: %w", err)
	}

	return &vehicle, nil
}

// FindAll finds all vehicles in Redis
func (r *DDDRedisVehicleRepository) FindAll(ctx context.Context) ([]*domain.Vehicle, error) {
	pattern := r.prefix + "*"

	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle keys: %w", err)
	}

	var vehicles []*domain.Vehicle
	for _, key := range keys {
		var data string
		err := r.client.Get(ctx, key, &data)
		if err != nil {
			continue // Skip if key expired
		}

		var vehicle domain.Vehicle
		err = json.Unmarshal([]byte(data), &vehicle)
		if err != nil {
			continue // Skip malformed data
		}

		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

// Delete deletes a vehicle by ID from Redis
func (r *DDDRedisVehicleRepository) Delete(ctx context.Context, id string) error {
	key := r.prefix + id

	err := r.client.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle from Redis: %w", err)
	}

	return nil
}

// FindByStatus finds vehicles by status in Redis
func (r *DDDRedisVehicleRepository) FindByStatus(ctx context.Context, status domain.VehicleStatus) ([]*domain.Vehicle, error) {
	allVehicles, err := r.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var vehicles []*domain.Vehicle
	for _, vehicle := range allVehicles {
		if vehicle.Status() == status {
			vehicles = append(vehicles, vehicle)
		}
	}

	return vehicles, nil
}

// Count returns the total number of vehicles in Redis
func (r *DDDRedisVehicleRepository) Count(ctx context.Context) (int, error) {
	pattern := r.prefix + "*"

	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return 0, fmt.Errorf("failed to count vehicle keys: %w", err)
	}

	return len(keys), nil
}
