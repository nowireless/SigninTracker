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

	r.HandleFunc("/people", app.GetPeopleCollection).Methods("GET")
	r.HandleFunc("/people", api.notImplemented).Methods("POST")

	r.HandleFunc("/people/{id}", app.GetPerson).Methods("GET")
	r.HandleFunc("/people/{id}", api.notImplemented).Methods("PATCH")
	r.HandleFunc("/people/{id}", api.notImplemented).Methods("PUT")
	r.HandleFunc("/people/{id}", api.notImplemented).Methods("DELETE")

	r.HandleFunc("/people/{id}/attendance", app.GetPersonAttendance).Methods("GET")

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
