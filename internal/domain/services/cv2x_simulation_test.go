package services_test

import (
	"testing"
	"time"

	"smart-outgoing-demo/internal/domain/services"
)

func TestNewCV2XSimulation(t *testing.T) {
	sim := services.NewCV2XSimulation()
	if sim == nil {
		t.Fatal("Expected non-nil simulation")
	}

	// Check if default RSUs are initialized
	rsus := sim.GetRSUs()
	if len(rsus) != 3 {
		t.Errorf("Expected 3 default RSUs, got %d", len(rsus))
	}
}

func TestRegisterVehicle(t *testing.T) {
	sim := services.NewCV2XSimulation()

	vehicle := &services.V2XVehicle{
		ID:         "test_vehicle_001",
		Type:       "car",
		Position:   &services.Position{Latitude: 22.6913, Longitude: 114.0448, Altitude: 0},
		Speed:      50,
		Heading:    0,
		V2XEnabled: true,
		LastSeen:   time.Now(),
	}

	err := sim.RegisterVehicle(vehicle)
	if err != nil {
		t.Fatalf("Failed to register vehicle: %v", err)
	}

	// Check if vehicle is registered
	vehicles := sim.GetVehicles()
	if len(vehicles) != 1 {
		t.Errorf("Expected 1 vehicle, got %d", len(vehicles))
	}

	if _, exists := vehicles["test_vehicle_001"]; !exists {
		t.Error("Vehicle not found in simulation")
	}
}

func TestGenerateRandomVehicles(t *testing.T) {
	sim := services.NewCV2XSimulation()

	count := 10
	centerLat := 22.6913
	centerLng := 114.0448
	radiusKm := 5.0

	vehicles, err := sim.GenerateRandomVehicles(count, centerLat, centerLng, radiusKm)
	if err != nil {
		t.Fatalf("Failed to generate random vehicles: %v", err)
	}

	if len(vehicles) != count {
		t.Errorf("Expected %d vehicles, got %d", count, len(vehicles))
	}

	// Check if all vehicles have valid positions
	for _, vehicle := range vehicles {
		if vehicle.Position == nil {
			t.Error("Vehicle position is nil")
		}
		if vehicle.Speed < 0 || vehicle.Speed > 150 {
			t.Errorf("Invalid vehicle speed: %f", vehicle.Speed)
		}
	}
}

func TestGenerateTrafficScenario(t *testing.T) {
	sim := services.NewCV2XSimulation()

	scenarios := []string{"congestion", "normal", "emergency"}

	for _, scenarioType := range scenarios {
		scenario, err := sim.GenerateTrafficScenario(scenarioType)
		if err != nil {
			t.Fatalf("Failed to generate scenario %s: %v", scenarioType, err)
		}

		if scenario == nil {
			t.Errorf("Scenario %s is nil", scenarioType)
		}

		if len(scenario.Vehicles) == 0 {
			t.Errorf("Scenario %s has no vehicles", scenarioType)
		}

		// Clear vehicles for next scenario
		sim.ClearAllVehicles()
	}

	// Test invalid scenario
	_, err := sim.GenerateTrafficScenario("invalid")
	if err == nil {
		t.Error("Expected error for invalid scenario type")
	}
}

func TestClearAllVehicles(t *testing.T) {
	sim := services.NewCV2XSimulation()

	// Add some vehicles
	_, err := sim.GenerateRandomVehicles(5, 22.6913, 114.0448, 5.0)
	if err != nil {
		t.Fatalf("Failed to generate vehicles: %v", err)
	}

	// Clear all vehicles
	err = sim.ClearAllVehicles()
	if err != nil {
		t.Fatalf("Failed to clear vehicles: %v", err)
	}

	// Check if all vehicles are cleared
	vehicles := sim.GetVehicles()
	if len(vehicles) != 0 {
		t.Errorf("Expected 0 vehicles after clear, got %d", len(vehicles))
	}
}

func TestUpdateVehiclePosition(t *testing.T) {
	sim := services.NewCV2XSimulation()

	vehicle := &services.V2XVehicle{
		ID:         "test_vehicle_001",
		Type:       "car",
		Position:   &services.Position{Latitude: 22.6913, Longitude: 114.0448, Altitude: 0},
		Speed:      50,
		Heading:    0,
		V2XEnabled: true,
		LastSeen:   time.Now(),
	}

	err := sim.RegisterVehicle(vehicle)
	if err != nil {
		t.Fatalf("Failed to register vehicle: %v", err)
	}

	// Update position
	newPosition := &services.Position{Latitude: 22.6920, Longitude: 114.0450, Altitude: 0}
	err = sim.UpdateVehiclePosition("test_vehicle_001", newPosition, 60, 45)
	if err != nil {
		t.Fatalf("Failed to update vehicle position: %v", err)
	}

	// Check if position was updated
	vehicles := sim.GetVehicles()
	updatedVehicle := vehicles["test_vehicle_001"]

	if updatedVehicle.Position.Latitude != 22.6920 {
		t.Errorf("Expected latitude 22.6920, got %f", updatedVehicle.Position.Latitude)
	}
	if updatedVehicle.Speed != 60 {
		t.Errorf("Expected speed 60, got %f", updatedVehicle.Speed)
	}
}

func TestGetStatistics(t *testing.T) {
	sim := services.NewCV2XSimulation()

	// Add some vehicles
	_, err := sim.GenerateRandomVehicles(5, 22.6913, 114.0448, 5.0)
	if err != nil {
		t.Fatalf("Failed to generate vehicles: %v", err)
	}

	stats := sim.GetStatistics()

	if stats["total_vehicles"] != 5 {
		t.Errorf("Expected 5 vehicles, got %v", stats["total_vehicles"])
	}
	if stats["total_rsus"] != 3 {
		t.Errorf("Expected 3 RSUs, got %v", stats["total_rsus"])
	}
}

func TestSendBSM(t *testing.T) {
	sim := services.NewCV2XSimulation()

	vehicle := &services.V2XVehicle{
		ID:         "test_vehicle_001",
		Type:       "car",
		Position:   &services.Position{Latitude: 22.6913, Longitude: 114.0448, Altitude: 0},
		Speed:      50,
		Heading:    0,
		V2XEnabled: true,
		LastSeen:   time.Now(),
	}

	err := sim.RegisterVehicle(vehicle)
	if err != nil {
		t.Fatalf("Failed to register vehicle: %v", err)
	}

	bsm := &services.BSM{
		VehicleID: "test_vehicle_001",
		Position:  &services.Position{Latitude: 22.6913, Longitude: 114.0448, Altitude: 0},
		Speed:     50,
		Heading:   0,
		Timestamp: time.Now(),
		Emergency: false,
	}

	message, err := sim.SendBSM("test_vehicle_001", bsm)
	if err != nil {
		t.Fatalf("Failed to send BSM: %v", err)
	}

	if message == nil {
		t.Error("Expected non-nil message")
	}
	if message.Type != "bsm" {
		t.Errorf("Expected message type 'bsm', got '%s'", message.Type)
	}
}

func TestSimulateTrafficHazard(t *testing.T) {
	sim := services.NewCV2XSimulation()

	position := &services.Position{Latitude: 22.6913, Longitude: 114.0448, Altitude: 0}
	duration := 30 * time.Minute

	denm, err := sim.SimulateTrafficHazard(position, "accident", "high", duration)
	if err != nil {
		t.Fatalf("Failed to simulate traffic hazard: %v", err)
	}

	if denm == nil {
		t.Error("Expected non-nil DENM")
	}
	if denm.EventType != "accident" {
		t.Errorf("Expected event type 'accident', got '%s'", denm.EventType)
	}
	if denm.Severity != "high" {
		t.Errorf("Expected severity 'high', got '%s'", denm.Severity)
	}
}
