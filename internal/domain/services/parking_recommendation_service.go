package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/events"
	"smart-outgoing-demo/internal/domain/repositories"
)

// ParkingRecommendationService provides intelligent parking recommendations
type ParkingRecommendationService struct {
	parkingLotRepo  repositories.ParkingLotRepository
	spaceRepo       repositories.ParkingSpaceRepository
	reservationRepo repositories.ParkingReservationRepository
	sessionRepo     repositories.ParkingSessionRepository
	eventBus        chan events.DomainEvent
	poolService     *ParkingPoolService // Three-tier parking pool service
}

// NewParkingRecommendationService creates a new parking recommendation service
func NewParkingRecommendationService(
	parkingLotRepo repositories.ParkingLotRepository,
	spaceRepo repositories.ParkingSpaceRepository,
	reservationRepo repositories.ParkingReservationRepository,
	eventBus chan events.DomainEvent,
) *ParkingRecommendationService {
	return &ParkingRecommendationService{
		parkingLotRepo:  parkingLotRepo,
		spaceRepo:       spaceRepo,
		reservationRepo: reservationRepo,
		eventBus:        eventBus,
		poolService:     NewParkingPoolService(),
	}
}

// FindBestParking finds the best parking options based on user preferences and current location
func (s *ParkingRecommendationService) FindBestParking(
	userLat, userLng float64,
	preferences *entities.UserParkingPreference,
	limit int,
) ([]*entities.ParkingRecommendation, error) {

	// Get all available parking lots
	lots, err := s.parkingLotRepo.FindAvailable()
	if err != nil {
		return nil, fmt.Errorf("failed to get parking lots: %w", err)
	}

	// Calculate distances and filter by preferences
	var recommendations []*entities.ParkingRecommendation

	for _, lot := range lots {
		// Calculate distance from user
		distance := calculateDistance(userLat, userLng, lot.Latitude, lot.Longitude)
		lot.Distance = distance

		// Check if lot meets user preferences
		if !s.meetsPreferences(lot, preferences) {
			continue
		}

		// Get available spaces for this lot
		spaces, err := s.spaceRepo.FindAvailableByLot(lot.ID)
		if err != nil {
			continue // Skip this lot if we can't get spaces
		}

		// Find best space for user
		bestSpace := s.findBestSpace(spaces, preferences)

		// Calculate recommendation score
		score := s.calculateScore(lot, bestSpace, preferences, distance)

		// Estimate travel time
		estimatedTime := s.estimateTravelTime(distance)

		// Calculate total price (estimated 2 hours)
		totalPrice := lot.PricePerHour * 2.0
		if bestSpace != nil && bestSpace.Price > 0 {
			totalPrice = bestSpace.Price * 2.0
		}

		// Generate route
		route := s.generateRouteToLot(userLat, userLng, lot, bestSpace)

		recommendation := &entities.ParkingRecommendation{
			ParkingLot:       lot,
			RecommendedSpace: bestSpace,
			Score:            score,
			Reasons:          s.generateReasons(lot, bestSpace, preferences),
			EstimatedTime:    estimatedTime,
			TotalPrice:       totalPrice,
			Route:            route,
		}

		recommendations = append(recommendations, recommendation)
	}

	// Sort by score (highest first)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit results
	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// ReserveSpace reserves a specific parking space
func (s *ParkingRecommendationService) ReserveSpace(
	userID, lotID, spaceID string,
	startTime, endTime time.Time,
) (*entities.ParkingReservation, error) {

	// Check if space is available
	space, err := s.spaceRepo.FindByID(spaceID)
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if !space.IsAvailable || space.IsReserved {
		return nil, fmt.Errorf("space is not available")
	}

	// Check for conflicting reservations
	conflicts, err := s.reservationRepo.FindConflicts(spaceID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check conflicts: %w", err)
	}

	if len(conflicts) > 0 {
		return nil, fmt.Errorf("space already reserved for the requested time")
	}

	// Get parking lot for pricing
	lot, err := s.parkingLotRepo.FindByID(lotID)
	if err != nil {
		return nil, fmt.Errorf("parking lot not found: %w", err)
	}

	// Calculate total price
	duration := endTime.Sub(startTime)
	hours := duration.Hours()
	pricePerHour := lot.PricePerHour
	if space.Price > 0 {
		pricePerHour = space.Price
	}
	totalPrice := hours * pricePerHour

	// Create reservation
	reservation := &entities.ParkingReservation{
		ID:           generateReservationID(),
		UserID:       userID,
		ParkingLotID: lotID,
		SpaceID:      spaceID,
		StartTime:    startTime,
		EndTime:      endTime,
		Status:       entities.StatusPending,
		TotalPrice:   totalPrice,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save reservation
	err = s.reservationRepo.Save(reservation)
	if err != nil {
		return nil, fmt.Errorf("failed to save reservation: %w", err)
	}

	// Mark space as reserved
	space.IsReserved = true
	err = s.spaceRepo.Update(space)
	if err != nil {
		// Rollback reservation
		s.reservationRepo.Delete(reservation.ID)
		return nil, fmt.Errorf("failed to reserve space: %w", err)
	}

	// TODO: Publish event when event system is implemented
	// if s.eventBus != nil {
	//     event := events.NewParkingEvent(...)
	//     s.eventBus <- event
	// }

	return reservation, nil
}

// StartParkingSession starts a parking session
func (s *ParkingRecommendationService) StartParkingSession(
	userID, lotID, spaceID string,
) (*entities.ParkingSession, error) {

	// Validate space
	space, err := s.spaceRepo.FindByID(spaceID)
	if err != nil {
		return nil, fmt.Errorf("space not found: %w", err)
	}

	if !space.IsAvailable {
		return nil, fmt.Errorf("space is not available")
	}

	// Create session
	session := &entities.ParkingSession{
		ID:           generateSessionID(),
		UserID:       userID,
		ParkingLotID: lotID,
		SpaceID:      spaceID,
		StartTime:    time.Now(),
		Status:       entities.SessionActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mark space as occupied
	space.IsAvailable = false
	space.IsReserved = false
	err = s.spaceRepo.Update(space)
	if err != nil {
		return nil, fmt.Errorf("failed to update space: %w", err)
	}

	// Save session
	err = s.sessionRepo.Save(session)
	if err != nil {
		// Rollback space status
		space.IsAvailable = true
		s.spaceRepo.Update(space)
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	// TODO: Publish event when event system is implemented
	// if s.eventBus != nil {
	//     event := events.NewParkingEvent(...)
	//     s.eventBus <- event
	// }

	return session, nil
}

// Helper methods

func (s *ParkingRecommendationService) meetsPreferences(
	lot *entities.ParkingLot,
	preferences *entities.UserParkingPreference,
) bool {
	// Check price
	if preferences.MaxPricePerHour > 0 && lot.PricePerHour > preferences.MaxPricePerHour {
		return false
	}

	// Check distance
	if preferences.PreferredDistance > 0 && lot.Distance > preferences.PreferredDistance {
		return false
	}

	// Check features
	if len(preferences.PreferredFeatures) > 0 {
		hasAllFeatures := true
		for _, requiredFeature := range preferences.PreferredFeatures {
			found := false
			for _, lotFeature := range lot.Features {
				if lotFeature == requiredFeature {
					found = true
					break
				}
			}
			if !found {
				hasAllFeatures = false
				break
			}
		}
		if !hasAllFeatures {
			return false
		}
	}

	return true
}

func (s *ParkingRecommendationService) findBestSpace(
	spaces []*entities.ParkingSpace,
	preferences *entities.UserParkingPreference,
) *entities.ParkingSpace {
	if len(spaces) == 0 {
		return nil
	}

	bestSpace := spaces[0]
	bestScore := 0.0

	for _, space := range spaces {
		score := s.calculateSpaceScore(space, preferences)
		if score > bestScore {
			bestScore = score
			bestSpace = space
		}
	}

	return bestSpace
}

func (s *ParkingRecommendationService) calculateSpaceScore(
	space *entities.ParkingSpace,
	preferences *entities.UserParkingPreference,
) float64 {
	score := 50.0 // Base score

	// Vehicle type match
	if space.Type == preferences.VehicleType {
		score += 20
	}

	// EV charging preference
	if preferences.PreferEV && space.Type == entities.SpaceTypeEV {
		score += 15
	}

	// Accessibility
	if !preferences.DisableAccessibility && space.Type == entities.SpaceTypeDisabled {
		score += 10
	}

	// VIP preference
	if space.Type == entities.SpaceTypeVIP {
		score += 5
	}

	return score
}

func (s *ParkingRecommendationService) calculateScore(
	lot *entities.ParkingLot,
	space *entities.ParkingSpace,
	preferences *entities.UserParkingPreference,
	distance float64,
) float64 {
	score := 0.0

	// Distance score (closer is better)
	distanceScore := math.Max(0, 100-distance*10) // 100 points, -10 per km
	score += distanceScore

	// Price score (cheaper is better)
	if preferences.MaxPricePerHour > 0 {
		priceScore := math.Max(0, 100-(lot.PricePerHour/preferences.MaxPricePerHour)*100)
		score += priceScore * 0.5 // Price is less important than distance
	}

	// Availability score (more spaces is better)
	availabilityScore := float64(lot.AvailableSpaces) / float64(lot.TotalSpaces) * 100
	score += availabilityScore * 0.3

	// Rating score
	ratingScore := lot.Rating * 20 // 5 stars = 100 points
	score += ratingScore * 0.2

	// Features match
	featureScore := 0.0
	if len(preferences.PreferredFeatures) > 0 {
		matchedFeatures := 0
		for _, requiredFeature := range preferences.PreferredFeatures {
			for _, lotFeature := range lot.Features {
				if lotFeature == requiredFeature {
					matchedFeatures++
					break
				}
			}
		}
		featureScore = float64(matchedFeatures) / float64(len(preferences.PreferredFeatures)) * 100
	}
	score += featureScore * 0.4

	return score
}

func (s *ParkingRecommendationService) generateReasons(
	lot *entities.ParkingLot,
	space *entities.ParkingSpace,
	preferences *entities.UserParkingPreference,
) []string {
	var reasons []string

	// Distance reason
	if lot.Distance < 1.0 {
		reasons = append(reasons, fmt.Sprintf("Only %.1f km away", lot.Distance))
	}

	// Price reason
	if lot.PricePerHour < 10.0 {
		reasons = append(reasons, fmt.Sprintf("Affordable at ¥%.1f/hour", lot.PricePerHour))
	}

	// Availability reason
	if lot.AvailableSpaces > 10 {
		reasons = append(reasons, fmt.Sprintf("%d spaces available", lot.AvailableSpaces))
	}

	// Rating reason
	if lot.Rating >= 4.0 {
		reasons = append(reasons, fmt.Sprintf("High rating: %.1f/5", lot.Rating))
	}

	// Features reasons
	for _, feature := range lot.Features {
		if contains(preferences.PreferredFeatures, feature) {
			reasons = append(reasons, fmt.Sprintf("Has %s", feature))
		}
	}

	// Space type reason
	if space != nil {
		switch space.Type {
		case entities.SpaceTypeEV:
			reasons = append(reasons, "EV charging available")
		case entities.SpaceTypeRegular:
			reasons = append(reasons, "Standard parking")
		case entities.SpaceTypeVIP:
			reasons = append(reasons, "VIP space")
		case entities.SpaceTypeCompact:
			reasons = append(reasons, "Compact space")
		case entities.SpaceTypeDisabled:
			reasons = append(reasons, "Disabled access")
		}
	}

	return reasons
}

func (s *ParkingRecommendationService) estimateTravelTime(distance float64) time.Duration {
	// Assume average speed of 30 km/h in city traffic
	avgSpeed := 30.0 // km/h
	hours := distance / avgSpeed
	return time.Duration(hours * float64(time.Hour))
}

func (s *ParkingRecommendationService) generateRouteToLot(
	userLat, userLng float64,
	lot *entities.ParkingLot,
	space *entities.ParkingSpace,
) *entities.ParkingRoute {
	// This would integrate with a real mapping service
	// For now, return a simple route
	steps := []entities.RouteStep{
		{
			Instruction: fmt.Sprintf("Drive %.1f km to %s", lot.Distance, lot.Name),
			Distance:    lot.Distance,
			Duration:    s.estimateTravelTime(lot.Distance),
			Direction:   "Towards destination",
		},
	}

	if space != nil {
		steps = append(steps, entities.RouteStep{
			Instruction: fmt.Sprintf("Go to space %s on %s", space.SpaceNumber, space.Floor),
			Distance:    0.1, // Estimated within parking lot
			Duration:    2 * time.Minute,
			Direction:   "Follow parking lot signs",
		})
	}

	return &entities.ParkingRoute{
		Steps:         steps,
		TotalDistance: lot.Distance,
		TotalTime:     s.estimateTravelTime(lot.Distance),
		Instructions:  fmt.Sprintf("Navigate to %s, %s", lot.Name, lot.Address),
	}
}

// Utility functions

func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Haversine formula for calculating distance between two points
	const earthRadius = 6371 // km

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func generateReservationID() string {
	return fmt.Sprintf("res_%d", time.Now().UnixNano())
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
