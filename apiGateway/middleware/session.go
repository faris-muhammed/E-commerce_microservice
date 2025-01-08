package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var secretKey = []byte("SECRETKEY") // Replace with a securely generated key

// GenerateSessionToken creates a signed session token containing userID and expiry time
func GenerateSessionToken(userID string) (string, error) {
	// Create a token: "<userID>:<expiryTime>"
	expiryTime := time.Now().Add(1 * time.Hour).Unix() // Token valid for 1 hour
	tokenData := fmt.Sprintf("%s:%d", userID, expiryTime)

	// Generate HMAC signature
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(tokenData))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Combine token data and signature
	sessionToken := fmt.Sprintf("%s:%s", tokenData, signature)
	return sessionToken, nil
}

// GetUserIDBySessionToken validates the token and extracts the userID
func GetUserIDBySessionToken(sessionToken string) (string, error) {
	// Split the token into its components: "<userID>:<expiryTime>:<signature>"
	parts := strings.Split(sessionToken, ":")
	if len(parts) != 3 {
		return "", errors.New("invalid session token format")
	}

	userID := parts[0]
	expiryTime := parts[1]
	signature := parts[2]

	// Validate the expiry time
	expiryTimeInt, err := strconv.ParseInt(expiryTime, 10, 64)
	if err != nil {
		return "", errors.New("invalid expiry time in session token")
	}
	if time.Now().Unix() > expiryTimeInt {
		return "", errors.New("session token expired")
	}

	// Recompute the signature and compare
	tokenData := fmt.Sprintf("%s:%s", userID, expiryTime)
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(tokenData))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", errors.New("invalid session token signature")
	}

	// Return the user ID
	return userID, nil
}
