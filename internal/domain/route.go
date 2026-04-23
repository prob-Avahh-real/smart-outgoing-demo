package domain

import (
	"time"
)

// Route 聚合根 - 车辆路线信息
type Route struct {
	id            string
	vehicleID     string
	waypoints     []Coordinates
	distance      float64
	estimatedTime time.Duration
	status        RouteStatus
	createdAt     time.Time
	updatedAt     time.Time
	version       int
}

// RouteStatus 路线状态
type RouteStatus int

const (
	RouteStatusPlanned RouteStatus = iota
	RouteStatusActive
	RouteStatusCompleted
	RouteStatusCancelled
)

// String 返回状态的字符串表示
func (rs RouteStatus) String() string {
	switch rs {
	case RouteStatusPlanned:
		return "planned"
	case RouteStatusActive:
		return "active"
	case RouteStatusCompleted:
		return "completed"
	case RouteStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// NewRoute 创建新的路线聚合
func NewRoute(id, vehicleID string, waypoints []Coordinates) *Route {
	now := time.Now()
	distance := calculateTotalDistance(waypoints)

	return &Route{
		id:            id,
		vehicleID:     vehicleID,
		waypoints:     waypoints,
		distance:      distance,
		estimatedTime: estimateTime(distance),
		status:        RouteStatusPlanned,
		createdAt:     now,
		updatedAt:     now,
		version:       1,
	}
}

// Getters
func (r *Route) ID() string                   { return r.id }
func (r *Route) VehicleID() string            { return r.vehicleID }
func (r *Route) Waypoints() []Coordinates     { return r.waypoints }
func (r *Route) Distance() float64            { return r.distance }
func (r *Route) EstimatedTime() time.Duration { return r.estimatedTime }
func (r *Route) Status() RouteStatus          { return r.status }
func (r *Route) CreatedAt() time.Time         { return r.createdAt }
func (r *Route) UpdatedAt() time.Time         { return r.updatedAt }
func (r *Route) Version() int                 { return r.version }

// Activate 激活路线
func (r *Route) Activate() {
	r.status = RouteStatusActive
	r.updatedAt = time.Now()
	r.version++
}

// Complete 完成路线
func (r *Route) Complete() {
	r.status = RouteStatusCompleted
	r.updatedAt = time.Now()
	r.version++
}

// Cancel 取消路线
func (r *Route) Cancel() {
	r.status = RouteStatusCancelled
	r.updatedAt = time.Now()
	r.version++
}

// UpdateProgress 更新路线进度（简化实现）
func (r *Route) UpdateProgress(currentIndex int) {
	// 这里可以添加更复杂的进度跟踪逻辑
	r.updatedAt = time.Now()
	r.version++
}

// calculateTotalDistance 计算路线总距离
func calculateTotalDistance(waypoints []Coordinates) float64 {
	if len(waypoints) < 2 {
		return 0
	}

	total := 0.0
	for i := 1; i < len(waypoints); i++ {
		total += waypoints[i-1].DistanceTo(waypoints[i])
	}
	return total
}

// estimateTime 估算行驶时间（简化实现）
func estimateTime(distance float64) time.Duration {
	// 假设平均速度为 10 m/s
	avgSpeed := 10.0
	seconds := distance / avgSpeed
	return time.Duration(seconds) * time.Second
}
