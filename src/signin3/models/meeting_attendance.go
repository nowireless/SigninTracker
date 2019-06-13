package models

// MeetingAttendance represents if a person was at an meeting, and from
// what times. Note: A person can have checked, but not checked out
type MeetingAttendance struct {
	PersonID  Link
	MeetingID Link
	InTime    string
	OutTime   *string
}
