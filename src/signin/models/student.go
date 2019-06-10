package models

type Student struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	GraduationYear int    `json:"graduationYear"`
	ID             string `json:"id"`
	SchoolID       string `json:"schoolID"`
}
