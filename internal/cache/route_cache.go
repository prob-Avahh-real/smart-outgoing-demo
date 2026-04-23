package cache

import (
	"sync"
	"time"
)

// RouteCache provides caching for vehicle routes and map data
type RouteCache struct {
	mapCache    *MapCache
	vehicleData map[string]*VehicleRouteData
	mu          sync.RWMutex
}

// VehicleRouteData contains cached vehicle route information
type VehicleRouteData struct {
	VehicleID    string
	Route        [][]float64
	LastUpdated  time.Time
	RouteHash    string // For cache invalidation
}

// NewRouteCache creates a new route cache
func NewRouteCache() *RouteCache {
	return &RouteCache{
		mapCache:    NewMapCache(),
		vehicleData: make(map[string]*VehicleRouteData),
	}
}

// CacheVehicleRoute stores vehicle route data in cache
func (rc *RouteCache) CacheVehicleRoute(vehicleID string, route [][]float64, hash string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	rc.vehicleData[vehicleID] = &VehicleRouteData{
		VehicleID:   vehicleID,
		Route:       route,
		LastUpdated: time.Now(),
		RouteHash:   hash,
	}

	// Cache with 5 minute TTL
	rc.mapCache.Set("route:"+vehicleID, route, 5*time.Minute)
}

// GetVehicleRoute retrieves cached vehicle route data
func (rc *RouteCache) GetVehicleRoute(vehicleID string) (*VehicleRouteData, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	data, exists := rc.vehicleData[vehicleID]
	if !exists {
		return nil, false
	}

	// Check if cache is still valid (5 minutes)
	if time.Since(data.LastUpdated) > 5*time.Minute {
		delete(rc.vehicleData, vehicleID)
		return nil, false
	}

	return data, true
}

// InvalidateVehicleRoute removes cached route data for a vehicle
func (rc *RouteCache) InvalidateVehicleRoute(vehicleID string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	delete(rc.vehicleData, vehicleID)
	rc.mapCache.Delete("route:" + vehicleID)
}

// CacheMapData stores general map data with TTL
func (rc *RouteCache) CacheMapData(key string, data interface{}, ttl time.Duration) {
	rc.mapCache.Set(key, data, ttl)
}

// GetMapData retrieves cached map data
func (rc *RouteCache) GetMapData(key string) (interface{}, bool) {
	return rc.mapCache.Get(key)
}

// Cleanup removes expired cache entries
func (rc *RouteCache) Cleanup() {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	now := time.Now()
	for id, data := range rc.vehicleData {
		if now.Sub(data.LastUpdated) > 5*time.Minute {
			delete(rc.vehicleData, id)
		}
	}

	rc.mapCache.Cleanup()
}

// GetCacheStats returns cache statistics
func (rc *RouteCache) GetCacheStats() map[string]int {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	return map[string]int{
		"vehicle_routes": len(rc.vehicleData),
		"total_items":    rc.mapCache.Size(),
	}
}
