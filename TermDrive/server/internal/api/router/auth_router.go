package router

import (
	"net/http"

	"github.com/AtahanPoyraz/TermDrive/server/internal/api/handler"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/middleware"
	"github.com/gorilla/mux"
)

// AuthRouter defines the interface for setting up authentication-related routes.
type AuthRouter interface {
	SetRouters(r *mux.Router)
}

// AuthRouterImpl is the concrete implementation of AuthRouter that configures authentication routes.
type AuthRouterImpl struct {
	authMiddleware middleware.AuthMiddleware
	authHandler    handler.AuthHandler
}

// NewAuthRouter initializes and returns an instance of AuthRouterImpl.
func NewAuthRouter(authMiddleware middleware.AuthMiddleware, authHandler handler.AuthHandler) AuthRouter {
	return &AuthRouterImpl{authMiddleware: authMiddleware, authHandler: authHandler}
}

// SetRouters configures the authentication-related routes in the provided router.
// This function sets up the following routes:
// - POST /api/v1/auth/sign-in for user sign-in
// - POST /api/v1/auth/sign-up for user registration
// - GET /api/v1/auth/me for fetching details of the authenticated user
// It also attaches the necessary middleware (e.g., authorization checks).
func (r *AuthRouterImpl) SetRouters(router *mux.Router) {
	signInRouter := router.Methods(http.MethodPost).Subrouter()
	signInRouter.HandleFunc("/api/v1/auth/sign-in", r.authHandler.SignInHandler)

	signUpRouter := router.Methods(http.MethodPost).Subrouter()
	signUpRouter.HandleFunc("/api/v1/auth/sign-up", r.authHandler.SignUpHandler)

	meRouter := router.Methods(http.MethodGet).Subrouter()
	meRouter.HandleFunc("/api/v1/auth/me", r.authHandler.MeHandler)
	meRouter.Use(r.authMiddleware.RequireAuthorize)
}
