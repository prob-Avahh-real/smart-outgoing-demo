package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"smart-outgoing-demo/internal/config"
	"smart-outgoing-demo/internal/store"
	"smart-outgoing-demo/internal/websocket"

	"github.com/gin-gonic/gin"
)

func TestGetConfig(t *testing.T) {
	cfg := &config.Config{
		AMapJsKey:        "test_key",
		AMapSecurityCode: "test_code",
		DefaultCenter:    []float64{114.0, 22.0},
		AdminToken:       "test_token",
	}

	handler := GetConfig(cfg)

	req := httptest.NewRequest("GET", "/api/config", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	if response["amap_js_key"] != "test_key" {
		t.Errorf("Expected amap_js_key 'test_key', got %v", response["amap_js_key"])
	}
}

func TestGetVehicles(t *testing.T) {
	vehicleStore := store.NewVehicleStore()

	// Add test vehicle
	vehicle := &store.Vehicle{
		ID:       "test-1",
		Name:     "Test Vehicle",
		StartLng: 114.0,
		StartLat: 22.0,
	}
	vehicleStore.Create(vehicle)

	handler := GetVehicles(vehicleStore)

	req := httptest.NewRequest("GET", "/api/vehicles", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []store.Vehicle
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 vehicle, got %d", len(response))
	}

	if response[0].Name != "Test Vehicle" {
		t.Errorf("Expected vehicle name 'Test Vehicle', got '%s'", response[0].Name)
	}
}

func TestCreateVehicle(t *testing.T) {
	vehicleStore := store.NewVehicleStore()
	hub := websocket.NewHub(vehicleStore)
	cfg := &config.Config{AdminToken: "test_token"}

	handler := CreateVehicle(vehicleStore, hub, cfg)

	requestBody := map[string]interface{}{
		"name":      "New Vehicle",
		"start_lng": 114.5,
		"start_lat": 22.5,
		"start_alt": 10.0,
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/vehicles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response store.Vehicle
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	if response.Name != "New Vehicle" {
		t.Errorf("Expected vehicle name 'New Vehicle', got '%s'", response.Name)
	}

	if response.StartLng != 114.5 {
		t.Errorf("Expected start_lng 114.5, got %f", response.StartLng)
	}
}

func TestImportVehiclesFromCSV(t *testing.T) {
	vehicleStore := store.NewVehicleStore()
	hub := websocket.NewHub(vehicleStore)
	cfg := &config.Config{AdminToken: "test_token"}

	handler := ImportVehiclesFromCSV(vehicleStore, hub, cfg)

	// Create CSV content
	csvContent := `name,lng,lat,alt
Vehicle1,114.1,22.1,100
Vehicle2,114.2,22.2,200`

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "test.csv")
	part.Write([]byte(csvContent))
	writer.Close()

	req := httptest.NewRequest("POST", "/api/vehicles/import", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	if response["success_count"] != float64(2) {
		t.Errorf("Expected success_count 2, got %v", response["success_count"])
	}

	if response["error_count"] != float64(0) {
		t.Errorf("Expected error_count 0, got %v", response["error_count"])
	}
}

func TestPlanRoute(t *testing.T) {
	cfg := &config.Config{AdminToken: "test_token"}
	handler := PlanRoute(cfg)

	requestBody := map[string]interface{}{
		"from": map[string]float64{"lng": 114.0, "lat": 22.0},
		"to":   map[string]float64{"lng": 115.0, "lat": 23.0},
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/api/algorithm/plan", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response JSON")
	}

	if response["path"] == nil {
		t.Error("Expected path in response")
	}

	if response["distance"] == nil {
		t.Error("Expected distance in response")
	}

	if response["route"] == nil {
		t.Error("Expected route in response")
	}
}
