package cache

import (
	"testing"
	"time"
)

func TestMapCache_SetAndGet(t *testing.T) {
	cache := NewMapCache()
	
	// Test setting and getting a value
	cache.Set("key1", "value1", 5*time.Minute)
	
	value, exists := cache.Get("key1")
	if !exists {
		t.Fatal("Value should exist")
	}
	
	if value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}
}

func TestMapCache_Expiration(t *testing.T) {
	cache := NewMapCache()
	
	// Set value with 1ms TTL
	cache.Set("expire", "value", 1*time.Millisecond)
	
	// Should exist immediately
	_, exists := cache.Get("expire")
	if !exists {
		t.Fatal("Value should exist immediately")
	}
	
	// Wait for expiration
	time.Sleep(10 * time.Millisecond)
	
	// Should be expired
	_, exists = cache.Get("expire")
	if exists {
		t.Fatal("Value should be expired")
	}
}

func TestMapCache_Delete(t *testing.T) {
	cache := NewMapCache()
	
	cache.Set("key1", "value1", 5*time.Minute)
	
	// Delete it
	cache.Delete("key1")
	
	// Should not exist
	_, exists := cache.Get("key1")
	if exists {
		t.Fatal("Value should not exist after deletion")
	}
}

func TestRouteCache(t *testing.T) {
	cache := NewRouteCache()
	
	// Cache route data
	route := [][]float64{{114.0, 22.0}, {115.0, 23.0}}
	cache.CacheVehicleRoute("vehicle1", route, "hash1")
	
	// Retrieve cached route
	data, exists := cache.GetVehicleRoute("vehicle1")
	if !exists {
		t.Fatal("Cached route should exist")
	}
	
	if len(data.Route) != 2 {
		t.Errorf("Expected 2 route points, got %d", len(data.Route))
	}
	
	// Test cache stats
	stats := cache.GetCacheStats()
	if stats["vehicle_routes"] != 1 {
		t.Errorf("Expected 1 cached route, got %d", stats["vehicle_routes"])
	}
}
