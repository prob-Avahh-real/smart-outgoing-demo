package handlers

import (
	"net/http"
	"smart-outgoing-demo/internal/abtesting"

	"github.com/gin-gonic/gin"
)

// ABTestingHandler handles A/B testing endpoints
type ABTestingHandler struct {
	manager *abtesting.ABTestingManager
}

// NewABTestingHandler creates a new A/B testing handler
func NewABTestingHandler() *ABTestingHandler {
	return &ABTestingHandler{
		manager: abtesting.NewABTestingManager(),
	}
}

// CreateExperiment creates a new A/B test experiment
func (h *ABTestingHandler) CreateExperiment(c *gin.Context) {
	var request struct {
		Name        string                    `json:"name" binding:"required"`
		Description string                    `json:"description"`
		Variants    []abtesting.Variant       `json:"variants" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	experiment, err := h.manager.CreateExperiment(request.Name, request.Description, request.Variants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, experiment)
}

// StartExperiment starts an experiment
func (h *ABTestingHandler) StartExperiment(c *gin.Context) {
	experimentID := c.Param("id")

	err := h.manager.StartExperiment(experimentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Experiment started"})
}

// AssignUser assigns a user to a variant
func (h *ABTestingHandler) AssignUser(c *gin.Context) {
	experimentID := c.Param("id")
	
	var request struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variant, err := h.manager.AssignUser(experimentID, request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"experiment_id": experimentID,
		"user_id":       request.UserID,
		"variant":       variant,
	})
}

// RecordConversion records a conversion
func (h *ABTestingHandler) RecordConversion(c *gin.Context) {
	experimentID := c.Param("id")
	
	var request struct {
		VariantID string `json:"variant_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.manager.RecordConversion(experimentID, request.VariantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversion recorded"})
}

// GetExperiment returns an experiment by ID
func (h *ABTestingHandler) GetExperiment(c *gin.Context) {
	experimentID := c.Param("id")

	experiment, err := h.manager.GetExperiment(experimentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, experiment)
}

// GetAllExperiments returns all experiments
func (h *ABTestingHandler) GetAllExperiments(c *gin.Context) {
	experiments := h.manager.GetAllExperiments()
	c.JSON(http.StatusOK, gin.H{"experiments": experiments})
}

// GetExperimentStats returns experiment statistics
func (h *ABTestingHandler) GetExperimentStats(c *gin.Context) {
	experimentID := c.Param("id")

	stats, err := h.manager.GetExperimentStats(experimentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
