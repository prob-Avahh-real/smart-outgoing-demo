package handlers

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"smart-outgoing-demo/internal/config"
	"smart-outgoing-demo/internal/domain/services"

	"github.com/gin-gonic/gin"
)

// TrafficHandler handles traffic-related HTTP requests
type TrafficHandler struct {
	schedulingEngine *services.TrafficSchedulingEngine
	densityAnalyzer  *services.TrafficDensityAnalyzer
	cv2xSimulation   *services.CV2XSimulation
	config           *config.Config
}

// NewTrafficHandler creates a new traffic handler
func NewTrafficHandler(cfg *config.Config) *TrafficHandler {
	return &TrafficHandler{
		schedulingEngine: services.NewTrafficSchedulingEngine(),
		densityAnalyzer:  services.NewTrafficDensityAnalyzer(),
		cv2xSimulation:   services.NewCV2XSimulation(),
		config:           cfg,
	}
}

// UpdateZoneDensityRequest represents a zone density update request
type UpdateZoneDensityRequest struct {
	ZoneID   string  `json:"zone_id" binding:"required"`
	Density  float64 `json:"density" binding:"required"`
	Vehicles int     `json:"vehicles"`
	Speed    float64 `json:"speed"`
}

// UpdateZoneDensity updates traffic density for a zone
func (h *TrafficHandler) UpdateZoneDensity(c *gin.Context) {
	var req UpdateZoneDensityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update density in analyzer
	err := h.densityAnalyzer.RecordDensity(req.ZoneID, req.Density, req.Vehicles, req.Speed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update density in scheduling engine
	err = h.schedulingEngine.UpdateZoneDensity(req.ZoneID, req.Density)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Zone density updated successfully",
		"zone_id": req.ZoneID,
		"density": req.Density,
	})
}

// GetZoneStatus returns status of all traffic zones
func (h *TrafficHandler) GetZoneStatus(c *gin.Context) {
	zones := h.schedulingEngine.GetZoneStatus()
	c.JSON(http.StatusOK, gin.H{
		"zones": zones,
		"count": len(zones),
	})
}

// AnalyzeTrafficPattern analyzes traffic pattern for a zone
func (h *TrafficHandler) AnalyzeTrafficPattern(c *gin.Context) {
	zoneID := c.Param("zone_id")
	if zoneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "zone_id is required"})
		return
	}

	analysis, err := h.schedulingEngine.AnalyzeTrafficPattern(zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// MakeSchedulingDecision makes a traffic scheduling decision
func (h *TrafficHandler) MakeSchedulingDecision(c *gin.Context) {
	zoneID := c.Param("zone_id")
	if zoneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "zone_id is required"})
		return
	}

	decision, err := h.schedulingEngine.MakeSchedulingDecision(zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, decision)
}

// GetSchedulingHistory returns scheduling decision history
func (h *TrafficHandler) GetSchedulingHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	history := h.schedulingEngine.GetSchedulingHistory(limit)
	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"count":   len(history),
	})
}

// AnalyzeDensity performs density analysis
func (h *TrafficHandler) AnalyzeDensity(c *gin.Context) {
	zoneID := c.Param("zone_id")
	if zoneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "zone_id is required"})
		return
	}

	result, err := h.densityAnalyzer.AnalyzeDensity(zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetDensityData returns density data for a zone
func (h *TrafficHandler) GetDensityData(c *gin.Context) {
	zoneID := c.Param("zone_id")
	if zoneID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "zone_id is required"})
		return
	}

	data, err := h.densityAnalyzer.GetZoneDensityData(zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetAllDensityData returns density data for all zones
func (h *TrafficHandler) GetAllDensityData(c *gin.Context) {
	data := h.densityAnalyzer.GetAllZoneData()
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"count": len(data),
	})
}

// SetAlertThreshold sets alert threshold for a zone
func (h *TrafficHandler) SetAlertThreshold(c *gin.Context) {
	var req struct {
		ZoneID    string  `json:"zone_id" binding:"required"`
		Threshold float64 `json:"threshold" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.densityAnalyzer.SetAlertThreshold(req.ZoneID, req.Threshold)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Alert threshold set successfully",
		"zone_id":   req.ZoneID,
		"threshold": req.Threshold,
	})
}

// GetAlertThresholds returns all alert thresholds
func (h *TrafficHandler) GetAlertThresholds(c *gin.Context) {
	thresholds := h.densityAnalyzer.GetAlertThresholds()
	c.JSON(http.StatusOK, thresholds)
}

// TrainPredictionModel trains the prediction model
func (h *TrafficHandler) TrainPredictionModel(c *gin.Context) {
	err := h.densityAnalyzer.TrainPredictionModel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prediction model trained successfully",
	})
}

// RegisterV2XVehicleRequest represents a V2X vehicle registration request
type RegisterV2XVehicleRequest struct {
	ID       string           `json:"id" binding:"required"`
	Type     string           `json:"type" binding:"required"`
	Position *PositionRequest `json:"position" binding:"required"`
	Speed    float64          `json:"speed"`
	Heading  float64          `json:"heading"`
}

// PositionRequest represents a position request
type PositionRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Altitude  float64 `json:"altitude"`
}

// RegisterV2XVehicle registers a V2X vehicle
func (h *TrafficHandler) RegisterV2XVehicle(c *gin.Context) {
	var req RegisterV2XVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicle := &services.V2XVehicle{
		ID:   req.ID,
		Type: req.Type,
		Position: &services.Position{
			Latitude:  req.Position.Latitude,
			Longitude: req.Position.Longitude,
			Altitude:  req.Position.Altitude,
		},
		Speed:      req.Speed,
		Heading:    req.Heading,
		V2XEnabled: true,
		LastSeen:   time.Now(),
	}

	err := h.cv2xSimulation.RegisterVehicle(vehicle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "V2X vehicle registered successfully",
		"vehicle": vehicle,
	})
}

// UnregisterV2XVehicle unregisters a V2X vehicle
func (h *TrafficHandler) UnregisterV2XVehicle(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	err := h.cv2xSimulation.UnregisterVehicle(vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "V2X vehicle unregistered successfully",
		"vehicle_id": vehicleID,
	})
}

// UpdateV2XVehiclePosition updates V2X vehicle position
func (h *TrafficHandler) UpdateV2XVehiclePosition(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	var req struct {
		Position *PositionRequest `json:"position" binding:"required"`
		Speed    float64          `json:"speed"`
		Heading  float64          `json:"heading"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	position := &services.Position{
		Latitude:  req.Position.Latitude,
		Longitude: req.Position.Longitude,
		Altitude:  req.Position.Altitude,
	}

	err := h.cv2xSimulation.UpdateVehiclePosition(vehicleID, position, req.Speed, req.Heading)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Vehicle position updated successfully",
		"vehicle_id": vehicleID,
	})
}

// GetV2XVehicles returns all V2X vehicles
func (h *TrafficHandler) GetV2XVehicles(c *gin.Context) {
	vehicles := h.cv2xSimulation.GetVehicles()
	c.JSON(http.StatusOK, gin.H{
		"vehicles": vehicles,
		"count":    len(vehicles),
	})
}

// GetV2XRSUs returns all RSUs
func (h *TrafficHandler) GetV2XRSUs(c *gin.Context) {
	rsus := h.cv2xSimulation.GetRSUs()
	c.JSON(http.StatusOK, gin.H{
		"rsus":  rsus,
		"count": len(rsus),
	})
}

// GetV2XMessages returns V2X message history
func (h *TrafficHandler) GetV2XMessages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	messages := h.cv2xSimulation.GetMessages(limit)
	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"count":    len(messages),
	})
}

// SendV2XBSM sends a Basic Safety Message
func (h *TrafficHandler) SendV2XBSM(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	var req struct {
		Position  *PositionRequest `json:"position" binding:"required"`
		Speed     float64          `json:"speed"`
		Heading   float64          `json:"heading"`
		Emergency bool             `json:"emergency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bsm := &services.BSM{
		VehicleID: vehicleID,
		Position: &services.Position{
			Latitude:  req.Position.Latitude,
			Longitude: req.Position.Longitude,
			Altitude:  req.Position.Altitude,
		},
		Speed:     req.Speed,
		Heading:   req.Heading,
		Timestamp: time.Now(),
		Emergency: req.Emergency,
	}

	message, err := h.cv2xSimulation.SendBSM(vehicleID, bsm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "BSM sent successfully",
		"bsm":     bsm,
		"msg_id":  message.ID,
	})
}

// SimulateTrafficHazard simulates a traffic hazard
func (h *TrafficHandler) SimulateTrafficHazard(c *gin.Context) {
	var req struct {
		Position  *PositionRequest `json:"position" binding:"required"`
		EventType string           `json:"event_type" binding:"required"`
		Severity  string           `json:"severity" binding:"required"`
		Duration  int              `json:"duration"` // minutes
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	position := &services.Position{
		Latitude:  req.Position.Latitude,
		Longitude: req.Position.Longitude,
		Altitude:  req.Position.Altitude,
	}

	duration := time.Duration(req.Duration) * time.Minute
	denm, err := h.cv2xSimulation.SimulateTrafficHazard(position, req.EventType, req.Severity, duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Traffic hazard simulated successfully",
		"denm":    denm,
	})
}

// GetV2XStatistics returns V2X simulation statistics
func (h *TrafficHandler) GetV2XStatistics(c *gin.Context) {
	stats := h.cv2xSimulation.GetStatistics()
	c.JSON(http.StatusOK, stats)
}

// GenerateRandomVehiclesRequest represents request for generating random vehicles
type GenerateRandomVehiclesRequest struct {
	Count     int     `json:"count" binding:"required"`
	CenterLat float64 `json:"center_lat"`
	CenterLng float64 `json:"center_lng"`
	RadiusKm  float64 `json:"radius_km"`
}

// GenerateRandomVehicles generates random V2X vehicles
func (h *TrafficHandler) GenerateRandomVehicles(c *gin.Context) {
	var req GenerateRandomVehiclesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use Longhua center as default if not provided
	if req.CenterLat == 0 {
		req.CenterLat = 22.6913
	}
	if req.CenterLng == 0 {
		req.CenterLng = 114.0448
	}
	if req.RadiusKm == 0 {
		req.RadiusKm = 5.0
	}

	vehicles, err := h.cv2xSimulation.GenerateRandomVehicles(req.Count, req.CenterLat, req.CenterLng, req.RadiusKm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Random vehicles generated successfully",
		"count":    len(vehicles),
		"vehicles": vehicles,
	})
}

// GenerateTrafficScenario generates a traffic scenario
func (h *TrafficHandler) GenerateTrafficScenario(c *gin.Context) {
	scenarioType := c.Param("type")
	if scenarioType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scenario type is required"})
		return
	}

	scenario, err := h.cv2xSimulation.GenerateTrafficScenario(scenarioType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Traffic scenario generated successfully",
		"scenario": scenario,
	})
}

// ClearVehicles clears all vehicles from simulation
func (h *TrafficHandler) ClearVehicles(c *gin.Context) {
	err := h.cv2xSimulation.ClearAllVehicles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All vehicles cleared successfully",
	})
}

// GetAppURL returns the current application URL for QR code generation on frontend
func (h *TrafficHandler) GetAppURL(c *gin.Context) {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// Use LAN IP address for mobile access
	// TODO: Make this configurable via environment variable
	host := "192.168.3.19:8080"

	url := scheme + "://" + host

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// getLocalIP returns the local LAN IP address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("[DEBUG] getLocalIP error: %v\n", err)
		return ""
	}
	fmt.Printf("[DEBUG] Found %d network interfaces\n", len(addrs))
	for _, addr := range addrs {
		fmt.Printf("[DEBUG] Checking address: %v\n", addr)
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// Return first non-loopback IPv4 address
				ip := ipnet.IP.String()
				fmt.Printf("[DEBUG] Found valid LAN IP: %s\n", ip)
				return ip
			}
		}
	}
	fmt.Printf("[DEBUG] No valid LAN IP found\n")
	return ""
}
