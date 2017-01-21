package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/aaronarduino/weekly-schedule/datastore"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

// TODO: create time block struct

type siteData struct {
	User     string
	Saturday datastore.Day
}

const (
	tmpFolder    string = "templates/"
	staticFolder string = "static/"
	port                = ":8080"
)

var (
	siteSessionName = "this"
	ErrEmtyValue    = errors.New("EmtyValue")
	globalDb        *bolt.DB
)

func main() {
	log.Println("Starting...")

	var err error
	globalDb, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Println(err)
	}
	defer globalDb.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", dashboard).Methods("GET")
	router.HandleFunc("/add", addEvent).Methods("GET")
	router.PathPrefix("/static").Handler(
		http.StripPrefix("/static",
			http.FileServer(http.Dir("static"))))

	//middleware := alice.New(requireLoginMiddleware, roleMiddleware).Then(router)

	log.Println("Started. Go to http://localhost:8080/ for site.")
	log.Fatal(http.ListenAndServe(port, router))
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to / request")
	sd := siteData{}

	var err error
	sd.Saturday, err = datastore.EventsForDay(globalDb, "2017", "1", "7")
	if err != nil {
		log.Println(err)
	}

	showPage("dashboard.html", sd, w, r)
}

// TODO: Create CRUD routes for time blocks

func addEvent(w http.ResponseWriter, r *http.Request) {
	// this is just a route to add a test event to the db
	// event is init'd below
	testEvent := datastore.Event{
		Name:      "Test Event",
		StartTime: "13:00",
		EndTime:   "14:00",
	}
	err := datastore.EventToDb(globalDb, "2017", "1", "7", testEvent)
	if err != nil {
		log.Println(err)
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "Test event added to db :)")
}
