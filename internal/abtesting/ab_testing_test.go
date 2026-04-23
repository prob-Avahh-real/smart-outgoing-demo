package abtesting

import (
	"fmt"
	"testing"
)

func TestABTestingManager_CreateExperiment(t *testing.T) {
	manager := NewABTestingManager()

	variants := []Variant{
		{ID: "control", Name: "Control", Weight: 50},
		{ID: "variant_a", Name: "Variant A", Weight: 50},
	}

	experiment, err := manager.CreateExperiment("Test Experiment", "Description", variants)
	if err != nil {
		t.Fatal("Failed to create experiment:", err)
	}

	if experiment.Name != "Test Experiment" {
		t.Errorf("Expected name 'Test Experiment', got '%s'", experiment.Name)
	}

	if experiment.Status != StatusDraft {
		t.Errorf("Expected status '%s', got '%s'", StatusDraft, experiment.Status)
	}

	if len(experiment.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(experiment.Variants))
	}
}

func TestABTestingManager_StartExperiment(t *testing.T) {
	manager := NewABTestingManager()

	variants := []Variant{
		{ID: "control", Name: "Control", Weight: 50},
		{ID: "variant_a", Name: "Variant A", Weight: 50},
	}

	experiment, _ := manager.CreateExperiment("Test", "Desc", variants)

	err := manager.StartExperiment(experiment.ID)
	if err != nil {
		t.Fatal("Failed to start experiment:", err)
	}

	// Verify status changed
	updated, _ := manager.GetExperiment(experiment.ID)
	if updated.Status != StatusRunning {
		t.Errorf("Expected status '%s', got '%s'", StatusRunning, updated.Status)
	}
}

func TestABTestingManager_AssignUser(t *testing.T) {
	manager := NewABTestingManager()

	variants := []Variant{
		{ID: "control", Name: "Control", Weight: 50},
		{ID: "variant_a", Name: "Variant A", Weight: 50},
	}

	experiment, _ := manager.CreateExperiment("Test", "Desc", variants)
	manager.StartExperiment(experiment.ID)

	variant, err := manager.AssignUser(experiment.ID, "user123")
	if err != nil {
		t.Fatal("Failed to assign user:", err)
	}

	if variant.ID == "" {
		t.Error("Variant ID should not be empty")
	}

	// Test consistent assignment
	variant2, err := manager.AssignUser(experiment.ID, "user123")
	if err != nil {
		t.Fatal("Failed to assign user again:", err)
	}

	if variant.ID != variant2.ID {
		t.Error("User should be assigned to same variant consistently")
	}
}

func TestABTestingManager_RecordConversion(t *testing.T) {
	manager := NewABTestingManager()

	variants := []Variant{
		{ID: "control", Name: "Control", Weight: 50},
		{ID: "variant_a", Name: "Variant A", Weight: 50},
	}

	experiment, _ := manager.CreateExperiment("Test", "Desc", variants)
	manager.StartExperiment(experiment.ID)

	variant, _ := manager.AssignUser(experiment.ID, "user123")

	// Record conversion
	err := manager.RecordConversion(experiment.ID, variant.ID)
	if err != nil {
		t.Fatal("Failed to record conversion:", err)
	}

	// Verify conversion recorded
	stats, _ := manager.GetExperimentStats(experiment.ID)
	variantStats := stats["variants"].([]map[string]interface{})

	var found bool
	for _, v := range variantStats {
		if v["id"] == variant.ID {
			if v["conversions"] != int64(1) {
				t.Errorf("Expected 1 conversion, got %v", v["conversions"])
			}
			found = true
			break
		}
	}

	if !found {
		t.Error("Variant not found in stats")
	}
}

func TestABTestingManager_WeightedSelection(t *testing.T) {
	manager := NewABTestingManager()

	// Create uneven weights
	variants := []Variant{
		{ID: "control", Name: "Control", Weight: 80},
		{ID: "variant_a", Name: "Variant A", Weight: 20},
	}

	experiment, _ := manager.CreateExperiment("Test", "Desc", variants)
	manager.StartExperiment(experiment.ID)

	// Test multiple users
	controlCount := 0
	variantACount := 0

	for i := 0; i < 100; i++ {
		userID := fmt.Sprintf("user%d", i)
		variant, _ := manager.AssignUser(experiment.ID, userID)

		if variant.ID == "control" {
			controlCount++
		} else if variant.ID == "variant_a" {
			variantACount++
		}
	}

	// Control should get roughly 80% of assignments
	controlRatio := float64(controlCount) / 100.0
	if controlRatio < 0.7 || controlRatio > 0.9 {
		t.Errorf("Expected control ratio around 0.8, got %.2f", controlRatio)
	}
}

func TestABTestingManager_ExperimentStats(t *testing.T) {
	manager := NewABTestingManager()

	variants := []Variant{
		{ID: "control", Name: "Control", Weight: 50},
		{ID: "variant_a", Name: "Variant A", Weight: 50},
	}

	experiment, _ := manager.CreateExperiment("Test", "Desc", variants)
	manager.StartExperiment(experiment.ID)

	// Record some conversions
	manager.RecordConversion(experiment.ID, "control")
	manager.RecordConversion(experiment.ID, "control")
	manager.RecordConversion(experiment.ID, "variant_a")

	stats, err := manager.GetExperimentStats(experiment.ID)
	if err != nil {
		t.Fatal("Failed to get stats:", err)
	}

	if stats["experiment_id"] != experiment.ID {
		t.Error("Experiment ID mismatch in stats")
	}

	variantStats := stats["variants"].([]map[string]interface{})
	if len(variantStats) != 2 {
		t.Errorf("Expected 2 variant stats, got %d", len(variantStats))
	}

	// Verify conversion counts
	controlFound := false
	variantAFound := false

	for _, v := range variantStats {
		id := v["id"].(string)
		conversions := v["conversions"].(int64)

		if id == "control" {
			if conversions != 2 {
				t.Errorf("Expected 2 conversions for control, got %d", conversions)
			}
			controlFound = true
		} else if id == "variant_a" {
			if conversions != 1 {
				t.Errorf("Expected 1 conversion for variant_a, got %d", conversions)
			}
			variantAFound = true
		}
	}

	if !controlFound || !variantAFound {
		t.Error("Not all variants found in stats")
	}
}
