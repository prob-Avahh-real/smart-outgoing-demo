package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"smart-outgoing-demo/internal/config"
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/services"
	"smart-outgoing-demo/internal/integration/parking"

	"github.com/gin-gonic/gin"
)

// ParkingHandler handles parking-related HTTP requests.
type ParkingHandler struct {
	parkingApp  *services.ParkingAppService
	config      *config.Config
	poolService *services.ParkingPoolService
	avpService  *services.AVPDispatchService
}

// NewParkingHandler creates a new parking handler.
func NewParkingHandler(parkingApp *services.ParkingAppService, cfg *config.Config) *ParkingHandler {
	return &ParkingHandler{
		parkingApp:  parkingApp,
		config:      cfg,
		poolService: services.NewParkingPoolService(),
		avpService:  services.NewAVPDispatchService(),
	}
}

func resolveUserID(c *gin.Context) string {
	userID := c.GetHeader("x-user-id")
	if userID == "" {
		return "demo_user"
	}
	return userID
}

func (h *ParkingHandler) cityBrainClient() *parking.CityBrainClient {
	cfg := &parking.CityBrainAPIConfig{
		BaseURL:      h.config.ParkingAPIBaseURL,
		AppID:        h.config.ParkingAPIAppID,
		AppSecret:    h.config.ParkingAPIAppSecret,
		ParkingLotNo: h.config.ParkingLotNo,
		PortNo:       h.config.PortNo,
		Timeout:      30 * time.Second,
		UseMock:      h.config.ParkingUseMock,
	}
	return parking.NewCityBrainClient(cfg)
}

// StartAVPRequest represents AVP auto-park task start payload.
type StartAVPRequest struct {
	VehicleID     string `json:"vehicle_id" binding:"required"`
	ParkingLotID  string `json:"parking_lot_id" binding:"required"`
	DropoffZone   string `json:"dropoff_zone" binding:"required"`
	TargetSpaceID string `json:"target_space_id,omitempty"`
}

// SummonAVPRequest represents AVP summon task payload.
type SummonAVPRequest struct {
	VehicleID    string `json:"vehicle_id" binding:"required"`
	ParkingLotID string `json:"parking_lot_id" binding:"required"`
	PickupZone   string `json:"pickup_zone" binding:"required"`
}

// FindParkingRequest represents a parking search request.
type FindParkingRequest struct {
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	MaxDistance float64 `json:"max_distance,omitempty"`
	MaxPrice    float64 `json:"max_price,omitempty"`
	Limit       int     `json:"limit,omitempty"`
}

// FindParking finds the best parking options based on user location and preferences.
func (h *ParkingHandler) FindParking(c *gin.Context) {
	var req FindParkingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recommendations, err := h.parkingApp.FindParking(req.Latitude, req.Longitude, req.MaxPrice, req.MaxDistance, req.Limit)
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

// ReserveSpaceRequest represents a reservation request.
type ReserveSpaceRequest struct {
	ParkingLotID string `json:"parking_lot_id" binding:"required"`
	SpaceID      string `json:"space_id" binding:"required"`
	StartTime    string `json:"start_time" binding:"required"`
	EndTime      string `json:"end_time" binding:"required"`
}

// ReserveSpace reserves a specific parking space.
func (h *ParkingHandler) ReserveSpace(c *gin.Context) {
	var req ReserveSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if endTime.Before(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}

	reservation, err := h.parkingApp.ReserveSpace(resolveUserID(c), req.ParkingLotID, req.SpaceID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"reservation": reservation,
		"message":     "Space reserved successfully",
	})
}

// StartParkingSessionRequest represents a parking session request.
type StartParkingSessionRequest struct {
	ParkingLotID string `json:"parking_lot_id" binding:"required"`
	SpaceID      string `json:"space_id" binding:"required"`
}

// StartParkingSession starts a parking session.
func (h *ParkingHandler) StartParkingSession(c *gin.Context) {
	var req StartParkingSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.parkingApp.StartSession(resolveUserID(c), req.ParkingLotID, req.SpaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session": session,
		"message": "Parking session started",
	})
}

// StartAVPTask creates a new autonomous valet parking task.
func (h *ParkingHandler) StartAVPTask(c *gin.Context) {
	var req StartAVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.avpService.StartAutoParkWithSessionCreation(
		h.parkingApp.MockService(),
		resolveUserID(c),
		req.VehicleID,
		req.ParkingLotID,
		req.DropoffZone,
		req.TargetSpaceID,
	)
	c.JSON(http.StatusCreated, gin.H{
		"task":            result.Task,
		"parking_session": result.Session,
		"message":         "AVP auto-park task started",
	})
}

// SummonAVPTask creates a summon task from parking slot to pickup zone.
func (h *ParkingHandler) SummonAVPTask(c *gin.Context) {
	var req SummonAVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := h.avpService.StartSummon(resolveUserID(c), req.VehicleID, req.ParkingLotID, req.PickupZone)
	c.JSON(http.StatusCreated, gin.H{
		"task":    task,
		"message": "AVP summon task started",
	})
}

// GetAVPTask returns AVP task status and progress.
func (h *ParkingHandler) GetAVPTask(c *gin.Context) {
	taskID := c.Param("id")
	task, ok := h.avpService.GetTask(taskID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// CancelAVPTask cancels an AVP task and triggers safe-stop semantics.
func (h *ParkingHandler) CancelAVPTask(c *gin.Context) {
	taskID := c.Param("id")
	task, err := h.avpService.CancelTask(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    task,
		"message": "AVP task cancelled and vehicle moved to minimum risk stop",
	})
}

// GetParkingLots returns all parking lots.
func (h *ParkingHandler) GetParkingLots(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	lots := h.parkingApp.ListLots(lat, lng)

	c.JSON(http.StatusOK, gin.H{
		"parking_lots": lots,
		"count":        len(lots),
	})
}

// GetParkingLot returns a specific parking lot.
func (h *ParkingHandler) GetParkingLot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"parking_lot": h.parkingApp.GetLot(c.Param("id")),
	})
}

// GetParkingSpaces returns available spaces in a parking lot.
func (h *ParkingHandler) GetParkingSpaces(c *gin.Context) {
	lotID := c.Param("id")
	spaces := h.parkingApp.GetSpaces(lotID)
	c.JSON(http.StatusOK, gin.H{
		"parking_spaces": spaces,
		"count":          len(spaces),
	})
}

// CityBrainEntryRequest represents vehicle entry request for City Brain API.
type CityBrainEntryRequest struct {
	PlateNo      string `json:"plate_no" binding:"required"`
	PortNo       string `json:"port_no" binding:"required"`
	ParkingLotNo string `json:"parking_lot_no" binding:"required"`
	VehicleType  string `json:"vehicle_type"`
}

// CityBrainExitRequest represents vehicle exit request for City Brain API.
type CityBrainExitRequest struct {
	PlateNo      string  `json:"plate_no" binding:"required"`
	PortNo       string  `json:"port_no" binding:"required"`
	ParkingLotNo string  `json:"parking_lot_no" binding:"required"`
	ParkingFee   float64 `json:"parking_fee"`
}

// ReportCityBrainEntry reports vehicle entry to City Brain API.
func (h *ParkingHandler) ReportCityBrainEntry(c *gin.Context) {
	var req CityBrainEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entryReq := &parking.ParkingEntryRequest{
		PlateNo:      req.PlateNo,
		EntryTime:    time.Now(),
		PortNo:       req.PortNo,
		ParkingLotNo: req.ParkingLotNo,
		OutOrderNo:   fmt.Sprintf("ORDER%d", time.Now().UnixNano()),
		VehicleType:  req.VehicleType,
	}

	resp, err := h.cityBrainClient().ReportVehicleEntry(entryReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ReportCityBrainExit reports vehicle exit to City Brain API.
func (h *ParkingHandler) ReportCityBrainExit(c *gin.Context) {
	var req CityBrainExitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exitReq := &parking.ParkingExitRequest{
		PlateNo:      req.PlateNo,
		ExitTime:     time.Now(),
		PortNo:       req.PortNo,
		ParkingLotNo: req.ParkingLotNo,
		OutOrderNo:   fmt.Sprintf("ORDER%d", time.Now().UnixNano()),
		ParkingFee:   req.ParkingFee,
	}

	resp, err := h.cityBrainClient().ReportVehicleExit(exitReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SendCityBrainHeartbeat sends heartbeat to City Brain API.
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

	heartbeatReq := &parking.HeartbeatRequest{
		ParkingLotNo:    h.config.ParkingLotNo,
		TotalSpaces:     req.TotalSpaces,
		AvailableSpaces: req.AvailableSpaces,
		HeartbeatTime:   time.Now(),
	}

	resp, err := h.cityBrainClient().SendHeartbeat(heartbeatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetParkingPoolStats returns statistics for all parking pools.
func (h *ParkingHandler) GetParkingPoolStats(c *gin.Context) {
	stats := h.poolService.GetPoolStatistics()
	c.JSON(http.StatusOK, gin.H{"pools": stats})
}

// AddLotToPoolRequest represents request to add a parking lot to a pool.
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

// AddLotToPool adds a parking lot to the appropriate pool.
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

// GetPoolRecommendationRequest represents request for pool-based recommendation.
type GetPoolRecommendationRequest struct {
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	MaxDistance float64 `json:"max_distance"`
}

// GetPoolRecommendation gets parking recommendation using three-tier pool logic.
func (h *ParkingHandler) GetPoolRecommendation(c *gin.Context) {
	var req GetPoolRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lot, reason, err := h.poolService.GetRecommendedParkingLot(req.Latitude, req.Longitude, req.MaxDistance)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"parking_lot": lot,
		"reason":      reason,
	})
}

// TriggerDiversionRequest represents request to trigger traffic diversion.
type TriggerDiversionRequest struct {
	SourceZone     string  `json:"source_zone" binding:"required"`
	CurrentDensity float64 `json:"current_density" binding:"required"`
}

// TriggerTrafficDiversion triggers traffic diversion based on current traffic density.
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
