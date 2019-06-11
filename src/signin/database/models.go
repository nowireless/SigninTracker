package database

import (
	"database/sql"
	"time"
)

// TODO: Need to decide what errors should cause an panic?
// TODO: Consider using Gorm?

type Person struct {
	db        *Database
	PersonID  int    //`db:"PersonID"`
	CheckInID string //`db:"CheckInID"`
	FirstName string //`db:"FirstName"`
	LastName  string //`db:"LastName"`

	Email       sql.NullString //`db:"Email"`
	Phone       sql.NullString //`db:"Phone"`
	SchoolEmail sql.NullString //`db:"SchoolEmail"`
	SchoolID    sql.NullString //`db:"SchoolID"`
}

func (p *Person) MentorOf() ([]Team, error) {
	// SELECT teams.teamid, compeition, number, name
	// FROM teams
	// INNER JOIN mentors m2 on teams.teamid = m2.teamid
	// WHERE m2.personid = 0;

	stmt, err := p.db.DB.Preparex(`
		SELECT teams.teamid, compeition, number, name
		FROM teams
		INNER JOIN mentors m2 on teams.teamid = m2.teamid
		WHERE m2.personid = $1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(p.PersonID)
	if err != nil {
		return nil, err
	}

	teams := []Team{}
	for rows.Next() {
		team := Team{db: p.db}
		err := rows.StructScan(&team)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (p *Person) StudentOf() ([]Team, error) {
	// SELECT teams.teamid, compeition, number, name
	// FROM teams
	// INNER JOIN students s2 on teams.teamid = s2.teamid
	// WHERE s2.personid = 0;

	stmt, err := p.db.DB.Preparex(`
		SELECT teams.teamid, compeition, number, name
		FROM teams
		INNER JOIN students s2 on teams.teamid = s2.teamid
		WHERE s2.personid = $1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(p.PersonID)
	if err != nil {
		return nil, err
	}

	teams := []Team{}
	for rows.Next() {
		team := Team{db: p.db}
		err := rows.StructScan(&team)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

type Relation string

const (
	RelationFather   Relation = "Father"
	RelationMother   Relation = "Mother"
	RelationGuardian Relation = "Guardian"
)

func (p *Person) ParentOf() ([]Person, error) {
	stmt, err := p.db.DB.Preparex(`	
		SELECT
			student.personid,
			student.checkinid,
			student.firstname,
			student.lastname,
			student.email,
			student.phone,
			student.schoolemail,
			student.schoolid
		FROM people as parent
			INNER JOIN parents p on parent.personid = p.parentid
			INNER JOIN students s on p.studentid = s.personid
			INNER JOIN people student on student.personid = s.personid
		WHERE parent.personid = $1
		ORDER BY student.lastname, student.firstname;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(p.PersonID)
	if err != nil {
		return nil, err
	}

	children := []Person{}
	for rows.Next() {
		child := Person{db: p.db}
		err := rows.StructScan(&child)
		if err != nil {
			return nil, err
		}

		children = append(children, child)
	}

	return children, nil
}

type Parent struct {
	Relation Relation
	Parent   *Person
}

func (p *Person) Parents() ([]Parent, error) {
	tx := p.db.DB.MustBegin()

	// Get parent IDS and relations
	stmt, err := tx.Preparex(`
		SELECT
			relation,
			parentid
		FROM parents
		WHERE studentid = $1;
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(p.PersonID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	type parentID struct {
		Relation Relation
		ParentID int
	}

	parentIDs := []parentID{}
	for rows.Next() {
		parentID := parentID{}
		err := rows.StructScan(&parentID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		parentIDs = append(parentIDs, parentID)
	}

	// Get Parents
	parents := []Parent{}
	for _, parentID := range parentIDs {
		parent := Parent{}
		parent.Relation = parentID.Relation
		parent.Parent, err = p.db.GetPersonTX(tx, parentID.ParentID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		parents = append(parents, parent)
	}

	tx.Commit()

	return parents, nil

}

type Meeting struct {
	db        *Database
	MeetingID int //`db:"MeetingID"`

	Date      time.Time //`db:"Date"`      // TODO use go type
	StartTime string    //`db:"StartTime"` // TODO use go type
	EndTime   string    //`db:"EndTime"`   // TODO use go type
	Location  string    //`db:"Location"`  // TODO use go type
}

func (m *Meeting) Teams() ([]Team, error) {
	stmt, err := m.db.DB.Preparex(`
		SELECT
			t.teamid,
			compeition,
			number,
			name
		FROM meetings
			INNER JOIN team_meetings m2 on meetings.meetingid = m2.meetingid
			INNER JOIN teams t on m2.teamid = t.teamid
		WHERE m2.meetingid = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(m.MeetingID)
	if err != nil {
		return nil, err
	}

	teams := []Team{}
	for rows.Next() {
		team := Team{db: m.db}
		err := rows.StructScan(&team)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (m *Meeting) Committed() ([]Person, error) {
	stmt, err := m.db.DB.Preparex(`	
		SELECT
			person.personid,
			person.checkinid,
			person.firstname,
			person.lastname,
			person.email,
			person.phone,
			person.schoolemail,
			person.schoolid
		FROM people as person
			INNER JOIN commitments c2 on person.personid = c2.personid
		WHERE c2.meetingid = $1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(m.MeetingID)
	if err != nil {
		return nil, err
	}

	people := []Person{}
	for rows.Next() {
		person := Person{db: m.db}
		err := rows.StructScan(&person)
		if err != nil {
			return nil, err
		}

		people = append(people, person)
	}

	return people, nil
}

// func (*Meeting) NoShows() []Person {

// }

func (m *Meeting) SignedIn() ([]SignIn, error) {
	stmt, err := m.db.DB.Preparex(`	
		SELECT
			personid,
			meetingid,
			intime
		FROM signed_in
	  	WHERE meetingid = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(m.MeetingID)
	if err != nil {
		return nil, err
	}

	signIns := []SignIn{}
	for rows.Next() {
		signIn := SignIn{db: m.db}
		err := rows.StructScan(&signIn)
		if err != nil {
			return nil, err
		}

		signIns = append(signIns, signIn)
	}

	return signIns, nil
}

func (m *Meeting) SignedOut() ([]SignOut, error) {
	stmt, err := m.db.DB.Preparex(`	
		SELECT
			personid,
			meetingid,
			outtime
		FROM signed_out
	  	WHERE meetingid = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(m.MeetingID)
	if err != nil {
		return nil, err
	}

	signOuts := []SignOut{}
	for rows.Next() {
		signOut := SignOut{db: m.db}
		err := rows.StructScan(&signOut)
		if err != nil {
			return nil, err
		}

		signOuts = append(signOuts, signOut)
	}

	return signOuts, nil
}

type Team struct {
	db     *Database
	TeamID int //`db:"TeamID"`

	Competition string `db:"compeition"`
	Number      int    //`db:"Number"`
	Name        string //`db:"Name"`
}

func (t *Team) Meetings() ([]Meeting, error) {
	stmt, err := t.db.DB.Preparex(`
		SELECT
			meeting.meetingid,
			meeting.date,
			meeting.starttime,
			meeting.endtime,
			meeting.location
		FROM meetings as meeting
			INNER JOIN team_meetings m2 on meeting.meetingid = m2.meetingid
			INNER JOIN teams t on m2.teamid = t.teamid
		WHERE t.teamid = $1
		ORDER BY date, starttime;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(t.TeamID)
	if err != nil {
		return nil, err
	}

	meetings := []Meeting{}
	for rows.Next() {
		meeting := Meeting{db: t.db}
		err := rows.StructScan(&meeting)
		if err != nil {
			return nil, err
		}

		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (t *Team) Students() ([]Person, error) {
	stmt, err := t.db.DB.Preparex(`
		SELECT
			people.personid,
			people.checkinid,
			people.firstname,
			people.lastname,
			people.email,
			people.schoolemail,
			people.schoolid
		FROM people
			INNER JOIN students s2 on people.personid = s2.personid
			INNER JOIN teams t on s2.teamid = t.teamid
		WHERE t.teamid = $1
		ORDER BY lastname, firstname;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(t.TeamID)
	if err != nil {
		return nil, err
	}

	people := []Person{}
	for rows.Next() {
		person := Person{db: t.db}
		err := rows.StructScan(&person)
		if err != nil {
			return nil, err
		}

		people = append(people, person)
	}

	return people, nil
}

func (t *Team) Mentors() ([]Person, error) {
	stmt, err := t.db.DB.Preparex(`
		SELECT
			people.personid,
			people.checkinid,
			people.firstname,
			people.lastname,
			people.email,
			people.schoolemail,
			people.schoolid
		FROM people
			INNER JOIN mentors m2 on people.personid = m2.personid
			INNER JOIN teams t on m2.teamid = t.teamid
		WHERE t.teamid = $1
		ORDER BY lastname, firstname;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Queryx(t.TeamID)
	if err != nil {
		return nil, err
	}

	people := []Person{}
	for rows.Next() {
		person := Person{db: t.db}
		err := rows.StructScan(&person)
		if err != nil {
			return nil, err
		}

		people = append(people, person)
	}

	return people, nil
}

type SignIn struct {
	db *Database

	PersonID  int    //`db:"PersonID"`
	MeetingID int    //`db:"MeetingID"`
	InTime    string //`db:"InTime"` // TODO use go type
}

func (s *SignIn) Person() (*Person, error) {
	return s.db.GetPerson(s.PersonID)
}

func (s *SignIn) Meeting() (*Meeting, error) {
	return s.db.GetMeeting(s.MeetingID)
}

type SignOut struct {
	db        *Database
	PersonID  int    //`db:"PersonID"`
	MeetingID int    //`db:"MeetingID"`
	OutTime   string //`db:"outTime"` // TODO use go type
}

func (s *SignOut) Person() (*Person, error) {
	return s.db.GetPerson(s.PersonID)
}

func (s *SignOut) Meeting() (*Meeting, error) {
	return s.db.GetMeeting(s.MeetingID)
}
