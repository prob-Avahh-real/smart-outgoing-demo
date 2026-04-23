package amap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RestClient provides access to AMap REST APIs
type RestClient struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewRestClient creates a new AMap REST API client
func NewRestClient(apiKey string) *RestClient {
	return &RestClient{
		APIKey:  apiKey,
		BaseURL: "https://restapi.amap.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GeocodingResponse represents the response from geocoding API
type GeocodingResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	Count     string `json:"count"`
	Geocodes  []Geocode `json:"geocodes"`
}

// Geocode represents a geocoding result
type Geocode struct {
	FormattedAddress string  `json:"formatted_address"`
	Location         Location `json:"location"`
	Level            string   `json:"level"`
	City             string   `json:"city"`
	District         string   `json:"district"`
}

// Location represents coordinates
type Location struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// DrivingResponse represents the response from driving directions API
type DrivingResponse struct {
	Status   string      `json:"status"`
	Info     string      `json:"info"`
	InfoCode string      `json:"infocode"`
	Routes   []Route     `json:"routes"`
}

// Route represents a driving route
type Route struct {
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Distance    string    `json:"distance"`
	Duration    string    `json:"duration"`
	Steps       []Step    `json:"steps"`
}

// Step represents a step in the route
type Step struct {
	Instruction string    `json:"instruction"`
	Distance    string    `json:"distance"`
	Duration    string    `json:"duration"`
	Polyline    string    `json:"polyline"`
}

// Geocode converts address to coordinates
func (c *RestClient) Geocode(address string) (*GeocodingResponse, error) {
	params := url.Values{}
	params.Set("key", c.APIKey)
	params.Set("address", address)
	params.Set("output", "json")

	url := fmt.Sprintf("%s/v3/geocode/geo?%s", c.BaseURL, params.Encode())

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result GeocodingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// Driving gets driving directions between two points
func (c *RestClient) Driving(origin, destination string) (*DrivingResponse, error) {
	params := url.Values{}
	params.Set("key", c.APIKey)
	params.Set("origin", origin)
	params.Set("destination", destination)
	params.Set("output", "json")
	params.Set("extensions", "all")

	url := fmt.Sprintf("%s/v3/direction/driving?%s", c.BaseURL, params.Encode())

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result DrivingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// IsValid checks if the REST API key is valid
func (c *RestClient) IsValid() bool {
	if c.APIKey == "" || c.APIKey == "75cde2597f0989d6e8fca0e7f69d98de" {
		return false
	}

	// Test with a simple geocoding request
	resp, err := c.Geocode("test")
	if err != nil {
		return false
	}

	return resp.Status == "1" && resp.InfoCode != "10009"
}
