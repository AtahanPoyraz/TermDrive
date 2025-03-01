package router

import (
	"net/http"

	"github.com/AtahanPoyraz/TermDrive/server/internal/api/handler"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/middleware"
	"github.com/gorilla/mux"
)

// StorageRouter defines the interface for setting up routes related to storage operations.
type StorageRouter interface {
	SetRouters(r *mux.Router)
}

// StorageRouterImpl is the concrete implementation of StorageRouter that configures the routes for storage operations.
type StorageRouterImpl struct {
	authMiddleware middleware.AuthMiddleware
	storageHandler handler.StorageHandler
}

// NewStorageRouter initializes and returns an instance of StorageRouterImpl.
func NewStorageRouter(authMiddleware middleware.AuthMiddleware, storageHandler handler.StorageHandler) StorageRouter {
	return &StorageRouterImpl{authMiddleware: authMiddleware, storageHandler: storageHandler}
}

// SetRouters configures the routes for the storage API in the provided router.
// It includes routes for the following operations:
// - Fetching a file: GET /api/v1/admin/storage/get
// - Creating a file: POST /api/v1/admin/storage/create
// - Updating a file: PATCH /api/admin/v1/storage/update
// - Deleting a file: DELETE /api/admin/v1/storage/delete
// - Uploading a file: POST /api/v1/storage/upload
// - Downloading a file: POST /api/v1/storage/download
// - Listing files: GET /api/v1/storage/list
// - Deleting a file: DELETE /api/v1/storage/delete
// Authorization middleware is applied to all routes, and admin role checks are enforced on some routes.
func (r *StorageRouterImpl) SetRouters(router *mux.Router) {
	fetchFileRouter := router.Methods(http.MethodGet).Subrouter()
	fetchFileRouter.HandleFunc("/api/v1/admin/storage/get", r.storageHandler.FetchFileHandler)
	fetchFileRouter.Use(r.authMiddleware.RequireAuthorize)
	fetchFileRouter.Use(r.authMiddleware.RequiredAdminRole)

	createFileRouter := router.Methods(http.MethodPost).Subrouter()
	createFileRouter.HandleFunc("/api/v1/admin/storage/create", r.storageHandler.CreateFileHandler)
	createFileRouter.Use(r.authMiddleware.RequireAuthorize)
	createFileRouter.Use(r.authMiddleware.RequiredAdminRole)

	updateFileRouter := router.Methods(http.MethodPatch).Subrouter()
	updateFileRouter.HandleFunc("/api/v1/admin/storage/update", r.storageHandler.UpdateFileHandler)
	updateFileRouter.Use(r.authMiddleware.RequireAuthorize)
	updateFileRouter.Use(r.authMiddleware.RequiredAdminRole)

	deleteFileRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteFileRouter.HandleFunc("/api/v1/admin/storage/delete", r.storageHandler.DeleteFileHandler)
	deleteFileRouter.Use(r.authMiddleware.RequireAuthorize)
	deleteFileRouter.Use(r.authMiddleware.RequiredAdminRole)

	uploadRouter := router.Methods(http.MethodPost).Subrouter()
	uploadRouter.HandleFunc("/api/v1/storage/upload", r.storageHandler.UploadHandler)
	uploadRouter.Use(r.authMiddleware.RequireAuthorize)

	downloadRouter := router.Methods(http.MethodGet).Subrouter()
	downloadRouter.HandleFunc("/api/v1/storage/download", r.storageHandler.DownloadHandler)
	downloadRouter.Use(r.authMiddleware.RequireAuthorize)

	listRouter := router.Methods(http.MethodGet).Subrouter()
	listRouter.HandleFunc("/api/v1/storage/list", r.storageHandler.ListHandler)
	listRouter.Use(r.authMiddleware.RequireAuthorize)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/api/v1/storage/delete", r.storageHandler.DeleteHandler)
	deleteRouter.Use(r.authMiddleware.RequireAuthorize)
}
