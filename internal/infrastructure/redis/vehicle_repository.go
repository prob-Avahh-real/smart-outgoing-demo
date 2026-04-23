package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
)

// RedisVehicleRepository implements VehicleRepository using Redis storage
type RedisVehicleRepository struct {
	client RedisClientInterface
}

// NewRedisVehicleRepository creates a new Redis vehicle repository
func NewRedisVehicleRepository(client RedisClientInterface) repositories.VehicleRepository {
	return &RedisVehicleRepository{
		client: client,
	}
}

// Save saves a vehicle aggregate
func (r *RedisVehicleRepository) Save(vehicle *entities.Vehicle) error {
	if vehicle == nil {
		return errors.New("vehicle cannot be nil")
	}

	ctx := context.Background()
	key := fmt.Sprintf("vehicle:%s", vehicle.ID)

	// Store with 24 hour TTL
	return r.client.Set(ctx, key, vehicle, 24*time.Hour)
}

// FindByID finds a vehicle by ID
func (r *RedisVehicleRepository) FindByID(id string) (*entities.Vehicle, error) {
	if id == "" {
		return nil, errors.New("vehicle ID cannot be empty")
	}

	ctx := context.Background()
	key := fmt.Sprintf("vehicle:%s", id)

	var vehicle entities.Vehicle
	err := r.client.Get(ctx, key, &vehicle)
	if err != nil {
		if err == ErrKeyNotFound {
			return nil, errors.New("vehicle not found")
		}
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	return &vehicle, nil
}

// FindAll finds all vehicles
func (r *RedisVehicleRepository) FindAll() ([]*entities.Vehicle, error) {
	ctx := context.Background()
	pattern := "vehicle:*"

	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle keys: %w", err)
	}

	var vehicles []*entities.Vehicle
	for _, key := range keys {
		var vehicle entities.Vehicle
		err := r.client.Get(ctx, key, &vehicle)
		if err != nil {
			// Log error but continue with other vehicles
			continue
		}
		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

// Delete deletes a vehicle by ID
func (r *RedisVehicleRepository) Delete(id string) error {
	if id == "" {
		return errors.New("vehicle ID cannot be empty")
	}

	ctx := context.Background()
	key := fmt.Sprintf("vehicle:%s", id)

	// Check if vehicle exists before deleting
	exists, err := r.client.Exists(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to check vehicle existence: %w", err)
	}

	if !exists {
		return errors.New("vehicle not found")
	}

	return r.client.Delete(ctx, key)
}

// Count returns the total number of vehicles
func (r *RedisVehicleRepository) Count() (int, error) {
	ctx := context.Background()
	pattern := "vehicle:*"

	keys, err := r.client.Keys(ctx, pattern)
	if err != nil {
		return 0, fmt.Errorf("failed to count vehicles: %w", err)
	}

	return len(keys), nil
}

// Exists checks if a vehicle exists
func (r *RedisVehicleRepository) Exists(id string) (bool, error) {
	if id == "" {
		return false, errors.New("vehicle ID cannot be empty")
	}

	ctx := context.Background()
	key := fmt.Sprintf("vehicle:%s", id)

	return r.client.Exists(ctx, key)
}
