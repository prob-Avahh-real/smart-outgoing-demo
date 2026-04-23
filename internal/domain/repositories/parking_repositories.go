package repositories

import (
	"smart-outgoing-demo/internal/domain/entities"
	"time"
)

// ParkingLotRepository defines the interface for parking lot data access
type ParkingLotRepository interface {
	Save(lot *entities.ParkingLot) error
	FindByID(id string) (*entities.ParkingLot, error)
	FindAll() ([]*entities.ParkingLot, error)
	FindAvailable() ([]*entities.ParkingLot, error)
	FindByLocation(lat, lng, radius float64) ([]*entities.ParkingLot, error)
	Update(lot *entities.ParkingLot) error
	Delete(id string) error
}

// ParkingSpaceRepository defines the interface for parking space data access
type ParkingSpaceRepository interface {
	Save(space *entities.ParkingSpace) error
	FindByID(id string) (*entities.ParkingSpace, error)
	FindAll() ([]*entities.ParkingSpace, error)
	FindByLot(lotID string) ([]*entities.ParkingSpace, error)
	FindAvailableByLot(lotID string) ([]*entities.ParkingSpace, error)
	FindByType(spaceType entities.SpaceType) ([]*entities.ParkingSpace, error)
	Update(space *entities.ParkingSpace) error
	Delete(id string) error
}

// ParkingReservationRepository defines the interface for reservation data access
type ParkingReservationRepository interface {
	Save(reservation *entities.ParkingReservation) error
	FindByID(id string) (*entities.ParkingReservation, error)
	FindByUser(userID string) ([]*entities.ParkingReservation, error)
	FindByLot(lotID string) ([]*entities.ParkingReservation, error)
	FindConflicts(spaceID string, startTime, endTime time.Time) ([]*entities.ParkingReservation, error)
	FindActive() ([]*entities.ParkingReservation, error)
	Update(reservation *entities.ParkingReservation) error
	Delete(id string) error
}

// ParkingSessionRepository defines the interface for parking session data access
type ParkingSessionRepository interface {
	Save(session *entities.ParkingSession) error
	FindByID(id string) (*entities.ParkingSession, error)
	FindByUser(userID string) ([]*entities.ParkingSession, error)
	FindByLot(lotID string) ([]*entities.ParkingSession, error)
	FindActive() ([]*entities.ParkingSession, error)
	Update(session *entities.ParkingSession) error
	Delete(id string) error
}

// UserParkingPreferenceRepository defines the interface for user preference data access
type UserParkingPreferenceRepository interface {
	Save(preference *entities.UserParkingPreference) error
	FindByUser(userID string) (*entities.UserParkingPreference, error)
	Update(preference *entities.UserParkingPreference) error
	Delete(userID string) error
}
