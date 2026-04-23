package entities

import "time"

// ParkingPoolLevel represents the three-tier parking pool classification
type ParkingPoolLevel string

const (
	PoolLevelCore      ParkingPoolLevel = "core"      // 核心商圈优先
	PoolLevelPeripheral ParkingPoolLevel = "peripheral" // 外围配套备用
	PoolLevelRoadside   ParkingPoolLevel = "roadside"  // 路边临停兜底
)

// ParkingPool represents a tiered parking pool for traffic diversion
type ParkingPool struct {
	ID          string           `json:"id"`
	Level       ParkingPoolLevel `json:"level"`       // core/peripheral/roadside
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Priority    int              `json:"priority"`    // 1=highest, 3=lowest
	Lots        []*ParkingLot    `json:"lots"`        // Parking lots in this pool
	TotalSpaces int              `json:"total_spaces"`
	FreeSpaces  int              `json:"free_spaces"`
	TurnoverRate float64         `json:"turnover_rate"` // 周转率 (vehicles/hour)
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// ParkingPoolStats represents statistics for a parking pool
type ParkingPoolStats struct {
	PoolID         string    `json:"pool_id"`
	OccupancyRate  float64   `json:"occupancy_rate"`  // 占用率
	TurnoverRate   float64   `json:"turnover_rate"`   // 周转率
	AvgWaitTime    float64   `json:"avg_wait_time"`   // 平均等待时间 (minutes)
	TotalVehicles  int       `json:"total_vehicles"`  // 当前车辆数
	LastUpdated    time.Time `json:"last_updated"`
}

// TrafficDiversionRule represents rules for traffic diversion
type TrafficDiversionRule struct {
	ID              string           `json:"id"`
	SourceZone      string           `json:"source_zone"`      // 拥堵区域
	TargetPoolLevel ParkingPoolLevel `json:"target_pool_level"` // 目标车场池级别
	TriggerDensity  float64          `json:"trigger_density"`  // 触发密度 (vehicles/km)
	MaxDistance     float64          `json:"max_distance"`     // 最大分流距离 (km)
	Enable          bool             `json:"enable"`
	Priority        int              `json:"priority"`
}

// ParkingPoolAssignment represents a parking lot assignment for a vehicle
type ParkingPoolAssignment struct {
	ID           string           `json:"id"`
	VehicleID    string           `json:"vehicle_id"`
	PoolID       string           `json:"pool_id"`
	LotID        string           `json:"lot_id"`
	Level        ParkingPoolLevel `json:"level"`
	Reason       string           `json:"reason"`        // 分配原因
	AssignedAt   time.Time        `json:"assigned_at"`
	CompletedAt  *time.Time       `json:"completed_at"`
	Status       string           `json:"status"`        // pending/active/completed/cancelled
}

// Helper methods

// GetAvailableLots returns available parking lots in the pool
func (p *ParkingPool) GetAvailableLots() []*ParkingLot {
	available := make([]*ParkingLot, 0)
	for _, lot := range p.Lots {
		if lot.IsAvailable() {
			available = append(available, lot)
		}
	}
	return available
}

// GetOccupancyRate returns the occupancy rate as a percentage
func (p *ParkingPool) GetOccupancyRate() float64 {
	if p.TotalSpaces == 0 {
		return 0
	}
	return float64(p.TotalSpaces-p.FreeSpaces) / float64(p.TotalSpaces) * 100
}

// UpdateStats updates pool statistics based on parking lots
func (p *ParkingPool) UpdateStats() {
	totalSpaces := 0
	freeSpaces := 0
	
	for _, lot := range p.Lots {
		totalSpaces += lot.TotalSpaces
		freeSpaces += lot.AvailableSpaces
	}
	
	p.TotalSpaces = totalSpaces
	p.FreeSpaces = freeSpaces
	p.UpdatedAt = time.Now()
}

// ShouldTriggerDiversion checks if traffic diversion should be triggered
func (r *TrafficDiversionRule) ShouldTriggerDiversion(currentDensity float64) bool {
	if !r.Enable {
		return false
	}
	return currentDensity >= r.TriggerDensity
}
