package models

type Person struct {
	// Read Only
	URI        string `meta:"readOnly"`
	DatabaseID int    `meta:"readOnly"`

	// Read/Write
	CheckInID *string `meta:"requiredOnCreate"`
	FirstName *string `meta:"requiredOnCreate"`
	LastName  *string `meta:"requiredOnCreate"`
	Phone     *string
	Email     *string

	Student *Student
	Mentor  *Mentor

	ParentOf []Link `meta:"readOnly"`
	Parents  []Link `meta:"readOnly"`

	Attendance Link `meta:"readOnly"` // Link to MeetingAttendance
}

func (p Person) GetDatabaseID() int {
	return p.DatabaseID
}

type Student struct {
	SchoolEmail *string
	SchoolID    *string
	Teams       []Link `meta:"readOnly"`
}

type Mentor struct {
	Teams []Link `meta:"readOnly"`
}
