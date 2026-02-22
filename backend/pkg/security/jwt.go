package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT generates a JWT token with the given claims and secret
func GenerateJWT(claims map[string]interface{}, secret string) (string, error) {
	// Create a new token with the specified claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString, secret string) (map[string]interface{}, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token has expired")
		}
	}

	// Convert jwt.MapClaims to map[string]interface{}
	result := make(map[string]interface{})
	for key, value := range claims {
		result[key] = value
	}

	return result, nil
}

// ExtractClaims extracts claims from a JWT token without validation (for debugging)
func ExtractClaims(tokenString string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	result := make(map[string]interface{})
	for key, value := range claims {
		result[key] = value
	}

	return result, nil
}

// GenerateAccessToken generates an access token with standard claims
func GenerateAccessToken(userID, username, role string, permissions []string, secret string, duration time.Duration) (string, error) {
	claims := map[string]interface{}{
		"user_id":     userID,
		"username":    username,
		"role":        role,
		"permissions": permissions,
		"exp":         time.Now().Add(duration).Unix(),
		"iat":         time.Now().Unix(),
		"type":        "access",
	}

	return GenerateJWT(claims, secret)
}

// GenerateRefreshToken generates a refresh token with minimal claims
func GenerateRefreshToken(userID string, secret string, duration time.Duration) (string, error) {
	claims := map[string]interface{}{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	return GenerateJWT(claims, secret)
}

