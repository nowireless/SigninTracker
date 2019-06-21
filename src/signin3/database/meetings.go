package database

import "signin3/models"

func (db *Database) GetMeetings() ([]models.Meeting, error) {
	panic("TODO")
}

func (db *Database) GetMeeting(id int) (*models.Meeting, error) {
	panic("TODO")
}

func (db *Database) CreateMeeting(mMeeting *models.Meeting) error {
	panic("TODO")
}

func (db *Database) UpdateMeeting(model models.Meeting) error {
	panic("TODO")
}

func (db *Database) DeleteMeeting(models *models.Meeting) error {
	panic("TODO")
}
