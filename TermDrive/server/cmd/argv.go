package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/dto"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model/enum"
	"github.com/AtahanPoyraz/TermDrive/server/internal/repository"
	"github.com/AtahanPoyraz/TermDrive/server/internal/service"
	"gorm.io/gorm"
)

var (
	option   string
	username string
	email    string
	password string
)

var userService service.UserService

// InitializeDependencies sets up the necessary dependencies for the command-line operation,
// including parsing command-line arguments and initializing the user service.
//
// Arguments:
// - dependencies: A pointer to the Dependencies struct containing application and system context.
func InitializeDependencies(dependencies *config.Dependencies) {
	flag.StringVar(&option, "option", "", "Operation to perform (e.g., createadmin)")
	flag.StringVar(&username, "username", "", "Admin username (required)")
	flag.StringVar(&email, "email", "", "Admin email address (required)")
	flag.StringVar(&password, "password", "", "Admin password (required)")

	flag.Parse()

	userService = service.NewUserService(dependencies.AppContext.Configuration, repository.NewUserRepository(dependencies.SystemContext.Database))
}

// ArgumentProcessing processes the provided command-line arguments and performs the requested operation.
//
// Arguments:
// - dependencies: A pointer to the Dependencies struct containing application and system context.
//
// Returns:
// - An error if the option is invalid or if there is an issue performing the requested operation.
func ArgumentProcessing(dependencies *config.Dependencies) error {
	InitializeDependencies(dependencies)
	switch strings.ToLower(option) {
	case "":
		return nil

	case "createadmin":
		if err := CreateAdmin(userService, username, email, password); err != nil {
			return err
		}

		dependencies.AppContext.Logger.Println("Admin user created successfully.")
		os.Exit(0)
		return nil

	default:
		return errors.New("error: Invalid option provided")
	}
}

// CreateAdmin creates a new admin user with the provided username, email, and password.
// It validates the input, checks if the username or email already exists, and then
// creates the user with an admin role.
//
// Arguments:
// - userService: The user service used to interact with user-related operations.
// - username: The username for the new admin user.
// - email: The email for the new admin user.
// - password: The password for the new admin user.
//
// Returns:
// - An error if the input is invalid or if there is an issue creating the admin user.
func CreateAdmin(userService service.UserService, username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("error: Admin username, email, and password are required")
	}

	matched, err := regexp.Match("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", []byte(email))
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("invalid email format")
	}

	_, err = userService.FetchUserByUsername(&dto.FetchUserByUsernameRequest{Username: username})
	if err == nil {
		return fmt.Errorf("error: Username '%s' is already taken", username)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	_, err = userService.FetchUserByEmail(&dto.FetchUserByEmailRequest{Email: email})
	if err == nil {
		return fmt.Errorf("error: Email '%s' is already registered", email)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return userService.CreateUser(&dto.CreateUserRequest{
		Username: username,
		Email:    email,
		Password: password,
		Role:     enum.ROLE_ADMIN,
		IsActive: true,
	})
}
