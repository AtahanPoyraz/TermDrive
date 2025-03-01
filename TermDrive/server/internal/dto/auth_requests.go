package dto

import (
	"github.com/google/uuid"
)

// MeRequest represents the request structure for fetching the currently authenticated user's information.
// It requires the user's ID as input.
//
// Fields:
// - UserId: The unique identifier of the user (required).
type MeRequest struct {
	UserId uuid.UUID `json:"user_id" validate:"required"`
}

// SignUpRequest represents the structure for the user registration request.
// It includes the necessary details for creating a new user account.
//
// Fields:
// - Username: The desired username of the new user (required, 3 to 32 characters).
// - Email: The email address of the new user (required, valid email format).
// - Password: The password for the new user account (required, 4 to 128 characters).
type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=128"`
}

// SignInRequest represents the structure for the user login request.
// It contains the required details for authenticating a user.
//
// Fields:
// - Email: The email address of the user attempting to sign in (required, valid email format).
// - Password: The password for the user account (required, 4 to 128 characters).
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=128"`
}
