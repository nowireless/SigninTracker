package database

import (
	"errors"
	"signin/models"
)

var ErrorIDNotExist = errors.New("ID doesn't exist")
var ErrorIDExists = errors.New("ID already exist")

type Database interface {
	GetStudents() ([]models.Student, error)
	GetStudentByID(id string) (models.Student, error)

	CreateStudent(student models.Student) error
	UpdateStudent(student models.Student) error
	DeleteStudent(student models.Student) error
	DeleteStudentByID(id string) error

	Lock()
	UnLock()

	Close()
}
