package dto

import (
	"github.com/AtahanPoyraz/TermDrive/server/internal/model/enum"
	"github.com/google/uuid"
)

// FetchUsersRequest represents a request to fetch a list of users.
// It includes pagination parameters for limiting the number of users
// returned and the offset for the results.
//
// Fields:
// - Limit: The maximum number of users to fetch (required).
// - Offset: The starting point for fetching users (required, must be >= 0).
type FetchUsersRequest struct {
	Limit  int `json:"limit" validate:"required"`
	Offset int `json:"offset" validate:"gte=0"`
}

// FetchUserByIdRequest represents a request to fetch a user by their unique ID.
//
// Fields:
// - UserId: The unique identifier of the user (required).
type FetchUserByIdRequest struct {
	UserId uuid.UUID `json:"userId" validate:"required"`
}

// FetchUserByEmailRequest represents a request to fetch a user by their email address.
//
// Fields:
// - Email: The email address of the user (required, must be in valid email format).
type FetchUserByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// FetchUserByUsernameRequest represents a request to fetch a user by their username.
//
// Fields:
// - Username: The username of the user (required, must be between 3 and 32 characters).
type FetchUserByUsernameRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
}

// CreateUserRequest represents a request to create a new user.
//
// Fields:
// - Username: The username of the user (required, between 3 and 32 characters).
// - Email: The email address of the user (required, must be in valid email format).
// - Password: The password of the user (required, between 4 and 128 characters).
// - Role: The role of the user (required, can be either 'ADMIN' or 'USER').
// - IsActive: Whether the user is active (required).
type CreateUserRequest struct {
	Username string    `json:"username" validate:"required,min=3,max=32"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=4,max=128"`
	Role     enum.Role `json:"role" validate:"required,oneof=ADMIN USER"`
	IsActive bool      `json:"isActive" validate:"required"`
}

// UpdateUserByIdRequest represents a request to update a user's information by their unique ID.
//
// Fields:
// - UserId: The unique ID of the user to update (required).
// - Email: The new email address for the user (optional, must be valid if provided).
// - Password: The new password for the user (optional, between 4 and 128 characters).
// - Role: The new role for the user (optional, can be 'ADMIN' or 'USER').
// - IsActive: The new active status for the user (optional).
type UpdateUserByIdRequest struct {
	UserId   uuid.UUID `json:"userId" validate:"required"`
	Email    string    `json:"email" validate:"omitempty,email"`
	Password string    `json:"password" validate:"omitempty,min=4,max=128"`
	Role     enum.Role `json:"role" validate:"omitempty,oneof=ADMIN USER"`
	IsActive bool      `json:"isActive"`
}

// DeleteUserByIdRequest represents a request to delete a user by their unique ID.
//
// Fields:
// - UserId: The unique ID of the user to delete (required).
type DeleteUserByIdRequest struct {
	UserId uuid.UUID `json:"userId" validate:"required"`
}
