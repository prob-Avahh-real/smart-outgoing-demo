package services

import (
	"fmt"
	"sync"
	"time"
)

// AVPTaskStatus represents lifecycle state for an AVP task.
type AVPTaskStatus string

const (
	AVPTaskStatusQueued      AVPTaskStatus = "queued"
	AVPTaskStatusDispatching AVPTaskStatus = "dispatching"
	AVPTaskStatusExecuting   AVPTaskStatus = "executing"
	AVPTaskStatusCompleted   AVPTaskStatus = "completed"
	AVPTaskStatusCancelled   AVPTaskStatus = "cancelled"
)

// AVPTaskType represents AVP mission type.
type AVPTaskType string

const (
	AVPTaskTypeAutoPark AVPTaskType = "auto_park"
	AVPTaskTypeSummon   AVPTaskType = "summon"
)

// AVPTask is a lightweight domain model for AVP orchestration.
type AVPTask struct {
	ID             string        `json:"id"`
	SessionID      string        `json:"session_id,omitempty"`
	UserID         string        `json:"user_id"`
	VehicleID      string        `json:"vehicle_id"`
	ParkingLotID   string        `json:"parking_lot_id"`
	SourceZone     string        `json:"source_zone"`
	TargetZone     string        `json:"target_zone"`
	TargetSpaceID  string        `json:"target_space_id,omitempty"`
	TaskType       AVPTaskType   `json:"task_type"`
	Status         AVPTaskStatus `json:"status"`
	SafetyMode     string        `json:"safety_mode"`
	Progress       int           `json:"progress"`
	LastCheckpoint string        `json:"last_checkpoint"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// AVPDispatchService provides in-memory AVP orchestration for demo.
type AVPDispatchService struct {
	mu    sync.RWMutex
	tasks map[string]*AVPTask
}

func NewAVPDispatchService() *AVPDispatchService {
	return &AVPDispatchService{
		tasks: make(map[string]*AVPTask),
	}
}

func (s *AVPDispatchService) StartAutoPark(userID, vehicleID, lotID, dropoffZone, targetSpaceID string) *AVPTask {
	return s.StartAutoParkWithSession("", userID, vehicleID, lotID, dropoffZone, targetSpaceID)
}

func (s *AVPDispatchService) StartAutoParkWithSession(sessionID, userID, vehicleID, lotID, dropoffZone, targetSpaceID string) *AVPTask {
	now := time.Now()
	task := &AVPTask{
		ID:             fmt.Sprintf("avp_%d", now.UnixNano()),
		SessionID:      sessionID,
		UserID:         userID,
		VehicleID:      vehicleID,
		ParkingLotID:   lotID,
		SourceZone:     dropoffZone,
		TargetZone:     "parking_slot_area",
		TargetSpaceID:  targetSpaceID,
		TaskType:       AVPTaskTypeAutoPark,
		Status:         AVPTaskStatusExecuting,
		SafetyMode:     "conservative",
		Progress:       35,
		LastCheckpoint: "left_dropoff_zone",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	s.mu.Lock()
	s.tasks[task.ID] = task
	s.mu.Unlock()
	return task
}

func (s *AVPDispatchService) StartSummon(userID, vehicleID, lotID, pickupZone string) *AVPTask {
	now := time.Now()
	task := &AVPTask{
		ID:             fmt.Sprintf("avp_%d", now.UnixNano()),
		UserID:         userID,
		VehicleID:      vehicleID,
		ParkingLotID:   lotID,
		SourceZone:     "parking_slot_area",
		TargetZone:     pickupZone,
		TaskType:       AVPTaskTypeSummon,
		Status:         AVPTaskStatusDispatching,
		SafetyMode:     "conservative",
		Progress:       20,
		LastCheckpoint: "route_allocated",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	s.mu.Lock()
	s.tasks[task.ID] = task
	s.mu.Unlock()
	return task
}

func (s *AVPDispatchService) GetTask(taskID string) (*AVPTask, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[taskID]
	if ok {
		s.advanceTaskLocked(task)
	}
	return task, ok
}

func (s *AVPDispatchService) CancelTask(taskID string) (*AVPTask, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}

	if task.Status == AVPTaskStatusCompleted || task.Status == AVPTaskStatusCancelled {
		return nil, fmt.Errorf("task already terminal")
	}

	task.Status = AVPTaskStatusCancelled
	task.Progress = 0
	task.LastCheckpoint = "minimum_risk_stop"
	task.UpdatedAt = time.Now()
	return task, nil
}

func (s *AVPDispatchService) advanceTaskLocked(task *AVPTask) {
	if task == nil {
		return
	}
	if task.Status == AVPTaskStatusCancelled || task.Status == AVPTaskStatusCompleted {
		return
	}

	elapsed := time.Since(task.CreatedAt)

	switch task.TaskType {
	case AVPTaskTypeAutoPark:
		if elapsed >= 90*time.Second {
			task.Status = AVPTaskStatusCompleted
			task.Progress = 100
			task.LastCheckpoint = "parked_in_target_space"
		} else if elapsed >= 25*time.Second {
			task.Status = AVPTaskStatusExecuting
			task.Progress = 70
			task.LastCheckpoint = "navigating_to_target_space"
		}
	case AVPTaskTypeSummon:
		if elapsed >= 80*time.Second {
			task.Status = AVPTaskStatusCompleted
			task.Progress = 100
			task.LastCheckpoint = "arrived_pickup_zone"
		} else if elapsed >= 20*time.Second {
			task.Status = AVPTaskStatusExecuting
			task.Progress = 60
			task.LastCheckpoint = "en_route_to_pickup_zone"
		}
	}

	task.UpdatedAt = time.Now()
}
