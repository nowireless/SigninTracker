package models

type Person struct {
	URI        string
	DatabaseID int

	CheckInID *string
	FirstName *string
	LastName  *string
	Phone     *string
	Email     *string

	Student *Student
	Mentor  *Mentor

	ParentOf []Link
	// Parents  []ParentRelation
	Parents []Link // TODO use ParentRelation

	Teams      []Link
	Committed  []Link
	Attendance Link // Link to MeetingAttendance
}

type Student struct {
	SchoolEmail *string
	SchoolID    *string
	Teams       []Link
}

type Mentor struct {
	Teams []Link
}

// TODO: Include this structure in people
type ParentRelation struct {
	Relation string
	Parent   Link
}
