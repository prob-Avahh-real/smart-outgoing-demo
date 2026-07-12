package services

import (
	"testing"
	"time"
)

func TestParkingAppService_FindParking_MockFallback(t *testing.T) {
	app := NewParkingAppService(nil, NewParkingMockService())
	recs, err := app.FindParking(22.6913, 114.0448, 20, 5, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 recommendations, got %d", len(recs))
	}
}

func TestParkingAppService_ListLotsAndSpaces(t *testing.T) {
	app := NewParkingAppService(nil, NewParkingMockService())

	lots := app.ListLots(22.6913, 114.0448)
	if len(lots) == 0 {
		t.Fatalf("expected parking lots")
	}
	if lots[0].Distance <= 0 {
		t.Fatalf("expected distance to be calculated")
	}

	spaces := app.GetSpaces("lot_1")
	if len(spaces) < 2 {
		t.Fatalf("expected at least 2 spaces")
	}
}

func TestParkingAppService_ReserveAndSession(t *testing.T) {
	app := NewParkingAppService(nil, NewParkingMockService())
	start := time.Now()
	end := start.Add(2 * time.Hour)

	reservation, err := app.ReserveSpace("user_1", "lot_1", "space_1", start, end)
	if err != nil {
		t.Fatalf("unexpected reserve error: %v", err)
	}
	if reservation.UserID != "user_1" {
		t.Fatalf("unexpected reservation user")
	}

	session, err := app.StartSession("user_1", "lot_1", "space_1")
	if err != nil {
		t.Fatalf("unexpected session error: %v", err)
	}
	if session.ParkingLotID != "lot_1" {
		t.Fatalf("unexpected session lot id")
	}
}
