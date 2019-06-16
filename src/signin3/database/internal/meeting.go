package internal

import "database/sql"

type Meeting struct {
	MeetingID int

	Date      sql.NullString
	StartTime sql.NullString
	EndTime   sql.NullString
}
