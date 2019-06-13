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
	row.StructScan(&person)
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
