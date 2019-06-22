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
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
)

type TeamHandlers struct {
	DB *database.Database
}

func (h *TeamHandlers) Collection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		people, err := h.DB.GetTeams()
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

		m := models.Team{}
		err = json.Unmarshal(requestBody, &m)
		if err != nil {
			log.Error(err)
			MalformedJSON(w, r)
			return
		}

		missing := tags.CheckRequiredOnCreate(m)
		if len(missing) > 0 {
			MissingRequiredOnCreate(w, r, missing)
			return
		}

		err = h.DB.CreateTeam(&m)
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

		result, err := h.DB.GetTeam(m.DatabaseID)
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

func (h *TeamHandlers) TeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	// TODO return 404 if the ID does not exist in the database

	switch r.Method {
	case http.MethodGet:
		log.Info("Getting team with id: ", id)

		// TODO add check for if the ID does not exist, add special error in database package.
		// ERRO[0708] sql: no rows in result set
		person, err := h.DB.GetTeam(int(id))
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
		person, err := h.DB.GetTeam(int(id))
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
		modified := models.Team{}
		err = json.Unmarshal(modifiedJSON, &modified)
		if err != nil {
			panic(err)
		}

		// Update person in database
		err = h.DB.UpdateTeam(modified)
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		result, err := h.DB.GetTeam(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// TODO End transaction

		writeStruct(w, http.StatusOK, result)
	case http.MethodDelete:

		// TODO Begin tranaction

		// Get person
		team, err := h.DB.GetTeam(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// Perform deletion
		err = h.DB.DeleteTeam(team)
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// TODO End transaction

		// Return deleted person
		writeStruct(w, http.StatusOK, team)
	default:
		MethodNotAllowed(w, r)
	}
}

func (h *TeamHandlers) Mentors(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func (h *TeamHandlers) RemoveMentor(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func (h *TeamHandlers) Students(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func (h *TeamHandlers) RemoveStudent(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func (h *TeamHandlers) Meetings(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}

func (h *TeamHandlers) MeetingID(w http.ResponseWriter, r *http.Request) {
	panic("TODO")
}