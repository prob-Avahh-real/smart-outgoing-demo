package services

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
)

// ParkingPoolService manages three-tier parking pool classification
type ParkingPoolService struct {
	mu    sync.RWMutex
	pools map[string]*entities.ParkingPool
	rules []*entities.TrafficDiversionRule
}

// NewParkingPoolService creates a new parking pool service
func NewParkingPoolService() *ParkingPoolService {
	service := &ParkingPoolService{
		pools: make(map[string]*entities.ParkingPool),
		rules: make([]*entities.TrafficDiversionRule, 0),
	}

	// Initialize default three-tier pools
	service.initializeDefaultPools()

	return service
}

// initializeDefaultPools initializes the three-tier parking pool structure
func (s *ParkingPoolService) initializeDefaultPools() {
	// Core business district (priority 1)
	corePool := &entities.ParkingPool{
		ID:           "pool_core",
		Level:        entities.PoolLevelCore,
		Name:         "核心商圈车场池",
		Description:  "优先推荐，位于核心商圈，周转率高",
		Priority:     1,
		Lots:         make([]*entities.ParkingLot, 0),
		TurnoverRate: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Peripheral supporting (priority 2)
	peripheralPool := &entities.ParkingPool{
		ID:           "pool_peripheral",
		Level:        entities.PoolLevelPeripheral,
		Name:         "外围配套车场池",
		Description:  "备用推荐，位于外围区域，容量充足",
		Priority:     2,
		Lots:         make([]*entities.ParkingLot, 0),
		TurnoverRate: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Roadside temporary (priority 3)
	roadsidePool := &entities.ParkingPool{
		ID:           "pool_roadside",
		Level:        entities.PoolLevelRoadside,
		Name:         "路边临停车场池",
		Description:  "兜底方案，路边临时停车位",
		Priority:     3,
		Lots:         make([]*entities.ParkingLot, 0),
		TurnoverRate: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	s.mu.Lock()
	s.pools[corePool.ID] = corePool
	s.pools[peripheralPool.ID] = peripheralPool
	s.pools[roadsidePool.ID] = roadsidePool
	s.mu.Unlock()
}

// ClassifyParkingLot classifies a parking lot into one of the three tiers
// Based on: location, capacity, turnover rate, price
func (s *ParkingPoolService) ClassifyParkingLot(lot *entities.ParkingLot) entities.ParkingPoolLevel {
	// Classification logic:
	// Core: within 2km of center, high turnover (>5/hour), premium location
	// Peripheral: 2-5km from center, medium turnover (2-5/hour)
	// Roadside: >5km or temporary parking, low turnover (<2/hour)

	centerLat, centerLng := 22.6913, 114.0448 // Shenzhen Longhua center
	distance := DistanceKM(lot.Latitude, lot.Longitude, centerLat, centerLng)

	// Calculate turnover rate (simplified)
	turnoverRate := s.calculateTurnoverRate(lot)

	if distance <= 2.0 && turnoverRate >= 5.0 {
		return entities.PoolLevelCore
	} else if distance <= 5.0 && turnoverRate >= 2.0 {
		return entities.PoolLevelPeripheral
	} else {
		return entities.PoolLevelRoadside
	}
}

// AddParkingLotToPool adds a parking lot to the appropriate pool
func (s *ParkingPoolService) AddParkingLotToPool(lot *entities.ParkingLot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	level := s.ClassifyParkingLot(lot)

	var targetPool *entities.ParkingPool
	for _, pool := range s.pools {
		if pool.Level == level {
			targetPool = pool
			break
		}
	}

	if targetPool == nil {
		return fmt.Errorf("no pool found for level: %s", level)
	}

	// Check if lot already exists
	for _, existingLot := range targetPool.Lots {
		if existingLot.ID == lot.ID {
			return fmt.Errorf("parking lot %s already exists in pool %s", lot.ID, targetPool.ID)
		}
	}

	targetPool.Lots = append(targetPool.Lots, lot)
	targetPool.UpdateStats()

	return nil
}

// GetRecommendedParkingLot gets the best parking lot recommendation based on:
// 1. Pool priority (core > peripheral > roadside)
// 2. Distance to user
// 3. Available spaces
// 4. Turnover rate
func (s *ParkingPoolService) GetRecommendedParkingLot(
	userLat, userLng float64,
	maxDistance float64,
) (*entities.ParkingLot, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Collect all pools sorted by priority
	pools := make([]*entities.ParkingPool, 0, len(s.pools))
	for _, pool := range s.pools {
		pools = append(pools, pool)
	}

	// Sort by priority (1=highest)
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Priority < pools[j].Priority
	})

	// Iterate through pools by priority
	for _, pool := range pools {
		availableLots := pool.GetAvailableLots()
		if len(availableLots) == 0 {
			continue
		}

		// Score each available lot
		candidates := make([]struct {
			lot   *entities.ParkingLot
			score float64
		}, 0)

		for _, lot := range availableLots {
			distance := DistanceKM(userLat, userLng, lot.Latitude, lot.Longitude)
			if maxDistance > 0 && distance > maxDistance {
				continue
			}

			// Scoring formula:
			// score = (distance_weight * distance) + (availability_weight * availability) + (turnover_weight * turnover)
			// Lower distance is better, higher availability is better, higher turnover is better
			distanceScore := distance / 10.0 // Normalize distance
			availabilityScore := float64(lot.AvailableSpaces) / float64(lot.TotalSpaces)
			turnoverScore := pool.TurnoverRate / 10.0 // Normalize turnover

			totalScore := distanceScore*0.5 + availabilityScore*0.3 + turnoverScore*0.2

			candidates = append(candidates, struct {
				lot   *entities.ParkingLot
				score float64
			}{
				lot:   lot,
				score: totalScore,
			})
		}

		// Sort by score (lowest score is best)
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].score < candidates[j].score
		})

		if len(candidates) > 0 {
			reason := fmt.Sprintf("推荐%s车场池：距离%.1fkm，可用车位%d个，周转率%.1f/h",
				pool.Name,
				DistanceKM(userLat, userLng, candidates[0].lot.Latitude, candidates[0].lot.Longitude),
				candidates[0].lot.AvailableSpaces,
				pool.TurnoverRate)
			return candidates[0].lot, reason, nil
		}
	}

	return nil, "", fmt.Errorf("no available parking lots found")
}

// GetPoolStatistics returns statistics for all pools
func (s *ParkingPoolService) GetPoolStatistics() map[string]*entities.ParkingPoolStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]*entities.ParkingPoolStats)

	for _, pool := range s.pools {
		pool.UpdateStats()
		stats[pool.ID] = &entities.ParkingPoolStats{
			PoolID:        pool.ID,
			OccupancyRate: pool.GetOccupancyRate(),
			TurnoverRate:  pool.TurnoverRate,
			TotalVehicles: pool.TotalSpaces - pool.FreeSpaces,
			LastUpdated:   pool.UpdatedAt,
		}
	}

	return stats
}

// TriggerTrafficDiversion triggers traffic diversion based on current density
func (s *ParkingPoolService) TriggerTrafficDiversion(sourceZone string, currentDensity float64) ([]*entities.ParkingLot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var triggeredRules []*entities.TrafficDiversionRule
	for _, rule := range s.rules {
		if rule.SourceZone == sourceZone && rule.ShouldTriggerDiversion(currentDensity) {
			triggeredRules = append(triggeredRules, rule)
		}
	}

	if len(triggeredRules) == 0 {
		return nil, fmt.Errorf("no diversion rules triggered for zone %s with density %.2f", sourceZone, currentDensity)
	}

	// Get parking lots from target pools
	var recommendedLots []*entities.ParkingLot
	for _, rule := range triggeredRules {
		for _, pool := range s.pools {
			if pool.Level == rule.TargetPoolLevel {
				available := pool.GetAvailableLots()
				recommendedLots = append(recommendedLots, available...)
				break
			}
		}
	}

	return recommendedLots, nil
}

// AddTrafficDiversionRule adds a traffic diversion rule
func (s *ParkingPoolService) AddTrafficDiversionRule(rule *entities.TrafficDiversionRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.rules = append(s.rules, rule)
	return nil
}

// calculateTurnoverRate calculates turnover rate for a parking lot (simplified)
func (s *ParkingPoolService) calculateTurnoverRate(lot *entities.ParkingLot) float64 {
	// Simplified calculation based on occupancy rate and price
	// In production, use historical data
	occupancyRate := lot.OccupancyRate()

	// Higher price usually means higher turnover
	priceFactor := lot.PricePerHour / 10.0

	// Base turnover rate + occupancy factor + price factor
	baseRate := 2.0
	occupancyFactor := occupancyRate / 100.0 * 3.0

	return baseRate + occupancyFactor + priceFactor
}
