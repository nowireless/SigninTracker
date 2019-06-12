package database

import (
	"database/sql"
	"fmt"
	"signin2/models"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func Connect(config Config) (*Database, error) {
	// See: https://github.com/jackc/pgx/blob/master/stdlib/sql.go
	connectionStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	log.Info("Connection String: ", connectionStr)

	db, err := sql.Open("pgx", connectionStr)
	if err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(db, "pgx")

	log.Info("Pinging Database")
	db.Ping()

	return &Database{DB: db, DBX: dbx}, nil
}

type Database struct {
	DB  *sql.DB
	DBX *sqlx.DB
}

func (db *Database) GetPeople() ([]*models.Person, error) {
	rows, err := db.DBX.Queryx(`
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

	people := []*models.Person{}
	for rows.Next() {
		person := Person{}
		rows.StructScan(&person)
		people = append(people, NewModelPerson(&person))
	}

	return people, nil
}

// func (db *Database) GetPerson(id int) (*models.Person, error) {

// }

// func (db *Database) GetMeetings() ([]models.Meeting, error) {

// }

// func (db *Database) GetMeeting(id int) (*models.Meeting, error) {

// }

// func (tb *Database) GetTeams() ([]models.Meeting, error) {

// }

// func (db *Database) GetTeam(id int) (*models.Team, error) {

// }
