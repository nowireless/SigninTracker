package app

import (
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

type MeetingHandlers struct {
	DB *database.Database
}

func (h *MeetingHandlers) Collection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		people, err := h.DB.GetMeetings()
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

		m := models.Meeting{}
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

		err = h.DB.CreateMeeting(&m)
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

		result, err := h.DB.GetMeeting(m.DatabaseID)
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

func (h *MeetingHandlers) MeetingID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	// TODO return 404 if the ID does not exist in the database

	switch r.Method {
	case http.MethodGet:
		log.Info("Getting meeting with id: ", id)

		// TODO add check for if the ID does not exist, add special error in database package.
		// ERRO[0708] sql: no rows in result set
		person, err := h.DB.GetMeeting(int(id))
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
		person, err := h.DB.GetMeeting(int(id))
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
		modified := models.Meeting{}
		err = json.Unmarshal(modifiedJSON, &modified)
		if err != nil {
			panic(err)
		}

		// Update person in database
		err = h.DB.UpdateMeeting(modified)
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		result, err := h.DB.GetMeeting(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// TODO End transaction

		writeStruct(w, http.StatusOK, result)
	case http.MethodDelete:

		// TODO Begin tranaction

		// Get person
		meeting, err := h.DB.GetMeeting(int(id))
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// Perform deletion
		err = h.DB.DeleteMeeting(meeting)
		if err != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		// TODO End transaction

		// Return deleted person
		writeStruct(w, http.StatusOK, meeting)
	default:
		MethodNotAllowed(w, r)
	}
}

func (h *MeetingHandlers) Attendance(w http.ResponseWriter, r *http.Request) {
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

	log.Info("Getting attendance for meeting id: ", id)

	attendance, err := h.DB.GetMeetingAttendance(id)
	if err != nil {
		InternalError(w, r, err, "Database error")
		return
	}

	writeStruct(w, http.StatusOK, attendance)
}

func (h *MeetingHandlers) Teams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		rows, err := h.DB.DB.Queryx(`
			SELECT teamid, kind
			FROM meetings
				INNER JOIN team_meetings m2 on meetings.meetingid = m2.meetingid
			WHERE m2.meetingid = $1;
		`, id)

		if rows.Err() != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		type teamMeeting struct {
			TeamID int
			Kind   string
		}

		results := []models.TeamMeeting{}
		for rows.Next() {
			tm := teamMeeting{}
			err := rows.StructScan(&tm)
			if err != nil {
				panic(err)
			}

			results = append(results, models.TeamMeeting{
				Team: models.Link{URI: fmt.Sprintf("/teams/%d", tm.TeamID)},
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
			TeamID int
			Kind   string
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
			request.TeamID, id, request.Kind,
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

func (h *MeetingHandlers) TeamID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	teamID, err := parseInt(vars["id"])
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
		log.Info("Deleting team meeting relationship: TeamID ", teamID, " MeetingID ", id)
		result, err := h.DB.DB.Exec(
			"DELETE FROM team_meetings WHERE teamid = $1 AND meetingid = $2",
			teamID, id,
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

func (h *MeetingHandlers) Commitments(w http.ResponseWriter, r *http.Request) {
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
		`INSERT INTO commitments(personid, meetingid) VALUES ($1, $2)`,
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

func (h *MeetingHandlers) RemoveCommitment(w http.ResponseWriter, r *http.Request) {
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

	personID, err := parseInt(vars["pid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse team id")
		return
	}

	log.Info("Deleting commitment relationship: PersonID ", personID, " MeetingID ", id)

	result, err := h.DB.DB.Exec(
		"DELETE FROM commitments WHERE personid = $1 AND meetingid = $2",
		personID, id,
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
		e := models.Error{Code: http.StatusBadRequest, Error: "No matching commitment relationship"}
		writeError(w, r, e)
		return
	}

	// Return nothing
	w.WriteHeader(http.StatusNoContent)
}

func (h *MeetingHandlers) SignIns(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		rows, err := h.DB.DB.Queryx(`
			SELECT personid,
				meetingid,
				intime
			FROM signed_in
			WHERE meetingid = $1;
		`, id)

		if rows.Err() != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		type signIn struct {
			PersonID int
			InTime   string
		}

		results := []models.SignIn{}
		for rows.Next() {
			si := signIn{}
			err := rows.StructScan(&si)
			if err != nil {
				panic(err)
			}

			results = append(results, models.SignIn{
				Person: &models.Link{URI: fmt.Sprintf("/people/%d", si.PersonID)},
				InTime: si.InTime,
			})
		}

		writeStruct(w, http.StatusOK, results)
	case http.MethodPost:
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		// Extract team id from request body
		type requestStruct struct {
			PersonID int
			InTime   string
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
			`INSERT INTO signed_in(personid, meetingid, intime)
				VALUES ($1, $2, $3);`,
			request.PersonID, id, request.InTime,
		)
		if err != nil {
			// TODO the following should maybe move to database package?
			// Create a custom error struct for friendlier error handle
			badRequestErrors := map[string]bool{}
			badRequestErrors["23503"] = true // foreign_key_violation - The parent ID does not exist
			// badRequestErrors["22P02"] = true // invalid_text_representation - The parent relation is invalid

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

func (h *MeetingHandlers) SignInsID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	personID, err := parseInt(vars["pid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse person id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		NotImplemented(w, r)
	case http.MethodPatch:
		NotImplemented(w, r)
	case http.MethodDelete:
		log.Info("Deleting signed in relationship: PersonID ", personID, " MeetingID ", id)

		result, err := h.DB.DB.Exec(
			"DELETE FROM signed_in WHERE personid = $1 AND meetingid = $2",
			personID, id,
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
			e := models.Error{Code: http.StatusBadRequest, Error: "No matching sign in relationship"}
			writeError(w, r, e)
			return
		}

		// Return nothing
		w.WriteHeader(http.StatusNoContent)
	default:
		MethodNotAllowed(w, r)
	}
}

func (h *MeetingHandlers) SignOuts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		rows, err := h.DB.DB.Queryx(`
			SELECT personid,
				meetingid,
				outtime
			FROM signed_out
			WHERE meetingid = $1;
		`, id)

		if rows.Err() != nil {
			InternalError(w, r, err, "Database error")
			return
		}

		type signOut struct {
			PersonID int
			OutTime  string
		}

		results := []models.SignOut{}
		for rows.Next() {
			so := signOut{}
			err := rows.StructScan(&so)
			if err != nil {
				panic(err)
			}

			results = append(results, models.SignOut{
				Person:  &models.Link{URI: fmt.Sprintf("/people/%d", so.PersonID)},
				OutTime: so.OutTime,
			})
		}

		writeStruct(w, http.StatusOK, results)
	case http.MethodPost:
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		// Extract team id from request body
		type requestStruct struct {
			PersonID int
			InTime   string
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
			`INSERT INTO signed_in(personid, meetingid, intime)
				VALUES ($1, $2, $3);`,
			request.PersonID, id, request.InTime,
		)
		if err != nil {
			// TODO the following should maybe move to database package?
			// Create a custom error struct for friendlier error handle
			badRequestErrors := map[string]bool{}
			badRequestErrors["23503"] = true // foreign_key_violation - The parent ID does not exist
			// badRequestErrors["22P02"] = true // invalid_text_representation - The parent relation is invalid

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

func (h *MeetingHandlers) SignOutsID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse id")
		return
	}

	personID, err := parseInt(vars["pid"])
	if err != nil {
		InternalError(w, r, err, "Unable to parse person id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		NotImplemented(w, r)
	case http.MethodPatch:
		NotImplemented(w, r)
	case http.MethodDelete:
		log.Info("Deleting signed out relationship: PersonID ", personID, " MeetingID ", id)

		result, err := h.DB.DB.Exec(
			"DELETE FROM signed_out WHERE personid = $1 AND meetingid = $2",
			personID, id,
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
			e := models.Error{Code: http.StatusBadRequest, Error: "No matching sign out relationship"}
			writeError(w, r, e)
			return
		}

		// Return nothing
		w.WriteHeader(http.StatusNoContent)
	default:
		MethodNotAllowed(w, r)
	}
}
