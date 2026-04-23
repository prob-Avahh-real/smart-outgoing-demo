package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"smart-outgoing-demo/internal/amap"
	"smart-outgoing-demo/internal/config"

	"github.com/gin-gonic/gin"
)

// AMapGeocode handles geocoding requests
func AMapGeocode(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
			return
		}

		client := amap.NewRestClient(cfg.AMapRestKey)

		// Check if REST API key is configured
		if !client.IsValid() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "AMap REST API key not configured or invalid",
				"info":  "Please set AMAP_REST_KEY environment variable with a valid REST API key",
			})
			return
		}

		result, err := client.Geocode(address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.Status != "1" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Geocoding failed",
				"info":      result.Info,
				"info_code": result.InfoCode,
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// AMapDriving handles driving directions requests
func AMapDriving(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Query("origin")
		destination := c.Query("destination")

		if origin == "" || destination == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "origin and destination parameters are required"})
			return
		}

		client := amap.NewRestClient(cfg.AMapRestKey)

		// Check if REST API key is configured
		if !client.IsValid() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "AMap REST API key not configured or invalid",
				"info":  "Please set AMAP_REST_KEY environment variable with a valid REST API key",
			})
			return
		}

		result, err := client.Driving(origin, destination)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.Status != "1" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Driving directions failed",
				"info":      result.Info,
				"info_code": result.InfoCode,
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// AMapValidate validates REST API key
func AMapValidate(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := amap.NewRestClient(cfg.AMapRestKey)

		isValid := client.IsValid()

		response := gin.H{
			"rest_api_key_configured": cfg.AMapRestKey != "" && cfg.AMapRestKey != "75cde2597f0989d6e8fca0e7f69d98de",
			"rest_api_key_valid":      isValid,
		}

		if !isValid {
			response["error"] = "REST API key validation failed"
			response["suggestion"] = "Obtain a REST API key from AMap console and set AMAP_REST_KEY"
		}

		c.JSON(http.StatusOK, response)
	}
}

// ParseCoordinates parses coordinate string to float64 array
func ParseCoordinates(coordStr string) ([]float64, error) {
	parts := strings.Split(coordStr, ",")
	if len(parts) != 2 {
		return nil, strconv.ErrSyntax
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return nil, err
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return nil, err
	}

	return []float64{lng, lat}, nil
}
