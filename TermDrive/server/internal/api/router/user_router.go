package router

import (
	"net/http"

	"github.com/AtahanPoyraz/TermDrive/server/internal/api/handler"
	"github.com/AtahanPoyraz/TermDrive/server/internal/api/middleware"
	"github.com/gorilla/mux"
)

// UserRouter defines the interface for setting up routes related to user operations.
type UserRouter interface {
	SetRouters(r *mux.Router)
}

// UserRouterImpl is the concrete implementation of UserRouter that configures the routes for user operations.
type UserRouterImpl struct {
	authMiddleware middleware.AuthMiddleware
	userHandler    handler.UserHandler
}

// NewUserRouter initializes and returns an instance of UserRouterImpl.
func NewUserRouter(authMiddleware middleware.AuthMiddleware, userHandler handler.UserHandler) UserRouter {
	return &UserRouterImpl{authMiddleware: authMiddleware, userHandler: userHandler}
}

// SetRouters configures the routes for the user API in the provided router.
// It includes routes for the following operations:
// - Fetching a user: GET /api/v1/admin/user/get
// - Creating a user: POST /api/v1/admin/user/create
// - Updating a user: PATCH /api/v1/admin/user/update
// - Deleting a user: DELETE /api/v1/admin/user/delete
// Authorization middleware is applied to all routes, and admin role checks are enforced on all routes.
func (r *UserRouterImpl) SetRouters(router *mux.Router) {
	fetchUserRouter := router.Methods(http.MethodGet).Subrouter()
	fetchUserRouter.HandleFunc("/api/v1/admin/user/get", r.userHandler.FetchUserHandler)
	fetchUserRouter.Use(r.authMiddleware.RequireAuthorize)
	fetchUserRouter.Use(r.authMiddleware.RequiredAdminRole)

	createUserRouter := router.Methods(http.MethodPost).Subrouter()
	createUserRouter.HandleFunc("/api/v1/admin/user/create", r.userHandler.CreateUserHandler)
	createUserRouter.Use(r.authMiddleware.RequireAuthorize)
	createUserRouter.Use(r.authMiddleware.RequiredAdminRole)

	updateUserRouter := router.Methods(http.MethodPatch).Subrouter()
	updateUserRouter.HandleFunc("/api/v1/admin/user/update", r.userHandler.UpdateUserByIdHandler)
	updateUserRouter.Use(r.authMiddleware.RequireAuthorize)
	updateUserRouter.Use(r.authMiddleware.RequiredAdminRole)

	deleteUserRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteUserRouter.HandleFunc("/api/v1/admin/user/delete", r.userHandler.DeleteUserByIdHandler)
	deleteUserRouter.Use(r.authMiddleware.RequireAuthorize)
	deleteUserRouter.Use(r.authMiddleware.RequiredAdminRole)
}
