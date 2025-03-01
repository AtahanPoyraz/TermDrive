package dto

import (
	"mime/multipart"
	"time"

	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/google/uuid"
)

// FetchFilesRequest represents a request to fetch a list of files with pagination.
// It includes parameters for limiting the number of files and specifying the offset.
//
// Fields:
// - Limit: The maximum number of files to fetch (required).
// - Offset: The starting point for fetching files (required, must be >= 0).
type FetchFilesRequest struct {
	Limit  int `json:"limit" validate:"required"`
	Offset int `json:"offset" validate:"gte=0"`
}

// FetchFileByNameAndUserIdRequest represents a request to fetch a file by its name and the user's unique ID.
//
// Fields:
// - FileName: The name of the file to fetch (required).
// - UserId: The unique ID of the user requesting the file (required).
type FetchFileByNameAndUserIdRequest struct {
	FileName string    `json:"fileName" validate:"required"`
	UserId   uuid.UUID `json:"userId" validate:"required"`
}

// FetchFileByPathAndUserIdRequest represents a request to fetch a file by its path and the user's unique ID.
//
// Fields:
// - FilePath: The path of the file to fetch (required).
// - UserId: The unique ID of the user requesting the file (required).
type FetchFileByPathAndUserIdRequest struct {
	FilePath string    `json:"filePath" validate:"required"`
	UserId   uuid.UUID `json:"userId" validate:"required"`
}

// FetchFilesByUserIDRequest represents a request to fetch all files owned by a user.
//
// Fields:
// - UserID: The unique ID of the user (required).
type FetchFilesByUserIDRequest struct {
	UserId uuid.UUID `json:"userId" validate:"required"`
}

// FetchFileByIdRequest represents a request to fetch a file by its unique ID.
//
// Fields:
// - FileId: The unique ID of the file (required).
type FetchFileByIdRequest struct {
	FileId uuid.UUID `json:"fileId" validate:"required"`
}

// FetchFileByNameRequest represents a request to fetch a file by its name.
//
// Fields:
// - FileName: The name of the file to fetch (required).
type FetchFileByNameRequest struct {
	FileName string `json:"fileName" validate:"required"`
}

// FetchFileByPathRequest represents a request to fetch a file by its path.
//
// Fields:
// - FilePath: The path of the file to fetch (required).
type FetchFileByPathRequest struct {
	FilePath string `json:"filePath" validate:"required"`
}

// CreateFileRequest represents a request to create a new file.
//
// Fields:
// - FileName: The name of the file (required).
// - FilePath: The path where the file is stored (required).
// - FileSize: The size of the file (required, must be >= 0).
// - MimeType: The MIME type of the file (required).
// - Permissions: The permissions for the file (required).
// - LastModified: The last modification time of the file (required).
// - UserId: The unique ID of the user who owns the file (required).
type CreateFileRequest struct {
	FileName     string    `json:"fileName" validate:"required"`
	FilePath     string    `json:"filePath" validate:"required"`
	FileSize     int64     `json:"fileSize" validate:"required,gte=0"`
	MimeType     string    `json:"mimeType" validate:"required"`
	Permissions  string    `json:"permissions" validate:"required"`
	LastModified time.Time `json:"lastModified" validate:"required"`
	UserId       uuid.UUID `json:"userId" validate:"required"`
}

// UpdateFileByIdRequest represents a request to update a file's information by its unique ID.
//
// Fields:
// - FileId: The unique ID of the file to update (required).
// - FileName: The new name of the file (optional).
// - FilePath: The new path of the file (optional).
// - FileSize: The new size of the file (optional, must be >= 0).
// - MimeType: The new MIME type of the file (optional).
// - Permissions: The new permissions for the file (optional).
// - LastModified: The new last modification time of the file (optional).
// - UserId: The unique ID of the user (optional).
type UpdateFileByIdRequest struct {
	FileId       uuid.UUID `json:"fileId" validate:"required"`
	FileName     string    `json:"fileName"`
	FilePath     string    `json:"filePath"`
	FileSize     int64     `json:"fileSize" validate:"gte=0"`
	MimeType     string    `json:"mimeType"`
	Permissions  string    `json:"permissions"`
	LastModified time.Time `json:"lastModified"`
	UserId       uuid.UUID `json:"userId"`
}

// DeleteFileByIdRequest represents a request to delete a file by its unique ID.
//
// Fields:
// - FileId: The unique ID of the file to delete (required).
type DeleteFileByIdRequest struct {
	FileId uuid.UUID `json:"fileId" validate:"required"`
}

// UploadRequest represents a request to upload a new file.
//
// Fields:
// - File: The file to upload (required).
// - FileName: The name of the file (required).
// - FilePath: The path where the file will be stored (required).
// - UserId: The unique ID of the user uploading the file (required).
type UploadRequest struct {
	File     multipart.File `json:"file" validate:"required"`
	FileName string         `json:"fileName" validate:"required"`
	FilePath string         `json:"filePath" validate:"required"`
	UserId   uuid.UUID      `json:"userId" validate:"required"`
}

// DownloadRequest represents a request to download a file.
//
// Fields:
// - FilePath: The path of the file to download (required).
// - UserId: The unique ID of the user requesting the file (required).
type DownloadRequest struct {
	FilePath string           `json:"filePath" validate:"required"`
	User     *model.UserModel `json:"user" validate:"required"`
}

// ListRequest represents a request to list the files in a specific directory.
//
// Fields:
// - Path: The path of the directory to list (required).
// - User: The user requesting the list (required).
type ListRequest struct {
	Path string           `json:"path" validate:"required"`
	User *model.UserModel `json:"user" validate:"required"`
}

// DeleteRequest represents a request to delete a file from a specific path.
//
// Fields:
// - Path: The path of the file to delete (required).
// - UserId: The user requesting the deletion (required).
type DeleteRequest struct {
	Path string           `json:"path" validate:"required"`
	User *model.UserModel `json:"user" validate:"required"`
}
