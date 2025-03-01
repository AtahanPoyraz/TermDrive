package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileModel represents the structure of a file record in the database.
// It stores the file metadata, such as file name, path, size, and associated user.
type FileModel struct {
	ID           uuid.UUID `gorm:"column:file_id;type:uuid;primaryKey;not null;unique;index" json:"file_id"`
	FileName     string    `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	FilePath     string    `gorm:"column:file_path;type:text;not null" json:"file_path"`
	FileSize     int64     `gorm:"column:file_size;type:bigint;not null" json:"file_size"`
	MimeType     string    `gorm:"column:mime_type;type:varchar(100);not null" json:"mimetype"`
	Permissions  string    `gorm:"column:permissions;type:varchar(255)" json:"permissions"`
	LastModified time.Time `gorm:"column:last_modified;autoUpdateTime:true" json:"last_modified"`

	UserId uuid.UUID `gorm:"column:user_id;type:uuid;not null;index" json:"user_id"`
	User   UserModel `gorm:"foreignKey:UserId;references:ID" json:"-"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime:true" json:"updated_at"`
}

// TableName defines the table name for the FileModel.
func (FileModel) TableName() string {
	return "files"
}

// BeforeSave is a hook that runs before saving the model to the database.
// It sets a new UUID for the file ID if it is not already set.
func (f *FileModel) BeforeSave(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}

	return nil
}
