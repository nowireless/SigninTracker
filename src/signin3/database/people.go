package database

import (
	"signin3/database/internal"
	"signin3/models"

	log "github.com/sirupsen/logrus"
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

func (db *Database) CreatePerson(mPerson *models.Person) error {
	person := internal.NewPerson(mPerson)

	row := db.DB.QueryRowx(`
		INSERT INTO people(checkinid, firstname, lastname, email, phone, schoolemail, schoolid)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING personid;`,
		person.CheckInID,
		person.FirstName,
		person.LastName,
		person.Email,
		person.Phone,
		person.SchoolEmail,
		person.SchoolID,
	)

	if err := row.Err(); err != nil {
		return err
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	mPerson.DatabaseID = id

	return nil
}

func (db *Database) UpdatePerson(model models.Person) error {
	person := internal.NewPerson(&model)

	log.Info("Updating person with ID: ", person.PersonID)

	_, err := db.DB.Exec(`
		UPDATE people SET
			checkinid    = $1,
			firstname    = $2,
			lastname     = $3,
			email        = $4,
			phone        = $5,
			schoolid     = $6,
			schoolemail  = $7
		WHERE personid = $8;
		`,
		person.CheckInID,
		person.FirstName,
		person.LastName,
		person.Email,
		person.Phone,
		person.SchoolID,
		person.SchoolEmail,
		person.PersonID,
	)

	return err
}

func (db *Database) DeletePerson(model *models.Person) error {
	// TODO: Some how mark/(remove database id) to indicate that the model no longer represents a row in the db
	person := internal.NewPerson(model)
	_, err := db.DB.Exec("DELETE FROM PEOPLE WHERE personid = $1", person.PersonID)
	return err
}
