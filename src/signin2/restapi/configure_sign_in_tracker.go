// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"signin2/restapi/operations"
	"signin2/restapi/operations/meetings"
	"signin2/restapi/operations/people"
	"signin2/restapi/operations/teams"
)

//go:generate swagger generate server --target ..\..\signin2 --name SignInTracker --spec ..\swagger.yaml

func configureFlags(api *operations.SignInTrackerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SignInTrackerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.MeetingsDeleteMeetingsIDHandler == nil {
		api.MeetingsDeleteMeetingsIDHandler = meetings.DeleteMeetingsIDHandlerFunc(func(params meetings.DeleteMeetingsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation meetings.DeleteMeetingsID has not yet been implemented")
		})
	}
	if api.PeopleDeletePeopleIDHandler == nil {
		api.PeopleDeletePeopleIDHandler = people.DeletePeopleIDHandlerFunc(func(params people.DeletePeopleIDParams) middleware.Responder {
			return middleware.NotImplemented("operation people.DeletePeopleID has not yet been implemented")
		})
	}
	if api.TeamsDeleteTeamsIDHandler == nil {
		api.TeamsDeleteTeamsIDHandler = teams.DeleteTeamsIDHandlerFunc(func(params teams.DeleteTeamsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation teams.DeleteTeamsID has not yet been implemented")
		})
	}
	if api.MeetingsGetMeetingsHandler == nil {
		api.MeetingsGetMeetingsHandler = meetings.GetMeetingsHandlerFunc(func(params meetings.GetMeetingsParams) middleware.Responder {
			return middleware.NotImplemented("operation meetings.GetMeetings has not yet been implemented")
		})
	}
	if api.MeetingsGetMeetingsIDHandler == nil {
		api.MeetingsGetMeetingsIDHandler = meetings.GetMeetingsIDHandlerFunc(func(params meetings.GetMeetingsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation meetings.GetMeetingsID has not yet been implemented")
		})
	}
	if api.PeopleGetPeopleHandler == nil {
		api.PeopleGetPeopleHandler = people.GetPeopleHandlerFunc(func(params people.GetPeopleParams) middleware.Responder {
			return middleware.NotImplemented("operation people.GetPeople has not yet been implemented")
		})
	}
	if api.PeopleGetPeopleIDHandler == nil {
		api.PeopleGetPeopleIDHandler = people.GetPeopleIDHandlerFunc(func(params people.GetPeopleIDParams) middleware.Responder {
			return middleware.NotImplemented("operation people.GetPeopleID has not yet been implemented")
		})
	}
	if api.TeamsGetTeamsHandler == nil {
		api.TeamsGetTeamsHandler = teams.GetTeamsHandlerFunc(func(params teams.GetTeamsParams) middleware.Responder {
			return middleware.NotImplemented("operation teams.GetTeams has not yet been implemented")
		})
	}
	if api.TeamsGetTeamsIDHandler == nil {
		api.TeamsGetTeamsIDHandler = teams.GetTeamsIDHandlerFunc(func(params teams.GetTeamsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation teams.GetTeamsID has not yet been implemented")
		})
	}
	if api.MeetingsPatchMeetingsIDHandler == nil {
		api.MeetingsPatchMeetingsIDHandler = meetings.PatchMeetingsIDHandlerFunc(func(params meetings.PatchMeetingsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation meetings.PatchMeetingsID has not yet been implemented")
		})
	}
	if api.PeoplePatchPeopleIDHandler == nil {
		api.PeoplePatchPeopleIDHandler = people.PatchPeopleIDHandlerFunc(func(params people.PatchPeopleIDParams) middleware.Responder {
			return middleware.NotImplemented("operation people.PatchPeopleID has not yet been implemented")
		})
	}
	if api.TeamsPatchTeamsIDHandler == nil {
		api.TeamsPatchTeamsIDHandler = teams.PatchTeamsIDHandlerFunc(func(params teams.PatchTeamsIDParams) middleware.Responder {
			return middleware.NotImplemented("operation teams.PatchTeamsID has not yet been implemented")
		})
	}
	if api.MeetingsPostMeetingsHandler == nil {
		api.MeetingsPostMeetingsHandler = meetings.PostMeetingsHandlerFunc(func(params meetings.PostMeetingsParams) middleware.Responder {
			return middleware.NotImplemented("operation meetings.PostMeetings has not yet been implemented")
		})
	}
	if api.PeoplePostPeopleHandler == nil {
		api.PeoplePostPeopleHandler = people.PostPeopleHandlerFunc(func(params people.PostPeopleParams) middleware.Responder {
			return middleware.NotImplemented("operation people.PostPeople has not yet been implemented")
		})
	}
	if api.TeamsPostTeamsHandler == nil {
		api.TeamsPostTeamsHandler = teams.PostTeamsHandlerFunc(func(params teams.PostTeamsParams) middleware.Responder {
			return middleware.NotImplemented("operation teams.PostTeams has not yet been implemented")
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
