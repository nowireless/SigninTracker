package database

import (
	"signin3/database/internal"
	"signin3/models"

	log "github.com/sirupsen/logrus"
)

func (db *Database) GetMeetings() ([]models.Meeting, error) {
	rows, err := db.DB.Queryx(`
		SELECT
			meetings.meetingid,
			meetings.date,
			meetings.starttime,
			meetings.endtime,
			meetings.location,
			array_agg(DISTINCT c2.personid) FILTER (WHERE c2.personid IS NOT NULL) as commitied
		FROM meetings
			FULL OUTER JOIN commitments c2 on meetings.meetingid = c2.meetingid
		GROUP BY meetings.meetingid;
	`)

	if err != nil {
		return nil, err
	}

	meetings := []models.Meeting{}
	for rows.Next() {
		meeting := internal.Meeting{}
		rows.StructScan(&meeting)
		meetings = append(meetings, meeting.Model())
	}

	return meetings, nil
}

func (db *Database) GetMeeting(id int) (*models.Meeting, error) {
	row := db.DB.QueryRowx(`
		SELECT
			meetings.meetingid,
			meetings.date,
			meetings.starttime,
			meetings.endtime,
			meetings.location,
			array_agg(DISTINCT c2.personid) FILTER (WHERE c2.personid IS NOT NULL) as commitied
		FROM meetings
			FULL OUTER JOIN commitments c2 on meetings.meetingid = c2.meetingid
		WHERE people.meetingid = $1
		GROUP BY meetings.meetingid;
	`, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	meeting := internal.Meeting{}
	err := row.StructScan(&meeting)
	if err != nil {
		return nil, err
	}
	m := meeting.Model()
	return &m, nil
}

func (db *Database) CreateMeeting(mMeeting *models.Meeting) error {
	meeting := internal.NewMeeting(mMeeting)

	row := db.DB.QueryRowx(`
		INSERT INTO meetings(date, starttime, endtime, location)
		VALUES ($1, $2, $3, $4);
	`, meeting.Date, meeting.StartTime, meeting.EndTime, meeting.Location)

	if err := row.Err(); err != nil {
		return err
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	mMeeting.DatabaseID = id

	return nil
}

func (db *Database) UpdateMeeting(model models.Meeting) error {
	meeting := internal.NewMeeting(&model)

	log.Info("Updating meeting with ID: ", meeting.MeetingID)

	_, err := db.DB.Exec(`
		UPDATE meetings SET
			date      = $1,
			starttime = $2,
			endtime   = $3,
			location  = $4,
		WHERE meetingid = $5;
		`,
		meeting.Date,
		meeting.StartTime,
		meeting.EndTime,
		meeting.Location,
	)

	return err
}

func (db *Database) DeleteMeeting(m *models.Meeting) error {
	// TODO: Some how mark/(remove database id) to indicate that the model no longer represents a row in the db
	_, err := db.DB.Exec("DELETE FROM meeting WHERE meetingid = $1", m.DatabaseID)
	return err
}
