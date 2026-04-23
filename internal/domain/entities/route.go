package entities

import (
	"crypto/md5"
	"fmt"
	"time"
)

// Route represents the route aggregate
type Route struct {
	ID          string      `json:"id"`
	VehicleID   string      `json:"vehicle_id"`
	Path        [][]float64 `json:"path"`
	Hash        string      `json:"hash"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// NewRoute creates a new route aggregate
func NewRoute(vehicleID string, path [][]float64) *Route {
	now := time.Now()
	hash := generateHash(path)
	
	return &Route{
		ID:        fmt.Sprintf("route_%s_%s", vehicleID, hash[:8]),
		VehicleID: vehicleID,
		Path:      path,
		Hash:      hash,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdatePath updates the route path and generates new hash
func (r *Route) UpdatePath(path [][]float64) {
	r.Path = path
	r.Hash = generateHash(path)
	r.UpdatedAt = time.Now()
}

// generateHash creates a hash for the route path
func generateHash(path [][]float64) string {
	if len(path) == 0 {
		return ""
	}
	
	data := ""
	for _, point := range path {
		data += fmt.Sprintf("%.6f,%.6f;", point[0], point[1])
	}
	
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}
