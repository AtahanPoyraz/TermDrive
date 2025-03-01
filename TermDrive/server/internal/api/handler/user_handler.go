package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/dto"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/AtahanPoyraz/TermDrive/server/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var userValidator = validator.New()

// UserHandler defines the interface for user-related HTTP handlers.
type UserHandler interface {
	FetchUserHandler(w http.ResponseWriter, r *http.Request)
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
	UpdateUserByIdHandler(w http.ResponseWriter, r *http.Request)
	DeleteUserByIdHandler(w http.ResponseWriter, r *http.Request)
}

// UserHandlerImpl is the implementation of the UserHandler interface.
type UserHandlerImpl struct {
	configuration *config.Configuration
	logger        *log.Logger
	userService   service.UserService
}

// NewUserHandler creates a new instance of UserHandlerImpl.
func NewUserHandler(configuration *config.Configuration, logger *log.Logger, userService service.UserService) UserHandler {
	return &UserHandlerImpl{
		configuration: configuration,
		logger:        logger,
		userService:   userService,
	}
}

// sendResponse sends the given response to the client with the appropriate HTTP status code.
// If encoding the response fails, an error is logged.
func (h *UserHandlerImpl) sendResponse(response *dto.GenericResponse, w http.ResponseWriter) {
	w.WriteHeader(response.StatusCode)
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		h.logger.Printf("Failed to encode response: %v.\n", encodeErr)
	}
}

// parseUrlQuery parses the URL query parameters and returns them as a map.
func (h *UserHandlerImpl) parseUrlQuery(r *http.Request) (map[string]interface{}, error) {
	params := make(map[string]interface{})

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params["limit"] = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params["offset"] = offset
		}
	}

	if userIDStr := r.URL.Query().Get("userId"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid userId format")
		}
		params["userId"] = userID
	}

	if username := r.URL.Query().Get("username"); username != "" {
		params["username"] = username
	}

	if email := r.URL.Query().Get("email"); email != "" {
		params["email"] = email
	}

	return params, nil
}

// FetchUserHandler handles requests for fetching user(s) based on query parameters.
func (h *UserHandlerImpl) FetchUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	params, err := h.parseUrlQuery(r)
	if err != nil {
		h.logger.Printf("Query parsing error: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, w)
		return
	}

	switch {
	case params["limit"] != nil && params["offset"] != nil:
		var request dto.FetchUsersRequest
		if request.Limit, ok = params["limit"].(int); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid limit format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid limit format",
			}, w)
			return
		}

		if request.Offset, ok = params["offset"].(int); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid offset format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid offset format",
			}, w)
			return
		}

		if err := userValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		users, err := h.userService.FetchUsers(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch users: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve users",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("Users retrieved successfully by %s.\n", user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "Users retrieved successfully",
			Data:       users,
		}, w)

	case params["userId"] != nil:
		var request dto.FetchUserByIdRequest
		if request.UserId, ok = params["userId"].(uuid.UUID); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid userId format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid userId format",
			}, w)
			return
		}

		if err := userValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		fetchedUser, err := h.userService.FetchUserById(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch user: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve user",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("User (%s) retrieved successfully by %s.\n", fetchedUser.Username, user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "User retrieved successfully",
			Data:       user,
		}, w)

	case params["username"] != nil:
		var request dto.FetchUserByUsernameRequest
		if request.Username, ok = params["username"].(string); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid username format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid username format",
			}, w)
			return
		}

		if err := userValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		fetchedUser, err := h.userService.FetchUserByUsername(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch user: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve user",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("User (%s) retrieved successfully by %s.\n", fetchedUser.Username, user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "User retrieved successfully",
			Data:       user,
		}, w)

	case params["email"] != nil:
		var request dto.FetchUserByEmailRequest
		if request.Email, ok = params["email"].(string); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid email format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid email format",
			}, w)
			return
		}

		if err := userValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		fetchedUser, err := h.userService.FetchUserByEmail(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch user: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve user",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("User (%s) retrieved successfully by %s.\n", fetchedUser.Username, user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "User retrieved successfully",
			Data:       user,
		}, w)
	}
}

// CreateUserHandler handles the creation of a new user.
func (h *UserHandlerImpl) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	var request dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Printf("Invalid request payload: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := userValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.userService.CreateUser(&request); err != nil {
		h.logger.Printf("Failed to create user: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to create user",
			Data:       err.Error(),
		}, w)
		return
	}

	h.logger.Printf("User (%s) created successfully by %s.\n", request.Username, user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "User created successfully",
	}, w)
}

// UpdateUserByIdHandler handles updating a user's details based on their ID.
func (h *UserHandlerImpl) UpdateUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	params, err := h.parseUrlQuery(r)
	if err != nil {
		h.logger.Printf("Query parsing error: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, w)
		return
	}

	var request dto.UpdateUserByIdRequest
	if request.UserId, ok = params["userId"].(uuid.UUID); !ok {
		h.logger.Printf("Validation failed: %v.\n", errors.New("invalid userId format"))
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid userId format",
		}, w)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Printf("Error decoding request payload: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := userValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	updatedUser, err := h.userService.FetchUserById(&dto.FetchUserByIdRequest{UserId: request.UserId})
	if err != nil {
		h.logger.Printf("Failed to fetch user: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "User not found.",
			}, w)
			return
		}

		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve user",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.userService.UpdateUserById(&request); err != nil {
		h.logger.Printf("Failed to update user: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update user",
			Data:       err.Error(),
		}, w)
		return
	}

	h.logger.Printf("User (%s) updated successfully by %s.\n", updatedUser.Username, user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "User updated successfully",
	}, w)
}

// DeleteUserByIdHandler handles the deletion of a user based on their ID.
func (h *UserHandlerImpl) DeleteUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	params, err := h.parseUrlQuery(r)
	if err != nil {
		h.logger.Printf("Query parsing error: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, w)
		return
	}

	var request dto.DeleteUserByIdRequest
	if request.UserId, ok = params["userId"].(uuid.UUID); !ok {
		h.logger.Printf("Validation failed: %v.\n", errors.New("invalid userId format"))
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid userId format",
		}, w)
		return
	}

	if err := userValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	deletedUser, err := h.userService.FetchUserById(&dto.FetchUserByIdRequest{UserId: request.UserId})
	if err != nil {
		h.logger.Printf("Failed to fetch user: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "User not found.",
			}, w)
			return
		}

		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve user",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.userService.DeleteUserById(&request); err != nil {
		h.logger.Printf("Failed to delete user: %v.\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "User not found.",
			}, w)
			return
		}

		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to delete user",
			Data:       err.Error(),
		}, w)
		return
	}

	h.logger.Printf("User (%s) deleted successfully by %s.\n", deletedUser.Username, user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "User deleted successfully",
	}, w)
}
