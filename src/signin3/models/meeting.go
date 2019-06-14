package models

type Meeting struct {
	databaseID int

	Day       *string
	StartTime *string
	EndTime   *string

	Committed []Link
	SignedIn  []Link
	SingedOut []Link

	Teams TeamMeeting

	Attendance Link
}

type TeamMeeting struct {
	Kind string
	Team Link
}
