// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"signin2/restapi/operations"
)

//go:generate swagger generate server --target ..\..\signin2 --name Signin --spec ..\swagger.yaml

func configureFlags(api *operations.SigninAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SigninAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.DeleteStudentsIDHandler == nil {
		api.DeleteStudentsIDHandler = operations.DeleteStudentsIDHandlerFunc(func(params operations.DeleteStudentsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation .DeleteStudentsID has not yet been implemented")
		})
	}
	if api.GetStudentsHandler == nil {
		api.GetStudentsHandler = operations.GetStudentsHandlerFunc(func(params operations.GetStudentsParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetStudents has not yet been implemented")
		})
	}
	if api.GetStudentsIDHandler == nil {
		api.GetStudentsIDHandler = operations.GetStudentsIDHandlerFunc(func(params operations.GetStudentsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation .GetStudentsID has not yet been implemented")
		})
	}
	if api.PatchStudentsIDHandler == nil {
		api.PatchStudentsIDHandler = operations.PatchStudentsIDHandlerFunc(func(params operations.PatchStudentsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation .PatchStudentsID has not yet been implemented")
		})
	}
	if api.PostStudentsHandler == nil {
		api.PostStudentsHandler = operations.PostStudentsHandlerFunc(func(params operations.PostStudentsParams) middleware.Responder {
			return middleware.NotImplemented("operation .PostStudents has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
