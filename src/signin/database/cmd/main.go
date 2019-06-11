package main

import (
	"fmt"
	"signin/database"

	"github.com/davecgh/go-spew/spew"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func main() {
	config := Config{
		User:     "signin",
		Password: "foobar",
		Host:     "localhost",
		Port:     5432,
		Database: "signin",
	}

	// See: https://github.com/jackc/pgx/blob/master/stdlib/sql.go
	connectionStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d database=%s sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	fmt.Println("Connection String: ", connectionStr)

	db, err := sqlx.Connect("pgx", connectionStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Pinging Database")
	db.Ping()

	database := &database.Database{DB: db}

	p, err := database.GetPerson(0)
	if err != nil {
		panic(err)
	}
	spew.Dump(p)

	teams, err := p.MentorOf()
	if err != nil {
		panic(err)
	}
	spew.Dump(teams)

	p, err = database.GetPerson(8)
	if err != nil {
		panic(err)
	}
	spew.Dump(p)

	teams, err = p.StudentOf()
	if err != nil {
		panic(err)
	}
	spew.Dump(teams)

	fmt.Println("Rick")
	p, err = database.GetPerson(5)
	if err != nil {
		panic(err)
	}
	spew.Dump(p)

	children, err := p.ParentOf()
	if err != nil {
		panic(err)
	}
	spew.Dump(children)

	fmt.Println("Trey")
	p, err = database.GetPerson(4)
	if err != nil {
		panic(err)
	}
	spew.Dump(p)

	parents, err := p.Parents()
	if err != nil {
		panic(err)
	}
	spew.Dump(parents)

	m, err := database.GetMeeting(0)
	if err != nil {
		panic(err)
	}
	spew.Dump(m)

	teams, err = m.Teams()
	if err != nil {
		panic(err)
	}
	spew.Dump(teams)

	people, err := m.Commited()
	if err != nil {
		panic(err)
	}
	spew.Dump(people)

	signedIn, err := m.SignedIn()
	spew.Dump(signedIn)
	for _, signin := range signedIn {
		fmt.Println(signin.MeetingID, signin.PersonID, signin.InTime)
		p, err := signin.Person()
		if err != nil {
			panic(err)
		}
		fmt.Println("Person:", p.FirstName, p.LastName)
		m, err := signin.Meeting()
		if err != nil {
			panic(err)
		}
		fmt.Println("Meeting:", m.Date, m.StartTime, m.EndTime)
	}

	signedOut, err := m.SignedOut()
	spew.Dump(signedOut)
	for _, signOut := range signedOut {
		fmt.Println(signOut.MeetingID, signOut.PersonID, signOut.OutTime)
		p, err := signOut.Person()
		if err != nil {
			panic(err)
		}
		fmt.Println("Person:", p.FirstName, p.LastName)
		m, err := signOut.Meeting()
		if err != nil {
			panic(err)
		}
		fmt.Println("Meeting:", m.Date, m.StartTime, m.EndTime)
	}

	team, err := database.GetTeam(1)
	if err != nil {
		panic(err)
	}
	spew.Dump(team)

	meetings, err := team.Meetings()
	if err != nil {
		panic(err)
	}
	spew.Dump(meetings)

	students, err := team.Students()
	if err != nil {
		panic(err)
	}
	spew.Dump(students)

	mentors, err := team.Mentors()
	if err != nil {
		panic(err)
	}
	spew.Dump(mentors)

}
