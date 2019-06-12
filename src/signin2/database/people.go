package database

import (
	"database/sql"
	"fmt"
	"signin2/models"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/pgtype"
)

// personDB is the internal struture used to fetch data from the database
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

func NewDBPerson(p *models.Person) *models.Person {
	panic("TOOD")
}

func NewModelPerson(p *Person) *models.Person {
	result := models.Person{}
	result.DatabaseID = &p.PersonID
	result.Checkinid = &p.CheckInID
	result.Name = &models.PersonName{}
	result.Name.First = &p.FirstName
	result.Name.Last = &p.LastName

	if p.Email.Valid {
		value := strfmt.Email(p.Email.String)
		result.Email = &value
	}
	if p.Phone.Valid {
		result.Phone = &p.Phone.String
	}
	if p.SchoolEmail.Valid || p.SchoolID.Valid || len(p.StudentOf.Elements) > 0 {
		student := &models.PersonStudent{}
		if p.SchoolEmail.Valid {
			value := strfmt.Email(p.SchoolEmail.String)
			student.SchoolEmail = &value
		}
		if p.SchoolID.Valid {
			value := p.SchoolID.String
			student.SchoolID = &value
		}

		if len(p.StudentOf.Elements) > 0 {
			student.Teams = []*models.IDRef{}
			for _, element := range p.StudentOf.Elements {
				if element.Status != pgtype.Present {
					panic("Element not present")
				}

				idRef := &models.IDRef{}
				idRef.CollectionURI = "/signin/v1/api/teams"
				idRef.DatabaseID = int64(element.Int)
				idRef.MetaURI = models.ID(fmt.Sprintf("%s/%d", idRef.CollectionURI, idRef.DatabaseID))

				student.Teams = append(student.Teams, idRef)
			}
		}

		result.Student = student
	}

	if len(p.MentorOf.Elements) > 0 {
		result.MentorOf = []*models.IDRef{}
		for _, element := range p.MentorOf.Elements {
			if element.Status != pgtype.Present {
				panic("Element not present")
			}

			idRef := &models.IDRef{}
			idRef.CollectionURI = "/signin/v1/api/teams"
			idRef.DatabaseID = int64(element.Int)
			idRef.MetaURI = models.ID(fmt.Sprintf("%s/%d", idRef.CollectionURI, idRef.DatabaseID))

			result.MentorOf = append(result.MentorOf, idRef)
		}
	}

	if len(p.Parents.Elements) > 0 {
		result.Parents = []*models.IDRef{}
		for _, element := range p.Parents.Elements {
			if element.Status != pgtype.Present {
				panic("Element not present")
			}

			idRef := &models.IDRef{}
			idRef.CollectionURI = "/signin/v1/api/people"
			idRef.DatabaseID = int64(element.Int)
			idRef.MetaURI = models.ID(fmt.Sprintf("%s/%d", idRef.CollectionURI, idRef.DatabaseID))

			result.Parents = append(result.Parents, idRef)
		}
	}

	if len(p.ParentOf.Elements) > 0 {
		result.ParentOf = []*models.PersonParentOfItems0{}
		for _, element := range p.ParentOf.Elements {
			if element.Status != pgtype.Present {
				panic("Element not present")
			}

			idRef := &models.IDRef{}
			idRef.CollectionURI = "/signin/v1/api/people"
			idRef.DatabaseID = int64(element.Int)
			idRef.MetaURI = models.ID(fmt.Sprintf("%s/%d", idRef.CollectionURI, idRef.DatabaseID))

			parent := &models.PersonParentOfItems0{}
			parent.Student = idRef
			// TODO: Suppot this field
			parent.Relation = models.PersonParentOfItems0RelationGuardian

			result.ParentOf = append(result.ParentOf, parent)
		}
	}

	return &result
}
