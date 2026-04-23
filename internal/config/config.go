package config

import (
	"bufio"
	"fmt"
	"os"
	"smart-outgoing-demo/pkg/security"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	Port              int
	AdminToken        string
	AMapJsKey         string
	AMapSecurityCode  string
	AMapRestKey       string
	DefaultCenter     []float64
	RateLimitRequests int
	RateLimitWindow   int
	EnableTLS         bool
	HTTPPort          int
	HTTPSPort         int
	TLSCertFile       string
	TLSKeyFile        string
	// Parking API Configuration
	ParkingAPIEnabled   bool
	ParkingAPIBaseURL   string
	ParkingAPIAppID     string
	ParkingAPIAppSecret string
	ParkingLotNo        string
	PortNo              string
	ParkingUseMock      bool
}

func Load() (*Config, error) {
	// Load .env file if exists
	if err := loadEnvFile(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg := &Config{
		Port:              getEnvInt("PORT", 8080),
		AdminToken:        getEnv("ADMIN_TOKEN", "demo_admin_token"),
		AMapJsKey:         getEnv("AMAP_JS_KEY", "45109d104b3c8d03a2c84175a7749241"),
		AMapSecurityCode:  getEnv("AMAP_SECURITY_CODE", "c552677838e5f5e71de92ce532c936bc"),
		AMapRestKey:       getEnv("AMAP_REST_KEY", "75cde2597f0989d6e8fca0e7f69d98de"),
		DefaultCenter:     []float64{114.0448, 22.6913},
		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvInt("RATE_LIMIT_WINDOW", 60),
		EnableTLS:         getEnvBool("ENABLE_TLS", false),
		HTTPPort:          getEnvInt("HTTP_PORT", 8080),
		HTTPSPort:         getEnvInt("HTTPS_PORT", 8443),
		TLSCertFile:       getEnv("TLS_CERT_FILE", "./certs/server.crt"),
		TLSKeyFile:        getEnv("TLS_KEY_FILE", "./certs/server.key"),
		// Parking API Configuration
		ParkingAPIEnabled:   getEnvBool("PARKING_API_ENABLED", true),
		ParkingAPIBaseURL:   getEnv("PARKING_API_BASE_URL", "https://api.citybrain.example.com"),
		ParkingAPIAppID:     getEnv("PARKING_API_APP_ID", ""),
		ParkingAPIAppSecret: getEnv("PARKING_API_APP_SECRET", ""),
		ParkingLotNo:        getEnv("PARKING_LOT_NO", "LOT001"),
		PortNo:              getEnv("PORT_NO", "PORT001"),
		ParkingUseMock:      getEnvBool("PARKING_USE_MOCK", true),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func loadEnvFile() error {
	file, err := os.Open(".env")
	if err != nil {
		return nil // .env file not found, skip
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}

// UpdateEnvFile dynamically updates the .env file
func UpdateEnvFile(updates map[string]string) error {
	// Read existing .env file
	var existingLines []string
	file, err := os.Open(".env")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			existingLines = append(existingLines, scanner.Text())
		}
	}

	// Create a map of existing keys
	existingKeys := make(map[string]int)
	for i, line := range existingLines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			existingKeys[strings.TrimSpace(parts[0])] = i
		}
	}

	// Update existing lines or add new ones
	for key, value := range updates {
		if lineIndex, exists := existingKeys[key]; exists {
			// Update existing line
			existingLines[lineIndex] = fmt.Sprintf("%s=%s", key, value)
			delete(existingKeys, key)
		} else {
			// Add new line
			existingLines = append(existingLines, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Write back to .env file
	file, err = os.Create(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range existingLines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	// Update environment variables immediately
	for key, value := range updates {
		os.Setenv(key, value)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		// Log warning about invalid value (in production, use proper logging)
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		// Log warning about invalid value (in production, use proper logging)
	}
	return defaultValue
}

// GetAuthConfig returns authentication configuration
func (c *Config) GetAuthConfig() *security.AuthConfig {
	environment := security.GetEnvironmentFromEnv()
	return security.NewAuthConfig(c.AdminToken, environment)
}

// Validate validates the configuration values
func (c *Config) Validate() error {
	// Validate port ranges
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid PORT: must be between 1 and 65535, got %d", c.Port)
	}
	if c.HTTPPort < 1 || c.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP_PORT: must be between 1 and 65535, got %d", c.HTTPPort)
	}
	if c.HTTPSPort < 1 || c.HTTPSPort > 65535 {
		return fmt.Errorf("invalid HTTPS_PORT: must be between 1 and 65535, got %d", c.HTTPSPort)
	}

	// Validate admin token
	if c.AdminToken == "" {
		return fmt.Errorf("ADMIN_TOKEN is required")
	}

	// Validate AMap configuration
	if c.AMapJsKey == "" {
		return fmt.Errorf("AMAP_JS_KEY is required")
	}
	if c.AMapSecurityCode == "" {
		return fmt.Errorf("AMAP_SECURITY_CODE is required")
	}

	// Validate default center coordinates
	if len(c.DefaultCenter) != 2 {
		return fmt.Errorf("default_center must have exactly 2 coordinates, got %d", len(c.DefaultCenter))
	}
	if c.DefaultCenter[0] < -180 || c.DefaultCenter[0] > 180 {
		return fmt.Errorf("invalid default_center longitude: must be between -180 and 180, got %f", c.DefaultCenter[0])
	}
	if c.DefaultCenter[1] < -90 || c.DefaultCenter[1] > 90 {
		return fmt.Errorf("invalid default_center latitude: must be between -90 and 90, got %f", c.DefaultCenter[1])
	}

	// Validate rate limiting
	if c.RateLimitRequests < 1 {
		return fmt.Errorf("RATE_LIMIT_REQUESTS must be positive, got %d", c.RateLimitRequests)
	}
	if c.RateLimitWindow < 1 {
		return fmt.Errorf("RATE_LIMIT_WINDOW must be positive, got %d", c.RateLimitWindow)
	}

	// Validate TLS configuration if enabled
	if c.EnableTLS {
		if c.TLSCertFile == "" {
			return fmt.Errorf("TLS_CERT_FILE is required when ENABLE_TLS is true")
		}
		if c.TLSKeyFile == "" {
			return fmt.Errorf("TLS_KEY_FILE is required when ENABLE_TLS is true")
		}
		// Check if cert and key files exist
		if _, err := os.Stat(c.TLSCertFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS certificate file does not exist: %s", c.TLSCertFile)
		}
		if _, err := os.Stat(c.TLSKeyFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS key file does not exist: %s", c.TLSKeyFile)
		}
	}

	return nil
}
