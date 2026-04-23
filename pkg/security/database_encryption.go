package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// DatabaseEncryption provides encryption for database storage
type DatabaseEncryption struct {
	encryptionKey []byte
}

// NewDatabaseEncryption creates a new database encryption instance
func NewDatabaseEncryption(key []byte) *DatabaseEncryption {
	if len(key) != 32 {
		// Generate a key if not provided or wrong length
		key = make([]byte, 32)
		rand.Read(key)
	}
	return &DatabaseEncryption{
		encryptionKey: key,
	}
}

// EncryptField encrypts a single database field
func (de *DatabaseEncryption) EncryptField(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(de.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptField decrypts a single database field
func (de *DatabaseEncryption) DecryptField(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(de.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherdata := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherdata, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptSensitiveData encrypts a map of sensitive data
func (de *DatabaseEncryption) EncryptSensitiveData(data map[string]string) (map[string]string, error) {
	encrypted := make(map[string]string)
	for key, value := range data {
		encryptedValue, err := de.EncryptField(value)
		if err != nil {
			return nil, err
		}
		encrypted[key] = encryptedValue
	}
	return encrypted, nil
}

// DecryptSensitiveData decrypts a map of sensitive data
func (de *DatabaseEncryption) DecryptSensitiveData(data map[string]string) (map[string]string, error) {
	decrypted := make(map[string]string)
	for key, value := range data {
		decryptedValue, err := de.DecryptField(value)
		if err != nil {
			return nil, err
		}
		decrypted[key] = decryptedValue
	}
	return decrypted, nil
}

// IsEncrypted checks if a value appears to be encrypted (base64 encoded)
func (de *DatabaseEncryption) IsEncrypted(value string) bool {
	if value == "" {
		return false
	}
	_, err := base64.URLEncoding.DecodeString(value)
	return err == nil
}
