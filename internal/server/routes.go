package server

import (
	"fmt"
	"log"
	"net/http"
	"pb/internal/domain"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// pages
	r.HandleFunc("/",s.handleIndex)
	r.HandleFunc("/add", s.addContactPage)

	// handlers
	r.HandleFunc("/contacts", s.handleContacts)
    r.HandleFunc("/contacts/{id}", s.handleContactsByID)

	return r 
}

// pages
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.templ.ExecuteTemplate(w, "base", nil)
}

func (s *Server) addContactPage(w http.ResponseWriter, r *http.Request) {
	s.templ.ExecuteTemplate(w, "contact-page", nil)
}

// handlers
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

	totalCount, err := s.db.Count()
	if err != nil {
		log.Println(err)
	}

	s.templ.ExecuteTemplate(w, "count", map[string]any{"Count": totalCount, "SwapOOB": true})
	s.templ.ExecuteTemplate(w, "list", contacts)
}

func (s *Server) getContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	log.Println("id", idParam)

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

	s.templ.ExecuteTemplate(w, "person", contact)
}

func (s *Server) postContact(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	tel := r.FormValue("tel")
	contact := domain.NewContact(name, tel)

	s.db.Store(contact)

	http.Redirect(w,r,"/add",http.StatusSeeOther)
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
		log.Println(id)
		log.Print("id delete or some shit")
	}

	totalCount, err := s.db.Count()
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w,r,"/",http.StatusSeeOther)
	s.templ.ExecuteTemplate(w, "count", map[string]any{"Count": totalCount, "SwapOOB": true})
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
