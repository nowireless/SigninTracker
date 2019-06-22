package internal

import (
	"database/sql"
	"fmt"
	"signin3/constants"
	"signin3/models"

	"github.com/jackc/pgx/pgtype"
)

type Meeting struct {
	MeetingID int64

	Date      string
	StartTime string
	EndTime   string

	Location sql.NullString

	Committed pgtype.Int4Array
}

func (m *Meeting) Model() models.Meeting {
	result := models.Meeting{}
	result.DatabaseID = int(m.MeetingID)
	result.URI = fmt.Sprintf("%s/%d", constants.MeetingsCollection, result.DatabaseID)

	result.Day = &m.Date
	result.StartTime = &m.StartTime
	result.EndTime = &m.EndTime
	result.Location = getNullString(m.Location)

	result.Committed = makeIDLinks(constants.PeopleCollection, m.Committed)
	result.SignedIn = models.Link{URI: fmt.Sprintf("%s/signedin", result.URI)}
	result.SignedOut = models.Link{URI: fmt.Sprintf("%s/signedout", result.URI)}
	result.Attendance = models.Link{URI: fmt.Sprintf("%s/attendance", result.URI)}
	result.Teams = models.Link{URI: fmt.Sprintf("%s/teams", result.URI)}

	return result
}

func NewMeeting(m *models.Meeting) Meeting {
	result := Meeting{}
	result.MeetingID = int64(m.DatabaseID)
	result.Date = *m.Day
	result.StartTime = *m.StartTime
	result.EndTime = *m.EndTime
	result.Location = setNullString(m.Location)
	return result
}
