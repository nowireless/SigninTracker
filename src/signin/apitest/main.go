package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"signin/database"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/shiena/ansicolor"
)

type Link struct {
	URI string `json:"@uri"`
}

type People struct {
	Members []Link
}

type ParentRelation struct {
	Relation database.Relation
	Link     Link
}

type Person struct {
	PersonID int

	CheckInID string

	Name struct {
		First string
		Last  string
	}
	Email       string
	PhoneNumber string
	SchoolID    string
	SchoolEmail string

	MentorOf  []Link
	StudentOf []Link
	ParentOf  []Link
	Parents   []ParentRelation
}

type Teams struct {
	Members []Team
}

type Team struct {
	TeamID int

	Competition string
	Number      int
	Name        string

	Students []Link
	Mentors  []Link
}

type Meetings struct {
	Members []Link
}

type SignIn struct {
	InTime string
	Person Link
}

type SignOut struct {
	OutTime string
	Person  Link
}

type Meeting struct {
	MeetingID int

	// TODO Should this just be a string?
	Date struct {
		Year  int
		Month int
		Day   int
	}
	// TODO Should this just be a string?
	StartTime struct {
		Hour   int
		Minute int
		Second int
	}
	// TODO Should this just be a string?
	EndTime struct {
		Hour   int
		Minute int
		Second int
	}
	Location string

	Teams     []Link
	Committed []Link
	// Use a map?
	SignedIn  []SignIn
	SignedOut []SignOut
}

func main() {
	if runtime.GOOS == "windows" {
		// Setup logging for windows
		log.SetFormatter(&log.TextFormatter{ForceColors: true})
		log.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
	}

	config := database.Config{
		User:     "signin",
		Password: "foobar",
		Host:     "localhost",
		Port:     5432,
		Database: "signin",
	}
	db, err := database.Connect(config)
	if err != nil {
		panic(err)
	}

	db.DB.Ping()

	r := mux.NewRouter()

	// Serve Client files
	// TODO Make client dir configurable
	clientDir := "./client"
	r.PathPrefix("/client/").Handler(http.StripPrefix("/client/", http.FileServer(http.Dir(clientDir))))

	/*
	 * Middleware Setup
	 */
	// Cross-Origin Resource Sharing middleware Setup
	r.Use(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	))

	// Logging
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Request ", r.Method, r.RequestURI, " from ", r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	})

	r.Handle("/api", timer(func(w http.ResponseWriter, req *http.Request) {
		result := map[string]interface{}{}
		result["Teams"] = Link{URI: "/api/teams"}
		result["Meetings"] = Link{URI: "/api/Meetings"}
		result["People"] = Link{URI: "/api/people"}

		body, _ := json.MarshalIndent(result, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	})).Methods("GET")

	r.Handle("/api/people", timer(func(w http.ResponseWriter, req *http.Request) {
		people, err := db.GetAllPeople()
		if err != nil {
			panic(err)
		}

		result := People{}
		for _, person := range people {
			result.Members = append(result.Members, Link{
				URI: fmt.Sprintf("/api/people/%d", person.PersonID),
			})
		}

		body, _ := json.MarshalIndent(result, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	})).Methods("GET")

	r.Handle("/api/people/{id}", timer(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		personIDRaw := vars["id"]
		personID, err := strconv.ParseInt(personIDRaw, 10, 64)
		if err != nil {
			panic(err)
		}

		person, err := db.GetPerson(int(personID))
		if err != nil {
			panic(err)
		}

		result := Person{}
		result.CheckInID = person.CheckInID
		result.PersonID = person.PersonID
		result.Name.First = person.FirstName
		result.Name.Last = person.LastName
		if person.Email.Valid {
			result.Email = person.Email.String
		}
		if person.Phone.Valid {
			result.PhoneNumber = person.Phone.String
		}
		if person.SchoolID.Valid {
			result.SchoolID = person.SchoolID.String
		}
		if person.SchoolEmail.Valid {
			result.SchoolEmail = person.SchoolEmail.String
		}

		teams, err := person.MentorOf()
		if err != nil {
			panic(err)
		}
		for _, team := range teams {
			link := Link{}
			link.URI = fmt.Sprintf("/api/teams/%d", team.TeamID)
			result.MentorOf = append(result.MentorOf, link)
		}

		teams, err = person.StudentOf()
		if err != nil {
			panic(err)
		}
		for _, team := range teams {
			link := Link{}
			link.URI = fmt.Sprintf("/api/teams/%d", team.TeamID)
			result.StudentOf = append(result.StudentOf, link)
		}

		children, err := person.ParentOf()
		if err != nil {
			panic(err)
		}
		for _, child := range children {
			link := Link{}
			link.URI = fmt.Sprintf("/api/person/%d", child.PersonID)
			result.ParentOf = append(result.ParentOf, link)
		}

		parents, err := person.Parents()
		if err != nil {
			panic(err)
		}
		for _, parent := range parents {
			relation := ParentRelation{}
			relation.Relation = parent.Relation
			relation.Link = Link{}
			relation.Link.URI = fmt.Sprintf("/api/person/%d", parent.Parent.PersonID)
			result.Parents = append(result.Parents, relation)
		}

		body, _ := json.MarshalIndent(result, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	})).Methods("GET")

	r.Handle("/api/teams", timer(func(w http.ResponseWriter, req *http.Request) {
		teams, err := db.GetAllTeams()
		if err != nil {
			panic(err)
		}

		result := People{}
		for _, team := range teams {
			result.Members = append(result.Members, Link{
				URI: fmt.Sprintf("/api/teams/%d", team.TeamID),
			})
		}

		body, _ := json.MarshalIndent(result, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	})).Methods("GET")

	r.Handle("/api/teams/{id}", timer(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		teamIdRaw := vars["id"]
		teamId, err := strconv.ParseInt(teamIdRaw, 10, 64)
		if err != nil {
			panic(err)
		}

		team, err := db.GetTeam(int(teamId))
		if err != nil {
			panic(err)
		}

		result := Team{}
		result.TeamID = team.TeamID
		result.Competition = team.Competition
		result.Number = team.Number
		result.Name = team.Name

		students, err := team.Students()
		if err != nil {
			panic(err)
		}
		for _, student := range students {
			link := Link{}
			link.URI = fmt.Sprintf("/api/people/%d", student.PersonID)
			result.Students = append(result.Students, link)
		}

		mentors, err := team.Mentors()
		if err != nil {
			panic(err)
		}
		for _, mentor := range mentors {
			link := Link{}
			link.URI = fmt.Sprintf("/api/people/%d", mentor.PersonID)
			result.Mentors = append(result.Mentors, link)
		}

		body, _ := json.MarshalIndent(result, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	})).Methods("GET")

	r.Handle("/api/meetings", timer(func(w http.ResponseWriter, req *http.Request) {
		meetings, err := db.GetAllMeetings()
		if err != nil {
			panic(err)
		}

		result := Meetings{}
		for _, meeting := range meetings {
			result.Members = append(result.Members, Link{
				URI: fmt.Sprintf("/api/meetings/%d", meeting.MeetingID),
			})
		}

		body, _ := json.MarshalIndent(result, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	})).Methods("GET")

	r.Handle("/api/meetings/{id}", timer(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		meetingIDRaw := vars["id"]
		meetingID, err := strconv.ParseInt(meetingIDRaw, 10, 64)
		if err != nil {
			panic(err)
		}

		meeting, err := db.GetMeeting(int(meetingID))
		if err != nil {
			panic(err)
		}

		result := Meeting{}
		result.MeetingID = meeting.MeetingID
		result.Location = meeting.Location
		result.Date.Year = meeting.Date.Year()
		result.Date.Month = int(meeting.Date.Month())
		result.Date.Day = meeting.Date.Day()
		startTime := strings.Split(meeting.StartTime, ":")
		if len(startTime) != 3 {
			panic("Error parsing start time: " + meeting.StartTime)
		}
		hour, err := strconv.ParseInt(startTime[0], 10, 32)
		if err != nil {
			panic(err)
		}
		minute, err := strconv.ParseInt(startTime[1], 10, 32)
		if err != nil {
			panic(err)
		}
		Second, err := strconv.ParseInt(startTime[2], 10, 32)
		if err != nil {
			panic(err)
		}
		result.StartTime.Hour = int(hour)
		result.StartTime.Minute = int(minute)
		result.StartTime.Second = int(Second)

		endTime := strings.Split(meeting.EndTime, ":")
		if len(startTime) != 3 {
			panic("Error parsing end time: " + meeting.StartTime)
		}
		hour, err = strconv.ParseInt(endTime[0], 10, 32)
		if err != nil {
			panic(err)
		}
		minute, err = strconv.ParseInt(endTime[1], 10, 32)
		if err != nil {
			panic(err)
		}
		Second, err = strconv.ParseInt(endTime[2], 10, 32)
		if err != nil {
			panic(err)
		}
		result.EndTime.Hour = int(hour)
		result.EndTime.Minute = int(minute)
		result.EndTime.Second = int(Second)

		teams, err := meeting.Teams()
		if err != nil {
			panic(err)
		}
		for _, team := range teams {
			link := Link{}
			link.URI = fmt.Sprintf("/api/teams/%d", team.TeamID)
			result.Teams = append(result.Teams, link)
		}

		committed, err := meeting.Committed()
		if err != nil {
			panic(err)
		}
		for _, person := range committed {
			link := Link{}
			link.URI = fmt.Sprintf("/api/people/%d", person.PersonID)
			result.Committed = append(result.Committed, link)
		}

		signedIn, err := meeting.SignedIn()
		if err != err {
			panic(err)
		}
		for _, signIn := range signedIn {
			s := SignIn{}
			s.InTime = signIn.InTime
			s.Person.URI = fmt.Sprintf("/api/people/%d", signIn.PersonID)
			result.SignedIn = append(result.SignedIn, s)
		}

		signedOut, err := meeting.SignedOut()
		if err != err {
			panic(err)
		}
		for _, signOut := range signedOut {
			s := SignOut{}
			s.OutTime = signOut.OutTime
			s.Person.URI = fmt.Sprintf("/api/people/%d", signOut.PersonID)
			result.SignedOut = append(result.SignedOut, s)
		}

		body, _ := json.MarshalIndent(result, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	})).Methods("GET")

	log.Info("Listening on :8081")
	http.ListenAndServe(":8081", r)
}

func timer(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		http.HandlerFunc(f).ServeHTTP(w, r)
		duration := time.Now().Sub(startTime)

		log.Info("Request ", r.Method, r.RequestURI, " from ", r.RemoteAddr, " took ", duration)
	})
}
