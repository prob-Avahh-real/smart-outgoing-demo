package parking

import (
	"fmt"
	"sync"
	"time"
)

// MockParkingService simulates a parking lot for local testing
type MockParkingService struct {
	mu              sync.RWMutex
	parkingLots     map[string]*MockParkingLot
	vehicles        map[string]*MockVehicle
	entries         []MockEntryRecord
	exits           []MockExitRecord
	heartbeatTimer  *time.Ticker
}

// MockParkingLot represents a mock parking lot
type MockParkingLot struct {
	ParkingLotNo    string
	Name            string
	Address         string
	Latitude        float64
	Longitude       float64
	TotalSpaces     int
	AvailableSpaces int
	Ports           []MockPort
	LastHeartbeat   time.Time
}

// MockPort represents a mock entrance/exit port
type MockPort struct {
	PortNo   string
	PortType string
	Latitude float64
	Longitude float64
}

// MockVehicle represents a vehicle in the mock system
type MockVehicle struct {
	PlateNo      string
	EntryTime    time.Time
	ExitTime     *time.Time
	ParkingLotNo string
	PortNo       string
	OutOrderNo   string
	Status       string // "parked", "exited"
}

// MockEntryRecord represents a vehicle entry record
type MockEntryRecord struct {
	OutOrderNo   string
	PlateNo      string
	EntryTime    time.Time
	ParkingLotNo string
	PortNo       string
}

// MockExitRecord represents a vehicle exit record
type MockExitRecord struct {
	OutOrderNo   string
	PlateNo      string
	ExitTime     time.Time
	ParkingFee   float64
}

// NewMockParkingService creates a new mock parking service
func NewMockParkingService() *MockParkingService {
	service := &MockParkingService{
		parkingLots: make(map[string]*MockParkingLot),
		vehicles:    make(map[string]*MockVehicle),
		entries:     make([]MockEntryRecord, 0),
		exits:       make([]MockExitRecord, 0),
	}

	// Initialize with default parking lots
	service.initializeDefaultLots()

	return service
}

// initializeDefaultLots initializes default parking lots for testing
func (s *MockParkingService) initializeDefaultLots() {
	defaultLots := []*MockParkingLot{
		{
			ParkingLotNo:    "LOT001",
			Name:            "龙华智行停车场A",
			Address:         "深圳市龙华区龙华街道",
			Latitude:        22.6913,
			Longitude:       114.0448,
			TotalSpaces:     100,
			AvailableSpaces: 50,
			Ports: []MockPort{
				{PortNo: "PORT001", PortType: "entry", Latitude: 22.6913, Longitude: 114.0448},
				{PortNo: "PORT002", PortType: "exit", Latitude: 22.6914, Longitude: 114.0449},
			},
			LastHeartbeat: time.Now(),
		},
		{
			ParkingLotNo:    "LOT002",
			Name:            "龙华智行停车场B",
			Address:         "深圳市龙华区民治街道",
			Latitude:        22.6923,
			Longitude:       114.0458,
			TotalSpaces:     80,
			AvailableSpaces: 30,
			Ports: []MockPort{
				{PortNo: "PORT003", PortType: "entry", Latitude: 22.6923, Longitude: 114.0458},
				{PortNo: "PORT004", PortType: "exit", Latitude: 22.6924, Longitude: 114.0459},
			},
			LastHeartbeat: time.Now(),
		},
	}

	for _, lot := range defaultLots {
		s.parkingLots[lot.ParkingLotNo] = lot
	}
}

// RegisterParkingLot registers a new parking lot
func (s *MockParkingService) RegisterParkingLot(registration *ParkingLotRegistration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.parkingLots[registration.ParkingLotNo]; exists {
		return fmt.Errorf("parking lot %s already exists", registration.ParkingLotNo)
	}

	ports := make([]MockPort, len(registration.Ports))
	for i, p := range registration.Ports {
		ports[i] = MockPort{
			PortNo:    p.PortNo,
			PortType:  p.PortType,
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
		}
	}

	lot := &MockParkingLot{
		ParkingLotNo:    registration.ParkingLotNo,
		Name:            registration.Name,
		Address:         registration.Address,
		Latitude:        registration.Latitude,
		Longitude:       registration.Longitude,
		TotalSpaces:     registration.TotalSpaces,
		AvailableSpaces: registration.TotalSpaces,
		Ports:           ports,
		LastHeartbeat:   time.Now(),
	}

	s.parkingLots[registration.ParkingLotNo] = lot
	return nil
}

// HandleVehicleEntry handles vehicle entry
func (s *MockParkingService) HandleVehicleEntry(req *ParkingEntryRequest) (*ParkingEntryResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if parking lot exists
	lot, exists := s.parkingLots[req.ParkingLotNo]
	if !exists {
		return &ParkingEntryResponse{
			Code:    1,
			Message: "Parking lot not found",
		}, nil
	}

	// Check if vehicle is already parked
	if vehicle, exists := s.vehicles[req.PlateNo]; exists && vehicle.Status == "parked" {
		return &ParkingEntryResponse{
			Code:    1,
			Message: "Vehicle already parked",
		}, nil
	}

	// Decrease available spaces
	if lot.AvailableSpaces <= 0 {
		return &ParkingEntryResponse{
			Code:    1,
			Message: "No available spaces",
		}, nil
	}
	lot.AvailableSpaces--

	// Record entry
	vehicle := &MockVehicle{
		PlateNo:      req.PlateNo,
		EntryTime:    req.EntryTime,
		ParkingLotNo: req.ParkingLotNo,
		PortNo:       req.PortNo,
		OutOrderNo:   req.OutOrderNo,
		Status:       "parked",
	}
	s.vehicles[req.PlateNo] = vehicle

	s.entries = append(s.entries, MockEntryRecord{
		OutOrderNo:   req.OutOrderNo,
		PlateNo:      req.PlateNo,
		EntryTime:    req.EntryTime,
		ParkingLotNo: req.ParkingLotNo,
		PortNo:       req.PortNo,
	})

	return &ParkingEntryResponse{
		Code:    0,
		Message: "success",
		Data: struct {
			OutOrderNo string `json:"out_order_no"`
			EntryTime  string `json:"entry_time"`
		}{
			OutOrderNo: req.OutOrderNo,
			EntryTime:  req.EntryTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

// HandleVehicleExit handles vehicle exit
func (s *MockParkingService) HandleVehicleExit(req *ParkingExitRequest) (*ParkingExitResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if vehicle exists
	vehicle, exists := s.vehicles[req.PlateNo]
	if !exists || vehicle.Status != "parked" {
		return &ParkingExitResponse{
			Code:    1,
			Message: "Vehicle not found or not parked",
		}, nil
	}

	// Check parking lot
	lot, exists := s.parkingLots[vehicle.ParkingLotNo]
	if !exists {
		return &ParkingExitResponse{
			Code:    1,
			Message: "Parking lot not found",
		}, nil
	}

	// Calculate parking fee (simple calculation: 5 CNY per hour)
	duration := req.ExitTime.Sub(vehicle.EntryTime).Hours()
	parkingFee := duration * 5.0

	// Update vehicle status
	exitTime := req.ExitTime
	vehicle.ExitTime = &exitTime
	vehicle.Status = "exited"

	// Increase available spaces
	lot.AvailableSpaces++

	// Record exit
	s.exits = append(s.exits, MockExitRecord{
		OutOrderNo: req.OutOrderNo,
		PlateNo:    req.PlateNo,
		ExitTime:   req.ExitTime,
		ParkingFee: parkingFee,
	})

	return &ParkingExitResponse{
		Code:    0,
		Message: "success",
		Data: struct {
			OutOrderNo string `json:"out_order_no"`
			ExitTime   string `json:"exit_time"`
			ParkingFee float64 `json:"parking_fee"`
		}{
			OutOrderNo: req.OutOrderNo,
			ExitTime:   req.ExitTime.Format("2006-01-02 15:04:05"),
			ParkingFee: parkingFee,
		},
	}, nil
}

// UpdateSpaceStatus updates parking space status
func (s *MockParkingService) UpdateSpaceStatus(parkingLotNo string, availableSpaces int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	lot, exists := s.parkingLots[parkingLotNo]
	if !exists {
		return fmt.Errorf("parking lot %s not found", parkingLotNo)
	}

	lot.AvailableSpaces = availableSpaces
	lot.LastHeartbeat = time.Now()

	return nil
}

// GetParkingLotStatus returns parking lot status
func (s *MockParkingService) GetParkingLotStatus(parkingLotNo string) (*MockParkingLot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lot, exists := s.parkingLots[parkingLotNo]
	if !exists {
		return nil, fmt.Errorf("parking lot %s not found", parkingLotNo)
	}

	return lot, nil
}

// GetAllParkingLots returns all parking lots
func (s *MockParkingService) GetAllParkingLots() []*MockParkingLot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lots := make([]*MockParkingLot, 0, len(s.parkingLots))
	for _, lot := range s.parkingLots {
		lots = append(lots, lot)
	}

	return lots
}

// GetEntryRecords returns entry records
func (s *MockParkingService) GetEntryRecords() []MockEntryRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := make([]MockEntryRecord, len(s.entries))
	copy(records, s.entries)
	return records
}

// GetExitRecords returns exit records
func (s *MockParkingService) GetExitRecords() []MockExitRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := make([]MockExitRecord, len(s.exits))
	copy(records, s.exits)
	return records
}

// StartHeartbeat starts automatic heartbeat
func (s *MockParkingService) StartHeartbeat(interval time.Duration) {
	if s.heartbeatTimer != nil {
		s.heartbeatTimer.Stop()
	}

	s.heartbeatTimer = time.NewTicker(interval)

	go func() {
		for range s.heartbeatTimer.C {
			s.mu.Lock()
			now := time.Now()
			for _, lot := range s.parkingLots {
				lot.LastHeartbeat = now
			}
			s.mu.Unlock()
		}
	}()
}

// StopHeartbeat stops automatic heartbeat
func (s *MockParkingService) StopHeartbeat() {
	if s.heartbeatTimer != nil {
		s.heartbeatTimer.Stop()
		s.heartbeatTimer = nil
	}
}
