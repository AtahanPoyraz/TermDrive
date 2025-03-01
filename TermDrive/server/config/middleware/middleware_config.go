package middleware

import (
	"net/http"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// SetConfiguration configures the provided router with middleware for handling CORS,
// as well as custom handlers for HTTP 404 (Not Found) and 405 (Method Not Allowed) errors.
// It uses the configuration provided to set the allowed origins, methods, and headers for CORS.
// The method also sets custom responses for unsupported routes and methods, returning HTTP 401 (Unauthorized)
// for both cases.
//
// Arguments:
// - configuration: A pointer to the Configuration struct that contains the CORS settings.
// - router: The router instance (mux.Router) to apply the middleware to.
//
// Returns:
// - A pointer to the modified router with the applied middleware and error handlers.

func SetConfiguration(configuration *config.Configuration, router *mux.Router) *mux.Router {
	router.Use(
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedOrigins(configuration.Cors.AllowedOrigins),
			handlers.AllowedMethods(configuration.Cors.AllowedMethods),
			handlers.AllowedHeaders(configuration.Cors.AllowedHeaders),
		),
	)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusUnauthorized) })
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusUnauthorized) })

	return router
}
