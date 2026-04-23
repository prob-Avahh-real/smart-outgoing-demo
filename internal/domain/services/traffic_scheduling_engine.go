package services

import (
	"fmt"
	"sync"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
)

// TrafficSchedulingEngine manages traffic flow scheduling and diversion
type TrafficSchedulingEngine struct {
	mu                sync.RWMutex
	zones             map[string]*TrafficZone
	diversionRules    []*entities.TrafficDiversionRule
	schedulingHistory []*SchedulingDecision
	monitoringEnabled bool
}

// TrafficZone represents a traffic monitoring zone
type TrafficZone struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	CenterLat         float64   `json:"center_lat"`
	CenterLng         float64   `json:"center_lng"`
	Radius            float64   `json:"radius"`          // km
	CurrentDensity    float64   `json:"current_density"` // vehicles/km
	HistoricalDensity []float64 `json:"historical_density"`
	PeakHourDensity   float64   `json:"peak_hour_density"`
	Status            string    `json:"status"` // normal/congested/severe
	LastUpdated       time.Time `json:"last_updated"`
}

// SchedulingDecision represents a traffic scheduling decision
type SchedulingDecision struct {
	ID              string                    `json:"id"`
	ZoneID          string                    `json:"zone_id"`
	DecisionType    string                    `json:"decision_type"` // diversion/rerouting/signal_control
	TargetPoolLevel entities.ParkingPoolLevel `json:"target_pool_level"`
	Reason          string                    `json:"reason"`
	Confidence      float64                   `json:"confidence"` // 0-1
	ExecutedAt      time.Time                 `json:"executed_at"`
	Effectiveness   float64                   `json:"effectiveness"` // 0-1
}

// NewTrafficSchedulingEngine creates a new traffic scheduling engine
func NewTrafficSchedulingEngine() *TrafficSchedulingEngine {
	engine := &TrafficSchedulingEngine{
		zones:             make(map[string]*TrafficZone),
		diversionRules:    make([]*entities.TrafficDiversionRule, 0),
		schedulingHistory: make([]*SchedulingDecision, 0),
		monitoringEnabled: true,
	}

	// Initialize default traffic zones for Longhua area
	engine.initializeDefaultZones()

	return engine
}

// initializeDefaultZones initializes default traffic monitoring zones
func (e *TrafficSchedulingEngine) initializeDefaultZones() {
	// Longhua Center Zone
	centerZone := &TrafficZone{
		ID:                "zone_longhua_center",
		Name:              "龙华中心区",
		CenterLat:         22.6913,
		CenterLng:         114.0448,
		Radius:            2.0,
		CurrentDensity:    0,
		HistoricalDensity: make([]float64, 0),
		PeakHourDensity:   0,
		Status:            "normal",
		LastUpdated:       time.Now(),
	}

	// Longhua North Zone
	northZone := &TrafficZone{
		ID:                "zone_longhua_north",
		Name:              "龙华北部区",
		CenterLat:         22.7100,
		CenterLng:         114.0500,
		Radius:            3.0,
		CurrentDensity:    0,
		HistoricalDensity: make([]float64, 0),
		PeakHourDensity:   0,
		Status:            "normal",
		LastUpdated:       time.Now(),
	}

	// Longhua South Zone
	southZone := &TrafficZone{
		ID:                "zone_longhua_south",
		Name:              "龙华南部区",
		CenterLat:         22.6700,
		CenterLng:         114.0400,
		Radius:            3.0,
		CurrentDensity:    0,
		HistoricalDensity: make([]float64, 0),
		PeakHourDensity:   0,
		Status:            "normal",
		LastUpdated:       time.Now(),
	}

	e.mu.Lock()
	e.zones[centerZone.ID] = centerZone
	e.zones[northZone.ID] = northZone
	e.zones[southZone.ID] = southZone
	e.mu.Unlock()
}

// UpdateZoneDensity updates the traffic density for a zone
func (e *TrafficSchedulingEngine) UpdateZoneDensity(zoneID string, density float64) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	zone, exists := e.zones[zoneID]
	if !exists {
		return fmt.Errorf("zone not found: %s", zoneID)
	}

	// Update current density
	zone.CurrentDensity = density

	// Update historical density (keep last 24 hours of data)
	zone.HistoricalDensity = append(zone.HistoricalDensity, density)
	if len(zone.HistoricalDensity) > 24 {
		zone.HistoricalDensity = zone.HistoricalDensity[1:]
	}

	// Update peak hour density
	if density > zone.PeakHourDensity {
		zone.PeakHourDensity = density
	}

	// Update zone status based on density
	zone.Status = e.determineZoneStatus(density)
	zone.LastUpdated = time.Now()

	return nil
}

// determineZoneStatus determines zone status based on density
func (e *TrafficSchedulingEngine) determineZoneStatus(density float64) string {
	if density < 50 {
		return "normal"
	} else if density < 100 {
		return "congested"
	}
	return "severe"
}

// AnalyzeTrafficPattern analyzes traffic patterns and predicts congestion
func (e *TrafficSchedulingEngine) AnalyzeTrafficPattern(zoneID string) (*TrafficPatternAnalysis, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	zone, exists := e.zones[zoneID]
	if !exists {
		return nil, fmt.Errorf("zone not found: %s", zoneID)
	}

	analysis := &TrafficPatternAnalysis{
		ZoneID:           zoneID,
		CurrentDensity:   zone.CurrentDensity,
		AverageDensity:   e.calculateAverageDensity(zone.HistoricalDensity),
		Trend:            e.calculateTrend(zone.HistoricalDensity),
		PeakHour:         e.predictPeakHour(zone.HistoricalDensity),
		PredictedDensity: e.predictNextHourDensity(zone.HistoricalDensity),
		RiskLevel:        e.calculateRiskLevel(zone),
		Recommendations:  e.generateRecommendations(zone),
		AnalyzedAt:       time.Now(),
	}

	return analysis, nil
}

// TrafficPatternAnalysis represents traffic pattern analysis results
type TrafficPatternAnalysis struct {
	ZoneID           string    `json:"zone_id"`
	CurrentDensity   float64   `json:"current_density"`
	AverageDensity   float64   `json:"average_density"`
	Trend            string    `json:"trend"`     // increasing/decreasing/stable
	PeakHour         int       `json:"peak_hour"` // 0-23
	PredictedDensity float64   `json:"predicted_density"`
	RiskLevel        string    `json:"risk_level"` // low/medium/high
	Recommendations  []string  `json:"recommendations"`
	AnalyzedAt       time.Time `json:"analyzed_at"`
}

// calculateAverageDensity calculates average density from historical data
func (e *TrafficSchedulingEngine) calculateAverageDensity(densities []float64) float64 {
	if len(densities) == 0 {
		return 0
	}

	sum := 0.0
	for _, d := range densities {
		sum += d
	}
	return sum / float64(len(densities))
}

// calculateTrend calculates traffic trend
func (e *TrafficSchedulingEngine) calculateTrend(densities []float64) string {
	if len(densities) < 3 {
		return "stable"
	}

	recent := densities[len(densities)-3:]
	if recent[2] > recent[1] && recent[1] > recent[0] {
		return "increasing"
	} else if recent[2] < recent[1] && recent[1] < recent[0] {
		return "decreasing"
	}
	return "stable"
}

// predictPeakHour predicts peak traffic hour
func (e *TrafficSchedulingEngine) predictPeakHour(densities []float64) int {
	// Simplified: return current hour + 2 (typical peak delay)
	return (time.Now().Hour() + 2) % 24
}

// predictNextHourDensity predicts density for next hour
func (e *TrafficSchedulingEngine) predictNextHourDensity(densities []float64) float64 {
	if len(densities) == 0 {
		return 0
	}

	// Simple linear prediction
	avg := e.calculateAverageDensity(densities)
	trend := e.calculateTrend(densities)

	if trend == "increasing" {
		return avg * 1.1
	} else if trend == "decreasing" {
		return avg * 0.9
	}
	return avg
}

// calculateRiskLevel calculates congestion risk level
func (e *TrafficSchedulingEngine) calculateRiskLevel(zone *TrafficZone) string {
	if zone.CurrentDensity > 100 {
		return "high"
	} else if zone.CurrentDensity > 50 {
		return "medium"
	}
	return "low"
}

// generateRecommendations generates traffic management recommendations
func (e *TrafficSchedulingEngine) generateRecommendations(zone *TrafficZone) []string {
	recommendations := make([]string, 0)

	switch zone.Status {
	case "normal":
		recommendations = append(recommendations, "维持当前交通信号配时")
		recommendations = append(recommendations, "监控车流变化趋势")
	case "congested":
		recommendations = append(recommendations, "启动一级分流预案")
		recommendations = append(recommendations, "引导车辆至外围配套停车场")
		recommendations = append(recommendations, "增加信号灯绿灯时长")
	case "severe":
		recommendations = append(recommendations, "启动二级分流预案")
		recommendations = append(recommendations, "引导车辆至路边临停区域")
		recommendations = append(recommendations, "实施交通管制")
		recommendations = append(recommendations, "通知相关部门协同疏导")
	}

	return recommendations
}

// MakeSchedulingDecision makes a traffic scheduling decision based on current conditions
func (e *TrafficSchedulingEngine) MakeSchedulingDecision(zoneID string) (*SchedulingDecision, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	zone, exists := e.zones[zoneID]
	if !exists {
		return nil, fmt.Errorf("zone not found: %s", zoneID)
	}

	// Analyze current conditions
	analysis, err := e.AnalyzeTrafficPattern(zoneID)
	if err != nil {
		return nil, err
	}

	// Make decision based on risk level
	var decision *SchedulingDecision
	switch analysis.RiskLevel {
	case "high":
		decision = e.createDiversionDecision(zone, analysis)
	case "medium":
		decision = e.createReroutingDecision(zone, analysis)
	default:
		decision = e.createMonitoringDecision(zone, analysis)
	}

	// Record decision
	e.schedulingHistory = append(e.schedulingHistory, decision)

	return decision, nil
}

// createDiversionDecision creates a traffic diversion decision
func (e *TrafficSchedulingEngine) createDiversionDecision(zone *TrafficZone, analysis *TrafficPatternAnalysis) *SchedulingDecision {
	decision := &SchedulingDecision{
		ID:              fmt.Sprintf("decision_%d", time.Now().UnixNano()),
		ZoneID:          zone.ID,
		DecisionType:    "diversion",
		TargetPoolLevel: entities.PoolLevelPeripheral,
		Reason:          fmt.Sprintf("区域车流密度过高(%.1f)，启动分流至外围停车场", zone.CurrentDensity),
		Confidence:      0.85,
		ExecutedAt:      time.Now(),
		Effectiveness:   0,
	}

	// If severe congestion, divert to roadside
	if zone.Status == "severe" {
		decision.TargetPoolLevel = entities.PoolLevelRoadside
		decision.Reason = fmt.Sprintf("区域严重拥堵(%.1f)，启动二级分流至路边临停", zone.CurrentDensity)
		decision.Confidence = 0.90
	}

	return decision
}

// createReroutingDecision creates a rerouting decision
func (e *TrafficSchedulingEngine) createReroutingDecision(zone *TrafficZone, analysis *TrafficPatternAnalysis) *SchedulingDecision {
	return &SchedulingDecision{
		ID:              fmt.Sprintf("decision_%d", time.Now().UnixNano()),
		ZoneID:          zone.ID,
		DecisionType:    "rerouting",
		TargetPoolLevel: entities.PoolLevelCore,
		Reason:          fmt.Sprintf("区域车流密度上升(%.1f)，建议优化路径规划", zone.CurrentDensity),
		Confidence:      0.70,
		ExecutedAt:      time.Now(),
		Effectiveness:   0,
	}
}

// createMonitoringDecision creates a monitoring decision
func (e *TrafficSchedulingEngine) createMonitoringDecision(zone *TrafficZone, analysis *TrafficPatternAnalysis) *SchedulingDecision {
	return &SchedulingDecision{
		ID:              fmt.Sprintf("decision_%d", time.Now().UnixNano()),
		ZoneID:          zone.ID,
		DecisionType:    "monitoring",
		TargetPoolLevel: entities.PoolLevelCore,
		Reason:          "区域车流正常，持续监控",
		Confidence:      0.95,
		ExecutedAt:      time.Now(),
		Effectiveness:   0,
	}
}

// UpdateDecisionEffectiveness updates the effectiveness of a decision
func (e *TrafficSchedulingEngine) UpdateDecisionEffectiveness(decisionID string, effectiveness float64) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, decision := range e.schedulingHistory {
		if decision.ID == decisionID {
			decision.Effectiveness = effectiveness
			return nil
		}
	}

	return fmt.Errorf("decision not found: %s", decisionID)
}

// GetZoneStatus returns the current status of all zones
func (e *TrafficSchedulingEngine) GetZoneStatus() map[string]*TrafficZone {
	e.mu.RLock()
	defer e.mu.RUnlock()

	status := make(map[string]*TrafficZone)
	for id, zone := range e.zones {
		status[id] = zone
	}
	return status
}

// GetSchedulingHistory returns the scheduling decision history
func (e *TrafficSchedulingEngine) GetSchedulingHistory(limit int) []*SchedulingDecision {
	e.mu.RLock()
	defer e.mu.RUnlock()

	history := e.schedulingHistory
	if limit > 0 && len(history) > limit {
		history = history[len(history)-limit:]
	}
	return history
}

// AddDiversionRule adds a traffic diversion rule
func (e *TrafficSchedulingEngine) AddDiversionRule(rule *entities.TrafficDiversionRule) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.diversionRules = append(e.diversionRules, rule)
	return nil
}

// GetDiversionRules returns all diversion rules
func (e *TrafficSchedulingEngine) GetDiversionRules() []*entities.TrafficDiversionRule {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.diversionRules
}
