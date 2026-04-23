package entities

import (
	"time"
)

// Vehicle represents the core vehicle aggregate root
type Vehicle struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Start     Coordinates `json:"start"`
	End       Coordinates `json:"end,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Coordinates is a value object representing geographic coordinates
type Coordinates struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
	Alt float64 `json:"alt,omitempty"`
}

// NewVehicle creates a new vehicle aggregate
func NewVehicle(id, name string, start Coordinates) *Vehicle {
	now := time.Now()
	return &Vehicle{
		ID:        id,
		Name:      name,
		Start:     start,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetDestination updates the vehicle's destination
func (v *Vehicle) SetDestination(end Coordinates) {
	v.End = end
	v.UpdatedAt = time.Now()
}

// UpdateName updates the vehicle's name
func (v *Vehicle) UpdateName(name string) {
	v.Name = name
	v.UpdatedAt = time.Now()
}
