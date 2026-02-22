package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Service interface defines JWT operations
type Service interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, role string) (string, error)
	ValidateAccessToken(tokenString string) (*Claims, error)
	ValidateRefreshToken(tokenString string) (*Claims, error)
	RefreshAccessToken(refreshToken string) (string, error)
	GetUserIDFromToken(tokenString string) (uuid.UUID, error)
	GetRoleFromToken(tokenString string) (string, error)
}

// JWTService implements the Service interface
type JWTService struct {
	secretKey        []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	issuer          string
}

// Claims represents JWT claims with all required fields
type Claims struct {
	UserID     string `json:"user_id"`     // Changed to string to match middleware expectations
	Role       string `json:"role"`
	Type       string `json:"type"`        // "access" or "refresh"
	Identifier string `json:"identifier"`  // Added missing field
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, accessTTL, refreshTTL time.Duration, issuer string) Service {
	return &JWTService{
		secretKey:       []byte(secretKey),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		issuer:         issuer,
	}
}

// GenerateAccessToken generates a new access token
func (j *JWTService) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := &Claims{
		UserID:     userID.String(),
		Role:       role,
		Type:       "access",
		Identifier: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateRefreshToken generates a new refresh token
func (j *JWTService) GenerateRefreshToken(userID uuid.UUID, role string) (string, error) {
	claims := &Claims{
		UserID:     userID.String(),
		Role:       role,
		Type:       "refresh",
		Identifier: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateAccessToken validates an access token
func (j *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return j.validateToken(tokenString, "access")
}

// ValidateRefreshToken validates a refresh token
func (j *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return j.validateToken(tokenString, "refresh")
}

// validateToken validates a token and checks its type
func (j *JWTService) validateToken(tokenString, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Type != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (j *JWTService) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", err
	}

	return j.GenerateAccessToken(userID, claims.Role)
}

// GetUserIDFromToken extracts user ID from token
func (j *JWTService) GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	claims, err := j.ValidateAccessToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(claims.UserID)
}

// GetRoleFromToken extracts role from token
func (j *JWTService) GetRoleFromToken(tokenString string) (string, error) {
	claims, err := j.ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.Role, nil
}

