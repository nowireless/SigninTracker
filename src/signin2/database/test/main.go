package main

import (
	"database/sql"
	"signin2/database"

	"github.com/jackc/pgx/pgtype"

	"github.com/davecgh/go-spew/spew"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := database.Config{
		User:     "signin",
		Password: "foobar",
		Host:     "localhost",
		Port:     5432,
		Database: "signin",
	}
	db, err := database.Connect(config)
	if err != nil {
		panic(err)
	}

	db.DB.Ping()

	log.Info("Ready")

	// conn, err := stdlib.AcquireConn(db.DB)
	// if err != nil {
	// 	panic(err)
	// }

	row := db.DBX.QueryRowx(`
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
		WHERE people.personid = 4
		GROUP BY people.personid;
	`)

	type personDB struct {
		PersonID    int
		CheckInID   int
		FirstName   string
		LastName    string
		Email       sql.NullString
		Phone       sql.NullString
		SchoolID    sql.NullString
		SchoolEmail sql.NullString

		MentorOf  pgtype.Int4Array `db:"mentor_of"`
		StudentOf pgtype.Int4Array `db:"student_of"`
		Parents   pgtype.Int4Array `db:"parents"`
		ParentOf  pgtype.Int4Array `db:"parent_of"`
	}

	m := personDB{}

	// row.Scan(
	// 	&m.PersonID,
	// 	&m.CheckInID,
	// 	&m.FirstName,
	// 	&m.LastName,
	// 	&m.Email,
	// 	&m.Phone,
	// 	&m.SchoolID,
	// 	&m.SchoolEmail,
	// 	&m.MentorOf,
	// 	&m.StudentOf,
	// 	&m.Parents,
	// 	&m.ParentOf,
	// )

	row.StructScan(&m)

	spew.Dump(m)
}
