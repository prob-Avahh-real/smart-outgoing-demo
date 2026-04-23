package redis

import (
	"context"
	"fmt"
	"time"
)

// RedisClientInterface defines the interface for Redis operations
// This allows us to work without the actual Redis dependency
type RedisClientInterface interface {
	Close() error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, target interface{}) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
}

// MockRedisClient implements RedisClientInterface for testing/fallback
type MockRedisClient struct {
	data map[string]interface{}
}

// NewMockRedisClient creates a new mock Redis client
func NewMockRedisClient() RedisClientInterface {
	return &MockRedisClient{
		data: make(map[string]interface{}),
	}
}

// Close closes the mock client
func (m *MockRedisClient) Close() error {
	m.data = nil
	return nil
}

// Set stores a value
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if m.data == nil {
		return ErrClientClosed
	}
	m.data[key] = value
	return nil
}

// Get retrieves a value
func (m *MockRedisClient) Get(ctx context.Context, key string, target interface{}) error {
	if m.data == nil {
		return ErrClientClosed
	}

	value, exists := m.data[key]
	if !exists {
		return ErrKeyNotFound
	}

	// Simple assignment for mock - in real implementation would use JSON unmarshaling
	switch v := target.(type) {
	case *string:
		*v = value.(string)
	case *interface{}:
		*v = value
	}

	return nil
}

// Delete removes keys
func (m *MockRedisClient) Delete(ctx context.Context, keys ...string) error {
	if m.data == nil {
		return ErrClientClosed
	}

	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}

// Exists checks if a key exists
func (m *MockRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	if m.data == nil {
		return false, ErrClientClosed
	}

	_, exists := m.data[key]
	return exists, nil
}

// Keys returns keys matching pattern (mock implementation)
func (m *MockRedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	if m.data == nil {
		return nil, ErrClientClosed
	}

	var keys []string
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys, nil
}

// Scan scans for keys (mock implementation)
func (m *MockRedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	if m.data == nil {
		return nil, 0, ErrClientClosed
	}

	keys, _ := m.Keys(ctx, match)
	return keys, 0, nil
}

// Custom errors
var (
	ErrClientClosed = fmt.Errorf("redis client is closed")
	ErrKeyNotFound  = fmt.Errorf("key not found")
)
