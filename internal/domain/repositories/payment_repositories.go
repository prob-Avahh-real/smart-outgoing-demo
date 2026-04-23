package repositories

import (
	"smart-outgoing-demo/internal/domain/entities"
	"time"
)

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	// Save saves a payment record
	Save(payment *entities.Payment) error
	
	// FindByID finds a payment by ID
	FindByID(id string) (*entities.Payment, error)
	
	// FindByUserID finds all payments for a user
	FindByUserID(userID string) ([]*entities.Payment, error)
	
	// FindByReservationID finds payment by reservation ID
	FindByReservationID(reservationID string) (*entities.Payment, error)
	
	// FindByStatus finds payments by status
	FindByStatus(status entities.PaymentStatus) ([]*entities.Payment, error)
	
	// FindByDateRange finds payments within a date range
	FindByDateRange(start, end time.Time) ([]*entities.Payment, error)
	
	// Update updates a payment record
	Update(payment *entities.Payment) error
	
	// Delete deletes a payment record
	Delete(id string) error
	
	// List lists all payments with pagination
	List(offset, limit int) ([]*entities.Payment, error)
	
	// Count returns the total count of payments
	Count() (int64, error)
}

// MockPaymentRepository provides a mock implementation for testing
type MockPaymentRepository struct {
	payments map[string]*entities.Payment
}

// NewMockPaymentRepository creates a new mock payment repository
func NewMockPaymentRepository() *MockPaymentRepository {
	return &MockPaymentRepository{
		payments: make(map[string]*entities.Payment),
	}
}

// Save saves a payment record
func (r *MockPaymentRepository) Save(payment *entities.Payment) error {
	r.payments[payment.ID] = payment
	return nil
}

// FindByID finds a payment by ID
func (r *MockPaymentRepository) FindByID(id string) (*entities.Payment, error) {
	payment, exists := r.payments[id]
	if !exists {
		return nil, ErrPaymentNotFound
	}
	return payment, nil
}

// FindByUserID finds all payments for a user
func (r *MockPaymentRepository) FindByUserID(userID string) ([]*entities.Payment, error) {
	var userPayments []*entities.Payment
	for _, payment := range r.payments {
		if payment.UserID == userID {
			userPayments = append(userPayments, payment)
		}
	}
	return userPayments, nil
}

// FindByReservationID finds payment by reservation ID
func (r *MockPaymentRepository) FindByReservationID(reservationID string) (*entities.Payment, error) {
	for _, payment := range r.payments {
		if payment.ReservationID == reservationID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

// FindByStatus finds payments by status
func (r *MockPaymentRepository) FindByStatus(status entities.PaymentStatus) ([]*entities.Payment, error) {
	var statusPayments []*entities.Payment
	for _, payment := range r.payments {
		if payment.Status == status {
			statusPayments = append(statusPayments, payment)
		}
	}
	return statusPayments, nil
}

// FindByDateRange finds payments within a date range
func (r *MockPaymentRepository) FindByDateRange(start, end time.Time) ([]*entities.Payment, error) {
	var rangePayments []*entities.Payment
	for _, payment := range r.payments {
		if payment.CreatedAt.After(start) && payment.CreatedAt.Before(end) {
			rangePayments = append(rangePayments, payment)
		}
	}
	return rangePayments, nil
}

// Update updates a payment record
func (r *MockPaymentRepository) Update(payment *entities.Payment) error {
	if _, exists := r.payments[payment.ID]; !exists {
		return ErrPaymentNotFound
	}
	r.payments[payment.ID] = payment
	return nil
}

// Delete deletes a payment record
func (r *MockPaymentRepository) Delete(id string) error {
	if _, exists := r.payments[id]; !exists {
		return ErrPaymentNotFound
	}
	delete(r.payments, id)
	return nil
}

// List lists all payments with pagination
func (r *MockPaymentRepository) List(offset, limit int) ([]*entities.Payment, error) {
	var allPayments []*entities.Payment
	for _, payment := range r.payments {
		allPayments = append(allPayments, payment)
	}
	
	// Simple pagination
	if offset >= len(allPayments) {
		return []*entities.Payment{}, nil
	}
	
	end := offset + limit
	if end > len(allPayments) {
		end = len(allPayments)
	}
	
	return allPayments[offset:end], nil
}

// Count returns the total count of payments
func (r *MockPaymentRepository) Count() (int64, error) {
	return int64(len(r.payments)), nil
}

// Error definitions
var (
	ErrPaymentNotFound = &RepositoryError{
		Code:    "PAYMENT_NOT_FOUND",
		Message: "Payment not found",
	}
	ErrInvalidPayment = &RepositoryError{
		Code:    "INVALID_PAYMENT",
		Message: "Invalid payment data",
	}
)

// RepositoryError represents a repository error
type RepositoryError struct {
	Code    string
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}
