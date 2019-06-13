package models

type Time string

type Link struct {
	URI string
}

// SignIn represents when a person signs into a meeting
type SignIn struct {
	Meeting *Link
	Person  *Link
	InTime  string
}

type SignOut struct {
	Meeting *Link
	Person  *Link
	OutTime string
}
