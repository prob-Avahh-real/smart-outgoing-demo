package services

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// TrafficDensityAnalyzer provides real-time traffic density analysis
type TrafficDensityAnalyzer struct {
	mu              sync.RWMutex
	zoneDensityData map[string]*DensityTimeSeries
	analysisWindow  time.Duration      // Time window for analysis
	updateInterval  time.Duration      // Update interval
	alertThresholds map[string]float64 // Alert thresholds by zone
	predictionModel *DensityPredictionModel
}

// DensityTimeSeries stores time series data for density
type DensityTimeSeries struct {
	ZoneID     string
	DataPoints []*DensityDataPoint
	MaxPoints  int
}

// DensityDataPoint represents a single density measurement
type DensityDataPoint struct {
	Timestamp time.Time
	Density   float64
	Vehicles  int
	Speed     float64 // Average speed km/h
}

// DensityPredictionModel predicts future density
type DensityPredictionModel struct {
	Trained     bool
	LastTrained time.Time
	ModelParams map[string]float64
	Accuracy    float64
}

// DensityAnalysisResult represents analysis results
type DensityAnalysisResult struct {
	ZoneID               string    `json:"zone_id"`
	CurrentDensity       float64   `json:"current_density"`
	AverageDensity       float64   `json:"average_density"`
	MaxDensity           float64   `json:"max_density"`
	MinDensity           float64   `json:"min_density"`
	Trend                string    `json:"trend"`       // increasing/decreasing/stable
	GrowthRate           float64   `json:"growth_rate"` // vehicles/hour
	PeakHour             int       `json:"peak_hour"`
	PredictedDensity     float64   `json:"predicted_density"`
	PredictionConfidence float64   `json:"prediction_confidence"`
	CongestionLevel      string    `json:"congestion_level"` // low/medium/high/severe
	AlertTriggered       bool      `json:"alert_triggered"`
	AlertMessage         string    `json:"alert_message"`
	Recommendations      []string  `json:"recommendations"`
	AnalyzedAt           time.Time `json:"analyzed_at"`
}

// NewTrafficDensityAnalyzer creates a new traffic density analyzer
func NewTrafficDensityAnalyzer() *TrafficDensityAnalyzer {
	analyzer := &TrafficDensityAnalyzer{
		zoneDensityData: make(map[string]*DensityTimeSeries),
		analysisWindow:  24 * time.Hour,
		updateInterval:  5 * time.Minute,
		alertThresholds: make(map[string]float64),
		predictionModel: &DensityPredictionModel{
			Trained:     false,
			ModelParams: make(map[string]float64),
			Accuracy:    0,
		},
	}

	// Set default alert thresholds
	analyzer.alertThresholds["zone_longhua_center"] = 80.0
	analyzer.alertThresholds["zone_longhua_north"] = 60.0
	analyzer.alertThresholds["zone_longhua_south"] = 60.0

	return analyzer
}

// RecordDensity records a density measurement
func (a *TrafficDensityAnalyzer) RecordDensity(zoneID string, density float64, vehicles int, speed float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Get or create time series
	series, exists := a.zoneDensityData[zoneID]
	if !exists {
		series = &DensityTimeSeries{
			ZoneID:     zoneID,
			DataPoints: make([]*DensityDataPoint, 0),
			MaxPoints:  288, // 24 hours * 12 points/hour (5-min intervals)
		}
		a.zoneDensityData[zoneID] = series
	}

	// Add new data point
	point := &DensityDataPoint{
		Timestamp: time.Now(),
		Density:   density,
		Vehicles:  vehicles,
		Speed:     speed,
	}

	series.DataPoints = append(series.DataPoints, point)

	// Trim old data points
	if len(series.DataPoints) > series.MaxPoints {
		series.DataPoints = series.DataPoints[len(series.DataPoints)-series.MaxPoints:]
	}

	return nil
}

// AnalyzeDensity performs comprehensive density analysis
func (a *TrafficDensityAnalyzer) AnalyzeDensity(zoneID string) (*DensityAnalysisResult, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	series, exists := a.zoneDensityData[zoneID]
	if !exists {
		return nil, fmt.Errorf("no data for zone: %s", zoneID)
	}

	if len(series.DataPoints) == 0 {
		return nil, fmt.Errorf("no data points for zone: %s", zoneID)
	}

	result := &DensityAnalysisResult{
		ZoneID:     zoneID,
		AnalyzedAt: time.Now(),
	}

	// Calculate statistics
	result.CurrentDensity = series.DataPoints[len(series.DataPoints)-1].Density
	result.AverageDensity = a.calculateAverage(series.DataPoints)
	result.MaxDensity = a.calculateMax(series.DataPoints)
	result.MinDensity = a.calculateMin(series.DataPoints)

	// Analyze trend
	result.Trend = a.analyzeTrend(series.DataPoints)
	result.GrowthRate = a.calculateGrowthRate(series.DataPoints)

	// Predict
	result.PeakHour = a.predictPeakHour(series.DataPoints)
	result.PredictedDensity, result.PredictionConfidence = a.predictDensity(series.DataPoints)

	// Determine congestion level
	result.CongestionLevel = a.determineCongestionLevel(result.CurrentDensity)

	// Check alert
	result.AlertTriggered, result.AlertMessage = a.checkAlert(zoneID, result.CurrentDensity)

	// Generate recommendations
	result.Recommendations = a.generateRecommendations(result)

	return result, nil
}

// calculateAverage calculates average density
func (a *TrafficDensityAnalyzer) calculateAverage(points []*DensityDataPoint) float64 {
	if len(points) == 0 {
		return 0
	}

	sum := 0.0
	for _, point := range points {
		sum += point.Density
	}
	return sum / float64(len(points))
}

// calculateMax calculates maximum density
func (a *TrafficDensityAnalyzer) calculateMax(points []*DensityDataPoint) float64 {
	if len(points) == 0 {
		return 0
	}

	max := points[0].Density
	for _, point := range points {
		if point.Density > max {
			max = point.Density
		}
	}
	return max
}

// calculateMin calculates minimum density
func (a *TrafficDensityAnalyzer) calculateMin(points []*DensityDataPoint) float64 {
	if len(points) == 0 {
		return 0
	}

	min := points[0].Density
	for _, point := range points {
		if point.Density < min {
			min = point.Density
		}
	}
	return min
}

// analyzeTrend analyzes density trend
func (a *TrafficDensityAnalyzer) analyzeTrend(points []*DensityDataPoint) string {
	if len(points) < 6 {
		return "stable"
	}

	// Compare recent 3 points with previous 3 points
	recent := points[len(points)-3:]
	previous := points[len(points)-6 : len(points)-3]

	recentAvg := a.calculateAverage(recent)
	previousAvg := a.calculateAverage(previous)

	changePercent := (recentAvg - previousAvg) / previousAvg * 100

	if changePercent > 10 {
		return "increasing"
	} else if changePercent < -10 {
		return "decreasing"
	}
	return "stable"
}

// calculateGrowthRate calculates growth rate (vehicles/hour)
func (a *TrafficDensityAnalyzer) calculateGrowthRate(points []*DensityDataPoint) float64 {
	if len(points) < 2 {
		return 0
	}

	latest := points[len(points)-1]
	earliest := points[0]

	duration := latest.Timestamp.Sub(earliest.Timestamp).Hours()
	if duration == 0 {
		return 0
	}

	return (latest.Density - earliest.Density) / duration
}

// predictPeakHour predicts peak traffic hour
func (a *TrafficDensityAnalyzer) predictPeakHour(points []*DensityDataPoint) int {
	if len(points) < 12 {
		return -1 // Not enough data
	}

	// Group by hour and find peak
	hourDensity := make(map[int][]float64)
	for _, point := range points {
		hour := point.Timestamp.Hour()
		hourDensity[hour] = append(hourDensity[hour], point.Density)
	}

	peakHour := 0
	maxAvgDensity := 0.0

	for hour, densities := range hourDensity {
		avg := 0.0
		for _, d := range densities {
			avg += d
		}
		avg /= float64(len(densities))

		if avg > maxAvgDensity {
			maxAvgDensity = avg
			peakHour = hour
		}
	}

	return peakHour
}

// predictDensity predicts density for next hour
func (a *TrafficDensityAnalyzer) predictDensity(points []*DensityDataPoint) (float64, float64) {
	if len(points) < 6 {
		return 0, 0 // Not enough data
	}

	// Simple linear regression prediction
	n := float64(len(points))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, point := range points {
		x := float64(i)
		y := point.Density
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate slope and intercept
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n

	// Predict next point
	nextX := float64(len(points))
	predicted := slope*nextX + intercept

	// Calculate confidence based on variance
	variance := 0.0
	for _, point := range points {
		x := float64(len(points) - 1) // Use latest index
		predictedAtX := slope*x + intercept
		variance += math.Pow(point.Density-predictedAtX, 2)
	}
	variance /= float64(len(points))

	// Confidence decreases with variance
	confidence := 1.0 / (1.0 + variance/100.0)
	if confidence > 1.0 {
		confidence = 1.0
	}

	return predicted, confidence
}

// determineCongestionLevel determines congestion level
func (a *TrafficDensityAnalyzer) determineCongestionLevel(density float64) string {
	if density < 30 {
		return "low"
	} else if density < 60 {
		return "medium"
	} else if density < 100 {
		return "high"
	}
	return "severe"
}

// checkAlert checks if alert should be triggered
func (a *TrafficDensityAnalyzer) checkAlert(zoneID string, density float64) (bool, string) {
	threshold, exists := a.alertThresholds[zoneID]
	if !exists {
		threshold = 80.0 // Default threshold
	}

	if density >= threshold {
		message := fmt.Sprintf("区域%s车流密度异常: %.1f (阈值: %.1f)", zoneID, density, threshold)
		return true, message
	}

	return false, ""
}

// generateRecommendations generates recommendations based on analysis
func (a *TrafficDensityAnalyzer) generateRecommendations(result *DensityAnalysisResult) []string {
	recommendations := make([]string, 0)

	switch result.CongestionLevel {
	case "low":
		recommendations = append(recommendations, "维持正常交通管理")
		recommendations = append(recommendations, "继续监控车流变化")
	case "medium":
		recommendations = append(recommendations, "增加交通监控频率")
		recommendations = append(recommendations, "准备分流预案")
		recommendations = append(recommendations, "优化信号灯配时")
	case "high":
		recommendations = append(recommendations, "启动一级分流预案")
		recommendations = append(recommendations, "引导车辆至外围停车场")
		recommendations = append(recommendations, "增加警力疏导")
	case "severe":
		recommendations = append(recommendations, "启动二级分流预案")
		recommendations = append(recommendations, "引导车辆至路边临停区域")
		recommendations = append(recommendations, "实施交通管制")
		recommendations = append(recommendations, "通知相关部门协同")
	}

	// Trend-based recommendations
	if result.Trend == "increasing" && result.GrowthRate > 10 {
		recommendations = append(recommendations, fmt.Sprintf("车流快速增长(%.1f辆/小时)，提前干预", result.GrowthRate))
	}

	// Prediction-based recommendations
	if result.PredictedDensity > result.CurrentDensity*1.2 && result.PredictionConfidence > 0.7 {
		recommendations = append(recommendations, fmt.Sprintf("预测1小时后车流将增长20%%以上，提前准备"))
	}

	return recommendations
}

// GetZoneDensityData returns density data for a zone
func (a *TrafficDensityAnalyzer) GetZoneDensityData(zoneID string) (*DensityTimeSeries, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	series, exists := a.zoneDensityData[zoneID]
	if !exists {
		return nil, fmt.Errorf("no data for zone: %s", zoneID)
	}

	return series, nil
}

// GetAllZoneData returns data for all zones
func (a *TrafficDensityAnalyzer) GetAllZoneData() map[string]*DensityTimeSeries {
	a.mu.RLock()
	defer a.mu.RUnlock()

	data := make(map[string]*DensityTimeSeries)
	for id, series := range a.zoneDensityData {
		data[id] = series
	}
	return data
}

// SetAlertThreshold sets alert threshold for a zone
func (a *TrafficDensityAnalyzer) SetAlertThreshold(zoneID string, threshold float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.alertThresholds[zoneID] = threshold
}

// GetAlertThresholds returns all alert thresholds
func (a *TrafficDensityAnalyzer) GetAlertThresholds() map[string]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()

	thresholds := make(map[string]float64)
	for id, threshold := range a.alertThresholds {
		thresholds[id] = threshold
	}
	return thresholds
}

// TrainPredictionModel trains the prediction model
func (a *TrafficDensityAnalyzer) TrainPredictionModel() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Simplified training: calculate model parameters from historical data
	allPoints := make([]*DensityDataPoint, 0)
	for _, series := range a.zoneDensityData {
		allPoints = append(allPoints, series.DataPoints...)
	}

	if len(allPoints) < 24 {
		return fmt.Errorf("insufficient data for training (need at least 24 points)")
	}

	// Calculate simple parameters
	a.predictionModel.ModelParams["avg_density"] = a.calculateAverage(allPoints)
	a.predictionModel.ModelParams["std_density"] = a.calculateStdDev(allPoints)
	a.predictionModel.ModelParams["trend_factor"] = a.calculateTrendFactor(allPoints)

	a.predictionModel.Trained = true
	a.predictionModel.LastTrained = time.Now()
	a.predictionModel.Accuracy = 0.75 // Estimated accuracy

	return nil
}

// calculateStdDev calculates standard deviation
func (a *TrafficDensityAnalyzer) calculateStdDev(points []*DensityDataPoint) float64 {
	if len(points) == 0 {
		return 0
	}

	avg := a.calculateAverage(points)
	variance := 0.0

	for _, point := range points {
		variance += math.Pow(point.Density-avg, 2)
	}
	variance /= float64(len(points))

	return math.Sqrt(variance)
}

// calculateTrendFactor calculates trend factor
func (a *TrafficDensityAnalyzer) calculateTrendFactor(points []*DensityDataPoint) float64 {
	if len(points) < 2 {
		return 0
	}

	first := points[0].Density
	last := points[len(points)-1].Density
	return (last - first) / first
}
