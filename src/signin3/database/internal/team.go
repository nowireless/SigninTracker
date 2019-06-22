package internal

import (
	"fmt"
	"signin3/models"
)

type Team struct {
	TeamID int64

	Competition string `db:"compeition"`
	Number      int
	Name        string
}

func (t *Team) Model() models.Team {
	result := models.Team{}
	result.DatabaseID = int(t.TeamID)
	result.URI = fmt.Sprintf("/teams/%d", result.DatabaseID)

	result.Competition = &t.Competition
	result.Number = &t.Number
	result.Name = &t.Name

	result.Mentors = models.Link{URI: fmt.Sprintf("%s/mentors", result.URI)}
	result.Students = models.Link{URI: fmt.Sprintf("%s/students", result.URI)}
	result.Meetings = models.Link{URI: fmt.Sprintf("%s/meetings", result.URI)}

	return result
}

func NewTeam(m *models.Team) Team {
	result := Team{}
	result.TeamID = int64(m.DatabaseID)
	result.Competition = *m.Competition
	result.Number = *m.Number
	result.Name = *m.Name

	return result
}
