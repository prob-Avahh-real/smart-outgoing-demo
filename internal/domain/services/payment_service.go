package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/repositories"
	"smart-outgoing-demo/pkg/security"
)

// PaymentService handles payment processing
type PaymentService struct {
	paymentRepo     repositories.PaymentRepository
	reservationRepo repositories.ParkingReservationRepository
	weChatGateway   *entities.WeChatPayGateway
	alipayGateway   *entities.AlipayGateway
	cardGateway     *entities.CreditCardGateway
	encryptionKey   []byte
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo repositories.PaymentRepository,
	reservationRepo repositories.ParkingReservationRepository,
) *PaymentService {
	// Generate encryption key for production use
	encryptionKey := make([]byte, 32) // AES-256 requires 32-byte key
	rand.Read(encryptionKey)

	return &PaymentService{
		paymentRepo:     paymentRepo,
		reservationRepo: reservationRepo,
		weChatGateway:   &entities.WeChatPayGateway{AppID: "mock_app_id", MerchantID: "mock_merchant_id", Encrypted: true},
		alipayGateway:   &entities.AlipayGateway{AppID: "mock_app_id", MerchantID: "mock_merchant_id", Encrypted: true},
		cardGateway:     &entities.CreditCardGateway{MerchantID: "mock_merchant_id", Environment: "sandbox", Encrypted: true},
		encryptionKey:   encryptionKey,
	}
}

// encryptSensitiveData encrypts sensitive payment information
func (s *PaymentService) encryptSensitiveData(data string) (string, error) {
	return security.Encrypt(data, s.encryptionKey)
}

// decryptSensitiveData decrypts sensitive payment information
func (s *PaymentService) decryptSensitiveData(encryptedData string) (string, error) {
	return security.Decrypt(encryptedData, s.encryptionKey)
}

// ProcessPayment processes a payment request
func (s *PaymentService) ProcessPayment(req *entities.PaymentRequest) (*entities.PaymentResponse, error) {
	// Validate reservation exists
	reservation, err := s.reservationRepo.FindByID(req.ReservationID)
	if err != nil {
		return nil, fmt.Errorf("reservation not found: %w", err)
	}

	// Validate amount matches reservation
	if req.Amount != reservation.TotalPrice {
		return nil, fmt.Errorf("payment amount does not match reservation total")
	}

	// Create payment record
	payment := &entities.Payment{
		ID:            s.generatePaymentID(),
		ReservationID: req.ReservationID,
		UserID:        reservation.UserID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Method:        req.Method,
		Status:        entities.PaymentStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save payment record
	err = s.paymentRepo.Save(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Process payment based on method
	var response *entities.PaymentResponse
	switch req.Method {
	case entities.PaymentMethodWeChat:
		response, err = s.processWeChatPayment(payment, req.PaymentInfo)
	case entities.PaymentMethodAlipay:
		response, err = s.processAlipayPayment(payment, req.PaymentInfo)
	case entities.PaymentMethodCreditCard:
		err := s.ProcessCreditCardPayment(payment, req.PaymentInfo)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", req.Method)
	}

	if err != nil {
		// Mark payment as failed
		payment.MarkFailed(err.Error())
		s.paymentRepo.Update(payment)
		return nil, fmt.Errorf("payment processing failed: %w", err)
	}

	// Update payment status
	payment.Status = response.Status
	if response.TransactionID != "" {
		payment.TransactionID = response.TransactionID
	}
	s.paymentRepo.Update(payment)

	return response, nil
}

// processWeChatPayment processes WeChat Pay payment
func (s *PaymentService) processWeChatPayment(payment *entities.Payment, paymentInfo interface{}) (*entities.PaymentResponse, error) {
	// Mock WeChat Pay implementation
	// In real implementation, would call WeChat Pay API

	paymentID := s.generateTransactionID()
	qrCode := fmt.Sprintf("weixin://wxpay/bizpayurl?pr=%s", paymentID)

	response := &entities.PaymentResponse{
		PaymentID:     payment.ID,
		Status:        entities.PaymentStatusProcessing,
		TransactionID: paymentID,
		QRCode:        qrCode,
		Message:       "Please scan QR code to complete payment",
	}

	// Simulate async payment completion
	go s.simulatePaymentCompletion(payment.ID, paymentID)

	return response, nil
}

// processAlipayPayment processes Alipay payment
func (s *PaymentService) processAlipayPayment(payment *entities.Payment, paymentInfo interface{}) (*entities.PaymentResponse, error) {
	// Mock Alipay implementation
	// In real implementation, would call Alipay API

	paymentID := s.generateTransactionID()
	paymentURL := fmt.Sprintf("https://openapi.alipay.com/gateway.do?app_id=mock&trade_no=%s", paymentID)

	response := &entities.PaymentResponse{
		PaymentID:     payment.ID,
		Status:        entities.PaymentStatusProcessing,
		TransactionID: paymentID,
		PaymentURL:    paymentURL,
		Message:       "Please complete payment on Alipay page",
	}

	// Simulate async payment completion
	go s.simulatePaymentCompletion(payment.ID, paymentID)

	return response, nil
}

// ProcessCreditCardPayment processes credit card payment with encryption
func (s *PaymentService) ProcessCreditCardPayment(payment *entities.Payment, paymentInfo interface{}) error {
	cardInfo, ok := paymentInfo.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid credit card info format")
	}

	// Validate card information
	cardNumber, ok := cardInfo["card_number"].(string)
	if !ok || cardNumber == "" {
		return fmt.Errorf("card number is required")
	}

	cardHolder, ok := cardInfo["card_holder_name"].(string)
	if !ok || cardHolder == "" {
		return fmt.Errorf("card holder name is required")
	}

	expiryMonth, ok := cardInfo["expiry_month"].(string)
	if !ok || expiryMonth == "" {
		return fmt.Errorf("expiry month is required")
	}

	expiryYear, ok := cardInfo["expiry_year"].(string)
	if !ok || expiryYear == "" {
		return fmt.Errorf("expiry year is required")
	}

	cvv, ok := cardInfo["cvv"].(string)
	if !ok || cvv == "" {
		return fmt.Errorf("CVV is required")
	}

	// Encrypt sensitive card information
	encryptedCardNumber, err := s.encryptSensitiveData(cardNumber)
	if err != nil {
		return fmt.Errorf("failed to encrypt card number: %w", err)
	}

	encryptedCVV, err := s.encryptSensitiveData(cvv)
	if err != nil {
		return fmt.Errorf("failed to encrypt CVV: %w", err)
	}

	// Store only encrypted data (in real implementation, use tokenization)
	// For mock, we just simulate success
	payment.Status = entities.PaymentStatusCompleted
	payment.TransactionID = s.generateTransactionID()
	payment.UpdatedAt = time.Now()

	// In production, store encrypted data securely
	_ = encryptedCardNumber
	_ = encryptedCVV

	return nil
}

// simulatePaymentCompletion simulates async payment completion for QR code and redirect payments
func (s *PaymentService) simulatePaymentCompletion(paymentID, transactionID string) {
	// Simulate payment processing delay
	time.Sleep(30 * time.Second)

	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		return
	}

	// Simulate 90% success rate
	if time.Now().UnixNano()%10 < 9 {
		payment.MarkCompleted(transactionID)
		s.paymentRepo.Update(payment)
	} else {
		payment.MarkFailed("Payment timeout or user cancelled")
		s.paymentRepo.Update(payment)
	}
}

// GetPaymentStatus retrieves payment status
func (s *PaymentService) GetPaymentStatus(paymentID string) (*entities.Payment, error) {
	return s.paymentRepo.FindByID(paymentID)
}

// RefundPayment processes a refund
func (s *PaymentService) RefundPayment(req *entities.RefundRequest) (*entities.RefundResponse, error) {
	payment, err := s.paymentRepo.FindByID(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	if !payment.CanRefund() {
		return nil, fmt.Errorf("payment cannot be refunded")
	}

	refundAmount := req.Amount
	if refundAmount <= 0 || refundAmount > payment.Amount {
		refundAmount = payment.CalculateRefundAmount()
	}

	// Mock refund processing
	refundID := s.generateRefundID()

	response := &entities.RefundResponse{
		RefundID:    refundID,
		PaymentID:   req.PaymentID,
		Amount:      refundAmount,
		Status:      "completed",
		ProcessedAt: time.Now(),
		Message:     "Refund processed successfully",
	}

	// In real implementation, would call payment gateway refund API
	// For now, just return mock response

	return response, nil
}

// GetUserPayments retrieves all payments for a user
func (s *PaymentService) GetUserPayments(userID string) ([]*entities.Payment, error) {
	return s.paymentRepo.FindByUserID(userID)
}

// GetPaymentByReservation retrieves payment for a specific reservation
func (s *PaymentService) GetPaymentByReservation(reservationID string) (*entities.Payment, error) {
	return s.paymentRepo.FindByReservationID(reservationID)
}

// Helper methods

func (s *PaymentService) generatePaymentID() string {
	return fmt.Sprintf("pay_%d_%s", time.Now().Unix(), s.generateRandomString(8))
}

func (s *PaymentService) generateTransactionID() string {
	return fmt.Sprintf("txn_%d_%s", time.Now().Unix(), s.generateRandomString(12))
}

func (s *PaymentService) generateRefundID() string {
	return fmt.Sprintf("ref_%d_%s", time.Now().Unix(), s.generateRandomString(8))
}

func (s *PaymentService) generateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
