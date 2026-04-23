package handlers

import (
	"net/http"
	"strconv"
	"time"

	"smart-outgoing-demo/internal/config"
	"smart-outgoing-demo/internal/domain/entities"
	"smart-outgoing-demo/internal/domain/services"

	"github.com/gin-gonic/gin"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService *services.PaymentService
	config         *config.Config
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(
	paymentService *services.PaymentService,
	cfg *config.Config,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		config:         cfg,
	}
}

// ProcessPaymentRequest represents a payment processing request
type ProcessPaymentRequest struct {
	ReservationID string                 `json:"reservation_id" binding:"required"`
	Amount        float64                `json:"amount" binding:"required"`
	Currency      string                 `json:"currency" binding:"required"`
	Method        entities.PaymentMethod `json:"method" binding:"required"`
	PaymentInfo   interface{}            `json:"payment_info"`
}

// ProcessPayment processes a payment
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate payment method
	switch req.Method {
	case entities.PaymentMethodWeChat, entities.PaymentMethodAlipay, entities.PaymentMethodCreditCard:
		// Valid methods
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported payment method"})
		return
	}

	// Create payment request
	paymentReq := &entities.PaymentRequest{
		ReservationID: req.ReservationID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Method:        req.Method,
		PaymentInfo:   req.PaymentInfo,
	}

	// Process payment
	if h.paymentService != nil {
		response, err := h.paymentService.ProcessPayment(paymentReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	} else {
		// Mock implementation
		mockResponse := &entities.PaymentResponse{
			PaymentID:     "mock_pay_" + strconv.FormatInt(int64(req.Amount*100), 10),
			Status:        entities.PaymentStatusCompleted,
			TransactionID: "mock_txn_" + strconv.FormatInt(int64(req.Amount*100), 10),
			Message:       "Payment processed successfully (mock)",
		}
		c.JSON(http.StatusOK, mockResponse)
	}
}

// GetPaymentStatus retrieves payment status
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment ID required"})
		return
	}

	if h.paymentService != nil {
		payment, err := h.paymentService.GetPaymentStatus(paymentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"payment": payment})
	} else {
		// Mock implementation
		now := time.Now()
		mockPayment := &entities.Payment{
			ID:            paymentID,
			ReservationID: "mock_reservation",
			UserID:        "demo_user",
			Amount:        30.0,
			Currency:      "CNY",
			Method:        entities.PaymentMethodWeChat,
			Status:        entities.PaymentStatusCompleted,
			TransactionID: "mock_txn_123",
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		c.JSON(http.StatusOK, gin.H{"payment": mockPayment})
	}
}

// RefundPaymentRequest represents a refund request
type RefundPaymentRequest struct {
	PaymentID string  `json:"payment_id" binding:"required"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

// RefundPayment processes a refund
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	var req RefundPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create refund request
	refundReq := &entities.RefundRequest{
		PaymentID: req.PaymentID,
		Amount:    req.Amount,
		Reason:    req.Reason,
	}

	// Process refund
	if h.paymentService != nil {
		response, err := h.paymentService.RefundPayment(refundReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	} else {
		// Mock implementation
		mockResponse := &entities.RefundResponse{
			RefundID:    "mock_refund_" + req.PaymentID,
			PaymentID:   req.PaymentID,
			Amount:      req.Amount,
			Status:      "completed",
			ProcessedAt: time.Now(),
			Message:     "Refund processed successfully (mock)",
		}
		c.JSON(http.StatusOK, mockResponse)
	}
}

// GetUserPayments retrieves all payments for a user
func (h *PaymentHandler) GetUserPayments(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = c.GetHeader("x-user-id")
	}
	if userID == "" {
		userID = "demo_user" // Default for demo
	}

	if h.paymentService != nil {
		payments, err := h.paymentService.GetUserPayments(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"payments": payments,
			"count":    len(payments),
		})
	} else {
		// Mock implementation
		mockPayments := []*entities.Payment{
			{
				ID:            "mock_pay_1",
				ReservationID: "res_123",
				UserID:        userID,
				Amount:        30.0,
				Currency:      "CNY",
				Method:        entities.PaymentMethodWeChat,
				Status:        entities.PaymentStatusCompleted,
				TransactionID: "txn_wechat_123",
				CreatedAt:     time.Now().Add(-1 * time.Hour),
				UpdatedAt:     time.Now().Add(-1 * time.Hour),
			},
			{
				ID:            "mock_pay_2",
				ReservationID: "res_456",
				UserID:        userID,
				Amount:        45.0,
				Currency:      "CNY",
				Method:        entities.PaymentMethodAlipay,
				Status:        entities.PaymentStatusPending,
				TransactionID: "",
				CreatedAt:     time.Now().Add(-30 * time.Minute),
				UpdatedAt:     time.Now().Add(-30 * time.Minute),
			},
		}
		c.JSON(http.StatusOK, gin.H{
			"payments": mockPayments,
			"count":    len(mockPayments),
		})
	}
}

// GetPaymentMethods returns available payment methods
func (h *PaymentHandler) GetPaymentMethods(c *gin.Context) {
	methods := []gin.H{
		{
			"id":          string(entities.PaymentMethodWeChat),
			"name":        "WeChat Pay",
			"description": "Pay with WeChat",
			"icon":        "wechat",
			"enabled":     true,
		},
		{
			"id":          string(entities.PaymentMethodAlipay),
			"name":        "Alipay",
			"description": "Pay with Alipay",
			"icon":        "alipay",
			"enabled":     true,
		},
		{
			"id":          string(entities.PaymentMethodCreditCard),
			"name":        "Credit Card",
			"description": "Pay with credit/debit card",
			"icon":        "credit-card",
			"enabled":     true,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"methods": methods,
		"count":   len(methods),
	})
}

// GetPaymentStats returns payment statistics
func (h *PaymentHandler) GetPaymentStats(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = c.GetHeader("x-user-id")
	}
	if userID == "" {
		userID = "demo_user"
	}

	// Mock statistics
	stats := gin.H{
		"total_payments":      15,
		"total_amount":        675.50,
		"successful_payments": 12,
		"failed_payments":     2,
		"pending_payments":    1,
		"refunded_amount":     45.00,
		"average_payment":     45.03,
		"most_used_method":    "wechat",
		"last_payment_date":   "2026-04-21T10:30:00Z",
	}

	c.JSON(http.StatusOK, stats)
}
