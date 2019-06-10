package app

import (
	"signin/database"
	"signin/models"
)

type App struct {
	db database.Database
}

func New() (*App, error) {
	return &App{
		db: database.NewFakeDB(),
	}, nil
}

func (a *App) Close() {
	a.db.Close()
}

func (a *App) GetStudents() ([]models.Student, error) {
	a.db.Lock()
	defer a.db.UnLock()

	return a.db.GetStudents()
}

func (a *App) GetStudentByID(id string) (models.Student, error) {
	a.db.Lock()
	defer a.db.UnLock()

	return a.db.GetStudentByID(id)
}

func (a *App) CreateStudent(student models.Student) error {
	a.db.Lock()
	defer a.db.UnLock()

	return a.db.CreateStudent(student)
}

func (a *App) UpdateStudent(student models.Student) error {
	a.db.Lock()
	defer a.db.UnLock()

	return a.db.UpdateStudent(student)
}

func (a *App) DeleteStudentByID(id string) error {
	a.db.Lock()
	defer a.db.UnLock()

	return a.db.DeleteStudentByID(id)
}
