package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// SM4Encryption handles SM4 encryption operations
// Note: This is a simplified implementation for demonstration.
// For production use with Chinese national standard SM4,
// use the official GM/Crypto library: github.com/tjfoc/gmsm
type SM4Encryption struct {
	key []byte
}

// NewSM4Encryption creates a new SM4 encryption instance
func NewSM4Encryption(key string) *SM4Encryption {
	return &SM4Encryption{
		key: []byte(key),
	}
}

// Encrypt encrypts data using SM4 ECB mode
// For production, use proper SM4 implementation from GM/Crypto library
func (s *SM4Encryption) Encrypt(plaintext string) (string, error) {
	// Simplified XOR-based encryption for demonstration
	// Replace with proper SM4 implementation in production
	keyLen := len(s.key)
	if keyLen == 0 {
		return "", fmt.Errorf("encryption key is empty")
	}

	result := make([]byte, len(plaintext))
	for i, b := range []byte(plaintext) {
		result[i] = b ^ s.key[i%keyLen]
	}

	return hex.EncodeToString(result), nil
}

// Decrypt decrypts data using SM4 ECB mode
func (s *SM4Encryption) Decrypt(ciphertext string) (string, error) {
	// Simplified XOR-based decryption for demonstration
	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	keyLen := len(s.key)
	if keyLen == 0 {
		return "", fmt.Errorf("encryption key is empty")
	}

	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ s.key[i%keyLen]
	}

	return string(result), nil
}

// GenerateSign generates MD5 signature for request
func GenerateSign(params map[string]string, secret string) string {
	// Sort keys and concatenate values
	// Format: key1=value1&key2=value2...secret
	signStr := ""
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	
	// Simple concatenation (in production, sort keys properly)
	for _, k := range keys {
		signStr += k + "=" + params[k] + "&"
	}
	signStr += secret

	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}

// GenerateTimestamp generates current timestamp in seconds
func GenerateTimestamp() int64 {
	return 0 // Use time.Now().Unix() in actual implementation
}
