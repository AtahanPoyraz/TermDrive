package model

import (
	"errors"
	"time"

	"github.com/AtahanPoyraz/TermDrive/server/internal/model/enum"
	"github.com/AtahanPoyraz/TermDrive/server/util/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserModel represents the structure of a user record in the database.
// It contains information about the user, including username, email, password, role, and status.
type UserModel struct {
	ID       uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey;not null;unique;index" json:"user_id"`
	Username string    `gorm:"column:username;type:VARCHAR(32);not null;unique;index" json:"username"`
	Email    string    `gorm:"column:email;type:VARCHAR(64);not null;unique;index" json:"email"`
	Password string    `gorm:"column:password;type:VARCHAR(128);not null" json:"password"`
	Role     enum.Role `gorm:"column:role;type:VARCHAR(16);not null;default:'USER';index" json:"role"`
	IsActive bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`

	Files []FileModel `gorm:"foreignKey:UserId;constraint:onDelete:CASCADE" json:"-"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:true" json:"updated_at"`
}

// TableName defines the table name for the UserModel.
func (UserModel) TableName() string {
	return "users"
}

// BeforeSave is a hook that runs before saving the model to the database.
// It validates the user's input and hashes the password before saving.
func (u *UserModel) BeforeSave(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if err := validator.IsValidUsername(u.Username); err != nil {
		return err
	}

	if err := validator.IsValidEmail(u.Email); err != nil {
		return err
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if !validator.IsValidPassword(u.Password) {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one digit, and one punctuation mark")
	}

	if !enum.Role(u.Role).IsValid() {
		return errors.New("invalid role")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}
