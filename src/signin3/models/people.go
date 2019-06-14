package models

type Person struct {
	// Read Only
	URI        string `meta:"readonly"`
	DatabaseID int    `meta:"readonly"`

	// Read/Write
	CheckInID *string `meta:"requiredOnCreate"`
	FirstName *string `meta:"requiredOnCreate"`
	LastName  *string `meta:"requiredOnCreate"`
	Phone     *string
	Email     *string

	Student *Student
	Mentor  *Mentor

	ParentOf []Link `meta:"readonly"`
	Parents  []Link `meta:"readonly"`

	Attendance Link `meta:"readonly"` // Link to MeetingAttendance
}

type Student struct {
	SchoolEmail *string
	SchoolID    *string
	Teams       []Link `meta:"readonly"`
}

type Mentor struct {
	Teams []Link `meta:"readonly"`
}
