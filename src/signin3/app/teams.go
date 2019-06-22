package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
			badRequestErrors := map[string]bool{}
			badRequestErrors["23503"] = true // foreign_key_violation - The team ID does not exist
			badRequestErrors["22P02"] = true // TODO
			badRequestErrors["23505"] = true // TODO

			if pgxErr, ok := err.(pgx.PgError); ok && badRequestErrors[pgxErr.Code] {
				e := models.Error{Code: http.StatusBadRequest, Error: pgxErr.Message}
				writeError(w, r, e)
				return
			}
			InternalError(w, r, err, "Database Error")
			return
		}

		log.Info("Created Team id: ", m.DatabaseID)

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
		if err == sql.ErrNoRows {
			NotFound(w, r)
			return
		} else if err != nil {
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
		if err == sql.ErrNoRows {
			NotFound(w, r)
			return
		} else if err != nil {
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
		if err == sql.ErrNoRows {
			NotFound(w, r)
			return
		} else if err != nil {
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
		PersonID int
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
		request.PersonID, id,
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

func (h *TeamHandlers) RemoveMentor(w http.ResponseWriter, r *http.Request) {
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

	meetingID, err := parseInt(vars["mid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse team id")
		return
	}

	result, err := h.DB.DB.Exec(
		"DELETE FROM mentors WHERE personid = $1 AND  teamid = $2",
		meetingID, id,
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

func (h *TeamHandlers) Students(w http.ResponseWriter, r *http.Request) {
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
		PersonID int
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
		request.PersonID, id,
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

func (h *TeamHandlers) RemoveStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		MethodNotAllowed(w, r)
		return
	}

	// Remove student from team

	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse person id")
		return
	}

	studentID, err := parseInt(vars["sid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse student id")
		return
	}

	log.Info("Deleting mentor relationship: PersonID ", studentID, " TeamID ", id)

	result, err := h.DB.DB.Exec(
		"DELETE FROM students WHERE personid = $1 AND teamid = $2",
		studentID, id,
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

func (h *TeamHandlers) Meetings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		rows, err := h.DB.DB.Queryx(`		
			SELECT
				meetingid,
				kind
			FROM (SELECT
				meetings.date,
				meetings.starttime,
				m2.meetingid,
				kind
				FROM meetings
				INNER JOIN team_meetings m2 on meetings.meetingid = m2.meetingid
				WHERE m2.teamid = $1
				ORDER BY meetings.date, meetings.starttime
			) AS _`, id)

		if rows.Err() != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		type teamMeeting struct {
			MeetingID int
			Kind      string
		}

		results := []models.TeamMeeting{}
		for rows.Next() {
			tm := teamMeeting{}
			err := rows.StructScan(&tm)
			if err != nil {
				panic(err)
			}

			results = append(results, models.TeamMeeting{
				Team: models.Link{URI: fmt.Sprintf("/meetings/%d", tm.MeetingID)},
				Kind: tm.Kind,
			})
		}

		writeStruct(w, http.StatusOK, results)
	case http.MethodPost:
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		type requestStruct struct {
			MeetingID int
			Kind      string
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
			`INSERT INTO team_meetings(teamid, meetingid, kind) VALUES ($1, $2, $3)`,
			id, request.MeetingID, request.Kind,
		)

		if err != nil {
			// TODO the following should maybe move to database package?
			// Create a custom error struct for friendlier error handle
			badRequestErrors := map[string]bool{}
			badRequestErrors["23503"] = true // foreign_key_violation - The team ID does not exist

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

	default:
		MethodNotAllowed(w, r)
	}
}

func (h *TeamHandlers) MeetingID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	meetingID, err := parseInt(vars["mid"])
	if err != nil {
		InternalError(w, r, err, "Unable to team id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		NotImplemented(w, r)
	case http.MethodPatch:
		NotImplemented(w, r)
	case http.MethodDelete:
		log.Info("Deleting team meeting relationship: TeamID ", id, " MeetingID ", meetingID)
		result, err := h.DB.DB.Exec(
			"DELETE FROM team_meetings WHERE teamid = $1 AND meetingid = $2",
			id, meetingID,
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
			e := models.Error{Code: http.StatusBadRequest, Error: "No matching meeting and team relationship"}
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
