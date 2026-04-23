package parking

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"smart-outgoing-demo/internal/crypto"
)

// CityBrainClient handles communication with City Brain Parking API
type CityBrainClient struct {
	config     *CityBrainAPIConfig
	httpClient *http.Client
	encryption *crypto.SM4Encryption
}

// NewCityBrainClient creates a new City Brain API client
func NewCityBrainClient(config *CityBrainAPIConfig) *CityBrainClient {
	// Create HTTP client with timeout
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For testing only
	}
	
	return &CityBrainClient{
		config: config,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   config.Timeout,
		},
		encryption: crypto.NewSM4Encryption(config.AppSecret),
	}
}

// ReportVehicleEntry reports vehicle entry to parking lot
func (c *CityBrainClient) ReportVehicleEntry(req *ParkingEntryRequest) (*ParkingEntryResponse, error) {
	if c.config.UseMock {
		return c.mockVehicleEntry(req)
	}

	// Prepare request data
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Encrypt data
	encryptedData, err := c.encryption.Encrypt(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Build API request
	apiReq := map[string]string{
		"app_id":         c.config.AppID,
		"timestamp":      fmt.Sprintf("%d", time.Now().Unix()),
		"encrypted_data": encryptedData,
	}

	// Generate signature
	sign := crypto.GenerateSign(apiReq, c.config.AppSecret)
	apiReq["sign"] = sign

	// Send request
	respBody, err := c.sendRequest("/barrier/parking", apiReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp ParkingEntryResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// ReportVehicleExit reports vehicle exit from parking lot
func (c *CityBrainClient) ReportVehicleExit(req *ParkingExitRequest) (*ParkingExitResponse, error) {
	if c.config.UseMock {
		return c.mockVehicleExit(req)
	}

	// Similar implementation to ReportVehicleEntry
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	encryptedData, err := c.encryption.Encrypt(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	apiReq := map[string]string{
		"app_id":         c.config.AppID,
		"timestamp":      fmt.Sprintf("%d", time.Now().Unix()),
		"encrypted_data": encryptedData,
	}

	sign := crypto.GenerateSign(apiReq, c.config.AppSecret)
	apiReq["sign"] = sign

	respBody, err := c.sendRequest("/barrier/away", apiReq)
	if err != nil {
		return nil, err
	}

	var resp ParkingExitResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// UploadSnapshot uploads vehicle snapshot image
func (c *CityBrainClient) UploadSnapshot(req *SnapshotUploadRequest) (*SnapshotUploadResponse, error) {
	if c.config.UseMock {
		return c.mockSnapshotUpload(req)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	encryptedData, err := c.encryption.Encrypt(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	apiReq := map[string]string{
		"app_id":         c.config.AppID,
		"timestamp":      fmt.Sprintf("%d", time.Now().Unix()),
		"encrypted_data": encryptedData,
	}

	sign := crypto.GenerateSign(apiReq, c.config.AppSecret)
	apiReq["sign"] = sign

	respBody, err := c.sendRequest("/barrier/uploadSnapshot", apiReq)
	if err != nil {
		return nil, err
	}

	var resp SnapshotUploadResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// SendHeartbeat sends parking lot heartbeat
func (c *CityBrainClient) SendHeartbeat(req *HeartbeatRequest) (*HeartbeatResponse, error) {
	if c.config.UseMock {
		return c.mockHeartbeat(req)
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	encryptedData, err := c.encryption.Encrypt(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	apiReq := map[string]string{
		"app_id":         c.config.AppID,
		"timestamp":      fmt.Sprintf("%d", time.Now().Unix()),
		"encrypted_data": encryptedData,
	}

	sign := crypto.GenerateSign(apiReq, c.config.AppSecret)
	apiReq["sign"] = sign

	respBody, err := c.sendRequest("/barrier/heartbeat", apiReq)
	if err != nil {
		return nil, err
	}

	var resp HeartbeatResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// sendRequest sends HTTP request to City Brain API
func (c *CityBrainClient) sendRequest(endpoint string, data map[string]string) ([]byte, error) {
	url := c.config.BaseURL + endpoint
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Mock implementations for testing
func (c *CityBrainClient) mockVehicleEntry(req *ParkingEntryRequest) (*ParkingEntryResponse, error) {
	return &ParkingEntryResponse{
		Code:    0,
		Message: "success",
		Data: struct {
			OutOrderNo string `json:"out_order_no"`
			EntryTime  string `json:"entry_time"`
		}{
			OutOrderNo: req.OutOrderNo,
			EntryTime:  req.EntryTime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (c *CityBrainClient) mockVehicleExit(req *ParkingExitRequest) (*ParkingExitResponse, error) {
	return &ParkingExitResponse{
		Code:    0,
		Message: "success",
		Data: struct {
			OutOrderNo string `json:"out_order_no"`
			ExitTime   string `json:"exit_time"`
			ParkingFee float64 `json:"parking_fee"`
		}{
			OutOrderNo: req.OutOrderNo,
			ExitTime:   req.ExitTime.Format("2006-01-02 15:04:05"),
			ParkingFee: req.ParkingFee,
		},
	}, nil
}

func (c *CityBrainClient) mockSnapshotUpload(req *SnapshotUploadRequest) (*SnapshotUploadResponse, error) {
	return &SnapshotUploadResponse{
		Code:    0,
		Message: "success",
		Data: struct {
			ImageURL string `json:"image_url"`
		}{
			ImageURL: fmt.Sprintf("https://mock.example.com/snapshots/%s.jpg", req.OutOrderNo),
		},
	}, nil
}

func (c *CityBrainClient) mockHeartbeat(req *HeartbeatRequest) (*HeartbeatResponse, error) {
	return &HeartbeatResponse{
		Code:    0,
		Message: "success",
	}, nil
}
