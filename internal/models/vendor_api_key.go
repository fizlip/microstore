// package models contains definitions of common types and models used accross the API.
package models

import (
    "time"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "crypto/rand"
)

// APIKey allows a user to use the API without the web api.
type APIKey struct {
    Key         string
    HashedKey   string
    Username    string
    Created     int64
    Partial     string
    secretKey   []byte
}

func (a *APIKey) generateSecretKey() ([]byte, error) {
    secretKey := make([]byte, 32)
    if _, err := rand.Read(secretKey); err != nil {
        return make([]byte, 0), err
    }
    return secretKey, nil
}

func getSecretKey() ([]byte, error) {
    return []byte("123"), nil
}

func ValidateKey(apiKey, storedHash string) bool {
    // Calculate HMAC of provided key
    calculatedHash := HMACKey(apiKey)

    // Compare hashes using hmac.Equal to prevent timing attacks
    return hmac.Equal([]byte(calculatedHash), []byte(storedHash))
}

// HMAC-SHA256 with a secret key (more secure for API key hashing)
func HMACKey(apiKey string) string {
    sk, err := getSecretKey()
    if err != nil {
        return ""
    }
    h := hmac.New(sha256.New, sk)
    h.Write([]byte(apiKey))
    return hex.EncodeToString(h.Sum(nil))
}

// GenerateAPIKey generates a new random API key
func NewAPIKey(username string) (APIKey, error) {
    // Generate random bytes for the API key
    keyBytes := make([]byte, 32)
    if _, err := rand.Read(keyBytes); err != nil {
        return APIKey{}, err
    }

    // Encode the key in base64
    apiKey := base64.URLEncoding.EncodeToString(keyBytes)

    // Hash the key for storage
    hashedKey := HMACKey(apiKey)

    return APIKey{
        Key: apiKey, 
        HashedKey: hashedKey, 
        Username: username, 
        Created: time.Now().UnixMilli(),
        Partial: "123",
    }, nil
}
