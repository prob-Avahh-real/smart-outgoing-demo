package parking

import "time"

// CityBrainAPIConfig holds configuration for City Brain API
type CityBrainAPIConfig struct {
	BaseURL     string
	AppID       string
	AppSecret   string
	ParkingLotNo string
	PortNo      string
	Timeout     time.Duration
	UseMock     bool // Use mock service for testing
}

// ParkingEntryRequest represents a vehicle entry request
type ParkingEntryRequest struct {
	PlateNo      string    `json:"plate_no"`
	EntryTime    time.Time `json:"entry_time"`
	PortNo       string    `json:"port_no"`
	ParkingLotNo string    `json:"parking_lot_no"`
	OutOrderNo   string    `json:"out_order_no"` // Unique order number
	VehicleType  string    `json:"vehicle_type"`
}

// ParkingExitRequest represents a vehicle exit request
type ParkingEntryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		OutOrderNo string `json:"out_order_no"`
		EntryTime  string `json:"entry_time"`
	} `json:"data"`
}

// ParkingExitRequest represents a vehicle exit request
type ParkingExitRequest struct {
	PlateNo      string    `json:"plate_no"`
	ExitTime     time.Time `json:"exit_time"`
	PortNo       string    `json:"port_no"`
	ParkingLotNo string    `json:"parking_lot_no"`
	OutOrderNo   string    `json:"out_order_no"`
	ParkingFee   float64   `json:"parking_fee"`
}

// ParkingExitResponse represents vehicle exit response
type ParkingExitResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		OutOrderNo string `json:"out_order_no"`
		ExitTime   string `json:"exit_time"`
		ParkingFee float64 `json:"parking_fee"`
	} `json:"data"`
}

// SnapshotUploadRequest represents snapshot image upload request
type SnapshotUploadRequest struct {
	OutOrderNo  string `json:"out_order_no"`
	ImageType   string `json:"image_type"` // entry/exit
	ImageData   string `json:"image_data"` // Base64 encoded
	UploadTime  time.Time `json:"upload_time"`
}

// SnapshotUploadResponse represents snapshot upload response
type SnapshotUploadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ImageURL string `json:"image_url"`
	} `json:"data"`
}

// HeartbeatRequest represents parking lot heartbeat request
type HeartbeatRequest struct {
	ParkingLotNo  string    `json:"parking_lot_no"`
	TotalSpaces   int       `json:"total_spaces"`
	AvailableSpaces int     `json:"available_spaces"`
	HeartbeatTime time.Time `json:"heartbeat_time"`
}

// HeartbeatResponse represents heartbeat response
type HeartbeatResponse struct {
	Code    int    `json:"code"`
	Message string `json:"data"`
}

// ParkingLotRegistration represents parking lot registration
type ParkingLotRegistration struct {
	ParkingLotNo string  `json:"parking_lot_no"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	TotalSpaces  int     `json:"total_spaces"`
	Ports        []Port  `json:"ports"`
}

// Port represents a parking lot entrance/exit port
type Port struct {
	PortNo   string  `json:"port_no"`
	PortType string  `json:"port_type"` // entry/exit
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
