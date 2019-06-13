package api

import (
	"encoding/json"
	"net/http"
	"signin3/app"
	"signin3/models"

	"github.com/gorilla/mux"
)

type Config struct {
	App app.Config
}

type API struct {
	App    *app.App
	Config Config
}

func NewAPI(config Config) *API {
	api := API{}
	api.Config = config
	return &api
}

func (api *API) Initialize(r *mux.Router) error {
	// Initialize Application
	app, err := app.NewApp(api.Config.App)
	if err != nil {
		return err
	}

	// Register handlers
	// Service Root
	r.HandleFunc("", api.notImplemented).Methods("GET")

	r.HandleFunc("/people", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/people/{id}", api.notImplemented).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/people/{id}/attendance", app.GetPersonAttendance).Methods("GET")
	r.HandleFunc("/people/{id}/mentors", api.notImplemented).Methods("POST")
	r.HandleFunc("/people/{id}/mentors/{tid}", api.notImplemented).Methods("DELETE")
	r.HandleFunc("/people/{id}/studentOf", api.notImplemented).Methods("POST")
	r.HandleFunc("/people/{id}/studentOf/{tid}", api.notImplemented).Methods("DELETE")
	r.HandleFunc("/people/{id}/parents", api.notImplemented).Methods("POST")
	r.HandleFunc("/people/{id}/parents/{pid}", api.notImplemented).Methods("DELETE")
	r.HandleFunc("/people/{id}/parentsOf", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/people/{id}/parentsOf/{sid}", api.notImplemented).Methods("DELETE")

	r.HandleFunc("/meetings", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/meetings/{id}", api.notImplemented).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/meetings/{id}/teams", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/meetings/{id}/teams/{tid}", api.notImplemented).Methods("DELETE")
	r.HandleFunc("/meetings/{id}/commitments", api.notImplemented).Methods("POST")
	r.HandleFunc("/meetings/{id}/commitments/{id}", api.notImplemented).Methods("DELETE")
	r.HandleFunc("/meetings/{id}/signins", api.notImplemented).Methods("POST")
	r.HandleFunc("/meetings/{id}/signins/{id}", api.notImplemented).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/meetings/{id}/signouts", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/meetings/{id}/signouts/{id}", api.notImplemented).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/meetings/{id}/attendance", api.notImplemented).Methods("GET")

	r.HandleFunc("/teams", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/teams/{id}", api.notImplemented).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/teams/{id}/mentors", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/teams/{id}/mentors/{mid}", api.notImplemented).Methods("DELETE")
	r.HandleFunc("/teams/{id}/students", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/teams/{id}/students/{sid}", api.notImplemented).Methods("DELETE")
	r.HandleFunc("/teams/{id}/meetings", api.notImplemented).Methods("GET", "POST")
	r.HandleFunc("/teams/{id}/meetings/{mid}", api.notImplemented).Methods("GET", "PATCH", "DELETE")

	return nil
}

func (api *API) notImplemented(w http.ResponseWriter, r *http.Request) {
	e := models.Error{Code: 501, Error: "Reqest not implemented"}
	body, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	w.Write(body)

}
