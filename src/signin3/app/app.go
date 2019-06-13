package app

import (
	"encoding/json"
	"net/http"
	"signin3/database"
	"signin3/models"
	"strconv"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	// Base URI for services
	APIBaseURI    string
	ClientBaseURI string

	Database database.Config
}

type App struct {
	DB *database.Database
}

func NewApp(config Config) (*App, error) {
	app := App{}

	var err error
	app.DB, err = database.Connect(config.Database)
	if err != nil {
		return nil, err
	}

	// Other initialization logic here...

	return &app, nil
}

func (app *App) GetPeopleCollection(w http.ResponseWriter, r *http.Request) {
	people, err := app.DB.GetPeople()
	if err != nil {
		log.Error(err)
		e := models.Error{Code: 500, Error: "Error accessing database"}
		app.InternalError(w, r, e)
		return
	}

	collection := map[string]interface{}{}
	collection["@uri"] = r.RequestURI
	collection["Members"] = people

	body, err := json.Marshal(collection)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

func (app *App) GetPerson(w http.ResponseWriter, r *http.Request) {
	idRaw := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idRaw, 10, 32)
	if err != nil {
		log.Error(err)
		e := models.Error{Code: http.StatusBadRequest, Error: "Unable to parse database id"}
		app.InternalError(w, r, e)
		return
	}

	log.Info("Getting attendance for person id: ", id)

	person, err := app.DB.GetPerson(int(id))
	if err != nil {
		log.Error(err)
		e := models.Error{Code: http.StatusInternalServerError, Error: "Database error"}
		app.InternalError(w, r, e)
		return
	}

	body, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

func (app *App) GetPersonAttendance(w http.ResponseWriter, r *http.Request) {
	idRaw := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idRaw, 10, 32)
	if err != nil {
		log.Error(err)
		e := models.Error{Code: http.StatusBadRequest, Error: "Unable to parse database id"}
		app.InternalError(w, r, e)
		return
	}

	log.Info("Getting attendance for person id: ", id)

	attendance, err := app.DB.GetPersonAttendances(int(id))
	if err != nil {
		log.Error(err)
		e := models.Error{Code: http.StatusInternalServerError, Error: "Database error"}
		app.InternalError(w, r, e)
		return
	}

	body, err := json.Marshal(attendance)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

func (app *App) InternalError(w http.ResponseWriter, r *http.Request, e models.Error) {
	body, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	w.Write(body)
}
