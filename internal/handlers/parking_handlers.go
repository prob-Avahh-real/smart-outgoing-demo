package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"smart-outgoing-demo/internal/config"
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/services"
	"smart-outgoing-demo/internal/integration/parking"

	"github.com/gin-gonic/gin"
)

// ParkingHandler handles parking-related HTTP requests
type ParkingHandler struct {
	recommendationService *services.ParkingRecommendationService
	config                *config.Config
	poolService           *services.ParkingPoolService
}

// NewParkingHandler creates a new parking handler
func NewParkingHandler(
	recommendationService *services.ParkingRecommendationService,
	cfg *config.Config,
) *ParkingHandler {
	return &ParkingHandler{
		recommendationService: recommendationService,
		config:                cfg,
		poolService:           services.NewParkingPoolService(),
	}
}

// FindParkingRequest represents a parking search request
type FindParkingRequest struct {
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	MaxDistance float64 `json:"max_distance,omitempty"`
	MaxPrice    float64 `json:"max_price,omitempty"`
	Limit       int     `json:"limit,omitempty"`
}

// FindParking finds the best parking options based on user location and preferences
func (h *ParkingHandler) FindParking(c *gin.Context) {
	var req FindParkingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If recommendation service is not available, return mock data
	if h.recommendationService == nil {
		mockRecommendations := generateMockRecommendations(req.Latitude, req.Longitude, req.Limit)
		c.JSON(http.StatusOK, gin.H{
			"recommendations": mockRecommendations,
			"count":           len(mockRecommendations),
			"search_location": gin.H{
				"latitude":  req.Latitude,
				"longitude": req.Longitude,
			},
		})
		return
	}

	// Get user preferences (for now, use defaults)
	preferences := &entities.UserParkingPreference{
		MaxPricePerHour:   req.MaxPrice,
		PreferredDistance: req.MaxDistance,
		VehicleType:       entities.SpaceTypeRegular,
		PreferredFeatures: []string{},
		PreferCovered:     false,
		PreferEV:          false,
	}

	// Find parking recommendations
	recommendations, err := h.recommendationService.FindBestParking(
		req.Latitude,
		req.Longitude,
		preferences,
		req.Limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
		"search_location": gin.H{
			"latitude":  req.Latitude,
			"longitude": req.Longitude,
		},
	})
}

// ReserveSpaceRequest represents a reservation request
type ReserveSpaceRequest struct {
	ParkingLotID string `json:"parking_lot_id" binding:"required"`
	SpaceID      string `json:"space_id" binding:"required"`
	StartTime    string `json:"start_time" binding:"required"` // ISO format
	EndTime      string `json:"end_time" binding:"required"`   // ISO format
}

// ReserveSpace reserves a specific parking space
func (h *ParkingHandler) ReserveSpace(c *gin.Context) {
	var req ReserveSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
		return
	}

	// Validate time range
	if endTime.Before(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}

	// Get user ID (for now, use a placeholder)
	userID := c.GetHeader("x-user-id")
	if userID == "" {
		userID = "demo_user"
	}

	// Debug: Check if service is nil
	if h.recommendationService != nil {
		fmt.Printf("Service is not nil, this is unexpected!\n")
		return
	}

	// Create reservation (mock implementation for now)
	reservation := &entities.ParkingReservation{
		ID:           fmt.Sprintf("res_%d", time.Now().UnixNano()),
		UserID:       userID,
		ParkingLotID: req.ParkingLotID,
		SpaceID:      req.SpaceID,
		StartTime:    startTime,
		EndTime:      endTime,
		Status:       entities.StatusConfirmed,
		TotalPrice:   30.0, // Mock price
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"reservation": reservation,
		"message":     "Space reserved successfully",
	})
}

// StartParkingSessionRequest represents a parking session request
type StartParkingSessionRequest struct {
	ParkingLotID string `json:"parking_lot_id" binding:"required"`
	SpaceID      string `json:"space_id" binding:"required"`
}

// StartParkingSession starts a parking session
func (h *ParkingHandler) StartParkingSession(c *gin.Context) {
	var req StartParkingSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID
	userID := c.GetHeader("x-user-id")
	if userID == "" {
		userID = "demo_user"
	}

	// Create parking session (mock implementation for now)
	session := &entities.ParkingSession{
		ID:           fmt.Sprintf("session_%d", time.Now().UnixNano()),
		UserID:       userID,
		ParkingLotID: req.ParkingLotID,
		SpaceID:      req.SpaceID,
		StartTime:    time.Now(),
		Status:       entities.SessionActive,
		TotalCost:    0.0, // Will be calculated when session ends
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"session": session,
		"message": "Parking session started (mock)",
	})
}

// GetParkingLots returns all parking lots
func (h *ParkingHandler) GetParkingLots(c *gin.Context) {
	// This would use the parking lot repository
	// For now, return mock data

	lats, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lngs, _ := strconv.ParseFloat(c.Query("lng"), 64)

	mockLots := []*entities.ParkingLot{
		{
			ID:              "lot_1",
			Name:            "CBD Central Parking",
			Address:         "123 Main Street",
			Latitude:        22.6913,
			Longitude:       114.0448,
			TotalSpaces:     200,
			AvailableSpaces: 45,
			PricePerHour:    15.0,
			Features:        []string{"covered", "24/7", "ev_charging"},
			Rating:          4.5,
			IsOpen:          true,
			LastUpdated:     time.Now(),
		},
		{
			ID:              "lot_2",
			Name:            "Shopping Mall Parking",
			Address:         "456 Shopping Ave",
			Latitude:        22.6950,
			Longitude:       114.0500,
			TotalSpaces:     150,
			AvailableSpaces: 12,
			PricePerHour:    10.0,
			Features:        []string{"covered", "security"},
			Rating:          4.2,
			IsOpen:          true,
			LastUpdated:     time.Now(),
		},
	}

	// Calculate distances if location provided
	if lats != 0 && lngs != 0 {
		for _, lot := range mockLots {
			lot.Distance = calculateDistance(lats, lngs, lot.Latitude, lot.Longitude)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"parking_lots": mockLots,
		"count":        len(mockLots),
	})
}

// GetParkingLot returns a specific parking lot
func (h *ParkingHandler) GetParkingLot(c *gin.Context) {
	lotID := c.Param("id")

	// Mock data - in real implementation, use repository
	mockLot := &entities.ParkingLot{
		ID:              lotID,
		Name:            "CBD Central Parking",
		Address:         "123 Main Street",
		Latitude:        22.6913,
		Longitude:       114.0448,
		TotalSpaces:     200,
		AvailableSpaces: 45,
		PricePerHour:    15.0,
		Features:        []string{"covered", "24/7", "ev_charging"},
		Rating:          4.5,
		IsOpen:          true,
		LastUpdated:     time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"parking_lot": mockLot,
	})
}

// GetParkingSpaces returns available spaces in a parking lot
func (h *ParkingHandler) GetParkingSpaces(c *gin.Context) {
	lotID := c.Param("id")

	// Mock data - in real implementation, use repository
	mockSpaces := []*entities.ParkingSpace{
		{
			ID:           "space_1",
			ParkingLotID: lotID,
			SpaceNumber:  "A-101",
			Type:         entities.SpaceTypeRegular,
			IsAvailable:  true,
			IsReserved:   false,
			Price:        15.0,
			Floor:        "A",
			Zone:         "North",
			LastUpdated:  time.Now(),
		},
		{
			ID:           "space_2",
			ParkingLotID: lotID,
			SpaceNumber:  "A-102",
			Type:         entities.SpaceTypeEV,
			IsAvailable:  true,
			IsReserved:   false,
			Price:        20.0,
			Floor:        "A",
			Zone:         "North",
			LastUpdated:  time.Now(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"parking_spaces": mockSpaces,
		"count":          len(mockSpaces),
	})
}

// generateMockRecommendations creates mock parking recommendations for demo
func generateMockRecommendations(userLat, userLng float64, limit int) []*entities.ParkingRecommendation {
	mockLots := []*entities.ParkingLot{
		{
			ID:              "lot_1",
			Name:            "CBD Central Parking",
			Address:         "123 Main Street",
			Latitude:        userLat + 0.001,
			Longitude:       userLng + 0.001,
			TotalSpaces:     200,
			AvailableSpaces: 45,
			PricePerHour:    15.0,
			Features:        []string{"covered", "24/7", "ev_charging"},
			Rating:          4.5,
			IsOpen:          true,
			LastUpdated:     time.Now(),
		},
		{
			ID:              "lot_2",
			Name:            "Shopping Mall Parking",
			Address:         "456 Shopping Ave",
			Latitude:        userLat - 0.001,
			Longitude:       userLng - 0.001,
			TotalSpaces:     150,
			AvailableSpaces: 12,
			PricePerHour:    10.0,
			Features:        []string{"covered", "security"},
			Rating:          4.2,
			IsOpen:          true,
			LastUpdated:     time.Now(),
		},
		{
			ID:              "lot_3",
			Name:            "Airport Parking",
			Address:         "789 Airport Road",
			Latitude:        userLat + 0.002,
			Longitude:       userLng - 0.002,
			TotalSpaces:     300,
			AvailableSpaces: 89,
			PricePerHour:    8.0,
			Features:        []string{"24/7", "security", "shuttle"},
			Rating:          4.0,
			IsOpen:          true,
			LastUpdated:     time.Now(),
		},
	}

	var recommendations []*entities.ParkingRecommendation

	for i, lot := range mockLots {
		// Calculate distance
		distance := calculateDistance(userLat, userLng, lot.Latitude, lot.Longitude)
		lot.Distance = distance

		// Create mock space
		space := &entities.ParkingSpace{
			ID:           fmt.Sprintf("space_%d", i+1),
			ParkingLotID: lot.ID,
			SpaceNumber:  fmt.Sprintf("A-%03d", i+1),
			Type:         entities.SpaceTypeRegular,
			IsAvailable:  true,
			IsReserved:   false,
			Price:        lot.PricePerHour,
			Floor:        "A",
			Zone:         "North",
			LastUpdated:  time.Now(),
		}

		// Create recommendation
		recommendation := &entities.ParkingRecommendation{
			ParkingLot:       lot,
			RecommendedSpace: space,
			Score:            85.0 - float64(i)*10, // Decreasing scores
			Reasons: []string{
				fmt.Sprintf("Only %.1f km away", distance),
				fmt.Sprintf("¥%.1f/hour", lot.PricePerHour),
				fmt.Sprintf("%d spaces available", lot.AvailableSpaces),
			},
			EstimatedTime: time.Duration(distance*10) * time.Minute,
			TotalPrice:    lot.PricePerHour * 2.0, // 2 hours estimate
			Route: &entities.ParkingRoute{
				Steps: []entities.RouteStep{
					{
						Instruction: fmt.Sprintf("Drive %.1f km to %s", distance, lot.Name),
						Distance:    distance,
						Duration:    time.Duration(distance*10) * time.Minute,
						Direction:   "Towards destination",
					},
				},
				TotalDistance: distance,
				TotalTime:     time.Duration(distance*10) * time.Minute,
				Instructions:  fmt.Sprintf("Navigate to %s, %s", lot.Name, lot.Address),
			},
		}

		recommendations = append(recommendations, recommendation)

		// Limit results
		if limit > 0 && len(recommendations) >= limit {
			break
		}
	}

	return recommendations
}

// Helper function
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Simple distance calculation (in km)
	const earthRadius = 6371.0

	dLat := (lat2 - lat1) * 3.14159 / 180.0
	dLng := (lng2 - lng1) * 3.14159 / 180.0

	a := (dLat/2)*(dLat/2) +
		(lat1*3.14159/180.0)*(lat2*3.14159/180.0)*
			(dLng/2)*(dLng/2)

	c := 2 * 3.14159 / 180.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// CityBrainEntryRequest represents vehicle entry request for City Brain API
type CityBrainEntryRequest struct {
	PlateNo      string `json:"plate_no" binding:"required"`
	PortNo       string `json:"port_no" binding:"required"`
	ParkingLotNo string `json:"parking_lot_no" binding:"required"`
	VehicleType  string `json:"vehicle_type"`
}

// CityBrainExitRequest represents vehicle exit request for City Brain API
type CityBrainExitRequest struct {
	PlateNo      string  `json:"plate_no" binding:"required"`
	PortNo       string  `json:"port_no" binding:"required"`
	ParkingLotNo string  `json:"parking_lot_no" binding:"required"`
	ParkingFee   float64 `json:"parking_fee"`
}

// ReportCityBrainEntry reports vehicle entry to City Brain API
func (h *ParkingHandler) ReportCityBrainEntry(c *gin.Context) {
	var req CityBrainEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate unique order number
	outOrderNo := fmt.Sprintf("ORDER%d", time.Now().UnixNano())

	// Build parking entry request
	entryReq := &parking.ParkingEntryRequest{
		PlateNo:      req.PlateNo,
		EntryTime:    time.Now(),
		PortNo:       req.PortNo,
		ParkingLotNo: req.ParkingLotNo,
		OutOrderNo:   outOrderNo,
		VehicleType:  req.VehicleType,
	}

	// Create City Brain client config
	config := &parking.CityBrainAPIConfig{
		BaseURL:      h.config.ParkingAPIBaseURL,
		AppID:        h.config.ParkingAPIAppID,
		AppSecret:    h.config.ParkingAPIAppSecret,
		ParkingLotNo: h.config.ParkingLotNo,
		PortNo:       h.config.PortNo,
		Timeout:      30 * time.Second,
		UseMock:      h.config.ParkingUseMock,
	}

	client := parking.NewCityBrainClient(config)
	resp, err := client.ReportVehicleEntry(entryReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ReportCityBrainExit reports vehicle exit to City Brain API
func (h *ParkingHandler) ReportCityBrainExit(c *gin.Context) {
	var req CityBrainExitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate unique order number
	outOrderNo := fmt.Sprintf("ORDER%d", time.Now().UnixNano())

	// Build parking exit request
	exitReq := &parking.ParkingExitRequest{
		PlateNo:      req.PlateNo,
		ExitTime:     time.Now(),
		PortNo:       req.PortNo,
		ParkingLotNo: req.ParkingLotNo,
		OutOrderNo:   outOrderNo,
		ParkingFee:   req.ParkingFee,
	}

	// Create City Brain client config
	config := &parking.CityBrainAPIConfig{
		BaseURL:      h.config.ParkingAPIBaseURL,
		AppID:        h.config.ParkingAPIAppID,
		AppSecret:    h.config.ParkingAPIAppSecret,
		ParkingLotNo: h.config.ParkingLotNo,
		PortNo:       h.config.PortNo,
		Timeout:      30 * time.Second,
		UseMock:      h.config.ParkingUseMock,
	}

	client := parking.NewCityBrainClient(config)
	resp, err := client.ReportVehicleExit(exitReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SendCityBrainHeartbeat sends heartbeat to City Brain API
func (h *ParkingHandler) SendCityBrainHeartbeat(c *gin.Context) {
	type HeartbeatRequest struct {
		TotalSpaces     int `json:"total_spaces" binding:"required"`
		AvailableSpaces int `json:"available_spaces" binding:"required"`
	}

	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build heartbeat request
	heartbeatReq := &parking.HeartbeatRequest{
		ParkingLotNo:    h.config.ParkingLotNo,
		TotalSpaces:     req.TotalSpaces,
		AvailableSpaces: req.AvailableSpaces,
		HeartbeatTime:   time.Now(),
	}

	// Create City Brain client config
	config := &parking.CityBrainAPIConfig{
		BaseURL:      h.config.ParkingAPIBaseURL,
		AppID:        h.config.ParkingAPIAppID,
		AppSecret:    h.config.ParkingAPIAppSecret,
		ParkingLotNo: h.config.ParkingLotNo,
		PortNo:       h.config.PortNo,
		Timeout:      30 * time.Second,
		UseMock:      h.config.ParkingUseMock,
	}

	client := parking.NewCityBrainClient(config)
	resp, err := client.SendHeartbeat(heartbeatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetParkingPoolStats returns statistics for all parking pools
func (h *ParkingHandler) GetParkingPoolStats(c *gin.Context) {
	stats := h.poolService.GetPoolStatistics()
	c.JSON(http.StatusOK, gin.H{"pools": stats})
}

// AddLotToPoolRequest represents request to add a parking lot to a pool
type AddLotToPoolRequest struct {
	ID           string   `json:"id" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Address      string   `json:"address" binding:"required"`
	Latitude     float64  `json:"latitude" binding:"required"`
	Longitude    float64  `json:"longitude" binding:"required"`
	TotalSpaces  int      `json:"total_spaces" binding:"required"`
	PricePerHour float64  `json:"price_per_hour"`
	Features     []string `json:"features"`
}

// AddLotToPool adds a parking lot to the appropriate pool
func (h *ParkingHandler) AddLotToPool(c *gin.Context) {
	var req AddLotToPoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lot := &entities.ParkingLot{
		ID:              req.ID,
		Name:            req.Name,
		Address:         req.Address,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		TotalSpaces:     req.TotalSpaces,
		AvailableSpaces: req.TotalSpaces,
		PricePerHour:    req.PricePerHour,
		Features:        req.Features,
		IsOpen:          true,
		LastUpdated:     time.Now(),
	}

	if err := h.poolService.AddParkingLotToPool(lot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parking lot added to pool successfully"})
}

// GetPoolRecommendationRequest represents request for pool-based recommendation
type GetPoolRecommendationRequest struct {
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	MaxDistance float64 `json:"max_distance"`
}

// GetPoolRecommendation gets parking recommendation using three-tier pool logic
func (h *ParkingHandler) GetPoolRecommendation(c *gin.Context) {
	var req GetPoolRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lot, reason, err := h.poolService.GetRecommendedParkingLot(
		req.Latitude,
		req.Longitude,
		req.MaxDistance,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"parking_lot": lot,
		"reason":      reason,
	})
}

// TriggerDiversionRequest represents request to trigger traffic diversion
type TriggerDiversionRequest struct {
	SourceZone     string  `json:"source_zone" binding:"required"`
	CurrentDensity float64 `json:"current_density" binding:"required"`
}

// TriggerTrafficDiversion triggers traffic diversion based on current traffic density
func (h *ParkingHandler) TriggerTrafficDiversion(c *gin.Context) {
	var req TriggerDiversionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lots, err := h.poolService.TriggerTrafficDiversion(req.SourceZone, req.CurrentDensity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommended_lots": lots,
		"count":            len(lots),
	})
}
