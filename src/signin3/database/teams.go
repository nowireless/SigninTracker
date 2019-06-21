package database

import "signin3/models"

func (db *Database) GetTeams() ([]models.Team, error) {
	panic("TODO")
}

func (db *Database) GetTeam(id int) (*models.Team, error) {
	panic("TODO")
}

func (db *Database) CreateTeam(mTeam *models.Team) error {
	panic("TODO")
}

func (db *Database) UpdateTeam(model models.Team) error {
	panic("TODO")
}

func (db *Database) DeleteTeam(models *models.Team) error {
	panic("TODO")
}
