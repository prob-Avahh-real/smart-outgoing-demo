package entities

import (
	"time"
)

// PaymentMethod represents different payment methods
type PaymentMethod string

const (
	PaymentMethodWeChat     PaymentMethod = "wechat"
	PaymentMethodAlipay     PaymentMethod = "alipay"
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodCash       PaymentMethod = "cash"
)

// PaymentStatus represents payment status
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

// Payment represents a payment transaction
type Payment struct {
	ID              string        `json:"id"`
	ReservationID   string        `json:"reservation_id"`
	UserID          string        `json:"user_id"`
	Amount          float64       `json:"amount"`
	Currency        string        `json:"currency"`
	Method          PaymentMethod `json:"method"`
	Status          PaymentStatus `json:"status"`
	TransactionID   string        `json:"transaction_id,omitempty"`
	GatewayResponse string        `json:"gateway_response,omitempty"`
	FailureReason   string        `json:"failure_reason,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
}

// PaymentGateway represents a payment gateway interface
type PaymentGateway interface {
	ProcessPayment(payment *Payment) error
	RefundPayment(paymentID string, amount float64) error
	GetPaymentStatus(paymentID string) (PaymentStatus, error)
}

// WeChatPayGateway represents WeChat Pay implementation
type WeChatPayGateway struct {
	AppID      string
	MerchantID string
	APIKey     string // Should be encrypted in production
	NotifyURL  string
	Encrypted  bool // Flag to indicate if API key is encrypted
}

// AlipayGateway represents Alipay implementation
type AlipayGateway struct {
	AppID      string
	MerchantID string
	PrivateKey string // Should be encrypted in production
	PublicKey  string
	NotifyURL  string
	GatewayURL string
	Encrypted  bool // Flag to indicate if private key is encrypted
}

// CreditCardGateway represents credit card payment implementation
type CreditCardGateway struct {
	MerchantID  string
	APIKey      string // Should be encrypted in production
	SecretKey   string // Should be encrypted in production
	Environment string // "sandbox" or "production"
	Encrypted   bool   // Flag to indicate if keys are encrypted
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	ReservationID string        `json:"reservation_id"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Method        PaymentMethod `json:"method"`
	PaymentInfo   interface{}   `json:"payment_info"` // Method-specific payment details
}

// WeChatPaymentInfo represents WeChat Pay specific info
type WeChatPaymentInfo struct {
	OpenID string `json:"open_id"`
}

// AlipayPaymentInfo represents Alipay specific info
type AlipayPaymentInfo struct {
	AlipayUserID string `json:"alipay_user_id"`
}

// CreditCardPaymentInfo represents credit card specific info
type CreditCardPaymentInfo struct {
	CardNumber     string `json:"card_number"`
	CardHolderName string `json:"card_holder_name"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	CVV            string `json:"cvv"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	PaymentID     string        `json:"payment_id"`
	Status        PaymentStatus `json:"status"`
	TransactionID string        `json:"transaction_id,omitempty"`
	PaymentURL    string        `json:"payment_url,omitempty"` // For redirect-based payments
	QRCode        string        `json:"qr_code,omitempty"`     // For QR code payments
	Message       string        `json:"message"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

// RefundResponse represents a refund response
type RefundResponse struct {
	RefundID    string    `json:"refund_id"`
	PaymentID   string    `json:"payment_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
	Message     string    `json:"message"`
}

// IsCompleted checks if payment is completed
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

// IsFailed checks if payment failed
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentStatusFailed
}

// CanRefund checks if payment can be refunded
func (p *Payment) CanRefund() bool {
	return p.Status == PaymentStatusCompleted && time.Since(*p.CompletedAt) <= 24*7*time.Hour // Within 7 days
}

// MarkCompleted marks payment as completed
func (p *Payment) MarkCompleted(transactionID string) {
	p.Status = PaymentStatusCompleted
	p.TransactionID = transactionID
	now := time.Now()
	p.CompletedAt = &now
	p.UpdatedAt = now
}

// MarkFailed marks payment as failed
func (p *Payment) MarkFailed(reason string) {
	p.Status = PaymentStatusFailed
	p.FailureReason = reason
	p.UpdatedAt = time.Now()
}

// CalculateRefundAmount calculates refundable amount (considering fees)
func (p *Payment) CalculateRefundAmount() float64 {
	// Simple refund calculation - in real system would consider gateway fees
	return p.Amount * 0.98 // 2% fee deduction
}
