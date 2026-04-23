package domain

import (
	"time"
)

// Vehicle 聚合根 - 核心车辆实体
type Vehicle struct {
	id          string
	name        string
	role        string
	coordinates Coordinates
	status      VehicleStatus
	destination *Coordinates
	createdAt   time.Time
	updatedAt   time.Time
	version     int // 乐观锁版本
}

// VehicleStatus 车辆状态
type VehicleStatus int

const (
	VehicleStatusIdle VehicleStatus = iota
	VehicleStatusMoving
	VehicleStatusArrived
	VehicleStatusMaintenance
)

// String 返回状态的字符串表示
func (vs VehicleStatus) String() string {
	switch vs {
	case VehicleStatusIdle:
		return "idle"
	case VehicleStatusMoving:
		return "moving"
	case VehicleStatusArrived:
		return "arrived"
	case VehicleStatusMaintenance:
		return "maintenance"
	default:
		return "unknown"
	}
}

// Coordinates 值对象 - 地理坐标
type Coordinates struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
	Altitude  float64 `json:"alt,omitempty"`
}

// NewVehicle 创建新的车辆聚合
func NewVehicle(id, name, role string, coords Coordinates) *Vehicle {
	now := time.Now()
	return &Vehicle{
		id:          id,
		name:        name,
		role:        role,
		coordinates: coords,
		status:      VehicleStatusIdle,
		createdAt:   now,
		updatedAt:   now,
		version:     1,
	}
}

// Getters
func (v *Vehicle) ID() string                { return v.id }
func (v *Vehicle) Name() string              { return v.name }
func (v *Vehicle) Role() string              { return v.role }
func (v *Vehicle) Coordinates() Coordinates  { return v.coordinates }
func (v *Vehicle) Status() VehicleStatus     { return v.status }
func (v *Vehicle) Destination() *Coordinates { return v.destination }
func (v *Vehicle) CreatedAt() time.Time      { return v.createdAt }
func (v *Vehicle) UpdatedAt() time.Time      { return v.updatedAt }
func (v *Vehicle) Version() int              { return v.version }

// SetDestination 设置目的地 - 领域业务逻辑
func (v *Vehicle) SetDestination(coords Coordinates) {
	v.destination = &coords
	v.status = VehicleStatusMoving
	v.updatedAt = time.Now()
	v.version++
}

// Arrive 到达目的地
func (v *Vehicle) Arrive() {
	v.status = VehicleStatusArrived
	v.destination = nil
	v.updatedAt = time.Now()
	v.version++
}

// SetMaintenance 设置维护状态
func (v *Vehicle) SetMaintenance() {
	v.status = VehicleStatusMaintenance
	v.updatedAt = time.Now()
	v.version++
}

// UpdatePosition 更新位置
func (v *Vehicle) UpdatePosition(coords Coordinates) {
	v.coordinates = coords
	v.updatedAt = time.Now()
	v.version++
}

// Equals 检查坐标相等
func (c Coordinates) Equals(other Coordinates) bool {
	const epsilon = 1e-9
	return abs(c.Longitude-other.Longitude) < epsilon &&
		abs(c.Latitude-other.Latitude) < epsilon &&
		abs(c.Altitude-other.Altitude) < epsilon
}

// DistanceTo 计算到另一个坐标的距离（简化版）
func (c Coordinates) DistanceTo(other Coordinates) float64 {
	// 简化的欧几里得距离计算
	dx := c.Longitude - other.Longitude
	dy := c.Latitude - other.Latitude
	dz := c.Altitude - other.Altitude
	return sqrt(dx*dx + dy*dy + dz*dz)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func sqrt(x float64) float64 {
	// 简化的平方根实现
	if x == 0 {
		return 0
	}
	guess := x
	for i := 0; i < 10; i++ {
		guess = (guess + x/guess) / 2
	}
	return guess
}
