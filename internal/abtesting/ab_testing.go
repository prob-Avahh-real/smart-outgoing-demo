package abtesting

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Experiment represents an A/B test experiment
type Experiment struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Variants    []Variant         `json:"variants"`
	Traffic     TrafficSplit      `json:"traffic"`
	CreatedAt   time.Time         `json:"created_at"`
	Status      ExperimentStatus  `json:"status"`
	Metrics     map[string]Metric `json:"metrics"`
}

// Variant represents a test variant
type Variant struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Weight      int                    `json:"weight"` // Traffic weight (0-100)
}

// TrafficSplit defines how traffic is distributed
type TrafficSplit struct {
	Type    string `json:"type"` // "weighted", "equal", "custom"
	Enabled bool   `json:"enabled"`
}

// ExperimentStatus represents experiment status
type ExperimentStatus string

const (
	StatusDraft     ExperimentStatus = "draft"
	StatusRunning   ExperimentStatus = "running"
	StatusPaused    ExperimentStatus = "paused"
	StatusCompleted ExperimentStatus = "completed"
)

// Metric represents experiment metrics
type Metric struct {
	Count       int64     `json:"count"`
	Conversion  float64   `json:"conversion"`
	LastUpdated  time.Time `json:"last_updated"`
}

// ABTestingManager manages A/B testing experiments
type ABTestingManager struct {
	mu          sync.RWMutex
	experiments map[string]*Experiment
	userAssignments map[string]string // user_id -> experiment_id
}

// NewABTestingManager creates a new A/B testing manager
func NewABTestingManager() *ABTestingManager {
	return &ABTestingManager{
		experiments:     make(map[string]*Experiment),
		userAssignments: make(map[string]string),
	}
}

// CreateExperiment creates a new A/B test experiment
func (m *ABTestingManager) CreateExperiment(name, description string, variants []Variant) (*Experiment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	experiment := &Experiment{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Variants:    variants,
		Traffic: TrafficSplit{
			Type:    "weighted",
			Enabled: true,
		},
		CreatedAt: time.Now(),
		Status:    StatusDraft,
		Metrics:   make(map[string]Metric),
	}

	// Initialize metrics for each variant
	for _, variant := range variants {
		experiment.Metrics[variant.ID] = Metric{
			LastUpdated: time.Now(),
		}
	}

	m.experiments[experiment.ID] = experiment
	return experiment, nil
}

// StartExperiment starts an experiment
func (m *ABTestingManager) StartExperiment(experimentID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	experiment, exists := m.experiments[experimentID]
	if !exists {
		return fmt.Errorf("experiment not found")
	}

	if experiment.Status != StatusDraft {
		return fmt.Errorf("experiment is not in draft status")
	}

	experiment.Status = StatusRunning
	return nil
}

// AssignUser assigns a user to a variant
func (m *ABTestingManager) AssignUser(experimentID, userID string) (*Variant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	experiment, exists := m.experiments[experimentID]
	if !exists {
		return nil, fmt.Errorf("experiment not found")
	}

	if experiment.Status != StatusRunning {
		return nil, fmt.Errorf("experiment is not running")
	}

	// Check if user is already assigned
	if assignment, exists := m.userAssignments[userID]; exists {
		if assignment == experimentID {
			// Return the same variant for consistency
			return m.selectVariant(experiment, userID), nil
		}
	}

	// Assign user to experiment
	m.userAssignments[userID] = experimentID

	// Select variant for user
	variant := m.selectVariant(experiment, userID)
	return variant, nil
}

// selectVariant selects a variant based on traffic weights
func (m *ABTestingManager) selectVariant(experiment *Experiment, userID string) *Variant {
	// Simple hash-based selection for consistency
	hash := hashUserID(userID)
	
	// Calculate total weight
	totalWeight := 0
	for _, variant := range experiment.Variants {
		totalWeight += variant.Weight
	}
	
	if totalWeight == 0 {
		// Equal distribution if no weights
		index := int(hash) % len(experiment.Variants)
		return &experiment.Variants[index]
	}
	
	// Weighted selection
	selection := int(hash) % totalWeight
	currentWeight := 0
	
	for _, variant := range experiment.Variants {
		currentWeight += variant.Weight
		if selection < currentWeight {
			return &variant
		}
	}
	
	// Fallback to first variant
	return &experiment.Variants[0]
}

// RecordConversion records a conversion for a variant
func (m *ABTestingManager) RecordConversion(experimentID, variantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	experiment, exists := m.experiments[experimentID]
	if !exists {
		return fmt.Errorf("experiment not found")
	}

	metric, exists := experiment.Metrics[variantID]
	if !exists {
		return fmt.Errorf("variant not found")
	}

	metric.Count++
	metric.LastUpdated = time.Now()
	experiment.Metrics[variantID] = metric

	return nil
}

// GetExperiment returns an experiment by ID
func (m *ABTestingManager) GetExperiment(experimentID string) (*Experiment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	experiment, exists := m.experiments[experimentID]
	if !exists {
		return nil, fmt.Errorf("experiment not found")
	}

	return experiment, nil
}

// GetAllExperiments returns all experiments
func (m *ABTestingManager) GetAllExperiments() []*Experiment {
	m.mu.RLock()
	defer m.mu.RUnlock()

	experiments := make([]*Experiment, 0, len(m.experiments))
	for _, exp := range m.experiments {
		experiments = append(experiments, exp)
	}

	return experiments
}

// GetExperimentStats returns statistics for an experiment
func (m *ABTestingManager) GetExperimentStats(experimentID string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	experiment, exists := m.experiments[experimentID]
	if !exists {
		return nil, fmt.Errorf("experiment not found")
	}

	stats := map[string]interface{}{
		"experiment_id": experiment.ID,
		"name":         experiment.Name,
		"status":       experiment.Status,
		"created_at":   experiment.CreatedAt,
		"variants":     make([]map[string]interface{}, 0),
	}

	totalConversions := int64(0)
	for _, metric := range experiment.Metrics {
		totalConversions += metric.Count
	}

	for _, variant := range experiment.Variants {
		metric := experiment.Metrics[variant.ID]
		conversionRate := float64(0)
		if totalConversions > 0 {
			conversionRate = float64(metric.Count) / float64(totalConversions) * 100
		}

		variantStats := map[string]interface{}{
			"id":               variant.ID,
			"name":             variant.Name,
			"conversions":      metric.Count,
			"conversion_rate":  conversionRate,
			"weight":           variant.Weight,
			"last_updated":     metric.LastUpdated,
		}

		stats["variants"] = append(stats["variants"].([]map[string]interface{}), variantStats)
	}

	return stats, nil
}

// generateID generates a unique ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// hashUserID creates a hash from user ID for consistent assignment
func hashUserID(userID string) int {
	hash := 0
	for _, char := range userID {
		hash = hash*31 + int(char)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash
}
