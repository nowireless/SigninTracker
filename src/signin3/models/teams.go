package models

type Team struct {
	DatabaseID int    `meta:"readOnly"`
	URI        string `meta:"readOnly"`

	Competition *string `meta:"requiredOnCreate"`
	Number      *int    `meta:"requiredOnCreate"`
	Name        *string `meta:"requiredOnCreate"`

	Mentors  Link `meta:"readOnly"`
	Students Link `meta:"readOnly"`
	Meetings Link `meta:"readOnly"`
}

func (t Team) GetDatabaseID() int {
	return t.DatabaseID
}
