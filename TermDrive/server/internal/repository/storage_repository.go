package repository

import (
	"errors"
	"time"

	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StorageRepository defines methods for interacting with file storage in the database.
type StorageRepository interface {
	FetchFiles(limit, offset int) ([]model.FileModel, error)
	FetchFileByFileId(id uuid.UUID) (model.FileModel, error)
	FetchFilesByUserId(id uuid.UUID) ([]model.FileModel, error)
	FetchFileByFileNameAndUserId(fileName string, userId uuid.UUID) (model.FileModel, error)
	FetchFileByFilePathAndUserId(filePath string, userId uuid.UUID) (model.FileModel, error)
	CreateFile(fileName, filePath string, fileSize int64, mimeType string, permissions string, lastModified time.Time, userId uuid.UUID) error
	UpdateFile(fileId uuid.UUID, fileName, filePath string, fileSize int64, mimeType string, permissions string, lastModified time.Time, userId uuid.UUID) error
	DeleteFileById(id uuid.UUID) error
	DeleteFilesByDirectory(directoryPath string) error
	DeleteFileByFilePath(filePath string) error
}

// StorageRepositoryImpl is the concrete implementation of StorageRepository using GORM.
type StorageRepositoryImpl struct {
	db *gorm.DB
}

// NewStorageRepository initializes and returns an instance of StorageRepository.
func NewStorageRepository(db *gorm.DB) StorageRepository {
	return &StorageRepositoryImpl{
		db: db,
	}
}

// FetchFiles retrieves a list of files with pagination based on the provided limit and offset.
// Returns a slice of FileModel objects and any error encountered.
func (r *StorageRepositoryImpl) FetchFiles(limit, offset int) ([]model.FileModel, error) {
	var files []model.FileModel
	if err := r.db.Debug().Preload("User").Limit(limit).Offset(offset).Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}

// FetchFileByFileId retrieves a file by its unique file ID (UUID).
// Returns the FileModel for the file and any error encountered.
func (r *StorageRepositoryImpl) FetchFileByFileId(id uuid.UUID) (model.FileModel, error) {
	var file model.FileModel
	if err := r.db.Debug().Preload("User").Where("file_id = ?", id).First(&file).Error; err != nil {
		return model.FileModel{}, err
	}

	return file, nil
}

// FetchFilesByUserId retrieves all files associated with a specific user by their user ID.
// Returns a slice of FileModel objects for the files and any error encountered.
func (r *StorageRepositoryImpl) FetchFilesByUserId(id uuid.UUID) ([]model.FileModel, error) {
	var files []model.FileModel
	if err := r.db.Debug().Preload("User").Where("user_id = ?", id).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// FetchFileByFileNameAndUserId retrieves a file by its name and the associated user ID.
// Returns the FileModel for the file and any error encountered.
func (r *StorageRepositoryImpl) FetchFileByFileNameAndUserId(fileName string, userId uuid.UUID) (model.FileModel, error) {
	var file model.FileModel
	if err := r.db.Debug().Preload("User").Where("file_name = ? AND user_id = ?", fileName, userId).First(&file).Error; err != nil {
		return model.FileModel{}, err
	}

	return file, nil
}

// FetchFileByFilePathAndUserId retrieves a file by its path and the associated user ID.
// Returns the FileModel for the file and any error encountered.
func (r *StorageRepositoryImpl) FetchFileByFilePathAndUserId(filePath string, userId uuid.UUID) (model.FileModel, error) {
	var file model.FileModel
	if err := r.db.Debug().Preload("User").Where("file_path = ? AND user_id = ?", filePath, userId).First(&file).Error; err != nil {
		return model.FileModel{}, err
	}

	return file, nil
}

// CreateFile adds a new file to the database with the specified attributes.
// Returns an error if the file creation fails.
func (r *StorageRepositoryImpl) CreateFile(fileName, filePath string, fileSize int64, mimeType string, permissions string, lastModified time.Time, userId uuid.UUID) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer tx.Rollback()

	if _, err := r.FetchFileByFileNameAndUserId(fileName, userId); err == nil {
		tx.Rollback()
		return errors.New("file with the same name already exists for this user")

	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return errors.New("error checking existing file name")
	}

	if _, err := r.FetchFileByFilePathAndUserId(filePath, userId); err == nil {
		tx.Rollback()
		return errors.New("file with the same path already exists for this user")

	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return errors.New("error checking existing file path")
	}

	file := model.FileModel{
		FileName:     fileName,
		FilePath:     filePath,
		FileSize:     fileSize,
		MimeType:     mimeType,
		Permissions:  permissions,
		LastModified: lastModified,
		UserId:       userId,
	}

	return tx.Debug().Create(&file).Commit().Error
}

// UpdateFile updates the attributes of an existing file in the database based on its file ID.
// Returns an error if the file update fails.
func (r *StorageRepositoryImpl) UpdateFile(fileId uuid.UUID, fileName, filePath string, fileSize int64, mimeType string, permissions string, lastModified time.Time, userId uuid.UUID) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer tx.Rollback()

	var existingFile model.FileModel
	if err := tx.Debug().Where("file_id = ?", fileId).First(&existingFile).Error; err != nil {
		return err
	}

	if fileName != "" && fileName != existingFile.FileName {
		existingFile.FileName = fileName
	}
	if filePath != "" && filePath != existingFile.FilePath {
		existingFile.FilePath = filePath
	}
	if fileSize != existingFile.FileSize {
		existingFile.FileSize = fileSize
	}
	if mimeType != "" && mimeType != existingFile.MimeType {
		existingFile.FileSize = fileSize
	}
	if permissions != "" && permissions != existingFile.Permissions {
		existingFile.Permissions = permissions
	}
	if lastModified != existingFile.LastModified {
		existingFile.LastModified = lastModified
	}

	if userId != uuid.Nil && userId != existingFile.UserId {
		existingFile.UserId = userId
	}

	return tx.Save(&existingFile).Commit().Error
}

// DeleteFileById removes a file from the database based on its unique file ID.
// Returns an error if the deletion fails.
func (r *StorageRepositoryImpl) DeleteFileById(id uuid.UUID) error {
	tx := r.db.Debug().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer tx.Rollback()

	if err := tx.Unscoped().Where("file_id = ?", id).Delete(&model.FileModel{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// DeleteFilesByDirectory removes all files in a specific directory, identified by its path.
// Returns an error if the deletion fails.
func (r *StorageRepositoryImpl) DeleteFilesByDirectory(directoryPath string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer tx.Rollback()

	if err := tx.Debug().Unscoped().Where("file_path LIKE ?", directoryPath+"%").Delete(&model.FileModel{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// DeleteFileByFilePath removes a file from the database based on its file path.
// Returns an error if the deletion fails.
func (r *StorageRepositoryImpl) DeleteFileByFilePath(filePath string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer tx.Rollback()

	if err := tx.Debug().Unscoped().Where("file_path = ?", filePath).Delete(&model.FileModel{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
