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

	// TODO return 404 if the ID does not exist in the database

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
	// Add person as mentor to team

	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Extract team id from request body
	type requestStruct struct {
		TeamID int
	}
	request := requestStruct{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		log.Error(err)
		MalformedJSON(w, r)
		return
	}

	// Exec query
	_, err = h.DB.DB.Exec(
		`INSERT INTO mentors(personid, teamid) VALUES ($1, $2)`,
		id, request.TeamID,
	)
	if err != nil {
		// TODO the following should maybe move to database package?
		// Create a custom error struct for friendlier error handle
		// Does the team id exist?
		// 23503 - foreign_key_violation
		if pgxErr, ok := err.(pgx.PgError); ok && pgxErr.Code == "23503" {
			e := models.Error{Code: http.StatusBadRequest, Error: pgxErr.Message}
			writeError(w, r, e)
			return
		}
		InternalError(w, r, err, "Database Error")
		return

	}

	// Return nothing
	w.WriteHeader(http.StatusNoContent)
}

func (h *PersonHandlers) MentorOfID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		MethodNotAllowed(w, r)
		return
	}

	// Remove mentor from team

	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse person id")
		return
	}

	teamId, err := parseInt(vars["tid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse team id")
		return
	}

	result, err := h.DB.DB.Exec(
		"DELETE FROM mentors WHERE personid = $1 AND  teamid = $2",
		id, teamId,
	)
	if err != nil {
		InternalError(w, r, err, "Database Error")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows == 0 {
		// No entry was deleted
		e := models.Error{Code: http.StatusBadRequest, Error: "No matching person and team mentor relationship"}
		writeError(w, r, e)
		return
	}

	// Return nothing
	w.WriteHeader(http.StatusNoContent)
}

func (h *PersonHandlers) StudentOf(w http.ResponseWriter, r *http.Request) {
	// Add person as student to team

	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Extract team id from request body
	type requestStruct struct {
		TeamID int
	}
	request := requestStruct{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		log.Error(err)
		MalformedJSON(w, r)
		return
	}

	// Exec query
	_, err = h.DB.DB.Exec(
		`INSERT INTO students(personid, teamid) VALUES ($1, $2)`,
		id, request.TeamID,
	)
	if err != nil {
		// TODO the following should maybe move to database package?
		// Create a custom error struct for friendlier error handle
		// Does the team id exist?
		// 23503 - foreign_key_violation
		if pgxErr, ok := err.(pgx.PgError); ok && pgxErr.Code == "23503" {
			e := models.Error{Code: http.StatusBadRequest, Error: pgxErr.Message}
			writeError(w, r, e)
			return
		}
		InternalError(w, r, err, "Database Error")
		return

	}

	// Return nothing
	w.WriteHeader(http.StatusNoContent)
}

func (h *PersonHandlers) StudentOfID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		MethodNotAllowed(w, r)
		return
	}

	// Remove mentor from team

	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse person id")
		return
	}

	teamID, err := parseInt(vars["tid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse team id")
		return
	}

	log.Info("Deleting mentor relationship: PersonID ", id, " TeamID ", teamID)

	result, err := h.DB.DB.Exec(
		"DELETE FROM students WHERE personid = $1 AND teamid = $2",
		id, teamID,
	)
	if err != nil {
		InternalError(w, r, err, "Database Error")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows == 0 {
		// No entry was deleted
		e := models.Error{Code: http.StatusBadRequest, Error: "No matching person and team student relationship"}
		writeError(w, r, e)
		return
	}

	// Return nothing
	w.WriteHeader(http.StatusNoContent)
}

func (h *PersonHandlers) Parents(w http.ResponseWriter, r *http.Request) {
	// Add person as mentor to team

	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Extract team id from request body
	type requestStruct struct {
		ParentID int
		Relation string
	}
	request := requestStruct{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		log.Error(err)
		MalformedJSON(w, r)
		return
	}

	// Exec query
	_, err = h.DB.DB.Exec(
		`INSERT INTO parents(studentid, parentid, relation) VALUES ($1, $2, $3)`,
		id, request.ParentID, request.Relation,
	)
	if err != nil {
		// TODO the following should maybe move to database package?
		// Create a custom error struct for friendlier error handle
		badRequestErrors := map[string]bool{}
		badRequestErrors["23503"] = true // foreign_key_violation - The parent ID does not exist
		badRequestErrors["22P02"] = true // invalid_text_representation - The parent relation is invalid

		if pgxErr, ok := err.(pgx.PgError); ok && badRequestErrors[pgxErr.Code] {
			e := models.Error{Code: http.StatusBadRequest, Error: pgxErr.Message}
			writeError(w, r, e)
			return
		}
		InternalError(w, r, err, "Database Error")
		return

	}

	// Return nothing
	w.WriteHeader(http.StatusNoContent)
}

func (h *PersonHandlers) ParentsID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	parentID, err := parseInt(vars["pid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse parent id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		NotImplemented(w, r)
	case http.MethodPatch:
		NotImplemented(w, r)
	case http.MethodDelete:
		log.Info("Deleting parent relationship: ParentID ", parentID, " StudentID ", id)
		result, err := h.DB.DB.Exec(
			"DELETE FROM parents WHERE parentid = $1 AND studentid = $2",
			parentID, id,
		)
		if err != nil {
			InternalError(w, r, err, "Database Error")
			return
		}

		rows, err := result.RowsAffected()
		if err != nil {
			panic(err)
		}
		if rows == 0 {
			// No entry was deleted
			e := models.Error{Code: http.StatusBadRequest, Error: "No matching person and team student relationship"}
			writeError(w, r, e)
			return
		}

		// Return nothing
		// TODO: Return deleted relationship?
		w.WriteHeader(http.StatusNoContent)
	default:
		MethodNotAllowed(w, r)
	}
}

func (h *PersonHandlers) ParentOf(w http.ResponseWriter, r *http.Request) {
	// Add person as mentor to team

	if r.Method != http.MethodPost {
		MethodNotAllowed(w, r)
		return
	}
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	// Extract team id from request body
	type requestStruct struct {
		StudentID      int
		ParentRelation string
	}
	request := requestStruct{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		log.Error(err)
		MalformedJSON(w, r)
		return
	}

	// Exec query
	_, err = h.DB.DB.Exec(
		`INSERT INTO parents(parentid, studentid, relation) VALUES ($1, $2, $3)`,
		id, request.StudentID, request.ParentRelation,
	)
	if err != nil {
		// TODO the following should maybe move to database package?
		// Create a custom error struct for friendlier error handle
		badRequestErrors := map[string]bool{}
		badRequestErrors["23503"] = true // foreign_key_violation - The parent ID does not exist
		badRequestErrors["22P02"] = true // invalid_text_representation - The parent relation is invalid

		if pgxErr, ok := err.(pgx.PgError); ok && badRequestErrors[pgxErr.Code] {
			e := models.Error{Code: http.StatusBadRequest, Error: pgxErr.Message}
			writeError(w, r, e)
			return
		}
		InternalError(w, r, err, "Database Error")
		return
	}

	// Return nothing
	w.WriteHeader(http.StatusNoContent)
}

func (h *PersonHandlers) ParentOfID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	studentID, err := parseInt(vars["sid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse student id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		NotImplemented(w, r)
	case http.MethodPatch:
		NotImplemented(w, r)
	case http.MethodDelete:
		log.Info("Deleting parent relationship: ParentID ", id, " StudentID ", studentID)
		result, err := h.DB.DB.Exec(
			"DELETE FROM parents WHERE parentid = $1 AND studentid = $2",
			id, studentID,
		)
		if err != nil {
			InternalError(w, r, err, "Database Error")
			return
		}

		rows, err := result.RowsAffected()
		if err != nil {
			panic(err)
		}
		if rows == 0 {
			// No entry was deleted
			e := models.Error{Code: http.StatusBadRequest, Error: "No matching person and team student relationship"}
			writeError(w, r, e)
			return
		}

		// Return nothing
		// TODO: Return deleted relationship?
		w.WriteHeader(http.StatusNoContent)
	default:
		MethodNotAllowed(w, r)
	}
}
