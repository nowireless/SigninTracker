package database

import (
	"signin3/database/internal"
	"signin3/models"
)

func (db *Database) GetPersonAttendances(personID int) ([]models.MeetingAttendance, error) {
	// What is a lateral join? See the following
	// https://medium.com/kkempin/postgresqls-lateral-join-bfd6bd0199df
	rows, err := db.DB.Queryx(`
		SELECT signed_in.meetingid, signed_in.personid, signed_in.intime, check_outs.outtime
		FROM signed_in
			LEFT JOIN LATERAL
		(SELECT meetingid, personid, outtime
		FROM signed_out
		WHERE signed_in.meetingid = signed_out.meetingid
		AND signed_in.personid = signed_out.personid) check_outs ON TRUE
		WHERE signed_in.personid = $1;
	`, personID)

	if err != nil {
		return nil, err
	}

	results := []models.MeetingAttendance{}
	for rows.Next() {
		result := internal.MeetingAttendance{}
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}

		results = append(results, result.Model())
	}

	return results, nil
}

func (db *Database) GetMeetingAttendance(meetingID int) ([]models.MeetingAttendance, error) {
	// What is a lateral join? See the following
	// https://medium.com/kkempin/postgresqls-lateral-join-bfd6bd0199df
	rows, err := db.DB.Queryx(`
		SELECT signed_in.meetingid, signed_in.personid, signed_in.intime, check_outs.outtime
		FROM signed_in
			LEFT JOIN LATERAL
		(SELECT meetingid, personid, outtime
		FROM signed_out
		WHERE signed_in.meetingid = signed_out.meetingid
		AND signed_in.personid = signed_out.personid) check_outs ON TRUE
		WHERE signed_in.meetingid = $1;
	`, meetingID)

	if err != nil {
		return nil, err
	}

	results := []models.MeetingAttendance{}
	for rows.Next() {
		result := internal.MeetingAttendance{}
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}

		results = append(results, result.Model())
	}

	return results, nil
}
