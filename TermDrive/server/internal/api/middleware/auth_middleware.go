package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/dto"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model/enum"
	"github.com/AtahanPoyraz/TermDrive/server/internal/service"
)

// AuthMiddleware defines the interface for authentication-related middleware.
type AuthMiddleware interface {
	RequiredAdminRole(next http.Handler) http.Handler
	RequireAuthorize(next http.Handler) http.Handler
}

// AuthMiddlewareImpl is the concrete implementation of AuthMiddleware.
type AuthMiddlewareImpl struct {
	logger      *log.Logger
	jwtService  service.JwtService
	userService service.UserService
}

// NewAuthMiddleware initializes and returns a new instance of AuthMiddlewareImpl.
func NewAuthMiddleware(logger *log.Logger, jwtService service.JwtService, userService service.UserService) AuthMiddleware {
	return &AuthMiddlewareImpl{logger: logger, jwtService: jwtService, userService: userService}
}

// sendResponse encodes and sends the response with the provided GenericResponse to the client.
func (h *AuthMiddlewareImpl) sendResponse(response *dto.GenericResponse, w http.ResponseWriter) {
	w.WriteHeader(response.StatusCode)
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		h.logger.Printf("Failed to encode response: %v.\n", encodeErr)
	}
}

// validateAndFetchUser validates the token and fetches the associated user from the database.
// It returns the user data if the token is valid or an error if validation fails.
func (m *AuthMiddlewareImpl) validateAndFetchUser(headerValue string) (model.UserModel, error) {
	var token string
	_, err := fmt.Sscanf(strings.Trim(headerValue, " "), "Bearer %s", &token)
	if err != nil {
		m.logger.Printf("Invalid token format: %v\n", err)
		return model.UserModel{}, errors.New("invalid authorization header format")
	}

	claims, err := m.jwtService.VerifyToken(token)
	if err != nil {
		m.logger.Printf("Token verification failed: %v\n", err)
		return model.UserModel{}, err
	}

	user, err := m.userService.FetchUserById(&dto.FetchUserByIdRequest{UserId: claims.UserId})
	if err != nil {
		m.logger.Printf("Error fetching user with ID %s: %v\n", claims.UserId, err)
		return model.UserModel{}, fmt.Errorf("error fetching user: %v", err)
	}

	if user.Username == "" {
		m.logger.Printf("User not found for ID: %s\n", claims.UserId)
		return model.UserModel{}, errors.New("user not found")
	}

	return user, nil
}

// RequiredAdminRole ensures that the user has an admin role before allowing access.
// If the user doesn't have the admin role, a forbidden response is sent.
func (m *AuthMiddlewareImpl) RequiredAdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := m.validateAndFetchUser(r.Header.Get("Authorization"))
		if err != nil {
			m.logger.Printf("Error validating header: %v.\n", err)
			m.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Error validating header",
			}, w)
			return
		}

		if user.Role != enum.ROLE_ADMIN {
			m.logger.Printf("Access denied: User (%s) does not have admin privileges.\n", user.Username)
			m.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusForbidden,
				Message:    "Access denied: Insufficient role privileges",
			}, w)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), config.UserContextKey, &user)))
	})
}

// RequireAuthorize ensures the user is authorized to access the route.
// It checks the validity of the token and whether the user's account is active.
func (m *AuthMiddlewareImpl) RequireAuthorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := m.validateAndFetchUser(r.Header.Get("Authorization"))
		if err != nil {
			m.logger.Printf("Error validating header: %v.\n", err)
			m.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    "Error validating header",
			}, w)
			return
		}

		if !user.IsActive {
			m.logger.Printf("Access denied: User (%s) account is inactive.\n", user.Username)
			m.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusForbidden,
				Message:    "User account is inactive",
			}, w)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), config.UserContextKey, &user)))
	})
}
