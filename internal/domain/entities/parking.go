package entities

import (
	"time"
)

// ParkingLot represents a parking facility
type ParkingLot struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	TotalSpaces     int       `json:"total_spaces"`
	AvailableSpaces int       `json:"available_spaces"`
	PricePerHour    float64   `json:"price_per_hour"`
	Features        []string  `json:"features"` // EV charging, covered, 24/7, etc
	Rating          float64   `json:"rating"`
	Distance        float64   `json:"distance"` // Distance from user
	IsOpen          bool      `json:"is_open"`
	LastUpdated     time.Time `json:"last_updated"`
}

// ParkingSpace represents an individual parking space
type ParkingSpace struct {
	ID           string    `json:"id"`
	ParkingLotID string    `json:"parking_lot_id"`
	SpaceNumber  string    `json:"space_number"`
	Type         SpaceType `json:"type"` // regular, compact, EV, disabled
	IsAvailable  bool      `json:"is_available"`
	IsReserved   bool      `json:"is_reserved"`
	Price        float64   `json:"price"`
	Floor        string    `json:"floor"`
	Zone         string    `json:"zone"`
	LastUpdated  time.Time `json:"last_updated"`
}

// SpaceType represents different types of parking spaces
type SpaceType string

const (
	SpaceTypeRegular  SpaceType = "regular"
	SpaceTypeCompact  SpaceType = "compact"
	SpaceTypeEV       SpaceType = "ev"
	SpaceTypeDisabled SpaceType = "disabled"
	SpaceTypeVIP      SpaceType = "vip"
)

// ParkingReservation represents a parking reservation
type ParkingReservation struct {
	ID           string            `json:"id"`
	UserID       string            `json:"user_id"`
	ParkingLotID string            `json:"parking_lot_id"`
	SpaceID      string            `json:"space_id"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time"`
	Status       ReservationStatus `json:"status"`
	TotalPrice   float64           `json:"total_price"`
	PaymentID    string            `json:"payment_id,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// ReservationStatus represents reservation status
type ReservationStatus string

const (
	StatusPending   ReservationStatus = "pending"
	StatusConfirmed ReservationStatus = "confirmed"
	StatusActive    ReservationStatus = "active"
	StatusCompleted ReservationStatus = "completed"
	StatusCancelled ReservationStatus = "cancelled"
)

// ParkingSession represents an active parking session
type ParkingSession struct {
	ID           string        `json:"id"`
	UserID       string        `json:"user_id"`
	ParkingLotID string        `json:"parking_lot_id"`
	SpaceID      string        `json:"space_id"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      *time.Time    `json:"end_time,omitempty"`
	Status       SessionStatus `json:"status"`
	TotalCost    float64       `json:"total_cost"`
	PaymentID    string        `json:"payment_id,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// SessionStatus represents parking session status
type SessionStatus string

const (
	SessionActive    SessionStatus = "active"
	SessionCompleted SessionStatus = "completed"
	SessionCancelled SessionStatus = "cancelled"
)

// UserParkingPreference represents user's parking preferences
type UserParkingPreference struct {
	UserID               string    `json:"user_id"`
	MaxPricePerHour      float64   `json:"max_price_per_hour"`
	PreferredDistance    float64   `json:"preferred_distance"`
	PreferredFeatures    []string  `json:"preferred_features"`
	VehicleType          SpaceType `json:"vehicle_type"`
	DisableAccessibility bool      `json:"disable_accessibility"`
	PreferCovered        bool      `json:"prefer_covered"`
	PreferEV             bool      `json:"prefer_ev"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ParkingRecommendation represents a parking recommendation
type ParkingRecommendation struct {
	ParkingLot       *ParkingLot   `json:"parking_lot"`
	RecommendedSpace *ParkingSpace `json:"recommended_space,omitempty"`
	Score            float64       `json:"score"`
	Reasons          []string      `json:"reasons"`
	EstimatedTime    time.Duration `json:"estimated_time"`
	TotalPrice       float64       `json:"total_price"`
	Route            *ParkingRoute `json:"route,omitempty"`
}

// ParkingRoute represents navigation route to parking space
type ParkingRoute struct {
	Steps         []RouteStep   `json:"steps"`
	TotalDistance float64       `json:"total_distance"`
	TotalTime     time.Duration `json:"total_time"`
	Instructions  string        `json:"instructions"`
}

// RouteStep represents a step in the navigation route
type RouteStep struct {
	Instruction string        `json:"instruction"`
	Distance    float64       `json:"distance"`
	Duration    time.Duration `json:"duration"`
	Direction   string        `json:"direction"`
}

// ParkingEvent represents parking-related events
type ParkingEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	UserID    string                 `json:"user_id"`
	LotID     string                 `json:"lot_id"`
	SpaceID   string                 `json:"space_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventType represents different parking event types
type EventType string

const (
	EventSpaceAvailable       EventType = "space_available"
	EventSpaceOccupied        EventType = "space_occupied"
	EventReservationCreated   EventType = "reservation_created"
	EventReservationCancelled EventType = "reservation_canceled"
	EventSessionStarted       EventType = "session_started"
	EventSessionEnded         EventType = "session_ended"
	EventPaymentCompleted     EventType = "payment_completed"
)

// Helper methods

// IsAvailable checks if a parking lot has available spaces
func (pl *ParkingLot) IsAvailable() bool {
	return pl.IsOpen && pl.AvailableSpaces > 0
}

// OccupancyRate returns the occupancy rate as a percentage
func (pl *ParkingLot) OccupancyRate() float64 {
	if pl.TotalSpaces == 0 {
		return 0
	}
	return float64(pl.TotalSpaces-pl.AvailableSpaces) / float64(pl.TotalSpaces) * 100
}

// IsActive checks if a reservation is currently active
func (r *ParkingReservation) IsActive() bool {
	now := time.Now()
	return r.Status == StatusConfirmed && now.After(r.StartTime) && now.Before(r.EndTime)
}

// Duration calculates the duration of a parking session
func (s *ParkingSession) Duration() time.Duration {
	if s.EndTime != nil {
		return s.EndTime.Sub(s.StartTime)
	}
	return time.Since(s.StartTime)
}

// CalculatePrice calculates total price for a reservation
func (r *ParkingReservation) CalculatePrice() float64 {
	duration := r.EndTime.Sub(r.StartTime)
	hours := duration.Hours()
	return hours * r.TotalPrice
}
