package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"signin3/database"
	"signin3/models"
	"signin3/tags"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/jackc/pgx"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type PersonHandlers struct {
	DB *database.Database
}

func (h *PersonHandlers) Collection(w http.ResponseWriter, r *http.Request) {
	ctx := handlerContext{w: w, r: r, db: h.DB}
	ctx.unmarshalJSON = func(raw []byte) (models.Model, error) {
		p := models.Person{}
		if err := json.Unmarshal(raw, &p); err != nil {
			return nil, err
		}
		return &p, nil
	}
	ctx.marshalJSON = func(obj models.Model) []byte {
		objJSON, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		return objJSON
	}

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
			// TODO the following should maybe move to database package?
			// Create a custom error struct for friendlier error handle
			if pgxErr, ok := err.(pgx.PgError); ok && pgxErr.Code == "23505" {
				e := models.Error{Code: http.StatusBadRequest, Error: pgxErr.Message}
				writeError(w, r, e)
				return
			}
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

		// TODO add check for if the ID does not exist, add special error in database package.
		// ERRO[0708] sql: no rows in result set
		person, err := h.DB.GetPerson(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		writeStruct(w, http.StatusOK, person)
	case http.MethodPatch:
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		// Check for malformed JSON
		patch := map[string]interface{}{}
		err = json.Unmarshal(requestBody, &patch)
		if err != nil {
			log.Error(err)
			MalformedJSON(w, r)
			return
		}

		// Check for known fields
		// TODO

		// Check the patch to see if it modifying readonly fields
		violations := tags.CheckPatchReadonly(models.Person{}, patch)
		if len(violations) > 0 {
			e := models.Error{
				Code:  http.StatusBadRequest,
				Error: "Setting read only fields: " + strings.Join(violations, ", "),
			}
			writeError(w, r, e)
			return
		}

		// TODO Begin transaction
		// Fetch current values of person
		person, err := h.DB.GetPerson(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// Convert source struct to JSON
		originalJSON, err := json.Marshal(person)
		if err != nil {
			panic(err)
		}

		// Apply JSON merge patch
		modifiedJSON, err := jsonpatch.MergePatch(originalJSON, requestBody)
		if err != nil {
			log.Error(err)
			// TODO Roll back transaction
			e := models.Error{Code: http.StatusBadRequest, Error: "Unable to apply merge patch"}
			writeError(w, r, e)
			return
		}

		// Convert dest JSON to person model
		modified := models.Person{}
		err = json.Unmarshal(modifiedJSON, &modified)
		if err != nil {
			panic(err)
		}

		// Update person in database
		err = h.DB.UpdatePerson(modified)
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		result, err := h.DB.GetPerson(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// TODO End transaction

		writeStruct(w, http.StatusOK, result)
	case http.MethodDelete:

		// TODO Begin tranaction

		// Get person
		person, err := h.DB.GetPerson(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// Perform deletion
		err = h.DB.DeletePerson(person)
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// TODO End transaction

		// Return deleted person
		writeStruct(w, http.StatusOK, person)
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
