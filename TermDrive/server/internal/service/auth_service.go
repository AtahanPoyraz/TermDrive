package service

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/dto"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model/enum"
	"github.com/AtahanPoyraz/TermDrive/server/internal/pb"
	"github.com/AtahanPoyraz/TermDrive/server/internal/repository"
	"github.com/AtahanPoyraz/TermDrive/server/util/fsutil"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService defines the methods for user authentication and management.
type AuthService interface {
	Me(request *dto.MeRequest) (model.UserModel, error)
	SignUp(request *pb.CreateUserRequest) error
	SignIn(request *dto.SignInRequest) (model.UserModel, error)
}

// AuthServiceImpl is the concrete implementation of AuthService.
type AuthServiceImpl struct {
	configuration  *config.Configuration
	userRepository repository.UserRepository
}

// NewAuthService initializes and returns a new instance of AuthServiceImpl with the provided configuration and user repository.
func NewAuthService(configuration *config.Configuration, userRepository repository.UserRepository) AuthService {
	return &AuthServiceImpl{
		configuration:  configuration,
		userRepository: userRepository,
	}
}

// Me retrieves the authenticated user's details based on their ID.
// Returns: The user model and an error if something goes wrong.
func (s *AuthServiceImpl) Me(request *dto.MeRequest) (model.UserModel, error) {
	return s.userRepository.FetchUserById(request.UserId)
}

// SignUp registers a new user and sets up their directory.
// Returns: An error if the user creation or directory setup fails.
func (s *AuthServiceImpl) SignUp(request *pb.CreateUserRequest) error {
	if err := s.userRepository.CreateUser(request.FirstName, request.Email, request.Password, enum.ROLE_USER, true); err != nil {
		return err
	}

	dirPath := filepath.Join(s.configuration.TermDrive.StoragePath, request.FirstName)
	if err := fsutil.Create(dirPath); err != nil {
		user, fetchErr := s.userRepository.FetchUserByUsername(request.FirstName)
		if fetchErr != nil {
			return fmt.Errorf("failed to create directory and could not fetch user for deletion: %w", fetchErr)
		}

		if delErr := s.userRepository.DeleteUserById(user.ID); delErr != nil {
			return fmt.Errorf("failed to create directory and could not delete user: %w", delErr)
		}

		return fmt.Errorf("failed to create directory, user has been deleted: %w", err)
	}

	return nil
}

// SignIn validates the user's email and password to authenticate them.
// Returns: The user model if authentication is successful, or an error if the credentials are invalid.
func (s *AuthServiceImpl) SignIn(request *dto.SignInRequest) (model.UserModel, error) {
	user, err := s.userRepository.FetchUserByEmail(request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.UserModel{}, errors.New("invalid credentials")
		}

		return model.UserModel{}, fmt.Errorf("an error occurred while fetching the user: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return model.UserModel{}, errors.New("invalid credentials")
	}

	return user, nil
}
