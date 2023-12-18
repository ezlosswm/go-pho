package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pb/internal/domain"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// mux.HandleFunc("/", s.hello)
	r.HandleFunc("/",s.handleIndex)
	r.HandleFunc("/contacts", s.handleContacts)
    r.HandleFunc("/contacts/{id}", s.handleContactsByID)

	return r 
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.templ.ExecuteTemplate(w, "index.html", nil)
}

func (s *Server)handleContacts(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case "GET":
		s.getContacts(w,r)
	case "POST": 
		s.postContact(w,r)
	}
}

func (s *Server) handleContactsByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.getContact(w,r)
	case "DELETE": 
		s.deleteContact(w,r)
	}
}

func (s *Server) getContacts(w http.ResponseWriter, r *http.Request) {
	contacts, err := s.db.GetContacts()
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusFound)
	resp, _ := json.Marshal(contacts)
	_,_ = w.Write(resp)
}

func (s *Server) getContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Error: id %v unable to convert to int", idParam)
	 	http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
	 	return
 	}

	 contact, err := s.db.GetContact(id)
	 if err != nil {
		 log.Printf("No contact with id %v found", id)
		 http.Error(w, "Contact not found", http.StatusNotFound)
		 return
	}

	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(contact)
	_, _ = w.Write(resp)
}

func (s *Server) postContact(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	tel := r.FormValue("tel")
	contact := domain.NewContact(name, tel)

	s.db.Store(contact)

	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(contact)
	_, _ = w.Write(resp)
}

func (s *Server) deleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		log.Fatal("err: unable to get id", id)
	}

	// TODO: work on validating ID provided
	// Example: id 30 does not exists
	// What do??????
	if err := s.db.Delete(id); err != nil {
		log.Fatal("id delete or some shit")
	}


	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal("deleted")
	_, _ = w.Write(resp)
}

func getID(r *http.Request) (int,error) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, fmt.Errorf("Error: id %v unable to convert to int", idParam)
 	}
	
	return id, nil
}
