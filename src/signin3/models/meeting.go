package models

type Meeting struct {
	DatabaseID int    `meta:"readOnly"`
	URI        string `meta:"readOnly"`

	Day       *string `meta:"requiredOnCreate"`
	StartTime *string `meta:"requiredOnCreate"`
	EndTime   *string `meta:"requiredOnCreate"`
	Location  *string

	Committed []Link `meta:"readOnly"`
	SignedIn  Link   `meta:"readOnly"`
	SignedOut Link   `meta:"readOnly"`

	Teams Link `meta:"readOnly"`

	Attendance Link `meta:"readOnly"` // Link to MeetingAttendance
}

func (m Meeting) GetDatabaseID() int {
	return m.DatabaseID
}

type TeamMeeting struct {
	Kind string
	Team Link
}
