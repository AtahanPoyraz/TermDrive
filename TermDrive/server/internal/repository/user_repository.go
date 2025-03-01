package repository

import (
	"errors"
	"strings"

	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository defines the methods for interacting with user data in the database.
type UserRepository interface {
	FetchUsers(limit, offset int) ([]model.UserModel, error)
	FetchUserById(id uuid.UUID) (model.UserModel, error)
	FetchUserByUsername(username string) (model.UserModel, error)
	FetchUserByEmail(email string) (model.UserModel, error)
	CreateUser(username, email, password string, role enum.Role, isActive bool) error
	UpdateUserById(id uuid.UUID, email, password string, role enum.Role, isActive bool) error
	DeleteUserById(id uuid.UUID) error
}

// UserRepositoryImpl is the concrete implementation of UserRepository using GORM.
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository initializes and returns an instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// FetchUsers retrieves a list of users with pagination parameters (limit and offset).
// Returns a slice of UserModel and an error if any.
func (r *UserRepositoryImpl) FetchUsers(limit, offset int) ([]model.UserModel, error) {
	var users []model.UserModel
	if err := r.db.Debug().Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FetchUserById retrieves a user record by their unique identifier.
// Returns a UserModel and an error if the user is not found or any other error occurs.
func (r *UserRepositoryImpl) FetchUserById(id uuid.UUID) (model.UserModel, error) {
	var user model.UserModel
	if err := r.db.Debug().Where("user_id = ?", id).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// FetchUserByUsername searches for a user by their username.
// Returns a UserModel and an error if the user is not found or any other error occurs.
func (r *UserRepositoryImpl) FetchUserByUsername(username string) (model.UserModel, error) {
	var user model.UserModel
	if err := r.db.Debug().Where("username = ?", username).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// FetchUserByEmail searches for a user by their email address.
// Returns a UserModel and an error if the user is not found or any other error occurs.
func (r *UserRepositoryImpl) FetchUserByEmail(email string) (model.UserModel, error) {
	var user model.UserModel
	if err := r.db.Debug().Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

// CreateUser inserts a new user record into the database with the provided details.
// Returns an error if the user could not be created due to a duplicate or other issue.
func (r *UserRepositoryImpl) CreateUser(username, email, password string, role enum.Role, isActive bool) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer tx.Rollback()

	user := &model.UserModel{
		Username: username,
		Email:    email,
		Password: password,
		Role:     role,
		IsActive: isActive,
	}

	if err := tx.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return errors.New("user could not be created due to existing data")
		}

		return err
	}

	return tx.Debug().Commit().Error
}

// UpdateUserById modifies the details of an existing user identified by their unique ID.
// Returns an error if the update could not be completed for any reason.
func (r *UserRepositoryImpl) UpdateUserById(id uuid.UUID, email, password string, role enum.Role, isActive bool) error {
	var user model.UserModel
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer tx.Rollback()

	if err := tx.Where("user_id = ?", id).First(&user).Error; err != nil {
		return err
	}

	if email != "" && email != user.Email {
		user.Email = email
	}
	if password != "" {
		user.Password = password
	}
	if role.IsValid() && role != user.Role {
		user.Role = role
	}
	if isActive != user.IsActive {
		user.IsActive = isActive
	}

	return tx.Debug().Save(&user).Commit().Error
}

// DeleteUserById removes a user permanently from the database by their unique ID.
// Returns an error if the deletion fails.
func (r *UserRepositoryImpl) DeleteUserById(id uuid.UUID) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer tx.Rollback()

	if err := tx.Debug().Unscoped().Where("user_id = ?", id).Delete(&model.UserModel{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
