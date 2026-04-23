package store

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewVehicleStore(t *testing.T) {
	store := NewVehicleStore()
	if store == nil {
		t.Fatal("NewVehicleStore returned nil")
	}
	if store.vehicles == nil {
		t.Error("vehicles map is nil")
	}
}

func TestCreateVehicle(t *testing.T) {
	store := NewVehicleStore()
	
	vehicle := &Vehicle{
		ID:       "test1",
		Name:     "Test Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	
	store.Create(vehicle)
	
	retrieved, exists := store.Get("test1")
	if !exists {
		t.Error("Vehicle not found after creation")
	}
	
	if retrieved.ID != "test1" {
		t.Errorf("Expected ID 'test1', got '%s'", retrieved.ID)
	}
	
	if retrieved.Name != "Test Vehicle" {
		t.Errorf("Expected name 'Test Vehicle', got '%s'", retrieved.Name)
	}
}

func TestGetVehicle(t *testing.T) {
	store := NewVehicleStore()
	
	vehicle := &Vehicle{
		ID:       "test1",
		Name:     "Test Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	
	store.Create(vehicle)
	
	retrieved, exists := store.Get("test1")
	if !exists {
		t.Error("Vehicle not found")
	}
	
	if retrieved.ID != "test1" {
		t.Errorf("Expected ID 'test1', got '%s'", retrieved.ID)
	}
	
	_, exists = store.Get("nonexistent")
	if exists {
		t.Error("Should not find nonexistent vehicle")
	}
}

func TestGetAllVehicles(t *testing.T) {
	store := NewVehicleStore()
	
	vehicles := store.GetAll()
	if len(vehicles) != 0 {
		t.Error("Expected empty store to return 0 vehicles")
	}
	
	store.Create(&Vehicle{ID: "test1", Name: "Vehicle 1", StartLng: 114.0, StartLat: 22.0})
	store.Create(&Vehicle{ID: "test2", Name: "Vehicle 2", StartLng: 114.1, StartLat: 22.1})
	
	vehicles = store.GetAll()
	if len(vehicles) != 2 {
		t.Errorf("Expected 2 vehicles, got %d", len(vehicles))
	}
}

func TestUpdateVehicle(t *testing.T) {
	store := NewVehicleStore()
	
	vehicle := &Vehicle{
		ID:       "test1",
		Name:     "Test Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	
	store.Create(vehicle)
	
	updated := store.Update("test1", func(v *Vehicle) {
		v.Name = "Updated Vehicle"
		v.EndLng = 114.5
		v.EndLat = 22.5
	})
	
	if !updated {
		t.Error("Update should return true for existing vehicle")
	}
	
	retrieved, _ := store.Get("test1")
	if retrieved.Name != "Updated Vehicle" {
		t.Errorf("Expected name 'Updated Vehicle', got '%s'", retrieved.Name)
	}
	
	if retrieved.EndLng != 114.5 {
		t.Errorf("Expected EndLng 114.5, got %f", retrieved.EndLng)
	}
}

func TestUpdateNonExistent(t *testing.T) {
	store := NewVehicleStore()
	
	updated := store.Update("nonexistent", func(v *Vehicle) {
		v.Name = "Should not update"
	})
	
	if updated {
		t.Error("Update should return false for nonexistent vehicle")
	}
}

func TestDeleteVehicle(t *testing.T) {
	store := NewVehicleStore()
	
	vehicle := &Vehicle{
		ID:       "test1",
		Name:     "Test Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	
	store.Create(vehicle)
	
	deleted := store.Delete("test1")
	if !deleted {
		t.Error("Delete should return true for existing vehicle")
	}
	
	_, exists := store.Get("test1")
	if exists {
		t.Error("Vehicle should not exist after deletion")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	store := NewVehicleStore()
	
	deleted := store.Delete("nonexistent")
	if deleted {
		t.Error("Delete should return false for nonexistent vehicle")
	}
}

func TestTimestamps(t *testing.T) {
	store := NewVehicleStore()
	
	beforeCreate := time.Now()
	
	vehicle := &Vehicle{
		ID:       "test1",
		Name:     "Test Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	
	store.Create(vehicle)
	
	retrieved, _ := store.Get("test1")
	
	if retrieved.CreatedAt.Before(beforeCreate) {
		t.Error("CreatedAt should be set to current time")
	}
	
	if retrieved.UpdatedAt.Before(beforeCreate) {
		t.Error("UpdatedAt should be set to current time")
	}
	
	time.Sleep(10 * time.Millisecond)
	
	store.Update("test1", func(v *Vehicle) {
		v.Name = "Updated"
	})
	
	retrieved, _ = store.Get("test1")
	
	if retrieved.UpdatedAt.Before(retrieved.CreatedAt) {
		t.Error("UpdatedAt should be updated after modification")
	}
}

func TestVehicleEdgeCases(t *testing.T) {
	store := NewVehicleStore()
	
	// Test vehicle with zero coordinates
	vehicle := &Vehicle{
		ID:       "zero",
		Name:     "Zero Vehicle",
		StartLng: 0,
		StartLat: 0,
	}
	store.Create(vehicle)
	
	retrieved, _ := store.Get("zero")
	if retrieved.StartLng != 0 || retrieved.StartLat != 0 {
		t.Error("Should accept zero coordinates")
	}
	
	// Test vehicle with negative coordinates
	vehicle = &Vehicle{
		ID:       "negative",
		Name:     "Negative Vehicle",
		StartLng: -180,
		StartLat: -90,
	}
	store.Create(vehicle)
	
	retrieved, _ = store.Get("negative")
	if retrieved.StartLng != -180 || retrieved.StartLat != -90 {
		t.Error("Should accept negative coordinates")
	}
	
	// Test vehicle with extreme coordinates
	vehicle = &Vehicle{
		ID:       "extreme",
		Name:     "Extreme Vehicle",
		StartLng: 180,
		StartLat: 90,
	}
	store.Create(vehicle)
	
	retrieved, _ = store.Get("extreme")
	if retrieved.StartLng != 180 || retrieved.StartLat != 90 {
		t.Error("Should accept extreme coordinates")
	}
	
	// Test vehicle with empty name
	vehicle = &Vehicle{
		ID:       "empty_name",
		Name:     "",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	store.Create(vehicle)
	
	retrieved, _ = store.Get("empty_name")
	if retrieved.Name != "" {
		t.Error("Should accept empty name")
	}
}

func TestConcurrentOperations(t *testing.T) {
	store := NewVehicleStore()
	
	// Test concurrent creates
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(i int) {
			vehicle := &Vehicle{
				ID:       fmt.Sprintf("concurrent_%d", i),
				Name:     fmt.Sprintf("Vehicle %d", i),
				StartLng: 114.0 + float64(i)*0.01,
				StartLat: 22.0 + float64(i)*0.01,
			}
			store.Create(vehicle)
			done <- true
		}(i)
	}
	
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify all vehicles were created
	vehicles := store.GetAll()
	if len(vehicles) != 10 {
		t.Errorf("Expected 10 vehicles after concurrent creates, got %d", len(vehicles))
	}
	
	// Test concurrent updates
	for i := 0; i < 10; i++ {
		go func(i int) {
			store.Update(fmt.Sprintf("concurrent_%d", i), func(v *Vehicle) {
				v.Name = fmt.Sprintf("Updated Vehicle %d", i)
			})
			done <- true
		}(i)
	}
	
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify all updates were applied
	vehicles = store.GetAll()
	for _, vehicle := range vehicles {
		if !strings.HasPrefix(vehicle.Name, "Updated") {
			t.Errorf("Vehicle name should be updated, got %s", vehicle.Name)
		}
	}
}

func TestStoreRobustness(t *testing.T) {
	store := NewVehicleStore()
	
	// Test with empty ID
	vehicle := &Vehicle{
		ID:       "",
		Name:     "Empty ID Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	store.Create(vehicle)
	
	_, exists := store.Get("")
	if !exists {
		t.Error("Should handle empty ID")
	}
	
	// Test with very long ID
	longID := strings.Repeat("a", 1000)
	vehicle = &Vehicle{
		ID:       longID,
		Name:     "Long ID Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	store.Create(vehicle)
	
	retrieved, exists := store.Get(longID)
	if !exists || retrieved.ID != longID {
		t.Error("Should handle very long IDs")
	}
	
	// Test with very long name
	longName := strings.Repeat("n", 1000)
	vehicle = &Vehicle{
		ID:       "long_name",
		Name:     longName,
		StartLng: 114.0,
		StartLat: 22.0,
	}
	store.Create(vehicle)
	
	retrieved, _ = store.Get("long_name")
	if retrieved.Name != longName {
		t.Error("Should handle very long names")
	}
	
	// Test update with no-op function
	store.Create(&Vehicle{
		ID:       "noop",
		Name:     "Noop Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	})
	
	updated := store.Update("noop", func(v *Vehicle) {
		// No-op
	})
	if !updated {
		t.Error("No-op update should still return true")
	}
}

func TestDataIntegrity(t *testing.T) {
	store := NewVehicleStore()
	
	// Test that data is not corrupted by rapid operations
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("integrity_%d", i)
		vehicle := &Vehicle{
			ID:       id,
			Name:     fmt.Sprintf("Vehicle %d", i),
			StartLng: 114.0 + float64(i)*0.001,
			StartLat: 22.0 + float64(i)*0.001,
		}
		store.Create(vehicle)
	}
	
	vehicles := store.GetAll()
	if len(vehicles) != 100 {
		t.Errorf("Expected 100 vehicles, got %d", len(vehicles))
	}
	
	// Verify data integrity
	for i := 0; i < 100; i++ {
		id := fmt.Sprintf("integrity_%d", i)
		vehicle, exists := store.Get(id)
		if !exists {
			t.Errorf("Vehicle %s should exist", id)
			continue
		}
		expectedLng := 114.0 + float64(i)*0.001
		expectedLat := 22.0 + float64(i)*0.001
		if vehicle.StartLng != expectedLng || vehicle.StartLat != expectedLat {
			t.Errorf("Vehicle %s data corrupted: expected (%f, %f), got (%f, %f)",
				id, expectedLng, expectedLat, vehicle.StartLng, vehicle.StartLat)
		}
	}
}
