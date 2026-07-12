package services

import (
	"time"

	"smart-outgoing-demo/internal/domain/entities"
)

// ParkingAppService is the single application entry for parking use-cases.
type ParkingAppService struct {
	recommendation *ParkingRecommendationService
	mock           *ParkingMockService
}

func NewParkingAppService(recommendation *ParkingRecommendationService, mock *ParkingMockService) *ParkingAppService {
	if mock == nil {
		mock = NewParkingMockService()
	}
	return &ParkingAppService{
		recommendation: recommendation,
		mock:           mock,
	}
}

func (s *ParkingAppService) FindParking(lat, lng, maxPrice, maxDistance float64, limit int) ([]*entities.ParkingRecommendation, error) {
	if s.recommendation == nil {
		return s.mock.GenerateRecommendations(lat, lng, limit), nil
	}

	preferences := &entities.UserParkingPreference{
		MaxPricePerHour:   maxPrice,
		PreferredDistance: maxDistance,
		VehicleType:       entities.SpaceTypeRegular,
		PreferredFeatures: []string{},
		PreferCovered:     false,
		PreferEV:          false,
	}
	return s.recommendation.FindBestParking(lat, lng, preferences, limit)
}

func (s *ParkingAppService) ReserveSpace(userID, lotID, spaceID string, startTime, endTime time.Time) (*entities.ParkingReservation, error) {
	if s.recommendation != nil {
		return s.recommendation.ReserveSpace(userID, lotID, spaceID, startTime, endTime)
	}
	return s.mock.CreateMockReservation(userID, lotID, spaceID, startTime, endTime), nil
}

func (s *ParkingAppService) StartSession(userID, lotID, spaceID string) (*entities.ParkingSession, error) {
	if s.recommendation != nil {
		return s.recommendation.StartParkingSession(userID, lotID, spaceID)
	}
	return s.mock.CreateMockSession(userID, lotID, spaceID), nil
}

func (s *ParkingAppService) ListLots(lat, lng float64) []*entities.ParkingLot {
	lots := s.mock.MockParkingLots()
	if lat != 0 && lng != 0 {
		for _, lot := range lots {
			lot.Distance = DistanceKM(lat, lng, lot.Latitude, lot.Longitude)
		}
	}
	return lots
}

func (s *ParkingAppService) GetLot(lotID string) *entities.ParkingLot {
	return s.mock.MockParkingLotByID(lotID)
}

func (s *ParkingAppService) GetSpaces(lotID string) []*entities.ParkingSpace {
	return s.mock.MockParkingSpacesByLot(lotID)
}

// MockService exposes the mock backend for AVP session bootstrap in demo mode.
func (s *ParkingAppService) MockService() *ParkingMockService {
	return s.mock
}
