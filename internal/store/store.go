package store

import (
	"crypto/md5"
	"fmt"
	"smart-outgoing-demo/internal/cache"
	"sync"
	"time"
)

type Vehicle struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	StartLng  float64     `json:"start_lng"`
	StartLat  float64     `json:"start_lat"`
	StartAlt  float64     `json:"start_alt,omitempty"` // Altitude in meters
	EndLng    float64     `json:"end_lng,omitempty"`
	EndLat    float64     `json:"end_lat,omitempty"`
	EndAlt    float64     `json:"end_alt,omitempty"` // Altitude in meters
	Route     [][]float64 `json:"route,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type VehicleStore struct {
	mu         sync.RWMutex
	vehicles   map[string]*Vehicle
	routeCache *cache.RouteCache
}

func NewVehicleStore() *VehicleStore {
	return &VehicleStore{
		vehicles:   make(map[string]*Vehicle),
		routeCache: cache.NewRouteCache(),
	}
}

func (vs *VehicleStore) GetAll() []*Vehicle {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	vehicles := make([]*Vehicle, 0, len(vs.vehicles))
	for _, v := range vs.vehicles {
		vehicles = append(vehicles, v)
	}
	return vehicles
}

func (vs *VehicleStore) Get(id string) (*Vehicle, bool) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	v, exists := vs.vehicles[id]
	return v, exists
}

func (vs *VehicleStore) Create(vehicle *Vehicle) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()
	vs.vehicles[vehicle.ID] = vehicle
}

func (vs *VehicleStore) Update(id string, updateFn func(*Vehicle)) bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if vehicle, exists := vs.vehicles[id]; exists {
		updateFn(vehicle)
		vehicle.UpdatedAt = time.Now()
		return true
	}
	return false
}

func (vs *VehicleStore) Delete(id string) bool {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	if _, exists := vs.vehicles[id]; exists {
		delete(vs.vehicles, id)
		// Invalidate cache for deleted vehicle
		vs.routeCache.InvalidateVehicleRoute(id)
		return true
	}
	return false
}

// GetCachedRoute returns cached route data for a vehicle
func (vs *VehicleStore) GetCachedRoute(vehicleID string) ([][]float64, bool) {
	if data, exists := vs.routeCache.GetVehicleRoute(vehicleID); exists {
		return data.Route, true
	}
	return nil, false
}

// CacheRoute stores route data in cache
func (vs *VehicleStore) CacheRoute(vehicle *Vehicle, route [][]float64) {
	// Generate hash for cache invalidation
	hash := vs.generateRouteHash(vehicle, route)
	vs.routeCache.CacheVehicleRoute(vehicle.ID, route, hash)
}

// generateRouteHash creates a hash for route cache invalidation
func (vs *VehicleStore) generateRouteHash(vehicle *Vehicle, route [][]float64) string {
	data := fmt.Sprintf("%s-%.6f-%.6f-%.6f-%.6f-%v",
		vehicle.ID, vehicle.StartLng, vehicle.StartLat,
		vehicle.EndLng, vehicle.EndLat, route)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// CleanupCache removes expired cache entries
func (vs *VehicleStore) CleanupCache() {
	vs.routeCache.Cleanup()
}

// GetCacheStats returns cache statistics
func (vs *VehicleStore) GetCacheStats() map[string]int {
	return vs.routeCache.GetCacheStats()
}

// GetRouteCache returns the route cache for external access
func (vs *VehicleStore) GetRouteCache() *cache.RouteCache {
	return vs.routeCache
}
