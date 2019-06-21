package app

import (
	"encoding/json"
	"net/http"
	"signin3/database"
	"signin3/models"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	// Base URI for services
	APIBaseURI    string
	ClientBaseURI string

	Database database.Config
}

type App struct {
	DB       *database.Database
	People   PersonHandlers
	Meetings MeetingHandlers
}

func NewApp(config Config) (*App, error) {
	app := App{}

	var err error
	app.DB, err = database.Connect(config.Database)
	if err != nil {
		return nil, err
	}

	// Other initialization logic here...
	app.People = PersonHandlers{DB: app.DB}
	app.Meetings = MeetingHandlers{DB: app.DB}

	return &app, nil
}

func InternalError(w http.ResponseWriter, r *http.Request, err error, errorMsg string) {
	log.Error(err)
	e := models.Error{
		Code:  http.StatusInternalServerError,
		Error: errorMsg,
	}
	writeError(w, r, e)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	e := models.Error{
		Code:  http.StatusNotImplemented,
		Error: "Reqest not implemented",
	}
	writeError(w, r, e)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	e := models.Error{
		Code:  http.StatusMethodNotAllowed,
		Error: "Method not allowed: " + r.Method,
	}
	writeError(w, r, e)
}

func MalformedJSON(w http.ResponseWriter, r *http.Request) {
	e := models.Error{
		Code:  http.StatusBadRequest,
		Error: "Malformed JSON in request body",
	}
	writeError(w, r, e)
}

func MissingRequiredOnCreate(w http.ResponseWriter, r *http.Request, missing []string) {
	e := models.Error{
		Code:  http.StatusBadRequest,
		Error: "Missing required on create fields: " + strings.Join(missing, ", "),
	}
	writeError(w, r, e)
}

func writeError(w http.ResponseWriter, r *http.Request, e models.Error) {
	body, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	writeJSON(w, e.Code, body)
}

func writeStruct(w http.ResponseWriter, status int, v interface{}) {
	body, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	writeJSON(w, status, body)
}

func writeJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func parseInt(idStr string) (int, error) {
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
