package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "enter data with the event title and description only")
	}

	json.Unmarshal(reqBody, &newEvent)
	newEvent.ID = strconv.Itoa(len(events) + 1)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	eid, err := strconv.Atoi(eventID)
	if err != nil {
		maxId := len(events)
		fmt.Fprintf(w, "please enter a valid ID between %v and %v", 1, maxId)
		return
	}

	json.NewEncoder(w).Encode(events[eid-1])
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	inputEid, err := strconv.Atoi(eventID)
	if err != nil {
		maxId := len(events)
		fmt.Fprintf(w, "please enter a valid ID between %v and %v", 1, maxId)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "enter data with event title and description only")
		return
	}

	var updatedEvent event
	json.Unmarshal(reqBody, &updatedEvent)

	e := events[inputEid-1]
	e.Title = updatedEvent.Title
	e.Description = updatedEvent.Description
	events[inputEid-1] = e

	json.NewEncoder(w).Encode(e)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	eId, err := strconv.Atoi(eventID)
	if err != nil {
		fmt.Fprintf(w, "enter a valid id")
		return
	}

	for i, e := range events {
		if i >= eId {
			e.ID = strconv.Itoa(i)
			events[i] = e
		}
	}

	events = append(events[:eId-1], events[eId:]...)
	fmt.Fprintf(w, "deleted event with ID %v", eventID)
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
