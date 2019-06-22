package database

import (
	"signin3/database/internal"
	"signin3/models"

	log "github.com/sirupsen/logrus"
)

func (db *Database) GetTeams() ([]models.Team, error) {
	rows, err := db.DB.Queryx(`
		SELECT teamid, compeition, number, name
		FROM teams
	`)
	if err != nil {
		return nil, err
	}

	teams := []models.Team{}
	for rows.Next() {
		team := internal.Team{}
		rows.StructScan(&team)
		teams = append(teams, team.Model())
	}

	return teams, nil
}

func (db *Database) GetTeam(id int) (*models.Team, error) {
	row := db.DB.QueryRowx(`
		SELECT teamid, compeition, number, name
		WHERE teamid = $1
		FROM teams
	`, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	team := internal.Team{}
	err := row.StructScan(&team)
	if err != nil {
		return nil, err
	}
	m := team.Model()
	return &m, nil
}

func (db *Database) CreateTeam(mTeam *models.Team) error {
	team := internal.NewTeam(mTeam)

	row := db.DB.QueryRowx(`
		INSERT INTO teams(compeition, number, name)
		VALUES ($1, $2, $3)`,
		team.Competition,
		team.Number,
		team.Name,
	)

	if err := row.Err(); err != nil {
		return err
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	mTeam.DatabaseID = id

	return nil
}

func (db *Database) UpdateTeam(model models.Team) error {
	team := internal.NewTeam(&model)

	log.Info("Updating team with ID: ", team.TeamID)

	_, err := db.DB.Exec(`
		UPDATE people SET
			compeition = $1,
			number     = $2,
			name       = $3,
		WHERE teamid   = $4;
		`,
		team.Competition,
		team.Number,
		team.Name,
		team.TeamID,
	)

	return err
}

func (db *Database) DeleteTeam(t *models.Team) error {
	// TODO: Some how mark/(remove database id) to indicate that the model no longer represents a row in the db
	_, err := db.DB.Exec("DELETE FROM teams WHERE teamid = $1", t.DatabaseID)
	return err
}
