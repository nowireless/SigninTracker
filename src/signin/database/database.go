package database

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
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

	fmt.Println("Connection String: ", connectionStr)

	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		return nil, err
	}

	fmt.Println("Pinging Database")
	db.Ping()

	return &Database{DB: db}, nil
}

type Database struct {
	DB *sqlx.DB
}

func (*Database) String() string {
	// TODO: This is a hack to stop spew from dumping the contents of database connection
	return ""
}

func (db *Database) GetAllPeople() ([]Person, error) {
	tx := db.DB.MustBegin()
	defer tx.Commit()

	return db.GetAllPeopleTx(tx)
}

func (db *Database) GetAllPeopleTx(tx *sqlx.Tx) ([]Person, error) {
	stmt, err := db.DB.Preparex(`
		SELECT personid,
			checkinid,
			firstname,
			lastname,
			email,
			phone,
			schoolemail,
			schoolid
		FROM people
		ORDER BY personid;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx()
	if err != nil {
		return nil, err
	}

	people := []Person{}
	for rows.Next() {
		person := Person{db: db}
		err := rows.StructScan(&person)
		if err != nil {
			return nil, err
		}

		people = append(people, person)
	}

	return people, nil

}

func (db *Database) GetPerson(id int) (*Person, error) {
	tx := db.DB.MustBegin()
	defer tx.Commit()

	return db.GetPersonTX(tx, id)
}

func (db *Database) GetPersonTX(tx *sqlx.Tx, id int) (*Person, error) {
	// SELECT personid, checkinid, firstname, lastname, email, phone, schoolemail, schoolid
	// FROM people
	// WHERE personid = 0;

	stmt, err := tx.Preparex(`
		SELECT personid, checkinid, firstname, lastname, email, phone, schoolemail, schoolid
		FROM people
		WHERE personid = $1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowx(id)
	if row.Err() != nil {
		return nil, err
	}

	person := Person{db: db}
	err = row.StructScan(&person)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (db *Database) GetAllMeetings() ([]Meeting, error) {
	tx := db.DB.MustBegin()
	defer tx.Commit()
	return db.GetAllMeetingsTx(tx)
}

func (db *Database) GetAllMeetingsTx(tx *sqlx.Tx) ([]Meeting, error) {
	stmt, err := tx.Preparex(`
		SELECT
			meetingid,
			date,
			starttime,
			endtime,
			location
		FROM meetings
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx()
	if err != nil {
		return nil, err
	}

	meetings := []Meeting{}

	for rows.Next() {
		meeting := Meeting{db: db}
		err = rows.StructScan(&meeting)
		if err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (db *Database) GetMeeting(id int) (*Meeting, error) {
	tx := db.DB.MustBegin()
	defer tx.Commit()
	return db.GetMeetingTx(tx, id)
}

func (db *Database) GetMeetingTx(tx *sqlx.Tx, id int) (*Meeting, error) {
	stmt, err := tx.Preparex(`
		SELECT meetingid, date, starttime, endtime, location
		FROM meetings
		WHERE meetingid = $1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowx(id)
	if row.Err() != nil {
		return nil, err
	}

	meeting := Meeting{db: db}
	err = row.StructScan(&meeting)
	if err != nil {
		return nil, err
	}

	return &meeting, nil
}

func (db *Database) GetAllTeams() ([]Team, error) {
	tx := db.DB.MustBegin()
	defer tx.Commit()

	return db.GetAllTeamsTx(tx)
}

func (db *Database) GetAllTeamsTx(tx *sqlx.Tx) ([]Team, error) {
	stmt, err := db.DB.Preparex(`
	SELECT
		teamid,
		compeition,
		number,
		name
	FROM teams
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx()
	if err != nil {
		return nil, err
	}

	teams := []Team{}
	for rows.Next() {
		team := Team{db: db}
		err = rows.StructScan(&team)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (db *Database) GetTeam(id int) (*Team, error) {
	tx := db.DB.MustBegin()
	defer tx.Commit()
	return db.GetTeamTx(tx, id)
}

func (db *Database) GetTeamTx(tx *sqlx.Tx, id int) (*Team, error) {
	stmt, err := db.DB.Preparex(`
	SELECT
		teamid,
		compeition,
		number,
		name
	FROM teams
	WHERE teamid = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowx(id)
	if row.Err() != nil {
		return nil, err
	}

	team := Team{db: db}
	err = row.StructScan(&team)
	if err != nil {
		return nil, err
	}

	return &team, nil

}
