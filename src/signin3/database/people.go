package database

import (
	"signin3/database/internal"
	"signin3/models"
)

func (db *Database) GetPeople() ([]models.Person, error) {
	rows, err := db.DB.Queryx(`
		SELECT
			people.personid,
			people.checkinid,
			people.firstname,
			people.lastname,
			people.email,
			people.phone,
			people.schoolid,
			people.schoolemail,
			array_agg(DISTINCT m2.teamid)    FILTER (WHERE m2.teamid    IS NOT NULL) as mentor_of,
			array_agg(DISTINCT s2.teamid)    FILTER (WHERE s2.teamid    IS NOT NULL) as student_of,
			array_agg(DISTINCT p.parentid)   FILTER (WHERE p.parentid   IS NOT NULL) as parents,
			array_agg(DISTINCT p2.studentid) FILTER (WHERE p2.studentid IS NOT NULL) as parent_of
		FROM people
			FULL OUTER JOIN mentors m2 on people.personid = m2.personid
			FULL OUTER JOIN students s2 on people.personid = s2.personid
			FULL OUTER JOIN parents p on people.personid = p.studentid
			FULL OUTER JOIN parents p2 on people.personid = p2.parentid
		GROUP BY people.personid;
	`)
	if err != nil {
		return nil, err
	}

	people := []models.Person{}
	for rows.Next() {
		person := internal.Person{}
		rows.StructScan(&person)
		people = append(people, person.Model())
	}

	return people, nil
}

func (db *Database) GetPerson(id int) (*models.Person, error) {
	row := db.DB.QueryRowx(`
		SELECT
			people.personid,
			people.checkinid,
			people.firstname,
			people.lastname,
			people.email,
			people.phone,
			people.schoolid,
			people.schoolemail,
			array_agg(DISTINCT m2.teamid)    FILTER (WHERE m2.teamid    IS NOT NULL) as mentor_of,
			array_agg(DISTINCT s2.teamid)    FILTER (WHERE s2.teamid    IS NOT NULL) as student_of,
			array_agg(DISTINCT p.parentid)   FILTER (WHERE p.parentid   IS NOT NULL) as parents,
			array_agg(DISTINCT p2.studentid) FILTER (WHERE p2.studentid IS NOT NULL) as parent_of
		FROM people
			FULL OUTER JOIN mentors m2 on people.personid = m2.personid
			FULL OUTER JOIN students s2 on people.personid = s2.personid
			FULL OUTER JOIN parents p on people.personid = p.studentid
			FULL OUTER JOIN parents p2 on people.personid = p2.parentid
		WHERE people.personid = $1
		GROUP BY people.personid;
	`, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	person := internal.Person{}
	err := row.StructScan(&person)
	if err != nil {
		return nil, err
	}
	m := person.Model()
	return &m, nil

}

func (db *Database) CreatePerson(*models.Person) error {
	panic("TODO")
}

func (db *Database) UpdatePerson(*models.Person) error {
	panic("TODO")
}

func (db *Database) DeletePerson(*models.Person) error {
	panic("TODO")
}

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
