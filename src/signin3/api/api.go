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

	r.HandleFunc("/people", app.People.Collection).Methods("GET", "POST")
	r.HandleFunc("/people/{id}", app.People.PersonID).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/people/{id}/attendance", app.People.Attendance).Methods("GET")
	r.HandleFunc("/people/{id}/mentorOf", app.People.MentorOf).Methods("POST")
	r.HandleFunc("/people/{id}/mentorOf/{tid}", app.People.MentorOfID).Methods("DELETE")
	r.HandleFunc("/people/{id}/studentOf", app.People.StudentOf).Methods("POST")
	r.HandleFunc("/people/{id}/studentOf/{tid}", app.People.StudentOfID).Methods("DELETE")
	r.HandleFunc("/people/{id}/parents", app.People.Parents).Methods("POST")
	r.HandleFunc("/people/{id}/parents/{pid}", app.People.ParentsID).Methods("DELETE")
	r.HandleFunc("/people/{id}/parentsOf", app.People.ParentOf).Methods("GET", "POST")
	r.HandleFunc("/people/{id}/parentsOf/{sid}", app.People.ParentOfID).Methods("DELETE")

	r.HandleFunc("/meetings", app.Meetings.Collection).Methods("GET", "POST")
	r.HandleFunc("/meetings/{id}", app.Meetings.MeetingID).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/meetings/{id}/teams", app.Meetings.Teams).Methods("GET", "POST")
	r.HandleFunc("/meetings/{id}/teams/{tid}", app.Meetings.RemoveTeam).Methods("DELETE")
	r.HandleFunc("/meetings/{id}/commitments", app.Meetings.Commitments).Methods("POST")
	r.HandleFunc("/meetings/{id}/commitments/{pid}", app.Meetings.RemoveCommitment).Methods("DELETE")
	r.HandleFunc("/meetings/{id}/signins", app.Meetings.SignIns).Methods("POST")
	r.HandleFunc("/meetings/{id}/signins/{pid}", app.Meetings.SignInsID).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/meetings/{id}/signouts", app.Meetings.SignOuts).Methods("GET", "POST")
	r.HandleFunc("/meetings/{id}/signouts/{pid}", app.Meetings.SignOutsID).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/meetings/{id}/attendance", app.Meetings.Attendance).Methods("GET")

	r.HandleFunc("/teams", app.Teams.Collection).Methods("GET", "POST")
	r.HandleFunc("/teams/{id}", app.Teams.TeamID).Methods("GET", "PATCH", "DELETE")
	r.HandleFunc("/teams/{id}/mentors", app.Teams.Mentors).Methods("GET", "POST")
	r.HandleFunc("/teams/{id}/mentors/{mid}", app.Teams.RemoveMentor).Methods("DELETE")
	r.HandleFunc("/teams/{id}/students", app.Teams.Students).Methods("GET", "POST")
	r.HandleFunc("/teams/{id}/students/{sid}", app.Teams.RemoveStudent).Methods("DELETE")
	r.HandleFunc("/teams/{id}/meetings", app.Teams.Meetings).Methods("GET", "POST")
	r.HandleFunc("/teams/{id}/meetings/{mid}", app.Teams.MeetingID).Methods("GET", "PATCH", "DELETE")

	return nil
}

func (api *API) notImplemented(w http.ResponseWriter, r *http.Request) {
	e := models.Error{Code: 501, Error: "Request not implemented"}
	body, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	w.Write(body)

}
