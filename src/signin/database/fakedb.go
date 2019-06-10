package database

import (
	"signin/models"
	"sync"
)

type FakeDB struct {
	students map[string]models.Student
	lock     *sync.Mutex
}

func NewFakeDB() *FakeDB {
	db := &FakeDB{
		students: make(map[string]models.Student),
		lock:     &sync.Mutex{},
	}

	db.students["test"] = models.Student{
		ID:             "test",
		FirstName:      "John",
		LastName:       "Doe",
		GraduationYear: 2015,
		SchoolID:       "doe123",
	}

	return db
}

func (db *FakeDB) GetStudents() ([]models.Student, error) {
	result := []models.Student{}
	for _, student := range db.students {
		result = append(result, student)
	}

	return result, nil
}

func (db *FakeDB) GetStudentByID(id string) (models.Student, error) {
	if student, ok := db.students[id]; ok {
		return student, nil
	}

	return models.Student{}, ErrorIDNotExist
}

func (db *FakeDB) CreateStudent(student models.Student) error {
	if _, present := db.students[student.ID]; present {
		return ErrorIDExists
	}

	db.students[student.ID] = student

	return nil
}

func (db *FakeDB) UpdateStudent(student models.Student) error {
	if _, present := db.students[student.ID]; present {
		db.students[student.ID] = student
		return nil
	}

	return ErrorIDNotExist

}

func (db *FakeDB) DeleteStudent(student models.Student) error {
	return db.DeleteStudentByID(student.ID)
}

func (db *FakeDB) DeleteStudentByID(id string) error {
	if _, present := db.students[id]; present {
		delete(db.students, id)
		return nil
	}

	return ErrorIDNotExist
}

func (db *FakeDB) Lock() {
	db.lock.Lock()
}

func (db *FakeDB) UnLock() {
	db.lock.Unlock()
}

func (db *FakeDB) Close() {

}
