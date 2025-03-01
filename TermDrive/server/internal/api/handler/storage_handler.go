package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/dto"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"github.com/AtahanPoyraz/TermDrive/server/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var storageValidator = validator.New()

// StorageHandler defines the interface for handling various file operations in the system.
// It includes methods for file creation, fetching, updating, deletion, and upload/download functionality.
type StorageHandler interface {
	FetchFileHandler(w http.ResponseWriter, r *http.Request)
	CreateFileHandler(w http.ResponseWriter, r *http.Request)
	UpdateFileHandler(w http.ResponseWriter, r *http.Request)
	DeleteFileHandler(w http.ResponseWriter, r *http.Request)
	UploadHandler(w http.ResponseWriter, r *http.Request)
	DownloadHandler(w http.ResponseWriter, r *http.Request)
	ListHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
}

// StorageHandlerImpl is the concrete implementation of the StorageHandler interface.
// It interacts with the storage service and manages file-related operations.
type StorageHandlerImpl struct {
	configuration  *config.Configuration
	logger         *log.Logger
	storageService service.StorageService
}

// NewStorageHandler creates and returns a new StorageHandlerImpl instance.
// It requires a configuration, logger, and a storage service to be provided.
func NewStorageHandler(configuration *config.Configuration, logger *log.Logger, storageService service.StorageService) StorageHandler {
	return &StorageHandlerImpl{
		configuration:  configuration,
		logger:         logger,
		storageService: storageService,
	}
}

// sendResponse sends the given response to the client with the appropriate HTTP status code.
// If encoding the response fails, an error is logged.
func (h *StorageHandlerImpl) sendResponse(response *dto.GenericResponse, w http.ResponseWriter) {
	w.WriteHeader(response.StatusCode)
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		h.logger.Printf("Failed to encode response: %v.\n", encodeErr)
	}
}

// parseUrlQuery parses the query parameters from the URL in the request and returns them as a map.
func (h *StorageHandlerImpl) parseUrlQuery(r *http.Request) (map[string]interface{}, error) {
	params := make(map[string]interface{})

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params["limit"] = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params["offset"] = offset
		}
	}

	if userIDStr := r.URL.Query().Get("userId"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid userId format")
		}
		params["userId"] = userID
	}

	if fileIDStr := r.URL.Query().Get("fileId"); fileIDStr != "" {
		fileID, err := uuid.Parse(fileIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid fileId format")
		}
		params["fileId"] = fileID
	}

	if filename := r.URL.Query().Get("fileName"); filename != "" {
		params["fileName"] = filename
	}

	if filepath := r.URL.Query().Get("filePath"); filepath != "" {
		params["filePath"] = filepath
	}

	return params, nil
}

// FetchFileHandler handles the HTTP request to fetch a file by its ID.
// It retrieves the file information based on the request parameters and responds accordingly.
func (h *StorageHandlerImpl) FetchFileHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	params, err := h.parseUrlQuery(r)
	if err != nil {
		h.logger.Printf("Query parsing error: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, w)
		return
	}

	switch {
	case params["limit"] != nil && params["offset"] != nil:
		var request dto.FetchFilesRequest
		if request.Limit, ok = params["limit"].(int); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid limit format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid limit format",
			}, w)
			return
		}

		if request.Offset, ok = params["offset"].(int); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid offset format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid offset format",
			}, w)
			return
		}

		if err := storageValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		files, err := h.storageService.FetchFiles(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch files: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve files",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("Files retrieved successfully by %s.\n", user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "Files retrieved successfully",
			Data:       files,
		}, w)

	case params["fileName"] != nil && params["userId"] != nil:
		var request dto.FetchFileByNameAndUserIdRequest
		if request.FileName, ok = params["fileName"].(string); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid fileName format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid fileName format",
			}, w)
			return
		}

		if request.UserId, ok = params["userId"].(uuid.UUID); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid userId format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid userId format",
			}, w)
			return
		}

		if err := storageValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		file, err := h.storageService.FetchFileByFileNameAndUserId(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch file: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve file",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("File (%s) retrieved successfully by %s.\n", file.FilePath, user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "File retrieved successfully",
			Data:       file,
		}, w)

	case params["filePath"] != nil && params["userId"] != nil:
		var request dto.FetchFileByPathAndUserIdRequest
		if request.FilePath, ok = params["filePath"].(string); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid filePath format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid filePath format",
			}, w)
			return
		}

		if request.UserId, ok = params["userId"].(uuid.UUID); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid userId format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid userId format",
			}, w)
			return
		}

		if err := storageValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		file, err := h.storageService.FetchFileByFilePathAndUserId(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch file: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve file",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("File (%s) retrieved successfully by %s.\n", file.FilePath, user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "File retrieved successfully",
			Data:       file,
		}, w)

	case params["userId"] != nil:
		var request dto.FetchFilesByUserIDRequest
		if request.UserId, ok = params["userId"].(uuid.UUID); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid userId format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid userId format",
			}, w)
			return
		}

		if err := storageValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		files, err := h.storageService.FetchFilesByUserId(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch files: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve files",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("Files retrieved successfully by %s.\n", user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "Files retrieved successfully",
			Data:       files,
		}, w)

	case params["fileId"] != nil:
		var request dto.FetchFileByIdRequest
		if request.FileId, ok = params["fileId"].(uuid.UUID); !ok {
			h.logger.Printf("Validation failed: %v.\n", errors.New("invalid fileId format"))
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid fileId format",
			}, w)
			return
		}

		if err := storageValidator.Struct(request); err != nil {
			h.logger.Printf("Validation failed: %v.\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid input data",
				Data:       err.Error(),
			}, w)
			return
		}

		file, err := h.storageService.FetchFileByFileId(&request)
		if err != nil {
			h.logger.Printf("Failed to fetch file: %v\n", err)
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to retrieve file",
				Data:       err.Error(),
			}, w)
			return
		}

		h.logger.Printf("File (%s) retrieved successfully by %s.\n", file.FilePath, user.Username)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusOK,
			Message:    "File retrieved successfully",
			Data:       file,
		}, w)

	default:
		err := errors.New("invalid query")
		h.logger.Printf("Query parsing error: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, w)
		return
	}
}

// CreateFileHandler handles the HTTP request to create a new file.
// It validates the input, creates the file, and responds with success or failure.
func (h *StorageHandlerImpl) CreateFileHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	var request dto.CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Printf("Invalid request payload: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := storageValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.storageService.CreateFile(&request); err != nil {
		h.logger.Printf("Failed to create file: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to create file",
			Data:       err.Error(),
		}, w)
		return
	}

	h.logger.Printf("File (%s) created successfully by %s.\n", request.FilePath, user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "File created successfully",
	}, w)
}

// UpdateFileHandler handles the HTTP request to update an existing file by its ID.
// It parses the URL parameters and the request body, validates the data, and updates the file accordingly.
func (h *StorageHandlerImpl) UpdateFileHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	params, err := h.parseUrlQuery(r)
	if err != nil {
		h.logger.Printf("Query parsing error: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, w)
		return
	}

	var request dto.UpdateFileByIdRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Printf("Invalid request payload: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request format",
			Data:       err.Error(),
		}, w)
		return
	}

	request.FileId = params["fileId"].(uuid.UUID)
	if request.FileId, ok = params["fileId"].(uuid.UUID); !ok {
		h.logger.Printf("Validation failed: %v.\n", errors.New("invalid uuid format"))
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid uuid format",
		}, w)
		return
	}

	if err := storageValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.storageService.UpdateFile(&request); err != nil {
		h.logger.Printf("Failed to update file: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to update file",
			Data:       err.Error(),
		}, w)
		return
	}

	file, err := h.storageService.FetchFileByFileId(&dto.FetchFileByIdRequest{FileId: request.FileId})
	if err != nil {
		h.logger.Printf("Failed to fetch file: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "file not found.",
			}, w)
			return
		}

		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve file",
			Data:       err.Error(),
		}, w)
		return
	}

	h.logger.Printf("File (%s) updated successfully by %s.\n", file.FilePath, user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "File updated successfully",
	}, w)
}

// DeleteFileHandler handles the HTTP request to delete a file by its ID.
// It retrieves the file using the ID, deletes it, and responds with the success or failure status.
func (h *StorageHandlerImpl) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	params, err := h.parseUrlQuery(r)
	if err != nil {
		h.logger.Printf("Query parsing error: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}, w)
		return
	}

	var request dto.DeleteFileByIdRequest
	if request.FileId, ok = params["fileId"].(uuid.UUID); !ok {
		h.logger.Printf("Validation failed: %v.\n", errors.New("invalid uuid format"))
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid uuid format",
		}, w)
		return
	}

	if err := storageValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	file, err := h.storageService.FetchFileByFileId(&dto.FetchFileByIdRequest{FileId: request.FileId})
	if err != nil {
		h.logger.Printf("Failed to fetch file: %v\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "file not found.",
			}, w)
			return
		}

		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve file",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.storageService.DeleteFile(&request); err != nil {
		h.logger.Printf("Failed to delete file: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to delete file",
			Data:       err.Error(),
		}, w)
		return
	}

	h.logger.Printf("File (%s) deleted successfully by %s.\n", file.FilePath, user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "File deleted successfully",
	}, w)
}

// UploadHandler handles the HTTP request to upload a file.
// It parses the multipart form data, validates the file path, and uploads the file.
func (h *StorageHandlerImpl) UploadHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	if err := r.ParseMultipartForm(0); err != nil {
		h.logger.Printf("Error parsing multipart form: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to parse multipart form",
		}, w)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Printf("Error retrieving the file: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Failed to retrieve file",
		}, w)
		return
	}
	defer file.Close()

	sanitizedPath := filepath.Clean(r.URL.Query().Get("path"))
	if strings.Contains(sanitizedPath, "..") || strings.HasPrefix(sanitizedPath, "/") {
		h.logger.Printf("Unauthorized directory traversal attempt by user %s. Context: %v\n", user.Username, r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid file path",
		}, w)
		return
	}

	basePath := filepath.Clean(filepath.Join(h.configuration.TermDrive.StoragePath, user.Username))
	destinationPath := filepath.Join(basePath, sanitizedPath)
	if !strings.HasPrefix(destinationPath, basePath) {
		h.logger.Printf("Unauthorized directory traversal attempt by user %s. Context: %v\n", user.Username, r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid file path",
		}, w)
		return
	}

	var request dto.UploadRequest
	request.File = file
	request.FileName = header.Filename
	request.FilePath = destinationPath
	request.UserId = user.ID
	if err := storageValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.storageService.Upload(&request); err != nil {
		h.logger.Printf("File upload failed: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "File upload failed",
		}, w)
		return
	}

	h.logger.Printf("file (%s) uploaded by %s.\n", header.Filename, user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "File uploaded successfully",
	}, w)
}

// DownloadHandler handles the HTTP request to download a file.
// It validates the file path and retrieves the file for download.
func (h *StorageHandlerImpl) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	sanitizedPath := filepath.Clean(r.URL.Query().Get("path"))
	if strings.Contains(sanitizedPath, "..") || strings.HasPrefix(sanitizedPath, "/") {
		h.logger.Printf("Unauthorized directory traversal attempt by user %s. Context: %v\n", user.Username, r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid file path",
		}, w)
		return
	}

	basePath := filepath.Clean(filepath.Join(h.configuration.TermDrive.StoragePath, user.Username))
	destinationPath := filepath.Join(basePath, sanitizedPath)
	if !strings.HasPrefix(destinationPath, basePath) {
		h.logger.Printf("Unauthorized directory traversal attempt by user %s. Context: %v\n", user.Username, r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid file path",
		}, w)
		return
	}

	var request dto.DownloadRequest
	request.FilePath = destinationPath
	request.User = user
	if err := storageValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed for user %s: %v\n", user.Username, err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	file, err := h.storageService.Download(&request)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "File not found",
			}, w)
			return
		}

		h.logger.Printf("File download failed for user %s: %v\n", user.Username, err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "File download failed",
		}, w)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		h.logger.Printf("Failed to get file (%s) stats for user %s: %v\n", filepath.Base(destinationPath), user.Username, err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "File download failed",
		}, w)
		return
	}

	h.logger.Printf("File (%s) downloaded by %s.\n", filepath.Base(destinationPath), user.Username)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(request.FilePath)))
	w.Header().Set("Content-Type", "application/octet-stream")

	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		var rangeStart int64
		if _, err := fmt.Sscanf(rangeHeader, "bytes=%d-", &rangeStart); err != nil || rangeStart >= fileStat.Size() {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		if _, err = file.Seek(rangeStart, io.SeekStart); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, fileStat.Size()-1, fileStat.Size()))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()-rangeStart))
		w.WriteHeader(http.StatusPartialContent)

	} else {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))
		w.WriteHeader(http.StatusOK)
	}

	buf := make([]byte, 5<<20)
	if _, err = io.CopyBuffer(w, file, buf); err != nil {
		h.logger.Printf("Error writing file to response for user %s: %v\n", user.Username, err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error writing file to response",
		}, w)
	}
}

// ListHandler handles the HTTP request to list the contents of a directory.
// It retrieves the list of files in the specified directory and returns the result.
func (h *StorageHandlerImpl) ListHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	sanitizedPath := filepath.Clean(r.URL.Query().Get("path"))
	if strings.Contains(sanitizedPath, "..") || strings.HasPrefix(sanitizedPath, "/") {
		h.logger.Printf("Unauthorized directory traversal attempt. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid file path",
		}, w)
		return
	}

	basePath := filepath.Clean(filepath.Join(h.configuration.TermDrive.StoragePath, user.Username))
	destinationPath := filepath.Join(basePath, sanitizedPath)
	if !strings.HasPrefix(destinationPath, basePath) {
		h.logger.Printf("Unauthorized directory traversal attempt. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid file path",
		}, w)
		return
	}

	var request dto.ListRequest
	request.Path = destinationPath
	request.User = user
	if err := storageValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	content, err := h.storageService.List(&request)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Directory not found",
				Data:       nil,
			}, w)
			return
		}

		h.logger.Printf("Directory contents listing failed: %v\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Directory contents listing failed",
		}, w)
		return
	}

	h.logger.Printf("Directory (%s) contents listed by %s.\n", filepath.Base(destinationPath), user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "Directory contents listed successfully",
		Data:       content,
	}, w)
}

// DeleteHandler handles the HTTP request to delete a file or folder.
// It checks for file/folder existence and removes it from the system.
func (h *StorageHandlerImpl) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*model.UserModel)
	if !ok || user == nil {
		h.logger.Printf("User not found or failed to retrieve. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User not found or session expired",
		}, w)
		return
	}

	sanitizedPath := filepath.Clean(r.URL.Query().Get("path"))
	if strings.Contains(sanitizedPath, "..") || strings.HasPrefix(sanitizedPath, "/") {
		h.logger.Printf("Unauthorized directory traversal attempt. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid file path",
		}, w)
		return
	}

	basePath := filepath.Clean(filepath.Join(h.configuration.TermDrive.StoragePath, user.Username))
	destinationPath := filepath.Join(basePath, sanitizedPath)
	if !strings.HasPrefix(destinationPath, basePath) {
		h.logger.Printf("Unauthorized directory traversal attempt. Context: %v\n", r.Context())
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid file path",
		}, w)
		return
	}

	var request dto.DeleteRequest
	request.Path = destinationPath
	request.User = user
	if err := storageValidator.Struct(request); err != nil {
		h.logger.Printf("Validation failed: %v.\n", err)
		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid input data",
			Data:       err.Error(),
		}, w)
		return
	}

	if err := h.storageService.Delete(&request); err != nil {
		h.logger.Printf("Contents deletion failed: %v\n", err)
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			h.sendResponse(&dto.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "File or directory not found",
			}, w)
			return
		}

		h.sendResponse(&dto.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Contents deletion failed",
		}, w)
		return
	}

	h.logger.Printf("(%s) content deleted by %s.\n", filepath.Base(destinationPath), user.Username)
	h.sendResponse(&dto.GenericResponse{
		StatusCode: http.StatusOK,
		Message:    "Content deleted successfully",
	}, w)
}
