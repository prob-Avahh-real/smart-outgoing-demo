package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"smart-outgoing-demo/internal/algorithm"
	"smart-outgoing-demo/internal/config"
	"smart-outgoing-demo/internal/store"
	"smart-outgoing-demo/internal/websocket"

	"github.com/gin-gonic/gin"
)

// UpdateConfig dynamically updates .env file configuration
func UpdateConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updates map[string]string
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := config.UpdateEnvFile(updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

// GetConfig returns the application's configuration to the client as JSON.
// It exposes only the necessary configuration for frontend functionality.
// Sensitive information like admin_token are not exposed for security reasons.
func GetConfig(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"amap_js_key":        cfg.AMapJsKey,        // Required for AMap JS API
			"amap_security_code": cfg.AMapSecurityCode, // Required for AMap security
			"default_center":     cfg.DefaultCenter,    // Public map center coordinates
			// Security: admin_token intentionally not exposed
		})
	}
}

// GetVehicles returns all vehicles
func GetVehicles(vehicleStore *store.VehicleStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if vehicleStore == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Vehicle store not initialized"})
			return
		}
		vehicles := vehicleStore.GetAll()
		c.JSON(http.StatusOK, vehicles)
	}
}

// CreateVehicle creates a new vehicle
func CreateVehicle(vehicleStore *store.VehicleStore, hub *websocket.Hub, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vehicle struct {
			Name     string  `json:"name"`
			StartLng float64 `json:"start_lng"`
			StartLat float64 `json:"start_lat"`
			StartAlt float64 `json:"start_alt,omitempty"` // Altitude in meters
		}

		if err := c.ShouldBindJSON(&vehicle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate vehicle data
		if vehicle.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle name is required"})
			return
		}
		if vehicle.StartLng < -180 || vehicle.StartLng > 180 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude: must be between -180 and 180"})
			return
		}
		if vehicle.StartLat < -90 || vehicle.StartLat > 90 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude: must be between -90 and 90"})
			return
		}

		newVehicle := &store.Vehicle{
			ID:       generateID(),
			Name:     vehicle.Name,
			StartLng: vehicle.StartLng,
			StartLat: vehicle.StartLat,
			StartAlt: vehicle.StartAlt,
		}

		vehicleStore.Create(newVehicle)
		hub.Broadcast()

		c.JSON(http.StatusCreated, newVehicle)
	}
}

// SetDestination sets the destination for a vehicle
func SetDestination(vehicleStore *store.VehicleStore, hub *websocket.Hub, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var destination struct {
			EndLng float64 `json:"end_lng"`
			EndLat float64 `json:"end_lat"`
			EndAlt float64 `json:"end_alt,omitempty"` // Altitude in meters
		}

		if err := c.ShouldBindJSON(&destination); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate destination coordinates
		if destination.EndLng < -180 || destination.EndLng > 180 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude: must be between -180 and 180"})
			return
		}
		if destination.EndLat < -90 || destination.EndLat > 90 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude: must be between -90 and 90"})
			return
		}

		updated := vehicleStore.Update(id, func(v *store.Vehicle) {
			v.EndLng = destination.EndLng
			v.EndLat = destination.EndLat
			v.EndAlt = destination.EndAlt
		})

		if !updated {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		hub.Broadcast()
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

// DeleteVehicle deletes a vehicle
func DeleteVehicle(vehicleStore *store.VehicleStore, hub *websocket.Hub, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		deleted := vehicleStore.Delete(id)
		if !deleted {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}

		hub.Broadcast()
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

// PlanRoute plans a route using the algorithm
func PlanRoute(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			From struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"from"`
			To struct {
				Lng float64 `json:"lng"`
				Lat float64 `json:"lat"`
			} `json:"to"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate coordinates
		if request.From.Lng < -180 || request.From.Lng > 180 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from longitude: must be between -180 and 180"})
			return
		}
		if request.From.Lat < -90 || request.From.Lat > 90 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from latitude: must be between -90 and 90"})
			return
		}
		if request.To.Lng < -180 || request.To.Lng > 180 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to longitude: must be between -180 and 180"})
			return
		}
		if request.To.Lat < -90 || request.To.Lat > 90 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to latitude: must be between -90 and 90"})
			return
		}

		// Create a simple graph for routing
		graph := algorithm.NewGraph()

		// Add nodes
		fromID := "from"
		toID := "to"
		graph.AddNode(fromID, request.From.Lng, request.From.Lat)
		graph.AddNode(toID, request.To.Lng, request.To.Lat)

		// Add edge
		distance := algorithm.CalculateDistance(
			request.From.Lng, request.From.Lat,
			request.To.Lng, request.To.Lat,
		)
		graph.AddEdge(fromID, toID, distance)

		// Create scheduler
		scheduler := algorithm.NewScheduler(graph)

		// Find shortest path
		path, totalDist := scheduler.Dijkstra(fromID, toID)

		// Generate route points
		route := buildRoutePoints(request.From.Lng, request.From.Lat, request.To.Lng, request.To.Lat)

		c.JSON(http.StatusOK, gin.H{
			"path":     path,
			"distance": totalDist,
			"route":    route,
		})
	}
}

// ScheduleTasks assigns tasks to vehicles
func ScheduleTasks(vehicleStore *store.VehicleStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if vehicleStore == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Vehicle store not initialized"})
			return
		}

		var request struct {
			Vehicles []string `json:"vehicles"`
			Tasks    []struct {
				ID       string  `json:"id"`
				Priority int     `json:"priority"`
				Lng      float64 `json:"lng"`
				Lat      float64 `json:"lat"`
				Duration int64   `json:"duration"`
			} `json:"tasks"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(request.Vehicles) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one vehicle is required"})
			return
		}
		if len(request.Tasks) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one task is required"})
			return
		}

		// Validate task coordinates
		for i, task := range request.Tasks {
			if task.ID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Task %d: ID is required", i)})
				return
			}
			if task.Lng < -180 || task.Lng > 180 {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Task %d: Invalid longitude: must be between -180 and 180", i)})
				return
			}
			if task.Lat < -90 || task.Lat > 90 {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Task %d: Invalid latitude: must be between -90 and 90", i)})
				return
			}
			if task.Duration <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Task %d: Duration must be positive", i)})
				return
			}
		}

		// Create graph
		graph := algorithm.NewGraph()

		// Add vehicle nodes
		for _, vehicleID := range request.Vehicles {
			if vehicle, exists := vehicleStore.Get(vehicleID); exists {
				graph.AddNode(vehicleID, vehicle.StartLng, vehicle.StartLat)
			}
		}

		// Add task nodes
		for _, task := range request.Tasks {
			graph.AddNode(task.ID, task.Lng, task.Lat)
		}

		// Create scheduler
		scheduler := algorithm.NewScheduler(graph)

		// Convert tasks
		tasks := make([]*algorithm.Task, len(request.Tasks))
		for i, task := range request.Tasks {
			tasks[i] = &algorithm.Task{
				ID:       task.ID,
				Priority: task.Priority,
				Location: algorithm.Node{
					ID:  task.ID,
					Lng: task.Lng,
					Lat: task.Lat,
				},
				Duration: task.Duration,
			}
		}

		// Schedule tasks
		routes := scheduler.ScheduleTasks(request.Vehicles, tasks)

		c.JSON(http.StatusOK, gin.H{
			"routes": routes,
		})
	}
}

// ImportVehiclesFromCSV imports vehicles from CSV file
func ImportVehiclesFromCSV(vehicleStore *store.VehicleStore, hub *websocket.Hub, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		reader := csv.NewReader(src)
		records, err := reader.ReadAll()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV: " + err.Error()})
			return
		}

		if len(records) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSV file is empty or missing header"})
			return
		}

		header := records[0]
		nameIndex := findColumnIndex(header, "name")
		lngIndex := findColumnIndex(header, "lng")
		latIndex := findColumnIndex(header, "lat")
		altIndex := findColumnIndex(header, "alt")

		if nameIndex == -1 || lngIndex == -1 || latIndex == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CSV must contain 'name', 'lng', and 'lat' columns"})
			return
		}

		successCount := 0
		errorCount := 0
		errors := []string{}

		for i, record := range records[1:] {
			if len(record) <= max(nameIndex, lngIndex, latIndex) {
				errorCount++
				errors = append(errors, fmt.Sprintf("Row %d: insufficient columns", i+2))
				continue
			}

			name := strings.TrimSpace(record[nameIndex])
			lng, err := strconv.ParseFloat(record[lngIndex], 64)
			if err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Row %d: invalid lng value", i+2))
				continue
			}

			lat, err := strconv.ParseFloat(record[latIndex], 64)
			if err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Row %d: invalid lat value", i+2))
				continue
			}

			alt := 0.0
			if altIndex != -1 && len(record) > altIndex {
				alt, _ = strconv.ParseFloat(record[altIndex], 64)
			}

			vehicle := &store.Vehicle{
				ID:        generateID(),
				Name:      name,
				StartLng:  lng,
				StartLat:  lat,
				StartAlt:  alt,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			vehicleStore.Create(vehicle)
			successCount++
		}

		hub.Broadcast()

		response := gin.H{
			"success_count": successCount,
			"error_count":   errorCount,
			"status":        "completed",
		}

		if errorCount > 0 {
			response["errors"] = errors
		}

		c.JSON(http.StatusOK, response)
	}
}

func findColumnIndex(header []string, columnName string) int {
	for i, col := range header {
		if strings.EqualFold(strings.TrimSpace(col), columnName) {
			return i
		}
	}
	return -1
}

// GetCacheStats returns cache statistics
func GetCacheStats(vehicleStore *store.VehicleStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if vehicleStore == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Vehicle store not initialized"})
			return
		}
		stats := vehicleStore.GetCacheStats()
		c.JSON(http.StatusOK, gin.H{
			"cache_stats": stats,
			"timestamp":   time.Now(),
		})
	}
}

// CleanupCache removes expired cache entries
func CleanupCache(vehicleStore *store.VehicleStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if vehicleStore == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Vehicle store not initialized"})
			return
		}
		vehicleStore.CleanupCache()
		c.JSON(http.StatusOK, gin.H{
			"message":   "Cache cleanup completed",
			"timestamp": time.Now(),
		})
	}
}

// max returns the maximum of integers
func max(a, b, c int) int {
	if a >= b && a >= c {
		return a
	}
	if b >= a && b >= c {
		return b
	}
	return c
}

func generateID() string {
	return strconv.FormatInt(int64(float64(1000000)*float64(1000000)), 10)
}

func buildRoutePoints(fromLng, fromLat, toLng, toLat float64) [][]float64 {
	// Simplified route generation - in production, use AMap API
	return [][]float64{
		{fromLng, fromLat},
		{toLng, toLat},
	}
}
