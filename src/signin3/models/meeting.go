package models

type Meeting struct {
	databaseID int

	Day       *string
	StartTime *string
	EndTime   *string
}
