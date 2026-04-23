package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// CV2XSimulation provides C-V2X (Cellular Vehicle-to-Everything) simulation
type CV2XSimulation struct {
	mu                sync.RWMutex
	vehicles          map[string]*V2XVehicle
	rsus              map[string]*RSU
	messages          []*V2XMessage
	simulationEnabled bool
	messageBufferSize int
}

// V2XVehicle represents a V2X-enabled vehicle
type V2XVehicle struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"` // car/bus/truck/emergency
	Position         *Position `json:"position"`
	Speed            float64   `json:"speed"`   // km/h
	Heading          float64   `json:"heading"` // degrees
	V2XEnabled       bool      `json:"v2x_enabled"`
	LastSeen         time.Time `json:"last_seen"`
	MessagesSent     int       `json:"messages_sent"`
	MessagesReceived int       `json:"messages_received"`
}

// Position represents geographic position
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

// RSU represents a Road Side Unit
type RSU struct {
	ID                string    `json:"id"`
	Position          *Position `json:"position"`
	Range             float64   `json:"range"`  // meters
	Type              string    `json:"type"`   // traffic_light/parking_sensor/speed_camera
	Status            string    `json:"status"` // active/inactive
	ConnectedVehicles []string  `json:"connected_vehicles"`
	LastUpdated       time.Time `json:"last_updated"`
}

// V2XMessage represents a V2X message
type V2XMessage struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // bsm/cam/denm/map/spat
	SourceID  string    `json:"source_id"`
	TargetID  string    `json:"target_id"` // empty for broadcast
	Payload   []byte    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
	Priority  int       `json:"priority"` // 0-7, 0=highest
	Received  bool      `json:"received"`
	Processed bool      `json:"processed"`
}

// BSM (Basic Safety Message) - Vehicle safety information
type BSM struct {
	VehicleID string    `json:"vehicle_id"`
	Position  *Position `json:"position"`
	Speed     float64   `json:"speed"`
	Heading   float64   `json:"heading"`
	Timestamp time.Time `json:"timestamp"`
	Emergency bool      `json:"emergency"`
}

// CAM (Cooperative Awareness Message) - Vehicle status
type CAM struct {
	VehicleID    string    `json:"vehicle_id"`
	Position     *Position `json:"position"`
	Speed        float64   `json:"speed"`
	Heading      float64   `json:"heading"`
	Acceleration float64   `json:"acceleration"`
	Timestamp    time.Time `json:"timestamp"`
}

// DENM (Decentralized Environmental Notification Message) - Hazard warning
type DENM struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"` // accident/road_work/weather
	Position  *Position `json:"position"`
	Radius    float64   `json:"radius"`   // meters
	Severity  string    `json:"severity"` // low/medium/high
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// NewCV2XSimulation creates a new C-V2X simulation
func NewCV2XSimulation() *CV2XSimulation {
	sim := &CV2XSimulation{
		vehicles:          make(map[string]*V2XVehicle),
		rsus:              make(map[string]*RSU),
		messages:          make([]*V2XMessage, 0),
		simulationEnabled: true,
		messageBufferSize: 1000,
	}

	// Initialize default RSUs for Longhua area
	sim.initializeDefaultRSUs()

	return sim
}

// initializeDefaultRSUs initializes default RSUs
func (s *CV2XSimulation) initializeDefaultRSUs() {
	// Traffic light RSU at intersection
	rsu1 := &RSU{
		ID:                "rsu_traffic_001",
		Position:          &Position{Latitude: 22.6913, Longitude: 114.0448, Altitude: 50},
		Range:             500,
		Type:              "traffic_light",
		Status:            "active",
		ConnectedVehicles: make([]string, 0),
		LastUpdated:       time.Now(),
	}

	// Parking sensor RSU
	rsu2 := &RSU{
		ID:                "rsu_parking_001",
		Position:          &Position{Latitude: 22.6920, Longitude: 114.0450, Altitude: 50},
		Range:             300,
		Type:              "parking_sensor",
		Status:            "active",
		ConnectedVehicles: make([]string, 0),
		LastUpdated:       time.Now(),
	}

	// Speed camera RSU
	rsu3 := &RSU{
		ID:                "rsu_speed_001",
		Position:          &Position{Latitude: 22.6900, Longitude: 114.0430, Altitude: 50},
		Range:             400,
		Type:              "speed_camera",
		Status:            "active",
		ConnectedVehicles: make([]string, 0),
		LastUpdated:       time.Now(),
	}

	s.mu.Lock()
	s.rsus[rsu1.ID] = rsu1
	s.rsus[rsu2.ID] = rsu2
	s.rsus[rsu3.ID] = rsu3
	s.mu.Unlock()
}

// RegisterVehicle registers a V2X-enabled vehicle
func (s *CV2XSimulation) RegisterVehicle(vehicle *V2XVehicle) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.vehicles[vehicle.ID]; exists {
		return fmt.Errorf("vehicle already registered: %s", vehicle.ID)
	}

	vehicle.V2XEnabled = true
	vehicle.LastSeen = time.Now()
	s.vehicles[vehicle.ID] = vehicle

	return nil
}

// UnregisterVehicle unregisters a V2X vehicle
func (s *CV2XSimulation) UnregisterVehicle(vehicleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.vehicles[vehicleID]; !exists {
		return fmt.Errorf("vehicle not found: %s", vehicleID)
	}

	delete(s.vehicles, vehicleID)

	// Remove from RSU connections
	for _, rsu := range s.rsus {
		for i, id := range rsu.ConnectedVehicles {
			if id == vehicleID {
				rsu.ConnectedVehicles = append(rsu.ConnectedVehicles[:i], rsu.ConnectedVehicles[i+1:]...)
				break
			}
		}
	}

	return nil
}

// UpdateVehiclePosition updates vehicle position
func (s *CV2XSimulation) UpdateVehiclePosition(vehicleID string, position *Position, speed, heading float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicle, exists := s.vehicles[vehicleID]
	if !exists {
		return fmt.Errorf("vehicle not found: %s", vehicleID)
	}

	vehicle.Position = position
	vehicle.Speed = speed
	vehicle.Heading = heading
	vehicle.LastSeen = time.Now()

	// Update RSU connections
	s.updateRSUConnections(vehicle)

	return nil
}

// updateRSUConnections updates which RSUs the vehicle is connected to
func (s *CV2XSimulation) updateRSUConnections(vehicle *V2XVehicle) {
	for _, rsu := range s.rsus {
		distance := calculateDistanceMeters(
			vehicle.Position.Latitude, vehicle.Position.Longitude,
			rsu.Position.Latitude, rsu.Position.Longitude,
		)

		connected := false
		for _, id := range rsu.ConnectedVehicles {
			if id == vehicle.ID {
				connected = true
				break
			}
		}

		if distance <= rsu.Range && !connected {
			rsu.ConnectedVehicles = append(rsu.ConnectedVehicles, vehicle.ID)
		} else if distance > rsu.Range && connected {
			for i, id := range rsu.ConnectedVehicles {
				if id == vehicle.ID {
					rsu.ConnectedVehicles = append(rsu.ConnectedVehicles[:i], rsu.ConnectedVehicles[i+1:]...)
					break
				}
			}
		}
	}
}

// SendBSM sends a Basic Safety Message
func (s *CV2XSimulation) SendBSM(vehicleID string, bsm *BSM) (*V2XMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicle, exists := s.vehicles[vehicleID]
	if !exists {
		return nil, fmt.Errorf("vehicle not found: %s", vehicleID)
	}

	payload, err := json.Marshal(bsm)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal BSM: %w", err)
	}

	message := &V2XMessage{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Type:      "bsm",
		SourceID:  vehicleID,
		TargetID:  "", // Broadcast
		Payload:   payload,
		Timestamp: time.Now(),
		Priority:  0, // Highest priority for safety messages
		Received:  false,
		Processed: false,
	}

	s.addMessage(message)
	vehicle.MessagesSent++

	return message, nil
}

// SendCAM sends a Cooperative Awareness Message
func (s *CV2XSimulation) SendCAM(vehicleID string, cam *CAM) (*V2XMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicle, exists := s.vehicles[vehicleID]
	if !exists {
		return nil, fmt.Errorf("vehicle not found: %s", vehicleID)
	}

	payload, err := json.Marshal(cam)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CAM: %w", err)
	}

	message := &V2XMessage{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Type:      "cam",
		SourceID:  vehicleID,
		TargetID:  "", // Broadcast
		Payload:   payload,
		Timestamp: time.Now(),
		Priority:  2,
		Received:  false,
		Processed: false,
	}

	s.addMessage(message)
	vehicle.MessagesSent++

	return message, nil
}

// SendDENM sends a Decentralized Environmental Notification Message
func (s *CV2XSimulation) SendDENM(rsuID string, denm *DENM) (*V2XMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.rsus[rsuID]
	if !exists {
		return nil, fmt.Errorf("RSU not found: %s", rsuID)
	}

	payload, err := json.Marshal(denm)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DENM: %w", err)
	}

	message := &V2XMessage{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Type:      "denm",
		SourceID:  rsuID,
		TargetID:  "", // Broadcast
		Payload:   payload,
		Timestamp: time.Now(),
		Priority:  1, // High priority for hazard warnings
		Received:  false,
		Processed: false,
	}

	s.addMessage(message)

	return message, nil
}

// addMessage adds a message to the buffer
func (s *CV2XSimulation) addMessage(message *V2XMessage) {
	s.messages = append(s.messages, message)

	// Trim old messages
	if len(s.messages) > s.messageBufferSize {
		s.messages = s.messages[len(s.messages)-s.messageBufferSize:]
	}
}

// ReceiveMessages receives messages for a vehicle
func (s *CV2XSimulation) ReceiveMessages(vehicleID string) []*V2XMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicle, exists := s.vehicles[vehicleID]
	if !exists {
		return nil
	}

	received := make([]*V2XMessage, 0)

	for _, message := range s.messages {
		// Skip own messages
		if message.SourceID == vehicleID {
			continue
		}

		// Check if message is in range
		if s.isMessageInRange(vehicle, message) {
			message.Received = true
			received = append(received, message)
			vehicle.MessagesReceived++
		}
	}

	return received
}

// isMessageInRange checks if a message is in range of the vehicle
func (s *CV2XSimulation) isMessageInRange(vehicle *V2XVehicle, message *V2XMessage) bool {
	// Get source position
	var sourcePosition *Position
	if sourceVehicle, exists := s.vehicles[message.SourceID]; exists {
		sourcePosition = sourceVehicle.Position
	} else if sourceRSU, exists := s.rsus[message.SourceID]; exists {
		sourcePosition = sourceRSU.Position
	} else {
		return false
	}

	// Calculate distance
	distance := calculateDistanceMeters(
		vehicle.Position.Latitude, vehicle.Position.Longitude,
		sourcePosition.Latitude, sourcePosition.Longitude,
	)

	// V2X communication range: 500-1000 meters
	return distance <= 1000
}

// GetVehicles returns all registered vehicles
func (s *CV2XSimulation) GetVehicles() map[string]*V2XVehicle {
	s.mu.RLock()
	defer s.mu.RUnlock()

	vehicles := make(map[string]*V2XVehicle)
	for id, vehicle := range s.vehicles {
		vehicles[id] = vehicle
	}
	return vehicles
}

// GetRSUs returns all RSUs
func (s *CV2XSimulation) GetRSUs() map[string]*RSU {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rsus := make(map[string]*RSU)
	for id, rsu := range s.rsus {
		rsus[id] = rsu
	}
	return rsus
}

// GetMessages returns message history
func (s *CV2XSimulation) GetMessages(limit int) []*V2XMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages := s.messages
	if limit > 0 && len(messages) > limit {
		messages = messages[len(messages)-limit:]
	}
	return messages
}

// GetNearbyVehicles returns vehicles near a position
func (s *CV2XSimulation) GetNearbyVehicles(position *Position, radius float64) []*V2XVehicle {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nearby := make([]*V2XVehicle, 0)

	for _, vehicle := range s.vehicles {
		distance := calculateDistanceMeters(
			position.Latitude, position.Longitude,
			vehicle.Position.Latitude, vehicle.Position.Longitude,
		)

		if distance <= radius {
			nearby = append(nearby, vehicle)
		}
	}

	return nearby
}

// GetNearbyRSUs returns RSUs near a position
func (s *CV2XSimulation) GetNearbyRSUs(position *Position, radius float64) []*RSU {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nearby := make([]*RSU, 0)

	for _, rsu := range s.rsus {
		distance := calculateDistanceMeters(
			position.Latitude, position.Longitude,
			rsu.Position.Latitude, rsu.Position.Longitude,
		)

		if distance <= radius {
			nearby = append(nearby, rsu)
		}
	}

	return nearby
}

// SimulateTrafficHazard simulates a traffic hazard
func (s *CV2XSimulation) SimulateTrafficHazard(position *Position, eventType, severity string, duration time.Duration) (*DENM, error) {
	denm := &DENM{
		EventID:   fmt.Sprintf("hazard_%d", time.Now().UnixNano()),
		EventType: eventType,
		Position:  position,
		Radius:    500, // 500 meters warning radius
		Severity:  severity,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(duration),
	}

	// Broadcast DENM from nearest RSU
	rsu := s.findNearestRSU(position)
	if rsu == nil {
		return nil, fmt.Errorf("no RSU found near position")
	}

	_, err := s.SendDENM(rsu.ID, denm)
	if err != nil {
		return nil, err
	}

	return denm, nil
}

// findNearestRSU finds the nearest RSU to a position
func (s *CV2XSimulation) findNearestRSU(position *Position) *RSU {
	var nearest *RSU
	minDistance := math.MaxFloat64

	for _, rsu := range s.rsus {
		distance := calculateDistanceMeters(
			position.Latitude, position.Longitude,
			rsu.Position.Latitude, rsu.Position.Longitude,
		)

		if distance < minDistance {
			minDistance = distance
			nearest = rsu
		}
	}

	return nearest
}

// GetStatistics returns simulation statistics
func (s *CV2XSimulation) GetStatistics() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalMessagesSent := 0
	totalMessagesReceived := 0

	for _, vehicle := range s.vehicles {
		totalMessagesSent += vehicle.MessagesSent
		totalMessagesReceived += vehicle.MessagesReceived
	}

	stats := map[string]interface{}{
		"total_vehicles":          len(s.vehicles),
		"total_rsus":              len(s.rsus),
		"total_messages":          len(s.messages),
		"total_messages_sent":     totalMessagesSent,
		"total_messages_received": totalMessagesReceived,
		"simulation_enabled":      s.simulationEnabled,
	}

	return stats
}

// GenerateRandomVehicles generates random V2X vehicles for simulation
func (s *CV2XSimulation) GenerateRandomVehicles(count int, centerLat, centerLng float64, radiusKm float64) ([]*V2XVehicle, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	vehicleTypes := []string{"car", "car", "car", "bus", "truck", "emergency"}
	vehicles := make([]*V2XVehicle, 0, count)

	for i := 0; i < count; i++ {
		// Generate random position within radius
		angle := rand.Float64() * 2 * math.Pi
		distance := rand.Float64() * radiusKm
		lat := centerLat + (distance*math.Cos(angle))/111.0
		lng := centerLng + (distance*math.Sin(angle))/(111.0*math.Cos(centerLat*math.Pi/180))

		// Random vehicle type
		typeIdx := rand.Intn(len(vehicleTypes))
		vehicleType := vehicleTypes[typeIdx]

		// Random speed and heading
		speed := 30 + rand.Float64()*70 // 30-100 km/h
		heading := rand.Float64() * 360 // 0-360 degrees

		vehicle := &V2XVehicle{
			ID:         fmt.Sprintf("v2x_v_%d_%d", time.Now().Unix(), i),
			Type:       vehicleType,
			Position:   &Position{Latitude: lat, Longitude: lng, Altitude: 0},
			Speed:      speed,
			Heading:    heading,
			V2XEnabled: true,
			LastSeen:   time.Now(),
		}

		s.vehicles[vehicle.ID] = vehicle
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

// GenerateTrafficScenario generates a traffic scenario with specific patterns
func (s *CV2XSimulation) GenerateTrafficScenario(scenarioType string) (*TrafficScenario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	scenario := &TrafficScenario{
		ID:          fmt.Sprintf("scenario_%d", time.Now().Unix()),
		Type:        scenarioType,
		CreatedAt:   time.Now(),
		Vehicles:    make([]*V2XVehicle, 0),
		Description: "",
	}

	centerLat := 22.6913 // Longhua center
	centerLng := 114.0448

	switch scenarioType {
	case "congestion":
		scenario.Description = "拥堵场景：高密度车辆，低速行驶"
		count := 50
		for i := 0; i < count; i++ {
			lat := centerLat + (rand.Float64()-0.5)*0.02
			lng := centerLng + (rand.Float64()-0.5)*0.02
			vehicle := &V2XVehicle{
				ID:         fmt.Sprintf("v2x_v_%d", i),
				Type:       "car",
				Position:   &Position{Latitude: lat, Longitude: lng, Altitude: 0},
				Speed:      10 + rand.Float64()*20, // 10-30 km/h (slow)
				Heading:    rand.Float64() * 360,
				V2XEnabled: true,
				LastSeen:   time.Now(),
			}
			s.vehicles[vehicle.ID] = vehicle
			scenario.Vehicles = append(scenario.Vehicles, vehicle)
		}

	case "normal":
		scenario.Description = "正常场景：中等密度，正常速度"
		count := 20
		for i := 0; i < count; i++ {
			lat := centerLat + (rand.Float64()-0.5)*0.03
			lng := centerLng + (rand.Float64()-0.5)*0.03
			vehicle := &V2XVehicle{
				ID:         fmt.Sprintf("v2x_v_%d", i),
				Type:       "car",
				Position:   &Position{Latitude: lat, Longitude: lng, Altitude: 0},
				Speed:      40 + rand.Float64()*40, // 40-80 km/h
				Heading:    rand.Float64() * 360,
				V2XEnabled: true,
				LastSeen:   time.Now(),
			}
			s.vehicles[vehicle.ID] = vehicle
			scenario.Vehicles = append(scenario.Vehicles, vehicle)
		}

	case "emergency":
		scenario.Description = "紧急场景：救护车优先，其他车辆避让"
		// Add emergency vehicle
		emergencyVehicle := &V2XVehicle{
			ID:         "v2x_emergency_001",
			Type:       "emergency",
			Position:   &Position{Latitude: centerLat, Longitude: centerLng, Altitude: 0},
			Speed:      80,
			Heading:    0,
			V2XEnabled: true,
			LastSeen:   time.Now(),
		}
		s.vehicles[emergencyVehicle.ID] = emergencyVehicle
		scenario.Vehicles = append(scenario.Vehicles, emergencyVehicle)

		// Add regular vehicles that should avoid
		for i := 0; i < 10; i++ {
			lat := centerLat + (rand.Float64()-0.5)*0.01
			lng := centerLng + (rand.Float64()-0.5)*0.01
			vehicle := &V2XVehicle{
				ID:         fmt.Sprintf("v2x_v_%d", i),
				Type:       "car",
				Position:   &Position{Latitude: lat, Longitude: lng, Altitude: 0},
				Speed:      20 + rand.Float64()*20,
				Heading:    rand.Float64() * 360,
				V2XEnabled: true,
				LastSeen:   time.Now(),
			}
			s.vehicles[vehicle.ID] = vehicle
			scenario.Vehicles = append(scenario.Vehicles, vehicle)
		}

	default:
		return nil, fmt.Errorf("unknown scenario type: %s", scenarioType)
	}

	return scenario, nil
}

// TrafficScenario represents a traffic simulation scenario
type TrafficScenario struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Vehicles    []*V2XVehicle `json:"vehicles"`
	CreatedAt   time.Time     `json:"created_at"`
}

// ClearAllVehicles removes all vehicles from simulation
func (s *CV2XSimulation) ClearAllVehicles() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vehicles = make(map[string]*V2XVehicle)
	return nil
}

// calculateDistanceMeters calculates distance between two coordinates in meters
func calculateDistanceMeters(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000 // meters

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
