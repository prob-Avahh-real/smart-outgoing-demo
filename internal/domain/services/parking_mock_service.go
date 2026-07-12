package services

import (
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
)

// ParkingMockService centralizes demo-only parking mock generation.
type ParkingMockService struct{}

func NewParkingMockService() *ParkingMockService {
	return &ParkingMockService{}
}

func (s *ParkingMockService) CreateMockReservation(userID, lotID, spaceID string, startTime, endTime time.Time) *entities.ParkingReservation {
	now := time.Now()
	return &entities.ParkingReservation{
		ID:           fmt.Sprintf("res_%d", now.UnixNano()),
		UserID:       userID,
		ParkingLotID: lotID,
		SpaceID:      spaceID,
		StartTime:    startTime,
		EndTime:      endTime,
		Status:       entities.StatusConfirmed,
		TotalPrice:   30.0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (s *ParkingMockService) CreateMockSession(userID, lotID, spaceID string) *entities.ParkingSession {
	now := time.Now()
	return &entities.ParkingSession{
		ID:           fmt.Sprintf("session_%d", now.UnixNano()),
		UserID:       userID,
		ParkingLotID: lotID,
		SpaceID:      spaceID,
		StartTime:    now,
		Status:       entities.SessionActive,
		TotalCost:    0.0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (s *ParkingMockService) GenerateRecommendations(userLat, userLng float64, limit int) []*entities.ParkingRecommendation {
	mockLots := []*entities.ParkingLot{
		{ID: "lot_1", Name: "CBD Central Parking", Address: "123 Main Street", Latitude: userLat + 0.001, Longitude: userLng + 0.001, TotalSpaces: 200, AvailableSpaces: 45, PricePerHour: 15.0, Features: []string{"covered", "24/7", "ev_charging"}, Rating: 4.5, IsOpen: true, LastUpdated: time.Now()},
		{ID: "lot_2", Name: "Shopping Mall Parking", Address: "456 Shopping Ave", Latitude: userLat - 0.001, Longitude: userLng - 0.001, TotalSpaces: 150, AvailableSpaces: 12, PricePerHour: 10.0, Features: []string{"covered", "security"}, Rating: 4.2, IsOpen: true, LastUpdated: time.Now()},
		{ID: "lot_3", Name: "Airport Parking", Address: "789 Airport Road", Latitude: userLat + 0.002, Longitude: userLng - 0.002, TotalSpaces: 300, AvailableSpaces: 89, PricePerHour: 8.0, Features: []string{"24/7", "security", "shuttle"}, Rating: 4.0, IsOpen: true, LastUpdated: time.Now()},
	}

	var recommendations []*entities.ParkingRecommendation
	for i, lot := range mockLots {
		distance := DistanceKM(userLat, userLng, lot.Latitude, lot.Longitude)
		lot.Distance = distance
		space := &entities.ParkingSpace{
			ID:           fmt.Sprintf("space_%d", i+1),
			ParkingLotID: lot.ID,
			SpaceNumber:  fmt.Sprintf("A-%03d", i+1),
			Type:         entities.SpaceTypeRegular,
			IsAvailable:  true,
			IsReserved:   false,
			Price:        lot.PricePerHour,
			Floor:        "A",
			Zone:         "North",
			LastUpdated:  time.Now(),
		}
		recommendations = append(recommendations, &entities.ParkingRecommendation{
			ParkingLot:       lot,
			RecommendedSpace: space,
			Score:            85.0 - float64(i)*10,
			Reasons: []string{
				fmt.Sprintf("Only %.1f km away", distance),
				fmt.Sprintf("¥%.1f/hour", lot.PricePerHour),
				fmt.Sprintf("%d spaces available", lot.AvailableSpaces),
			},
			EstimatedTime: time.Duration(distance*10) * time.Minute,
			TotalPrice:    lot.PricePerHour * 2.0,
			Route: &entities.ParkingRoute{
				Steps:         []entities.RouteStep{{Instruction: fmt.Sprintf("Drive %.1f km to %s", distance, lot.Name), Distance: distance, Duration: time.Duration(distance*10) * time.Minute, Direction: "Towards destination"}},
				TotalDistance: distance,
				TotalTime:     time.Duration(distance*10) * time.Minute,
				Instructions:  fmt.Sprintf("Navigate to %s, %s", lot.Name, lot.Address),
			},
		})
		if limit > 0 && len(recommendations) >= limit {
			break
		}
	}
	return recommendations
}

func (s *ParkingMockService) MockParkingLots() []*entities.ParkingLot {
	now := time.Now()
	return []*entities.ParkingLot{
		{
			ID:              "lot_1",
			Name:            "CBD Central Parking",
			Address:         "123 Main Street",
			Latitude:        22.6913,
			Longitude:       114.0448,
			TotalSpaces:     200,
			AvailableSpaces: 45,
			PricePerHour:    15.0,
			Features:        []string{"covered", "24/7", "ev_charging"},
			Rating:          4.5,
			IsOpen:          true,
			LastUpdated:     now,
		},
		{
			ID:              "lot_2",
			Name:            "Shopping Mall Parking",
			Address:         "456 Shopping Ave",
			Latitude:        22.6950,
			Longitude:       114.0500,
			TotalSpaces:     150,
			AvailableSpaces: 12,
			PricePerHour:    10.0,
			Features:        []string{"covered", "security"},
			Rating:          4.2,
			IsOpen:          true,
			LastUpdated:     now,
		},
	}
}

func (s *ParkingMockService) MockParkingLotByID(lotID string) *entities.ParkingLot {
	for _, lot := range s.MockParkingLots() {
		if lot.ID == lotID {
			return lot
		}
	}
	// Keep demo compatibility for unknown IDs.
	lot := s.MockParkingLots()[0]
	lot.ID = lotID
	return lot
}

func (s *ParkingMockService) MockParkingSpacesByLot(lotID string) []*entities.ParkingSpace {
	now := time.Now()
	return []*entities.ParkingSpace{
		{
			ID:           "space_1",
			ParkingLotID: lotID,
			SpaceNumber:  "A-101",
			Type:         entities.SpaceTypeRegular,
			IsAvailable:  true,
			IsReserved:   false,
			Price:        15.0,
			Floor:        "A",
			Zone:         "North",
			LastUpdated:  now,
		},
		{
			ID:           "space_2",
			ParkingLotID: lotID,
			SpaceNumber:  "A-102",
			Type:         entities.SpaceTypeEV,
			IsAvailable:  true,
			IsReserved:   false,
			Price:        20.0,
			Floor:        "A",
			Zone:         "North",
			LastUpdated:  now,
		},
	}
}
