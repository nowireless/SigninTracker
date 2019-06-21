package models

type Team struct {
	DatabaseID int
	URI        string

	Competition string
	Number      int
	Name        string
}

func (t Team) GetDatabaseID() int {
	return t.DatabaseID
}
