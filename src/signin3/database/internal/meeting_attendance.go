package internal

import (
	"database/sql"
	"signin3/models"
)

// MeetingAttendance represents if a person was at an meeting, and from
// what times. Note: A person can have checked, but not checked out
type MeetingAttendance struct {
	PersonID  int
	MeetingID int
	InTime    string
	OutTime   sql.NullString
}

func (ma *MeetingAttendance) Model() *models.MeetingAttendance {
	result := models.MeetingAttendance{}
	result.PersonID = makeLink("/people", ma.PersonID)
	result.MeetingID = makeLink("/meetings", ma.MeetingID)
	result.InTime = ma.InTime
	result.OutTime = getNullString(ma.OutTime)

	return &result
}
