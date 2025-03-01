package manager

import (
	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/handler"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/middleware"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/router"
	"github.com/AtahanPoyraz/TermDrive/server/internal/repository"
	"github.com/AtahanPoyraz/TermDrive/server/internal/service"
)

// StorageManager is responsible for managing all storage-related components of the application.
// It initializes services, handlers, middleware, and routers necessary for file storage functionality.
// The manager coordinates the user service, JWT service, storage service, storage handler,
// authentication middleware, and router, enabling seamless integration of file storage functionality
// into the application.
//
// Fields:
// - Dependencies: A pointer to the Dependencies struct, which contains the application context and system context.
// - UserService: The service responsible for managing user-related operations, such as retrieving user data.
// - JwtService: The service responsible for handling JWT-related operations, such as token validation.
// - StorageService: The service responsible for managing storage-related operations, such as file uploads, retrievals, and deletions.
// - StorageHandler: The handler that processes HTTP requests related to file storage, such as uploading and retrieving files.
// - AuthMiddleware: The middleware that ensures authentication checks for protected routes.
// - StorageRouter: The router responsible for routing HTTP requests related to file storage.
type StorageManager struct {
	Dependencies *config.Dependencies

	UserService    service.UserService
	JwtService     service.JwtService
	StorageService service.StorageService

	// Handlers
	StorageHandler handler.StorageHandler

	// Middleware
	AuthMiddleware middleware.AuthMiddleware

	// Routers
	StorageRouter router.StorageRouter
}

// NewStorageManager initializes a new StorageManager instance by setting up the necessary services, handlers,
// middleware, and routers. It accepts a pointer to the Dependencies struct to obtain application and system contexts.
//
// Arguments:
// - dependencies: A pointer to the Dependencies struct, containing the app context and system context.
//
// Returns:
// - A pointer to the newly initialized StorageManager.
func NewStorageManager(dependencies *config.Dependencies) *StorageManager {
	manager := &StorageManager{Dependencies: dependencies}
	manager.InitServices()
	manager.InitHandlers()
	manager.InitMiddlewares()
	manager.InitRouters()
	return manager
}

// InitServices initializes all storage-related services, including the user service, JWT service, and storage service,
// using the dependencies provided by the StorageManager.
//
// This method configures services that are essential for storage functionality, such as user data management,
// authentication, and file operations.
func (m *StorageManager) InitServices() {
	m.UserService = service.NewUserService(m.Dependencies.AppContext.Configuration, repository.NewUserRepository(m.Dependencies.SystemContext.Database))
	m.JwtService = service.NewJwtService(m.Dependencies.AppContext.Configuration.Security.JwtSecretKey)
	m.StorageService = service.NewStorageService(m.Dependencies.AppContext.Configuration, repository.NewStorageRepository(m.Dependencies.SystemContext.Database))
}

// InitHandlers initializes the storage handler, which processes HTTP requests related to file storage,
// such as uploading, retrieving, and deleting files. The handler is configured using the storage service.
func (m *StorageManager) InitHandlers() {
	m.StorageHandler = handler.NewStorageHandler(m.Dependencies.AppContext.Configuration, m.Dependencies.AppContext.Logger, m.StorageService)
}

// InitMiddlewares initializes the authentication middleware, which checks for valid authentication tokens
// in incoming requests for storage-related operations.
func (m *StorageManager) InitMiddlewares() {
	m.AuthMiddleware = middleware.NewAuthMiddleware(m.Dependencies.AppContext.Logger, m.JwtService, m.UserService)
}

// InitRouters initializes the storage router, which is responsible for routing HTTP requests related to file storage,
// such as upload and retrieve requests, to the appropriate handler methods.
func (m *StorageManager) InitRouters() {
	m.StorageRouter = router.NewStorageRouter(m.AuthMiddleware, m.StorageHandler)
}

// RegisterRoutes registers all storage-related routes with the provided router.
// It sets up the routes required for storage operations (e.g., upload, retrieve, delete) to be served by the application.
func (m *StorageManager) RegisterRoutes() {
	m.StorageRouter.SetRouters(m.Dependencies.SystemContext.Router)
}
