package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"signin3/database"
	"signin3/models"
	"signin3/tags"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type PersonHandlers struct {
	DB *database.Database
}

func (h *PersonHandlers) Collection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		people, err := h.DB.GetPeople()
		if err != nil {
			InternalError(w, r, err, "Error accessing database")
			return
		}

		collection := map[string]interface{}{}
		collection["@uri"] = r.RequestURI
		collection["Members"] = people

		writeStruct(w, http.StatusOK, collection)
	case http.MethodPost:
		// TODO: Somwhere check headers of the POST request for JSON
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		p := models.Person{}
		err = json.Unmarshal(requestBody, &p)
		if err != nil {
			log.Error(err)
			MalformedJSON(w, r)
			return
		}

		missing := tags.CheckRequiredOnCreate(p)
		if len(missing) > 0 {
			MissingRequiredOnCreate(w, r, missing)
			return
		}

		err = h.DB.CreatePerson(&p)
		if err != nil {
			InternalError(w, r, err, "Database Error")
			return
		}

		result, err := h.DB.GetPerson(p.DatabaseID)
		if err != nil {
			InternalError(w, r, err, "Database Error")
			return
		}

		// Does a header need to created to store the new resources location?
		writeStruct(w, http.StatusCreated, result)
	default:
		MethodNotAllowed(w, r)
	}
}

func (h *PersonHandlers) PersonID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Info("Getting person with id: ", id)

		person, err := h.DB.GetPerson(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		writeStruct(w, http.StatusOK, person)
	case http.MethodPatch:
		NotImplemented(w, r)
	case http.MethodDelete:
		NotImplemented(w, r)
	default:
		MethodNotAllowed(w, r)
	}
}

func (h *PersonHandlers) Attendance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotAllowed(w, r)
		return
	}

	// URI variables
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse database id")
		return
	}

	log.Info("Getting attendance for person id: ", id)

	attendance, err := h.DB.GetPersonAttendances(id)
	if err != nil {
		InternalError(w, r, err, "Database error")
		return
	}

	writeStruct(w, http.StatusOK, attendance)
}

func (h *PersonHandlers) MentorOf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}

	NotImplemented(w, r)
}

func (h *PersonHandlers) MentorOfID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		MethodNotAllowed(w, r)
		return
	}

	NotImplemented(w, r)
}

func (h *PersonHandlers) StudentOf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}

	NotImplemented(w, r)
}

func (h *PersonHandlers) StudentOfID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		MethodNotAllowed(w, r)
		return
	}

	NotImplemented(w, r)
}

func (h *PersonHandlers) Parents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}

	NotImplemented(w, r)
}

func (h *PersonHandlers) ParentsID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		MethodNotAllowed(w, r)
		return
	}

	NotImplemented(w, r)
}

func (h *PersonHandlers) ParentOf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}

	NotImplemented(w, r)
}

func (h *PersonHandlers) ParentOfID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		MethodNotAllowed(w, r)
		return
	}

	NotImplemented(w, r)
}

// func (app *App) getPeopleCollection(w http.ResponseWriter, r *http.Request) {
// 	people, err := app.DB.GetPeople()
// 	if err != nil {
// 		log.Error(err)
// 		e := models.Error{Code: 500, Error: "Error accessing database"}
// 		app.InternalError(w, r, e)
// 		return
// 	}

// 	collection := map[string]interface{}{}
// 	collection["@uri"] = r.RequestURI
// 	collection["Members"] = people

// 	body, err := json.Marshal(collection)
// 	if err != nil {
// 		panic(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(body)
// }

// func (app *App) getPerson(w http.ResponseWriter, r *http.Request) {
// 	idRaw := mux.Vars(r)["id"]
// 	id, err := strconv.ParseInt(idRaw, 10, 32)
// 	if err != nil {
// 		log.Error(err)
// 		e := models.Error{Code: http.StatusBadRequest, Error: "Unable to parse database id"}
// 		app.InternalError(w, r, e)
// 		return
// 	}

// 	log.Info("Getting person with id: ", id)

// 	person, err := app.DB.GetPerson(int(id))
// 	if err != nil {
// 		log.Error(err)
// 		e := models.Error{Code: http.StatusInternalServerError, Error: "Database error"}
// 		app.InternalError(w, r, e)
// 		return
// 	}

// 	body, err := json.Marshal(person)
// 	if err != nil {
// 		panic(err)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write(body)
// }

// func (app *App) getPersonAttendance(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, err := getInt(vars["id"])
// 	if err != nil {
// 		InternalError(w, r, err, "Unable to parse database id")
// 		return
// 	}

// 	log.Info("Getting attendance for person id: ", id)

// 	attendance, err := app.DB.GetPersonAttendances(id)
// 	if err != nil {
// 		InternalError(w, r, err, "Database error")
// 		return
// 	}

// 	body, err := json.Marshal(attendance)
// 	if err != nil {
// 		panic(err)
// 	}

// 	writeJSON(w, http.StatusOK, body)
// }
