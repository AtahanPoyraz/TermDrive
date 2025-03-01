package manager

import (
	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/handler"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/middleware"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/router"
	"github.com/AtahanPoyraz/TermDrive/server/internal/repository"
	"github.com/AtahanPoyraz/TermDrive/server/internal/service"
)

// AuthManager is responsible for managing all authentication-related components of the application.
// It initializes services, handlers, middleware, and routers necessary for authentication.
// The manager coordinates the JWT service, user service, authentication service, authentication handler,
// authentication middleware, and router, allowing easy integration of authentication functionality
// into the application.
//
// Fields:
// - Dependencies: A pointer to the Dependencies struct, which contains the application context and system context.
// - JwtService: The service responsible for handling JWT-related operations, such as token generation and validation.
// - UserService: The service responsible for managing user data and user-related operations.
// - AuthService: The service that handles authentication logic, such as login and token verification.
// - AuthHandler: The handler that processes HTTP requests related to authentication.
// - AuthMiddleware: The middleware that enforces authentication checks on protected routes.
// - AuthRouter: The router responsible for routing authentication-related HTTP requests.
type AuthManager struct {
	Dependencies *config.Dependencies

	// Services
	JwtService  service.JwtService
	UserService service.UserService
	AuthService service.AuthService

	// Handlers
	AuthHandler handler.AuthHandler

	// Middleware
	AuthMiddleware middleware.AuthMiddleware

	// Routers
	AuthRouter router.AuthRouter
}

// NewAuthManager initializes a new AuthManager instance by setting up the necessary services, handlers,
// middleware, and routers. It accepts a pointer to the Dependencies struct to obtain application and system contexts.
//
// Arguments:
// - dependencies: A pointer to the Dependencies struct, containing the app context and system context.
//
// Returns:
// - A pointer to the newly initialized AuthManager.
func NewAuthManager(dependencies *config.Dependencies) *AuthManager {
	manager := &AuthManager{Dependencies: dependencies}
	manager.InitServices()
	manager.InitHandlers()
	manager.InitMiddlewares()
	manager.InitRouters()
	return manager
}

// InitServices initializes all authentication-related services, including the JWT service, user service, and
// authentication service, using the dependencies provided by the AuthManager.
//
// This method configures services that are essential for authentication functionality, such as token generation
// and user verification.
func (m *AuthManager) InitServices() {
	m.JwtService = service.NewJwtService(m.Dependencies.AppContext.Configuration.Security.JwtSecretKey)
	m.UserService = service.NewUserService(m.Dependencies.AppContext.Configuration, repository.NewUserRepository(m.Dependencies.SystemContext.Database))
	m.AuthService = service.NewAuthService(m.Dependencies.AppContext.Configuration, repository.NewUserRepository(m.Dependencies.SystemContext.Database))
}

// InitHandlers initializes the authentication handler, which processes HTTP requests related to authentication,
// such as login, registration, and token verification. The handler is configured using the JWT service,
// user service, and authentication service.
func (m *AuthManager) InitHandlers() {
	m.AuthHandler = handler.NewAuthHandler(m.Dependencies.AppContext.Configuration, m.Dependencies.AppContext.Logger, m.AuthService, m.JwtService)
}

// InitMiddlewares initializes the authentication middleware, which is responsible for handling authentication-related
// logic, such as checking the presence of valid tokens in incoming requests. This middleware is necessary for
// protecting routes that require user authentication.
func (m *AuthManager) InitMiddlewares() {
	m.AuthMiddleware = middleware.NewAuthMiddleware(m.Dependencies.AppContext.Logger, m.JwtService, m.UserService)
}

// InitRouters initializes the authentication router, which is responsible for routing HTTP requests related
// to authentication, such as login and token refresh, to the appropriate handler methods.
func (m *AuthManager) InitRouters() {
	m.AuthRouter = router.NewAuthRouter(m.AuthMiddleware, m.AuthHandler)
}

// RegisterRoutes registers all authentication-related routes with the provided router.
// It sets up the routes required for authentication (e.g., login, token refresh) to be served by the application.
func (m *AuthManager) RegisterRoutes() {
	m.AuthRouter.SetRouters(m.Dependencies.SystemContext.Router)
}
