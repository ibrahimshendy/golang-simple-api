package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Note struct {
	Title       string
	Description string
	CreatedOn   time.Time
}

var notes = make(map[string]Note)

var id = 0

func main() {
	r := mux.NewRouter().StrictSlash(false)

	r.HandleFunc("/api/v1/notes", listAllNotes).Methods("GET")
	r.HandleFunc("/api/v1/note/create", addNote).Methods("POST")
	r.HandleFunc("/api/v1/note/{id}", deleteNote).Methods("POST")
	r.HandleFunc("/api/v1/note/{id}", updateNote).Methods("PUT")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Server Listening On Localhost:8080")
	server.ListenAndServe()
}

func listAllNotes(w http.ResponseWriter, r *http.Request) {
	var list_notes []Note

	for _, v := range notes {
		list_notes = append(list_notes, v)
	}

	j, err := json.Marshal(list_notes)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func addNote(w http.ResponseWriter, r *http.Request) {
	var note Note

	err := json.NewDecoder(r.Body).Decode(&note)

	if err != nil {
		panic(err)
	}

	note.CreatedOn = time.Now()
	id++
	k := strconv.Itoa(id)
	notes[k] = note

	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]

	var note Note

	if _, ok := notes[k]; ok {
		err := json.NewDecoder(r.Body).Decode(&note)

		if err != nil {
			panic(err)
		}
		notes[k] = note
	} else {
		log.Printf("Could not find key of Note %s to delete", k)
		return
	}

	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	// Remove from Store
	if _, ok := notes[k]; ok {
		//delete existing item
		delete(notes, k)
	} else {
		log.Printf("Could not find key of Note %s to delete", k)
	}

	w.WriteHeader(http.StatusNoContent)
}
