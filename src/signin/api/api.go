package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"signin/app"
	"signin/database"
	"signin/models"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type API struct {
	Config Config
	App    *app.App
}

func New(a *app.App) (*API, error) {
	return &API{
		App: a,
		Config: Config{
			Port: 8080,
		},
	}, nil
}

func (api *API) Init(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		response := map[string]interface{}{}
		response["API"] = "Hello World!"

		body, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			log.Panic(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}).Methods("GET")

	r.HandleFunc("/students", func(w http.ResponseWriter, req *http.Request) {
		students, err := api.App.GetStudents()
		if err != nil {
			log.Panic(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Add sorting support

		body, err := json.MarshalIndent(students, "", "  ")
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}).Methods("GET")

	r.HandleFunc("/students", func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Panic(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Body: " + string(body))

		var student models.Student
		err = json.Unmarshal(body, &student)
		if err != nil {
			log.Panic(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check for requried parameters
		if student.FirstName == "" || student.LastName == "" || student.ID == "" {
			log.Warn("Bad request missing required parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = api.App.CreateStudent(student)
		if err == database.ErrorIDExists {
			// Duplicate Student
			// TODO figure out correct code
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else if err != nil {
			log.Panic(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}).Methods("POST")

}
