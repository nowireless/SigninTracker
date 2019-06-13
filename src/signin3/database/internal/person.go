package internal

import (
	"database/sql"
	"fmt"
	"signin3/models"

	"github.com/jackc/pgx/pgtype"
)

// Person is the internal struture used to fetch data from the database
type Person struct {
	PersonID    int64
	CheckInID   string
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

func (p *Person) Model() models.Person {
	result := models.Person{}
	result.DatabaseID = int(p.PersonID)
	result.URI = fmt.Sprintf("/people/%d", result.DatabaseID)

	result.CheckInID = &p.CheckInID
	result.FirstName = &p.FirstName
	result.LastName = &p.LastName
	result.Email = getNullString(p.Email)
	result.Phone = getNullString(p.Phone)

	if p.SchoolEmail.Valid || p.SchoolID.Valid || len(p.StudentOf.Elements) > 0 {
		student := models.Student{}
		student.SchoolID = getNullString(p.SchoolID)
		student.SchoolEmail = getNullString(p.SchoolEmail)
		student.Teams = makeIDLinks("/teams", p.StudentOf)

		result.Student = &student
	}

	if len(p.MentorOf.Elements) > 0 {
		mentor := models.Mentor{}
		mentor.Teams = makeIDLinks("/teams", p.MentorOf)

		result.Mentor = &mentor
	}

	result.ParentOf = makeIDLinks("/people", p.ParentOf)
	result.Parents = makeIDLinks("/people", p.Parents)

	// Static link to array of
	result.Attendance = models.Link{URI: fmt.Sprintf("%s/attendance", result.URI)}

	return result
}
