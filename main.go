package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO: create time block struct
type siteData struct {
	User     string
	Fullname string
	Name     string
	ObsID    string
}

const (
	tmpFolder    string = "templates/"
	staticFolder string = "static/"
	port                = ":8080"
)

var (
	siteSessionName = "this"
	ErrEmtyValue    = errors.New("EmtyValue")
)

func main() {
	log.Println("Starting...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", dashboard).Methods("GET")
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	//middleware := alice.New(requireLoginMiddleware, roleMiddleware).Then(router)

	log.Println("Started. Go to http://localhost:8080/ for site.")
	log.Fatal(http.ListenAndServe(port, router))
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	log.Println("Responding to / request")

	showPage("dashboard.html", siteData{}, w, r)
}

// TODO: Create CRUD routes for time blocks
