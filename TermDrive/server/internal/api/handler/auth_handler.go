package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/dto"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/AtahanPoyraz/TermDrive/server/internal/service"
	"github.com/go-playground/validator/v10"
)

var authValidator = validator.New()

// AuthHandler defines the methods for handling authentication-related requests
type AuthHandler interface {
	MeHandler(w http.ResponseWriter, r *http.Request)
	SignUpHandler(w http.ResponseWriter, r *http.Request)
	SignInHandler(w http.ResponseWriter, r *http.Request)
}

// AuthHandlerImpl implements the AuthHandler interface and handles authentication logic
type AuthHandlerImpl struct {
	configuration *config.Configuration
	logger        *log.Logger
	authService   service.AuthService
	jwtService    service.JwtService
}

// NewAuthHandler creates and returns a new instance of AuthHandlerImpl
func NewAuthHandler(configuration *config.Configuration, logger *log.Logger, authService service.AuthService, jwtService service.JwtService) AuthHandler {
	return &AuthHandlerImpl{
		configuration: configuration,
		logger:        logger,
		authService:   authService,
		jwtService:    jwtService,
	}
}

// sendResponse sends the given response to the client with the appropriate HTTP status code.
// If encoding the response fails, an error is logged.
func (h *AuthHandlerImpl) sendResponse(response *dto.GenericResponse, w http.ResponseWriter) {
	w.WriteHeader(response.StatusCode)
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		h.logger.Printf("Failed to encode response: %v.\n", encodeErr)
	}
}

// MeHandler handles the "Me" endpoint, which retrieves the current user's information.
func (h *AuthHandlerImpl) MeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("Error: User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "User not found or failed to retrieve",
		}, w)
		return
	}

	h.logger.Printf("User %s successfully retrieved.\n", user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "User successfully retrieved",
		Data:       user,
	}, w)
}

// SignUpHandler handles user registration requests, validating and creating a new user.
func (h *AuthHandlerImpl) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var request dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Printf("Invalid request payload: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       err.Error(),
		}, w)
	}

	if err := authValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.authService.SignUp(&request); err != nil {
		h.logger.Printf("User registration failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "User registration failed",
			Data:       err.Error(),
		}, w)
		return
	}

	h.logger.Printf("User %s successfully registered.\n", request.Email)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusCreated,
		Message:    "User registered successfully. Please sign in to continue.",
	}, w)
}

// SignInHandler handles user sign-in requests and generates a JWT token for the user.
func (h *AuthHandlerImpl) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var request dto.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Printf("Invalid request payload: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := authValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	user, err := h.authService.SignIn(&request)
	if err != nil {
		h.logger.Printf("Sign-in failed for %s: %v.\n", request.Email, err)
		response := dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid credentials",
			Data:       err.Error(),
		}

		w.WriteHeader(response.StatusCode)
		if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
			h.logger.Printf("Failed to encode error response: %v.\n", encodeErr)
		}
		return
	}

	value, err := h.jwtService.CreateToken(user.ID, user.Role, user.IsActive)
	if err != nil {
		h.logger.Printf("Login failed for %s: %v.\n", request.Email, err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Something went wrong when creating the token",
			Data:       err.Error(),
		}, w)
		return
	}

	h.logger.Printf("User %s successfully signed in.\n", request.Email)
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", value))
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "User signed in successfully",
		Data:       fmt.Sprintf("Authorization: Bearer %s", value),
	}, w)
}
