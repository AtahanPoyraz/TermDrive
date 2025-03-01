package manager

import (
	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/handler"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/middleware"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/router"
	"github.com/AtahanPoyraz/TermDrive/server/internal/repository"
	"github.com/AtahanPoyraz/TermDrive/server/internal/service"
)

// UserManager is responsible for managing all user-related components of the application.
// It handles initialization of services, handlers, middleware, and routers necessary for user operations,
// including user registration, login, and profile management.
//
// Fields:
// - Dependencies: A pointer to the Dependencies struct, which contains the application context and system context.
// - JwtService: The service responsible for handling JWT-related operations, such as token generation and validation.
// - UserService: The service responsible for managing user-related operations, such as user registration and authentication.
// - UserHandler: The handler that processes HTTP requests related to user operations, such as registration and login.
// - AuthMiddleware: The middleware that ensures authentication checks for routes that require user validation.
// - UserRouter: The router responsible for routing HTTP requests related to user operations.
type UserManager struct {
	Dependencies *config.Dependencies

	// Services
	JwtService  service.JwtService
	UserService service.UserService

	// Handlers
	UserHandler handler.UserHandler

	// Middleware
	AuthMiddleware middleware.AuthMiddleware

	// Routers
	UserRouter router.UserRouter
}

// NewUserManager initializes a new UserManager instance by setting up the necessary services, handlers,
// middleware, and routers. It accepts a pointer to the Dependencies struct to obtain application and system contexts.
//
// Arguments:
// - dependencies: A pointer to the Dependencies struct, containing the app context and system context.
//
// Returns:
// - A pointer to the newly initialized UserManager.
func NewUserManager(dependencies *config.Dependencies) *UserManager {
	manager := &UserManager{Dependencies: dependencies}
	manager.InitServices()
	manager.InitHandlers()
	manager.InitMiddlewares()
	manager.InitRouters()
	return manager
}

// InitServices initializes all user-related services, including the JWT service and the user service,
// using the dependencies provided by the UserManager.
//
// This method configures services essential for user functionality, such as user registration, authentication,
// and token management.
func (m *UserManager) InitServices() {
	m.JwtService = service.NewJwtService(m.Dependencies.AppContext.Configuration.Security.JwtSecretKey)
	m.UserService = service.NewUserService(m.Dependencies.AppContext.Configuration, repository.NewUserRepository(m.Dependencies.SystemContext.Database))
}

// InitHandlers initializes the user handler, which processes HTTP requests related to user operations,
// such as registration, login, and profile management. The handler is configured using the user service.
func (m *UserManager) InitHandlers() {
	m.UserHandler = handler.NewUserHandler(m.Dependencies.AppContext.Configuration, m.Dependencies.AppContext.Logger, m.UserService)
}

// InitMiddlewares initializes the authentication middleware, which checks for valid authentication tokens
// in incoming requests for user-related operations. This ensures that routes requiring authentication are protected.
func (m *UserManager) InitMiddlewares() {
	m.AuthMiddleware = middleware.NewAuthMiddleware(m.Dependencies.AppContext.Logger, m.JwtService, m.UserService)
}

// InitRouters initializes the user router, which is responsible for routing HTTP requests related to user operations,
// such as user registration, login, and profile management, to the appropriate handler methods.
func (m *UserManager) InitRouters() {
	m.UserRouter = router.NewUserRouter(m.AuthMiddleware, m.UserHandler)
}

// RegisterRoutes registers all user-related routes with the provided router.
// It sets up the routes required for user operations (e.g., registration, login, profile management) to be served by the application.
func (m *UserManager) RegisterRoutes() {
	m.UserRouter.SetRouters(m.Dependencies.SystemContext.Router)
}
