package services

import (
	"testing"
	"time"
)

func TestAVPDispatchService_StartAutoParkWithSession(t *testing.T) {
	svc := NewAVPDispatchService()
	task := svc.StartAutoParkWithSession("session_1", "user_1", "vehicle_1", "lot_1", "dropoff_a", "space_101")

	if task.ID == "" {
		t.Fatalf("expected task ID")
	}
	if task.SessionID != "session_1" {
		t.Fatalf("expected session_id=session_1, got %s", task.SessionID)
	}
	if task.TaskType != AVPTaskTypeAutoPark {
		t.Fatalf("expected auto park task type")
	}
	if task.Status != AVPTaskStatusExecuting {
		t.Fatalf("expected initial status executing, got %s", task.Status)
	}
}

func TestAVPDispatchService_AutoProgressToCompleted(t *testing.T) {
	svc := NewAVPDispatchService()
	task := svc.StartSummon("user_1", "vehicle_1", "lot_1", "pickup_a")

	// Simulate elapsed time so that lifecycle advancement reaches completed.
	task.CreatedAt = time.Now().Add(-2 * time.Minute)
	got, ok := svc.GetTask(task.ID)
	if !ok {
		t.Fatalf("expected task to exist")
	}
	if got.Status != AVPTaskStatusCompleted {
		t.Fatalf("expected completed status, got %s", got.Status)
	}
	if got.Progress != 100 {
		t.Fatalf("expected progress 100, got %d", got.Progress)
	}
}

func TestAVPDispatchService_CancelTask(t *testing.T) {
	svc := NewAVPDispatchService()
	task := svc.StartAutoPark("user_1", "vehicle_1", "lot_1", "dropoff_a", "space_101")

	cancelled, err := svc.CancelTask(task.ID)
	if err != nil {
		t.Fatalf("unexpected cancel error: %v", err)
	}
	if cancelled.Status != AVPTaskStatusCancelled {
		t.Fatalf("expected cancelled status, got %s", cancelled.Status)
	}
	if cancelled.LastCheckpoint != "minimum_risk_stop" {
		t.Fatalf("expected minimum risk stop checkpoint, got %s", cancelled.LastCheckpoint)
	}
}
