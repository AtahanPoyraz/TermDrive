package service

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/dto"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/AtahanPoyraz/TermDrive/server/internal/repository"
	"github.com/AtahanPoyraz/TermDrive/server/util/fsutil"
	"gorm.io/gorm"
)

// UserService defines the methods for managing user-related operations such as fetching, creating, updating, and deleting users.
type UserService interface {
	FetchUsers(request *dto.FetchUsersRequest) ([]model.UserModel, error)
	FetchUserById(request *dto.FetchUserByIdRequest) (model.UserModel, error)
	FetchUserByUsername(request *dto.FetchUserByUsernameRequest) (model.UserModel, error)
	FetchUserByEmail(request *dto.FetchUserByEmailRequest) (model.UserModel, error)
	CreateUser(request *dto.CreateUserRequest) error
	UpdateUserById(request *dto.UpdateUserByIdRequest) error
	DeleteUserById(request *dto.DeleteUserByIdRequest) error
}

// UserServiceImpl is the concrete implementation of the UserService interface.
type UserServiceImpl struct {
	configuration  *config.Configuration
	userRepository repository.UserRepository
}

// NewUserService initializes and returns a new instance of UserServiceImpl with the provided configuration and user repository.
func NewUserService(configuration *config.Configuration, userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{configuration: configuration, userRepository: userRepository}
}

// FetchUsers retrieves a list of users based on the provided request parameters (limit, offset).
// Returns: A list of users and an error if something goes wrong.
func (s *UserServiceImpl) FetchUsers(request *dto.FetchUsersRequest) ([]model.UserModel, error) {
	return s.userRepository.FetchUsers(request.Limit, request.Offset)
}

// FetchUsers retrieves a list of users based on the provided request parameters (limit, offset).
// Returns: A list of users and an error if something goes wrong.
func (s *UserServiceImpl) FetchUserById(request *dto.FetchUserByIdRequest) (model.UserModel, error) {
	return s.userRepository.FetchUserById(request.UserId)
}

// FetchUserByUsername fetches a single user by their username.
// Returns: The user model and an error if something goes wrong.
func (s *UserServiceImpl) FetchUserByUsername(request *dto.FetchUserByUsernameRequest) (model.UserModel, error) {
	return s.userRepository.FetchUserByUsername(request.Username)
}

// FetchUserByEmail fetches a single user by their email address.
// Returns: The user model and an error if something goes wrong.
func (s *UserServiceImpl) FetchUserByEmail(request *dto.FetchUserByEmailRequest) (model.UserModel, error) {
	return s.userRepository.FetchUserByEmail(request.Email)
}

// CreateUser creates a new user in the system and sets up their directory.
// Returns: An error if the creation or directory setup fails.
func (s *UserServiceImpl) CreateUser(request *dto.CreateUserRequest) error {
	if err := s.userRepository.CreateUser(
		request.Username,
		request.Email,
		request.Password,
		request.Role,
		request.IsActive,
	); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	directoryPath := filepath.Join(s.configuration.TermDrive.StoragePath, request.Username)
	if err := fsutil.Create(directoryPath); err != nil {
		user, fetchErr := s.userRepository.FetchUserByUsername(request.Username)
		if fetchErr != nil {
			return fmt.Errorf("failed to create directory and could not fetch user for deletion: %v", fetchErr)
		}

		if delErr := s.userRepository.DeleteUserById(user.ID); delErr != nil {
			return fmt.Errorf("failed to create directory and could not delete user: %v", delErr)
		}

		return fmt.Errorf("failed to create directory, user has been deleted: %v", err)
	}

	return nil
}

// UpdateUserById updates the details of an existing user by their ID.
// Returns: An error if the update fails.
func (s *UserServiceImpl) UpdateUserById(request *dto.UpdateUserByIdRequest) error {
	return s.userRepository.UpdateUserById(
		request.UserId,
		request.Email,
		request.Password,
		request.Role,
		request.IsActive,
	)
}

// UpdateUserById updates the details of an existing user by their ID.
// Returns: An error if the update fails.
func (s *UserServiceImpl) DeleteUserById(request *dto.DeleteUserByIdRequest) error {
	user, err := s.userRepository.FetchUserById(request.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return fmt.Errorf("failed to fetch user: %v", err)
	}

	if err := s.userRepository.DeleteUserById(request.UserId); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	directoryPath := filepath.Join(s.configuration.TermDrive.StoragePath, user.Username)
	if err := fsutil.Delete(directoryPath); err != nil {
		return fmt.Errorf("user deleted, but failed to delete directory: %v", err)
	}

	return nil
}
