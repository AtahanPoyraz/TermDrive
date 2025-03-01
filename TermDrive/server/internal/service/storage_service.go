package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/dto"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/AtahanPoyraz/TermDrive/server/internal/repository"
	"github.com/AtahanPoyraz/TermDrive/server/util/fsutil"
	"gorm.io/gorm"
)

// StorageService defines the methods for file management and storage operations.
type StorageService interface {
	FetchFiles(request *dto.FetchFilesRequest) ([]model.FileModel, error)
	FetchFileByFileNameAndUserId(request *dto.FetchFileByNameAndUserIdRequest) (model.FileModel, error)
	FetchFileByFilePathAndUserId(request *dto.FetchFileByPathAndUserIdRequest) (model.FileModel, error)
	FetchFilesByUserId(request *dto.FetchFilesByUserIDRequest) ([]model.FileModel, error)
	FetchFileByFileId(request *dto.FetchFileByIdRequest) (model.FileModel, error)
	CreateFile(request *dto.CreateFileRequest) error
	UpdateFile(request *dto.UpdateFileByIdRequest) error
	DeleteFile(request *dto.DeleteFileByIdRequest) error

	Upload(request *dto.UploadRequest) error
	Download(request *dto.DownloadRequest) (*os.File, error)
	List(request *dto.ListRequest) ([]string, error)
	Delete(request *dto.DeleteRequest) error
}

// StorageServiceImpl is the concrete implementation of the StorageService interface.
type StorageServiceImpl struct {
	configuration     *config.Configuration
	storageRepository repository.StorageRepository
}

// NewStorageService initializes and returns a new instance of StorageServiceImpl with the provided configuration and repository.
func NewStorageService(configuration *config.Configuration, storageRepository repository.StorageRepository) StorageService {
	return &StorageServiceImpl{
		configuration:     configuration,
		storageRepository: storageRepository,
	}
}

// FetchFiles fetches a list of files with the specified limit and offset.
// Returns: A slice of FileModel and an error if any issue occurs while fetching the files.
func (s *StorageServiceImpl) FetchFiles(request *dto.FetchFilesRequest) ([]model.FileModel, error) {
	return s.storageRepository.FetchFiles(request.Limit, request.Offset)

}

// FetchFileByFileNameAndUserId fetches a file by its name and the user ID.
// Returns: A FileModel and an error if the file is not found or any issue occurs.
func (s *StorageServiceImpl) FetchFileByFileNameAndUserId(request *dto.FetchFileByNameAndUserIdRequest) (model.FileModel, error) {
	return s.storageRepository.FetchFileByFileNameAndUserId(request.FileName, request.UserId)
}

// FetchFileByFilePathAndUserId fetches a file by its path and the user ID.
// Returns: A FileModel and an error if the file is not found or any issue occurs.
func (s *StorageServiceImpl) FetchFileByFilePathAndUserId(request *dto.FetchFileByPathAndUserIdRequest) (model.FileModel, error) {
	return s.storageRepository.FetchFileByFilePathAndUserId(request.FilePath, request.UserId)
}

// FetchFilesByUserId fetches all files for a specific user.
// Returns: A slice of FileModel and an error if any issue occurs while fetching the files.
func (s *StorageServiceImpl) FetchFilesByUserId(request *dto.FetchFilesByUserIDRequest) ([]model.FileModel, error) {
	return s.storageRepository.FetchFilesByUserId(request.UserId)
}

// FetchFileByFileId fetches a file by its unique ID.
// Returns: A FileModel and an error if the file is not found or any issue occurs.
func (s *StorageServiceImpl) FetchFileByFileId(request *dto.FetchFileByIdRequest) (model.FileModel, error) {
	return s.storageRepository.FetchFileByFileId(request.FileId)
}

// CreateFile creates a new file record in the database.
// Returns: An error if there is an issue while creating the file.
func (s *StorageServiceImpl) CreateFile(request *dto.CreateFileRequest) error {
	return s.storageRepository.CreateFile(
		request.FileName,
		request.FilePath,
		request.FileSize,
		request.MimeType,
		request.Permissions,
		request.LastModified,
		request.UserId,
	)
}

// UpdateFile updates an existing file's details.
// Returns: An error if there is an issue while updating the file.
func (s *StorageServiceImpl) UpdateFile(request *dto.UpdateFileByIdRequest) error {
	return s.storageRepository.UpdateFile(
		request.FileId,
		request.FileName,
		request.FilePath,
		request.FileSize,
		request.MimeType,
		request.Permissions,
		request.LastModified,
		request.UserId,
	)
}

// DeleteFile deletes a file by its ID.
// Returns: An error if there is an issue while deleting the file.
func (s *StorageServiceImpl) DeleteFile(request *dto.DeleteFileByIdRequest) error {
	if _, err := s.storageRepository.FetchFileByFileId(request.FileId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("file does not exist")
		}

		return nil
	}
	return s.storageRepository.DeleteFileById(request.FileId)
}

// Upload handles the file upload process. It writes the file to the specified path, extracts relevant file
// Returns: An error if there is an issue while upload the file.
func (s *StorageServiceImpl) Upload(request *dto.UploadRequest) error {
	if err := fsutil.WriteFile(context.Background(), request.File, request.FilePath, s.configuration.TermDrive.UploadSize); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fileSize, err := fsutil.ExtractFileSize(request.FilePath)
	if err != nil {
		if deleteErr := fsutil.Delete(request.FilePath); deleteErr != nil {
			return fmt.Errorf("failed to delete directory after extracting file size error: %v", deleteErr)
		}

		return fmt.Errorf("failed to extract file size: %v", err)
	}

	mimeType, err := fsutil.ExtractMimeType(request.FilePath)
	if err != nil {
		if deleteErr := fsutil.Delete(request.FilePath); deleteErr != nil {
			return fmt.Errorf("failed to delete directory after extracting mime type error: %v", deleteErr)
		}

		return fmt.Errorf("failed to extract mime type: %v", err)
	}

	permissions, err := fsutil.ExtractFilePermissions(request.FilePath)
	if err != nil {
		if deleteErr := fsutil.Delete(request.FilePath); deleteErr != nil {
			return fmt.Errorf("failed to delete directory after extracting file permissions error: %v", deleteErr)
		}

		return fmt.Errorf("failed to extract file permissions: %v", err)
	}

	lastModified, err := fsutil.ExtractFileLastModified(request.FilePath)
	if err != nil {
		if deleteErr := fsutil.Delete(request.FilePath); deleteErr != nil {
			return fmt.Errorf("failed to delete directory after extracting last modified date error: %v", deleteErr)
		}

		return fmt.Errorf("failed to extract last modified date: %v", err)
	}

	file, err := s.storageRepository.FetchFileByFilePathAndUserId(request.FilePath, request.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.storageRepository.CreateFile(
				request.FileName,
				request.FilePath,
				fileSize,
				mimeType,
				permissions,
				lastModified,
				request.UserId,
			); err != nil {
				if deleteErr := fsutil.Delete(request.FilePath); deleteErr != nil {
					return fmt.Errorf("failed to delete directory after database creation error: %v", deleteErr)
				}

				return fmt.Errorf("failed to create new file in the database: %v", err)
			}

			return nil
		}

		return fmt.Errorf("failed to fetch file: %v", err)
	}

	if err := s.storageRepository.UpdateFile(
		file.ID,
		request.FileName,
		request.FilePath,
		fileSize,
		mimeType,
		permissions,
		lastModified,
		request.UserId,
	); err != nil {
		if deleteErr := fsutil.Delete(request.FilePath); deleteErr != nil {
			return fmt.Errorf("failed to delete directory after database update error: %v", deleteErr)
		}

		return fmt.Errorf("failed to update file in the database: %v", err)
	}

	return nil
}

// Download retrieves the content of a file specified by the given path. It returns the file if it exists.
// Returns: An error if there is an issue while download the file.
func (s *StorageServiceImpl) Download(request *dto.DownloadRequest) (*os.File, error) {
	file, err := fsutil.ReadFile(context.Background(), request.FilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}

		return nil, fmt.Errorf("failed to retrieve content: %w", err)
	}

	return file, nil
}

// List retrieves the list of files or directories located at the specified path for the given user. It
// Returns: An error if there is an issue while list the file.
func (s *StorageServiceImpl) List(request *dto.ListRequest) ([]string, error) {
	content, err := fsutil.List(context.WithValue(context.Background(), config.UserContextKey, request.User), request.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}

		return nil, fmt.Errorf("failed to retrieve content: %w", err)
	}

	return content, nil
}

// Delete handles the deletion of files and directories.
// It deletes the contents inside the directory if it's a directory, or just the file if it's a file.
func (s *StorageServiceImpl) Delete(request *dto.DeleteRequest) error {
	stat, err := os.Stat(request.Path)
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			if strings.HasSuffix(pathErr.Error(), "no such file or directory") {
				return &os.PathError{
					Op:   pathErr.Op,
					Path: pathErr.Path,
					Err:  fmt.Errorf("no such file or directory"),
				}
			}
		}
		return fmt.Errorf("failed to access path: %v", err.Error())
	}

	if filepath.Base(request.Path) == request.User.Username {
		if err := s.storageRepository.DeleteFilesByDirectory(request.Path); err != nil {
			return fmt.Errorf("failed to delete files in directory: %v", err)
		}

		entries, err := os.ReadDir(request.Path)
		if err != nil {
			return fmt.Errorf("failed to read directory: %v", err)
		}

		if len(entries) == 0 {
			return &os.PathError{
				Op:   "delete",
				Path: request.Path,
				Err:  fmt.Errorf("no such file or directory"),
			}
		}

		for _, entry := range entries {
			fullPath := filepath.Join(request.Path, entry.Name())
			if err := fsutil.Delete(fullPath); err != nil {
				return fmt.Errorf("failed to delete %s: %v", fullPath, err)
			}
		}

	} else if stat.IsDir() {
		if err := s.storageRepository.DeleteFilesByDirectory(request.Path); err != nil {
			return fmt.Errorf("failed to delete files in directory: %v", err)
		}

		if err := fsutil.Delete(request.Path); err != nil {
			return fmt.Errorf("failed to delete directory: %v", err)
		}

	} else {
		if err := s.storageRepository.DeleteFileByFilePath(request.Path); err != nil {
			return fmt.Errorf("failed to delete file: %v", err)
		}

		if err := fsutil.Delete(request.Path); err != nil {
			return fmt.Errorf("failed to delete file: %v", err)
		}
	}

	return nil
}
