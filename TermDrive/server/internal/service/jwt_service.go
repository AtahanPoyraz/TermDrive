package service

import (
	"errors"
	"time"

	"github.com/AtahanPoyraz/TermDrive/server/internal/model/enum"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthClaims represents the structure of the claims in the JWT token.
// It includes user information like UserID, Role, and IsActive status.
type AuthClaims struct {
	jwt.RegisteredClaims
	UserId   uuid.UUID `json:"userId"`
	UserRole enum.Role `json:"userRole"`
	IsActive bool      `json:"isActive"`
}

// JwtService defines the interface for creating and verifying JWT tokens.
type JwtService interface {
	CreateToken(userId uuid.UUID, role enum.Role, isActive bool) (string, error)
	VerifyToken(tokenString string) (AuthClaims, error)
}

// JwtServiceImpl is the concrete implementation of the JwtService interface.
type JwtServiceImpl struct {
	secretKey []byte
}

// NewJwtService initializes and returns a new instance of JwtServiceImpl with the given secret key.
func NewJwtService(secretKey string) JwtService {
	return &JwtServiceImpl{
		secretKey: []byte(secretKey),
	}
}

// CreateToken generates a JWT token for the provided user details, including user ID, role, and account status.
// The token will be valid for 6 hours from the time of creation.
func (s *JwtServiceImpl) CreateToken(userId uuid.UUID, role enum.Role, isActive bool) (string, error) {
	claims := AuthClaims{
		UserId:   userId,
		UserRole: role,
		IsActive: isActive,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// VerifyToken validates a given JWT token string.
// It parses the token and verifies its signature using the secret key.
// If the token is valid, it returns the decoded claims; otherwise, it returns an error.
func (s *JwtServiceImpl) VerifyToken(tokenString string) (AuthClaims, error) {
	var claims AuthClaims

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		var errMsg string
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			errMsg = "Token has expired"
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			errMsg = "Invalid token signature"
		default:
			errMsg = "Invalid token"
		}
		return AuthClaims{}, errors.New(errMsg)
	}

	return claims, nil
}
